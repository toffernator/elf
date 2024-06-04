package core

import "errors"

var (
	ErrWishlistDoesNotExist = errors.New("A wishlist with that 'Id' does not exist")
)

type Wishlist struct {
	Id       int       `json:"id"`
	OwnerId  int       `json:"ownerId"`
	Name     string    `json:"name"`
	Image    string    `json:"image"`
	Products []Product `json:"products"`
}
