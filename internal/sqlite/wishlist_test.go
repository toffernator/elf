package sqlite_test

import (
	"context"
	"elf/internal/sqlite"
	"fmt"
	"testing"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
	"github.com/stretchr/testify/assert"
)

var db = sqlx.MustConnect("sqlite", ":memory:")
var wishlists = sqlite.NewWishlistStore(db)

func seed() {
	goose.SetDialect("sqlite3")
	err := goose.Up(db.DB, "../../db/migrations")
	if err != nil {
		panic(fmt.Errorf("Migrating the database failed with: %w", err))
	}
	err = goose.Up(db.DB, "../../db/seeds")
	if err != nil {
		panic(fmt.Errorf("Seeding the database failed with: %w", err))
	}
}

var readyByIdTests = []struct {
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

func TestReadById(t *testing.T) {
	seed()
	for _, tt := range readyByIdTests {
		t.Run(fmt.Sprintf("Read %d", tt.input), func(t *testing.T) {
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
