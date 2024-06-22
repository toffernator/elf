package rest

import (
	"context"
	"crypto/rand"
	"elf/internal/auth/auth0"
	"elf/internal/core"
	restcontext "elf/internal/rest_context"
	"encoding/base64"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	state, err := generateState(s.Config.OAuth.StateLength)
	if err != nil {
		return err
	}

	encodedStateValue, err := s.SecureCookies.Encode(s.Config.OAuth.StateCookieName, state)
	if err != nil {
		return err
	}
	c := &http.Cookie{
		Name:    s.Config.OAuth.StateCookieName,
		Value:   encodedStateValue,
		Path:    "/",
		Expires: time.Now().Add(time.Minute * 10),
		// TODO: What does it require to uncomment this line, certs?
		// Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, c)

	http.Redirect(w, r, s.Authenticator.Config.AuthCodeURL(state), http.StatusTemporaryRedirect)
	return nil
}

func generateState(n int) (string, error) {
	data := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) error {
	state, err := generateState(s.Config.OAuth.StateLength)
	if err != nil {
		return err
	}

	encodedStateValue, err := s.SecureCookies.Encode(s.Config.OAuth.StateCookieName, state)
	if err != nil {
		return err
	}
	c := &http.Cookie{
		Name:    s.Config.OAuth.StateCookieName,
		Value:   encodedStateValue,
		Path:    "/",
		Expires: time.Now().Add(time.Minute * 10),
		// Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, c)

	http.Redirect(w, r, s.Authenticator.Config.AuthCodeURL(state), http.StatusTemporaryRedirect)
	return nil
}

func (s *Server) LoginCallback(w http.ResponseWriter, r *http.Request) (err error) {
	var req LoginCallBackRequest
	err = s.decodeAndValidateLoginCallBackRequest(&req, r)
	if err != nil {
		return err
	}

	if req.State != req.ExpectedState {
		return NewAuthenticationError("OAuth State Mismatch")
	}

	token, err := s.Authenticator.Config.Exchange(r.Context(), req.Code)
	if err != nil {
		return err
	}

	idToken, err := s.Authenticator.VerifyIDToken(r.Context(), token)
	if err != nil {
		return err
	}

	var p auth0.Profile
	if err := idToken.Claims(&p); err != nil {
		return err
	}

	user, err := s.Users.Create(r.Context(), core.UserCreateParams{
		Name: p.Name,
	})
	if err != nil {
		return err
	}

	session, _ := s.Sessions.Get(r, s.Config.Auth.SessionCookieName)
	session.Values[s.Config.Auth0.SessionCookieAccessTokenKey] = token.AccessToken
	session.Values[s.Config.Auth.SessionCookieUserKey] = user
	if err = session.Save(r, w); err != nil {
		return err
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	return nil
}

type LoginCallBackRequest struct {
	Code          string `validate:"required" form:"code"`
	State         string `validate:"required" form:"state"`
	ExpectedState string `validate:"required"`
}

func (l *LoginCallBackRequest) Validate() (err error) {
	err = validate.Struct(&l)
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	return core.ValidationErrorsFromValidatorErrors(errs)
}

func (s *Server) decodeAndValidateLoginCallBackRequest(req *LoginCallBackRequest, r *http.Request) (err error) {
	err = Decode(req, r)
	if err != nil {
		return err
	}

	expectedStateCookie, err := r.Cookie(s.Config.OAuth.StateCookieName)
	if err != nil {
		return NewAuthenticationError("OAuth state mismatch")
	}
	if err = expectedStateCookie.Valid(); err != nil {
		return err
	}

	var expectedState string
	if err = s.SecureCookies.Decode(s.Config.OAuth.StateCookieName, expectedStateCookie.Value, &expectedState); err != nil {
		// TODO: DO I want a AuthenticationError or a DecodeError here?
		return err
	}

	req.ExpectedState = expectedState

	return nil
}

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) error {
	http.Redirect(w, r, s.Authenticator.LogoutUrl(s.Config.Auth0.LogoutCallbackUrl), http.StatusMovedPermanently)
	return nil
}

func (s *Server) LogoutCallback(w http.ResponseWriter, r *http.Request) error {
	session, _ := s.Sessions.Get(r, s.Config.Auth.SessionCookieName)
	delete(session.Values, s.Config.Auth.SessionCookieName)
	err := session.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) AddUserToContext(w http.ResponseWriter, r *http.Request, next http.Handler) error {
	if !s.Config.IsDevelop() {
		session, err := s.Sessions.Get(r, s.Config.Auth.SessionCookieName)
		if err != nil {
			return err
		}

		user, ok := session.Values[s.Config.Auth.SessionCookieUserKey].(core.User)
		if !ok {
			// TODO: better error
			return errors.New("Cannot cast the user in the session to an auth.AuthenticatedUser")
		}

		ctxWithUser := context.WithValue(r.Context(), restcontext.UserKey, user)
		rWithUser := r.WithContext(ctxWithUser)
		next.ServeHTTP(w, rWithUser)
	} else {
		user := core.User{
			Id:   1,
			Name: "Dev User",
		}
		ctxWithUser := context.WithValue(r.Context(), restcontext.UserKey, user)
		rWithUser := r.WithContext(ctxWithUser)
		slog.Info("AddUserToContext injected development user", "user", user)

		next.ServeHTTP(w, rWithUser)
	}

	return nil
}
