package sqlite_test

import (
	"context"
	"elf/internal/core"
	"elf/internal/store"
	"elf/internal/store/sqlite"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserCreate(t *testing.T) {
	var tests = map[string]struct {
		// The "expected" user is derived from the params
		input       core.UserCreateParams
		expectedErr error
	}{
		"A user": {
			input:       core.UserCreateParams{Name: "A user"},
			expectedErr: nil,
		},
	}
	for name, tt := range tests {
		name := name
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db := setupSqlite(name)
			users := sqlite.NewUserStore(db)
			t.Cleanup(func() { db.Close() })

			actual, err := users.Create(context.TODO(), tt.input)
			if err != nil {
				t.Errorf("%s failed with error: %v", t.Name(), err)
				t.FailNow()
			}

			require.Equal(t, tt.input.Name, actual.Name)
		})
	}
}

func TestUserRead(t *testing.T) {
	var tests = map[string]struct {
		input       int64
		expected    core.User
		expectedErr error
	}{
		"User with ID 1": {
			input:       1,
			expected:    core.User{Id: 1, Name: "test user 1"},
			expectedErr: nil,
		},
		"User with an ID that does not exist": {
			input:       math.MaxInt64,
			expectedErr: store.NewEntityDoesNotExistError("user", math.MaxInt64),
		},
	}
	for name, tt := range tests {
		name := name
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db := setupSqlite(name)
			users := sqlite.NewUserStore(db)
			t.Cleanup(func() { db.Close() })

			actual, err := users.Read(context.TODO(), tt.input)
			if err != nil {
				t.Errorf("%s failed with error: %v", t.Name(), err)
				t.FailNow()
			}

			require.Equal(t, tt.expected, actual)
		})
	}
}
