package core

import (
	"database/sql"
	"errors"
)

var (
	ErrWishlistDoesNotExist = errors.New("A wishlist with that 'Id' does not exist")
)

type Wishlist struct {
	Id       int            `json:"id" db:"id"`
	OwnerId  int            `json:"ownerId" db:"owner_id"`
	Name     string         `json:"name" db:"name"`
	Image    sql.NullString `json:"image" db:"image"`
	Products []Product      `json:"products"`
}
