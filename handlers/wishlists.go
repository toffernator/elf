package handlers

import (
	"elf/internal/config"
	"elf/internal/core"
	"elf/middleware"
	components "elf/views/wishlist"
	"net/http"
)

func NewWishlist(cfg *config.Config, srvcs *WishlistServices) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		r.ParseForm()

		owner, err := middleware.GetUser(r.Context())
		if err != nil {
			return err
		}

		wl, err := srvcs.Wishlists.Create(r.FormValue("name"), owner.Id)
		if err != nil {
			return err
		}

		return Render(w, r, components.Wishlist(wl))
	}
}

type WishlistServices struct {
	Wishlists WishlistCreator
}

type WishlistCreator interface {
	Create(name string, ownerId int, products ...core.Product) (core.Wishlist, error)
}
