package rest

import (
	"elf/internal/core"
	"elf/internal/rest/views/components"
	"net/http"

	"github.com/a-h/templ"
)

func (s *Server) HandleHome(w http.ResponseWriter, r *http.Request) error {
	if isModalOpen := r.URL.Query().Has("openModal"); isModalOpen {
		switch r.URL.Query().Get("openModal") {
		case "newWishlist":
			return Render(w, r, components.Modal())
		default:
			return &Error{StatusCode: http.StatusUnprocessableEntity, Reason: "unsupported modal"}
		}
	}

	u, err := GetUser(r.Context())
	if err != nil {
		return err
	}

	_, err = s.Wishlists.ReadBy(r.Context(), core.WishlistReadByParams{
		OwnerId: u.Id,
	})

	if err != nil {
		return err
	}

	return Render(w, r, templ.NopComponent)
}
