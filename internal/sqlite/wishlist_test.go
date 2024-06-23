package sqlite_test

import (
	"context"
	"database/sql"
	"elf/internal/core"
	"elf/internal/sqlite"
	"fmt"
	"math"
	"testing"

	_ "github.com/glebarez/go-sqlite"
	"github.com/stretchr/testify/require"
)

// TODO: Tests that are expected to fail (Read non-existent id, create invalid wishlist, ...)

var wishlists = sqlite.NewWishlistStore(db)

// The expected values is derived from the values of the the params's fields.
var wishlistCreateTests = []struct {
	params      core.WishlistCreateParams
	expectedErr error
}{
	{
		params: core.WishlistCreateParams{
			OwnerId:  1,
			Name:     "A wishlist created without products",
			Products: []core.ProductCreateParams{},
		},
	},
	{
		params: core.WishlistCreateParams{
			OwnerId: 1,
			Name:    "A wishlist created with products",
			Products: []core.ProductCreateParams{
				{Name: "A product"},
				{Name: "Another product"},
			},
		},
	},
	{
		params: core.WishlistCreateParams{
			OwnerId: math.MaxInt64,
			Name:    "A wishlist with a non-sensical owner id",
		},
		expectedErr: sqlite.NewEntityDoesNotExistError("user", math.MaxInt64),
	},
}

func TestWishlistCreate(t *testing.T) {
	seed()
	for _, tt := range wishlistCreateTests {
		t.Run(fmt.Sprintf("Create Wishlist %+v", tt), func(t *testing.T) {
			actual, err := wishlists.Create(context.TODO(), tt.params)
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr, "The actual error is %v.", err)
				return
			} else {
				require.NoError(t, err)
			}

			actualWithProducts, err := wishlists.Read(context.TODO(), actual.Id)
			require.NoError(t, err)

			require.Equal(t, tt.params.Name, actual.Name)
			require.Equal(t, tt.params.OwnerId, actual.OwnerId)
			require.Len(t, actualWithProducts.Products, len(tt.params.Products))
		})
	}
}

var wishlistReadTests = []struct {
	input              int64
	expectedName       string
	expectedOwnerId    int64
	expectedProductLen int
	expectedErr        error
}{
	{
		input:              1,
		expectedName:       "test wishlist 1 belonging to user with id 1",
		expectedOwnerId:    1,
		expectedProductLen: 1,
	},
	{
		input:       math.MaxInt64,
		expectedErr: sql.ErrNoRows,
	},
}

func TestWishlistRead(t *testing.T) {
	seed()
	for _, tt := range wishlistReadTests {
		t.Run(fmt.Sprintf("Read Wishlist %d", tt.input), func(t *testing.T) {
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

var wishlistReadByTests = []struct {
	params   core.WishlistReadByParams
	expected []core.Wishlist
}{
	{
		params: core.WishlistReadByParams{OwnerId: 1},
		expected: []core.Wishlist{
			{Id: 1, OwnerId: 1, Name: "test wishlist 1 belonging to user with id 1"},
			{Id: 2, OwnerId: 1, Name: "test wishlist 2 belonging to user with id 1"},
		},
	},
	{
		params: core.WishlistReadByParams{OwnerId: 2},
		expected: []core.Wishlist{
			{Id: 3, OwnerId: 2, Name: "test wishlist 3 belonging to user with id 2"},
		},
	},
	{
		params:   core.WishlistReadByParams{OwnerId: 3},
		expected: []core.Wishlist{},
	},
}

func TestWishlistReadBy(t *testing.T) {
	seed()
	for _, tt := range wishlistReadByTests {
		t.Run(fmt.Sprintf("ReadBy Wishlist %+v", tt.params), func(t *testing.T) {
			actual, err := wishlists.ReadBy(context.TODO(), tt.params)
			if err != nil {
				t.Errorf("%s failed with the error: %s", t.Name(), err)
			}

			require.ElementsMatch(t, actual, tt.expected)
		})
	}
}
