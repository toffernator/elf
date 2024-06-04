package handlers

import (
	"crypto/rand"
	"elf/internal/auth"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log/slog"
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

func Login(authenticator *auth.Authenticator, secureCookies *securecookie.SecureCookie) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state, err := generateState(OAUTH_STATE_LENGTH)
		if err != nil {
			slog.Error("Failed to generate state when logging in", "err", err.Error())
			http.Error(w, "We did an oopsie!", http.StatusInternalServerError)
			return
		}

		if encoded, err := secureCookies.Encode(OAUTH_STATE_COOKIE_NAME, state); err == nil {
			cookie := &http.Cookie{
				Name:    OAUTH_STATE_COOKIE_NAME,
				Value:   encoded,
				Path:    "/",
				Expires: time.Now().Add(OAUTH_STATE_COOKIE_EXPIRATION),
				// Secure:   true,
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)
		} else {
			slog.Error("Failed to encode secure cookie", "err", err.Error())
			http.Error(w, "We did an oopsie!", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, authenticator.Config.AuthCodeURL(state), http.StatusTemporaryRedirect)
	}
}

func LoginCallback(authenticator *auth.Authenticator, store sessions.Store, secureCookies *securecookie.SecureCookie, users auth.AuthenticatedUserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !r.URL.Query().Has("state") {
			slog.Error("The 'state' query parameter is not present.", "url", r.URL.String())
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		state := r.URL.Query().Get("state")

		expectedStateCookie, err := r.Cookie(OAUTH_STATE_COOKIE_NAME)
		if err != nil {
			slog.Error("The 'state' cookie is not set.", "err", err.Error())
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if err = expectedStateCookie.Valid(); err != nil {
			slog.Error("The 'state' cookie is set, however, it's invalid.", "err", err.Error())
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		var expectedState string
		if err = secureCookies.Decode(OAUTH_STATE_COOKIE_NAME, expectedStateCookie.Value, &expectedState); err != nil {
			slog.Error("Failed to decode secure cookie", "err", err.Error())
			http.Error(w, "We did an oopsie!", http.StatusInternalServerError)
			return
		}

		if state != expectedState {
			slog.Error("The provided state does not match the expected state.")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !r.URL.Query().Has("code") {
			slog.Error("The 'code' query parameter is not present.", "url", r.URL.String())
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		code := r.URL.Query().Get("code")

		token, err := authenticator.Config.Exchange(r.Context(), code)
		if err != nil {
			slog.Error("The authorization code failed to be exchanged for a token.", "err", err.Error())
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		idToken, err := authenticator.VerifyIDToken(r.Context(), token)
		if err != nil {
			slog.Error("The ID token failed verification", "err", err.Error())
			http.Error(w, "We did an oopsie!", http.StatusInternalServerError)
			return
		}

		var profile auth.Profile
		if err := idToken.Claims(&profile); err != nil {
			slog.Error("The JSON of the ID token's claims could not be unmarshalled.", "err", err.Error())
			http.Error(w, "We did an oopsie!", http.StatusInternalServerError)
			return
		}

		user, err := users.Create(profile)
		if errors.Is(err, auth.ErrDuplicateSub) {
			slog.Error("An authenticated user with that 'Sub' already exists", "Profile", profile)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		session, _ := store.Get(r, AUTH0_SESSION_NAME)
		session.Values[AUTH0_SESSION_ACCESS_TOKEN_KEY] = token.AccessToken
		session.Values[AUTH0_SESSION_USER_KEY] = *user
		if err = session.Save(r, w); err != nil {
			slog.Error("The session store failed to save.", "err", err.Error())
			http.Error(w, "We did an oopsie!", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

func generateState(n int) (string, error) {
	data := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func Logout(a *auth.Authenticator, store sessions.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, AUTH0_SESSION_NAME)
		slog.Info("Pre delete from session.Values", "session.Values", session.Values)
		delete(session.Values, AUTH0_SESSION_USER_KEY)
		slog.Info("Pre delete from session.Values", "session.Values", session.Values)

		http.Redirect(w, r, a.LogoutUrl(), http.StatusMovedPermanently)
	}
}
