package rest

import (
	"elf/internal/core"
	"elf/internal/rest/views/components"
	"elf/internal/rest/views/home"
	restcontext "elf/internal/rest_context"
	"net/http"
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

	u, err := restcontext.GetUser(r.Context())
	if err != nil {
		return err
	}

	wls, err := s.Wishlists.ReadBy(r.Context(), core.WishlistReadByParams{
		OwnerId: u.Id,
	})

	if err != nil {
		return err
	}

	return Render(w, r, home.Index(wls))
}
