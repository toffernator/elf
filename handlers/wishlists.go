package handlers

import (
	"elf/internal/core"
	"elf/middleware"
	components "elf/views/wishlist"
	"net/http"
)

type WishlistCreator interface {
	Create(name string, ownerId int, products ...core.Product) core.Wishlist
}

func NewWishlist(wls WishlistCreator) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		r.ParseForm()

		owner, err := middleware.GetUser(r.Context())
		if err != nil {
			return err
		}

		wl := wls.Create(r.FormValue("name"), owner.User.Id)

		return Render(w, r, components.Wishlist(wl))
	}
}
