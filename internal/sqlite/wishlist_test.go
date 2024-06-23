package sqlite_test

import (
	"context"
	"elf/internal/core"
	"elf/internal/sqlite"
	"fmt"
	"testing"

	_ "github.com/glebarez/go-sqlite"
	"github.com/stretchr/testify/assert"
)

// TODO: Tests that are expected to fail (Read non-existent id, create invalid wishlist)

var wishlists = sqlite.NewWishlistStore(db)

// The expected values is derived from the values of the the input's fields.
var wishlistCreateTests = []core.WishlistCreateParams{
	{
		OwnerId:  1,
		Name:     "A wishlist created without products",
		Products: []core.ProductCreateParams{},
	},
	{
		OwnerId: 1,
		Name:    "A wishlist created with products",
		Products: []core.ProductCreateParams{
			{Name: "A product"},
			{Name: "Another product"},
		},
	},
}

func TestWishlistCreate(t *testing.T) {
	seed()
	for _, tt := range wishlistCreateTests {
		t.Run(fmt.Sprintf("Create Wishlist %+v", tt), func(t *testing.T) {
			actual, err := wishlists.Create(context.TODO(), tt)
			if err != nil {
				t.Errorf("%s failed with error: %v", t.Name(), err)
			}
			actualWithProducts, err := wishlists.Read(context.TODO(), actual.Id)
			if err != nil {
				t.Errorf("%s failed with error: %v", t.Name(), err)
			}

			assert.Equal(t, tt.Name, actual.Name)
			assert.Equal(t, tt.OwnerId, actual.OwnerId)
			assert.Len(t, actualWithProducts.Products, len(tt.Products))
		})
	}
}

var wishlistReadTests = []struct {
	input              int64
	expectedName       string
	expectedOwnerId    int64
	expectedProductLen int
}{
	{
		input:              1,
		expectedName:       "test wishlist 1 belonging to user with id 1",
		expectedOwnerId:    1,
		expectedProductLen: 1,
	},
}

func TestWishlistRead(t *testing.T) {
	seed()
	for _, tt := range wishlistReadTests {
		t.Run(fmt.Sprintf("Read Wishlist %d", tt.input), func(t *testing.T) {
			actual, err := wishlists.Read(context.TODO(), tt.input)
			if err != nil {
				t.Errorf("%s failed with error: %v", t.Name(), err)
			}

			assert.Equal(t, tt.input, actual.Id)
			assert.Equal(t, tt.expectedName, actual.Name)
			assert.Equal(t, tt.expectedOwnerId, actual.OwnerId)
			assert.Len(t, actual.Products, tt.expectedProductLen)
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

			assert.ElementsMatch(t, actual, tt.expected)
		})
	}
}
