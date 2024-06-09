package middleware

import (
	"context"
	"elf/internal/config"
	"elf/internal/core"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gorilla/sessions"
)

type MiddlewareFunc func(w http.ResponseWriter, r *http.Request, next http.Handler) error

type contextKey int

const UserKey contextKey = 0

type AuthServices struct {
	Sessions sessions.Store
}

func AddUserToContext(cfg *config.Config, srvcs *AuthServices) MiddlewareFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.Handler) error {
		session, err := srvcs.Sessions.Get(r, cfg.Auth.SessionCookieName)
		if err != nil {
			return err
		}

		user, ok := session.Values[cfg.Auth.SessionCookieUserKey].(core.User)
		if !ok && cfg.IsProduction() {
			// TODO: better error
			return errors.New("Cannot cast the user in the session to an auth.AuthenticatedUser")
		} else if !ok && cfg.IsDevelop() {
			user = core.User{
				Id: 1,
			}
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
	// TODO: Better (api) error
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
