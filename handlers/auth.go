package handlers

import (
	"crypto/rand"
	"elf/internal/auth"
	"elf/internal/config"
	"elf/internal/core"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

func init() {
	gob.Register(core.User{})
}

type UserCreator interface {
	Create(sub string, name string) (core.User, error)
}

type AuthServices struct {
	SecureCookie  *securecookie.SecureCookie
	Sessions      sessions.Store
	Authenticator *auth.Authenticator

	Users UserCreator
}

func Login(cfg *config.Config, srvcs *AuthServices) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		state, err := generateState(cfg.OAuth.StateLength)
		if err != nil {
			return err
		}

		encodedStateValue, err := srvcs.SecureCookie.Encode(cfg.OAuth.StateCookieName, state)
		if err != nil {
			return err
		}
		c := &http.Cookie{
			Name:    cfg.OAuth.StateCookieName,
			Value:   encodedStateValue,
			Path:    "/",
			Expires: time.Now().Add(time.Minute * 10),
			// Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, c)

		http.Redirect(w, r, srvcs.Authenticator.Config.AuthCodeURL(state), http.StatusTemporaryRedirect)
		return nil
	}
}

func LoginCallback(cfg *config.Config, srvcs *AuthServices) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		es := map[Field]FieldError{}
		if !r.URL.Query().Has("state") {
			es["state"] = FieldError{
				Location: QUERY_PARAM_LOCATION,
				Value:    "",
				Reason:   "must be set",
			}
		}
		if !r.URL.Query().Has("code") {
			es["code"] = FieldError{
				Location: QUERY_PARAM_LOCATION,
				Value:    "",
				Reason:   "must be set",
			}
		}
		if len(es) > 0 {
			return ValidationErrors(es)
		}
		state := r.URL.Query().Get("state")
		code := r.URL.Query().Get("code")

		expectedStateCookie, err := r.Cookie(cfg.OAuth.StateCookieName)
		if err != nil {
			return fmt.Errorf("%s not present because: %w", cfg.OAuth.StateCookieName, err)
		}
		if err = expectedStateCookie.Valid(); err != nil {
			return err
		}

		var expectedState string
		if err = srvcs.SecureCookie.Decode(cfg.OAuth.StateCookieName, expectedStateCookie.Value, &expectedState); err != nil {
			return err
		}
		if state != expectedState {
			return errors.New("The provided state does not match the expected state.")
		}

		token, err := srvcs.Authenticator.Config.Exchange(r.Context(), code)
		if err != nil {
			return err
		}

		idToken, err := srvcs.Authenticator.VerifyIDToken(r.Context(), token)
		if err != nil {
			return err
		}

		var p auth.Profile
		if err := idToken.Claims(&p); err != nil {
			return err
		}

		user, err := srvcs.Users.Create(p.Sub, p.Name)
		if err != nil {
			return err
		}

		session, _ := srvcs.Sessions.Get(r, cfg.Auth.SessionCookieName)
		session.Values[cfg.Auth0.SessionCookieAccessTokenKey] = token.AccessToken
		session.Values[cfg.Auth.SessionCookieUserKey] = user
		if err = session.Save(r, w); err != nil {
			return err
		}

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return nil
	}
}

func generateState(n int) (string, error) {
	data := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func Logout(cfg *config.Config, srvcs *AuthServices) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		http.Redirect(w, r, srvcs.Authenticator.LogoutUrl(cfg.Auth0.LogoutCallbackUrl), http.StatusMovedPermanently)
		return nil
	}
}

func LogoutCallback(cfg *config.Config, srvcs *AuthServices) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		session, _ := srvcs.Sessions.Get(r, cfg.Auth.SessionCookieName)
		delete(session.Values, cfg.Auth.SessionCookieName)
		err := session.Save(r, w)
		if err != nil {
			return err
		}
		return nil
	}
}
