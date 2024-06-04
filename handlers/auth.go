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

func Login(authenticator *auth.Authenticator, secureCookies *securecookie.SecureCookie, oauthStateLength int, oauthStateCookieName string) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		state, err := generateState(oauthStateLength)
		if err != nil {
			return err
		}

		encodedStateValue, err := secureCookies.Encode(oauthStateCookieName, state)
		if err != nil {
			return err
		}
		c := &http.Cookie{
			Name:    oauthStateCookieName,
			Value:   encodedStateValue,
			Path:    "/",
			Expires: time.Now().Add(time.Minute * 10),
			// Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, c)

		http.Redirect(w, r, authenticator.Config.AuthCodeURL(state), http.StatusTemporaryRedirect)
		return nil
	}
}

func LoginCallback(authenticator *auth.Authenticator, store sessions.Store, secureCookies *securecookie.SecureCookie, users auth.AuthenticatedUserStore, oauthStateCookieName string, sessionCookieName string, sessionCookieUserKey string, sessionCookieAccessTokenKey string) HTTPHandler {
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

		expectedStateCookie, err := r.Cookie(oauthStateCookieName)
		if err != nil {
			return fmt.Errorf("%s not present because: %w", oauthStateCookieName, err)
		}
		if err = expectedStateCookie.Valid(); err != nil {
			return err
		}

		var expectedState string
		if err = secureCookies.Decode(oauthStateCookieName, expectedStateCookie.Value, &expectedState); err != nil {
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

		session, _ := store.Get(r, sessionCookieName)
		session.Values[sessionCookieAccessTokenKey] = token.AccessToken
		session.Values[sessionCookieUserKey] = *user
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

func Logout(a *auth.Authenticator, callbackUrl string) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		http.Redirect(w, r, a.LogoutUrl(callbackUrl), http.StatusMovedPermanently)
		return nil
	}
}

func LogoutCallback(store sessions.Store, sessionCookieName string, sessionCookieUserKey string) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		session, _ := store.Get(r, sessionCookieName)
		delete(session.Values, sessionCookieUserKey)
		err := session.Save(r, w)
		if err != nil {
			return err
		}
		return nil
	}
}
