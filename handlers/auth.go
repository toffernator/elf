package handlers

import (
	"crypto/rand"
	"elf/internal/auth"
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
	gob.Register(auth.AuthenticatedUser{})
}

const (
	OAUTH_STATE_LENGTH                    = 16
	OAUTH_STATE_COOKIE_NAME               = "oauthstate"
	AUTH0_SESSION_NAME                    = "auth0"
	AUTH0_SESSION_USER_KEY         string = "user"
	AUTH0_SESSION_ACCESS_TOKEN_KEY        = "access_token"
)

var (
	OAUTH_STATE_COOKIE_EXPIRATION = time.Minute * 10
	ErrNoUserInSession            = errors.New(fmt.Sprintf("There is no value in the session '%s' assocated with the key '%s'.", AUTH0_SESSION_NAME, AUTH0_SESSION_USER_KEY))
)

func Login(authenticator *auth.Authenticator, secureCookies *securecookie.SecureCookie) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		state, err := generateState(OAUTH_STATE_LENGTH)
		if err != nil {
			return err
		}

		encodedStateValue, err := secureCookies.Encode(OAUTH_STATE_COOKIE_NAME, state)
		if err != nil {
			return err
		}
		c := &http.Cookie{
			Name:    OAUTH_STATE_COOKIE_NAME,
			Value:   encodedStateValue,
			Path:    "/",
			Expires: time.Now().Add(OAUTH_STATE_COOKIE_EXPIRATION),
			// Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, c)

		http.Redirect(w, r, authenticator.Config.AuthCodeURL(state), http.StatusTemporaryRedirect)
		return nil
	}
}

func LoginCallback(authenticator *auth.Authenticator, store sessions.Store, secureCookies *securecookie.SecureCookie, users auth.AuthenticatedUserStore) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		queryParameterErrors := map[string]string{}
		if !r.URL.Query().Has("state") {
			queryParameterErrors["state"] = ""
		}
		if !r.URL.Query().Has("code") {
			queryParameterErrors["code"] = ""
		}
		if len(queryParameterErrors) != 0 {
			return InvalidQueryParameters(queryParameterErrors)
		}
		state := r.URL.Query().Get("state")
		code := r.URL.Query().Get("code")

		expectedStateCookie, err := r.Cookie(OAUTH_STATE_COOKIE_NAME)
		if err != nil {
			return err
		}
		if err = expectedStateCookie.Valid(); err != nil {
			return err
		}

		var expectedState string
		if err = secureCookies.Decode(OAUTH_STATE_COOKIE_NAME, expectedStateCookie.Value, &expectedState); err != nil {
			return err
		}
		if state != expectedState {
			return errors.New("The provided state does not match the expected state.")
		}

		token, err := authenticator.Config.Exchange(r.Context(), code)
		if err != nil {
			return err
		}

		idToken, err := authenticator.VerifyIDToken(r.Context(), token)
		if err != nil {
			return err
		}

		var profile auth.Profile
		if err := idToken.Claims(&profile); err != nil {
			return err
		}

		user, err := users.Create(profile)
		if errors.Is(err, auth.ErrDuplicateSub) {
			return err
		}

		session, _ := store.Get(r, AUTH0_SESSION_NAME)
		session.Values[AUTH0_SESSION_ACCESS_TOKEN_KEY] = token.AccessToken
		session.Values[AUTH0_SESSION_USER_KEY] = *user
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

func Logout(a *auth.Authenticator) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		http.Redirect(w, r, a.LogoutUrl(), http.StatusMovedPermanently)
		return nil
	}
}

func LogoutCallback(store sessions.Store) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		session, _ := store.Get(r, AUTH0_SESSION_NAME)
		delete(session.Values, AUTH0_SESSION_USER_KEY)
		err := session.Save(r, w)
		if err != nil {
			return err
		}
		return nil
	}
}
