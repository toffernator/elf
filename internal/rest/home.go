package rest

import (
	"elf/internal/core"
	"elf/internal/rest/views/home"
	restcontext "elf/internal/rest_context"
	"net/http"
)

func (s *Server) HandleHome(w http.ResponseWriter, r *http.Request) error {
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
