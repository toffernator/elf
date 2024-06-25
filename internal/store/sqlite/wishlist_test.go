package sqlite_test

import (
	"context"
	"database/sql"
	"elf/internal/core"
	"elf/internal/sqlite"
	"math"
	"testing"

	_ "github.com/glebarez/go-sqlite"
	"github.com/stretchr/testify/require"
)

// TODO: Tests that are expected to fail (Read non-existent id, create invalid wishlist, ...)

func TestWishlistCreate(t *testing.T) {
	tests := map[string]struct {
		// The "expected" wishlist is derived from the params
		params      core.WishlistCreateParams
		expectedErr error
	}{
		"Create a wishlist": {
			params: core.WishlistCreateParams{
				OwnerId: 1,
				Name:    "A wishlist without products",
			},
		},
		"Create a wishlist owned by a user that does not exist": {
			params: core.WishlistCreateParams{
				OwnerId: math.MaxInt64,
				Name:    "A wishlist owned by a user that does not exist",
			},
			expectedErr: sqlite.NewEntityDoesNotExistError("user", math.MaxInt64),
		},
	}

	for name, tt := range tests {
		name := name
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db := setupSqlite(name)
			wishlists := sqlite.NewWishlistStore(db)
			t.Cleanup(func() { db.Close() })

			actual, err := wishlists.Create(context.TODO(), tt.params)
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
				return
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.params.Name, actual.Name)
			require.Equal(t, tt.params.OwnerId, actual.OwnerId)
		})
	}
}

func TestWishlistRead(t *testing.T) {
	tests := map[string]struct {
		input              int64
		expectedName       string
		expectedOwnerId    int64
		expectedProductLen int
		expectedErr        error
	}{
		"Wishlist with ID 1": {
			input:              1,
			expectedName:       "test wishlist 1 belonging to user with id 1",
			expectedOwnerId:    1,
			expectedProductLen: 1,
		},
		"Wishlist with an ID that does not exist": {
			input:       math.MaxInt64,
			expectedErr: sql.ErrNoRows,
		},
	}
	for name, tt := range tests {
		tt := tt
		name := name
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db := setupSqlite(name)
			wishlists := sqlite.NewWishlistStore(db)
			t.Cleanup(func() { db.Close() })

			actual, err := wishlists.Read(context.TODO(), tt.input)
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr, err)
				return
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.input, actual.Id)
			require.Equal(t, tt.expectedName, actual.Name)
			require.Equal(t, tt.expectedOwnerId, actual.OwnerId)
			require.Len(t, actual.Products, tt.expectedProductLen)
		})
	}
}

func TestWishlistReadBy(t *testing.T) {
	var wishlistReadByTests = map[string]struct {
		params   core.WishlistReadByParams
		expected []core.Wishlist
	}{
		"An OwnerId with several wishlists": {
			params: core.WishlistReadByParams{OwnerId: 1},
			expected: []core.Wishlist{
				{Id: 1, OwnerId: 1, Name: "test wishlist 1 belonging to user with id 1"},
				{Id: 2, OwnerId: 1, Name: "test wishlist 2 belonging to user with id 1"},
			},
		},
		"An OwnerId with one wishlist": {
			params: core.WishlistReadByParams{OwnerId: 2},
			expected: []core.Wishlist{
				{Id: 3, OwnerId: 2, Name: "test wishlist 3 belonging to user with id 2"},
			},
		},
		"An OwnerId with no wishlists": {
			params:   core.WishlistReadByParams{OwnerId: 3},
			expected: []core.Wishlist{},
		},
	}

	for name, tt := range wishlistReadByTests {
		name := name
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			db := setupSqlite(name)
			wishlists := sqlite.NewWishlistStore(db)
			t.Cleanup(func() { db.Close() })

			actual, err := wishlists.ReadBy(context.TODO(), tt.params)
			if err != nil {
				t.Errorf("%s failed with the error: %s", t.Name(), err)
			}

			require.ElementsMatch(t, actual, tt.expected)
		})
	}
}
