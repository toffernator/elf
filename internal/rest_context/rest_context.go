package restcontext

import (
	"context"
	"elf/internal/core"
	"errors"
	"log/slog"
)

func GetUser(ctx context.Context) (core.User, error) {
	if user, ok := ctx.Value(UserKey).(core.User); ok {
		slog.Info("GetUser found a user in the context", "user", user)
		return user, nil
	}
	// TODO: Better (api) error
	return core.User{}, errors.New("Unauthenticated")
}

type contextKey int

const UserKey contextKey = 0
