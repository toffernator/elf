package handlers

import (
	"elf/internal/config"
	"elf/middleware"
	"elf/views/components"
	"elf/views/home"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) error {
	http.Redirect(w, r, "/home", http.StatusSeeOther)
	return nil
}

func Home(cfg *config.Config, srvcs *HomeServices) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
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

		ws, err := srvcs.Wishlists.ReadByOwner(r.Context(), u.Id)
		if err != nil {
			return err
		}

		return Render(w, r, home.Index(ws))
	}
}

type HomeServices struct {
	Wishlists WishlistReader
}
