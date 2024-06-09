package handlers

import (
	"context"
	"elf/internal/core"
)

type WishlistServices struct {
	WishlistCreator WishlistCreator
	WishlistReader  WishlistReader
	WishlistUpdater WishlistUpdater
}

type WishlistReader interface {
	ReadById(ctx context.Context, id int) (core.Wishlist, error)
	ReadByOwner(ctx context.Context, id int) (ws []core.Wishlist, err error)
}

type WishlistUpdater interface {
	AddProduct(ctx context.Context, id int, p core.Product) (core.Wishlist, error)
}

type WishlistCreator interface {
	Create(name string, ownerId int, products ...core.Product) (core.Wishlist, error)
}
