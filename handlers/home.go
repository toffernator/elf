package handlers

import (
	"context"
	"elf/internal/core"
	"elf/middleware"
	"elf/views/components"
	"elf/views/home"
	"net/http"
)

type Services struct {
	WlCreator WishlistCreator
	WlReader  WishlistReader
}

type WishlistReader interface {
	ReadById(ctx context.Context, id int) (w core.Wishlist, err error)
	ReadByOwner(ctx context.Context, id int) (ws []core.Wishlist, err error)
}

func Index(w http.ResponseWriter, r *http.Request) error {
	http.Redirect(w, r, "/home", http.StatusSeeOther)
	return nil
}

// HandleHome assumes that the user is authenticated
func (s *Services) HandleHome(w http.ResponseWriter, r *http.Request) error {
	if isModalOpen := r.URL.Query().Has("openModal"); isModalOpen {
		switch r.URL.Query().Get("openModal") {
		case "newWishlist":
			return Render(w, r, components.Modal())
		default:
			return ApiError{StatusCode: http.StatusUnprocessableEntity, Msg: "unsupported modal"}
		}
	}

	u, err := middleware.GetUser(r.Context())
	if err != nil {
		return err
	}

	ws, err := s.WlReader.ReadByOwner(r.Context(), u.Id)
	if err != nil {
		return err
	}

	return Render(w, r, home.Index(ws))
}
