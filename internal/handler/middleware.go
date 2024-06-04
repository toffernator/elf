package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/toffernator/elf/auth"
)

func EnsureAuthenticated(next http.Handler, store sessions.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := ComputeAuthenticatedUser(r, store)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}

		ctxWithUser := context.WithValue(r.Context(), AuthenticatedUserKey, user)
		rWithUser := r.WithContext(ctxWithUser)
		next.ServeHTTP(w, rWithUser)
	}
}

// ComputeAuthenticatedUser returns the *AuthenticatedUser contained by the
// store if any. Otherwise, it returns an ErrNoProfileInSession.
func ComputeAuthenticatedUser(r *http.Request, store sessions.Store) (*auth.AuthenticatedUser, error) {
	session, err := store.Get(r, AUTH0_SESSION_NAME)
	if err != nil {
		return nil, err
	}

	user, found := session.Values[AUTH0_SESSION_USER_KEY].(auth.AuthenticatedUser)
	if !found {
		return nil, ErrNoUserInSession
	}
	return &user, nil
}

type contextKey int

const AuthenticatedUserKey contextKey = 0

func LogRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("", "method", r.Method, "url", r.URL.String())
		next.ServeHTTP(w, r)
	})
}
