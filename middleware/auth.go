package middleware

import (
	"context"
	"elf/internal/auth"
	"elf/internal/handler"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gorilla/sessions"
)

type MiddlewareFunc func(w http.ResponseWriter, r *http.Request, next http.Handler) error

type contextKey int

const AuthenticatedUserKey contextKey = 0

func AddUserToContext(store sessions.Store, sessionCookieName string) MiddlewareFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.Handler) error {
		session, err := store.Get(r, sessionCookieName)
		if err != nil {
			return err
		}

		user, ok := session.Values[handler.AUTH0_SESSION_USER_KEY].(auth.AuthenticatedUser)
		if !ok {
			return errors.New("Cannot cast the user in the session to an auth.AuthenticatedUser")
		}

		ctxWithUser := context.WithValue(r.Context(), AuthenticatedUserKey, user)
		rWithUser := r.WithContext(ctxWithUser)
		next.ServeHTTP(w, rWithUser)

		return nil
	}
}

func GetUser(ctx context.Context) (auth.AuthenticatedUser, error) {
	if user, ok := ctx.Value(AuthenticatedUserKey).(auth.AuthenticatedUser); ok {
		return user, nil
	}
	return auth.AuthenticatedUser{}, errors.New("Unauthenticated")
}

func Make(m MiddlewareFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := m(w, r, next); err != nil {
				slog.Error("Middleware handler error", "err", err, "path", r.URL.Path)
			}
		})
	}
}
