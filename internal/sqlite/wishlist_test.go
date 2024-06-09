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
var wishlists = sqlite.NewWishlist(db)

var insertUser, _ = db.Preparex(`INSERT INTO user (name)
    VALUES($1)`)

var insertWishlist, _ = db.Preparex(`INSERT INTO wishlist (owner_id, name)
    VALUES($1, $2)`)
var insertProduct, _ = db.Preparex(`INSERT INTO product (name, url, price, currency, belongs_to_id)
    VALUES($1, $2, $3, $4, $5)`)

func seed() {
	goose.SetDialect("sqlite3")
	err := goose.Up(db.DB, "../../db/migrations")
	if err != nil {
		panic(fmt.Errorf("goose.Up error: %w", err))
	}

	_, err = insertUser.Exec("randy bobandy")
	if err != nil {
		panic(fmt.Errorf("Insert user panic: %w", err))
	}

	_, err = insertWishlist.Exec(1, "birthday")
	if err != nil {
		panic(fmt.Errorf("Insert wishlist panic: %w", err))
	}

	_, err = insertProduct.Exec("iPad", "www.example.com", 100, "eur", 1)
	if err != nil {
		panic(fmt.Errorf("Insert product panic: %w", err))
	}

	_, err = insertProduct.Exec("Macbook", "www.example.com", 200, "eur", 1)
	if err != nil {
		panic(fmt.Errorf("Insert product panic: %w", err))
	}

}

var readyByIdTests = []struct {
	input              int
	expectedName       string
	expectedOwnerId    int
	expectedProductLen int
}{
	{
		input:              1,
		expectedName:       "birthday",
		expectedOwnerId:    1,
		expectedProductLen: 2,
	},
}

func TestReadById(t *testing.T) {
	seed()
	for _, tt := range readyByIdTests {
		t.Run(fmt.Sprintf("Read %d", tt.input), func(t *testing.T) {
			actual, err := wishlists.ReadById(context.Background(), tt.input)
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
