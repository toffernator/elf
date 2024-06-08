package store_test

import (
	"context"
	"elf/internal/store"
	"fmt"
	"testing"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

var db = sqlx.MustConnect("sqlite", ":memory:")
var wls = store.NewWishlist(db)

var insertUser, _ = db.Preparex(`INSERT INTO user (name)
    VALUES($1)`)

var insertWishlist, _ = db.Preparex(`INSERT INTO wishlist (owner_id, name)
    VALUES($1, $2)`)
var insertProduct, _ = db.Preparex(`INSERT INTO product (name, url, price, currency, belongs_to_id)
    VALEUS($1, $2, $3, $4, $5)`)

func seed() {
	insertUser.Exec("randy bobandy")
	insertWishlist.Exec(1, "birthday")
	insertProduct.Exec("iPad", "www.example.com", 100, "eur", 1)
	insertProduct.Exec("Macbook", "www.example.com", 200, "eur", 1)
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
			actual, err := wls.ReadById(context.Background(), tt.input)
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
