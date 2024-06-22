package sqlite_test

import (
	"context"
	"elf/internal/sqlite"
	"fmt"
	"testing"

	_ "github.com/glebarez/go-sqlite"
	"github.com/stretchr/testify/assert"
)

var wishlists = sqlite.NewWishlistStore(db)

var wishlistReadByIdTests = []struct {
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
	{
		input:              2,
		expectedName:       "test wishlist 2 belonging to user with id 1",
		expectedOwnerId:    1,
		expectedProductLen: 2,
	},
	{
		input:              3,
		expectedName:       "test wishlist 3 belonging to user with id 2",
		expectedOwnerId:    2,
		expectedProductLen: 0,
	},
}

func TestWishlistReadById(t *testing.T) {
	seed()
	for _, tt := range wishlistReadByIdTests {
		t.Run(fmt.Sprintf("Read Wishlist %d", tt.input), func(t *testing.T) {
			actual, err := wishlists.Read(context.Background(), tt.input)
			if err != nil {
				t.Errorf("%s failed with error: %v", t.Name(), err)
				t.FailNow()
			}

			assert.Equal(t, tt.input, actual.Id)
			assert.Equal(t, tt.expectedName, actual.Name)
			assert.Equal(t, tt.expectedOwnerId, actual.OwnerId)
			assert.Len(t, actual.Products, tt.expectedProductLen)
		})
	}
}
