package middleware

import (
	"context"
	"elf/internal/core"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gorilla/sessions"
)

type MiddlewareFunc func(w http.ResponseWriter, r *http.Request, next http.Handler) error

type contextKey int

const UserKey contextKey = 0

func AddUserToContext(store sessions.Store, sessionCookieName string, sessionCookieUserKey string) MiddlewareFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.Handler) error {
		session, err := store.Get(r, sessionCookieName)
		if err != nil {
			return err
		}

		user, ok := session.Values[sessionCookieUserKey].(core.User)
		if !ok {
			return errors.New("Cannot cast the user in the session to an auth.AuthenticatedUser")
		}

		ctxWithUser := context.WithValue(r.Context(), UserKey, user)
		rWithUser := r.WithContext(ctxWithUser)
		next.ServeHTTP(w, rWithUser)

		return nil
	}
}

func GetUser(ctx context.Context) (core.User, error) {
	if user, ok := ctx.Value(UserKey).(core.User); ok {
		return user, nil
	}
	return core.User{}, errors.New("Unauthenticated")
}

func Make(m MiddlewareFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := m(w, r, next); err != nil {
				slog.Error("Middleware handler error", "err", err, "path", r.URL.Path)
				next.ServeHTTP(w, r)
			}
		})
	}
}
