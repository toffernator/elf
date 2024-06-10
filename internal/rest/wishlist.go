package rest

import (
	"elf/internal/core"
	components "elf/internal/rest/views/wishlist"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
)

type WishlistCreateReq struct {
	OwnerId int64
	Name    string `form:"name"`
	Image   string `form:"image"`
}

func (s *Server) HandleWishlistCreate(w http.ResponseWriter, r *http.Request) (err error) {
	var req WishlistCreateReq
	err = decodeWishlistCreateReq(&req, r)
	if err != nil {
		return err
	}
	slog.Info("HandleWishlistCreate is called", "args", req)

	wl, err := s.Wishlists.Create(r.Context(), core.WishlistCreateParams{
		OwnerId: req.OwnerId,
		Name:    req.Name,
		Image:   req.Image,
	})
	if err != nil {
		return err
	}

	return Render(w, r, components.Wishlist(wl))
}

func decodeWishlistCreateReq(req *WishlistCreateReq, r *http.Request) (err error) {
	err = Decode(&req, r)
	if err != nil {
		return err
	}

	u, err := GetUser(r.Context())
	if err != nil {
		return err
	}
	slog.Info("decodeWishlistCreateReq got the user from the request's context.", "user", u)

	req.OwnerId = u.Id
	return
}

type WishlistReadByReq struct {
	OwnerId int64 `form:"ownerId"`
}

func (s *Server) HandleWishlistReadBy(w http.ResponseWriter, r *http.Request) (err error) {
	var req WishlistReadByReq
	err = Decode(&req, r)
	if err != nil {
		return err
	}

	_, err = s.Wishlists.ReadBy(r.Context(), core.WishlistReadByParams{
		OwnerId: req.OwnerId,
	})
	if err != nil {
		return err
	}

	// TODO: Need a component for a collection of wishlists
	return Render(w, r, templ.NopComponent)
}

type WishlistReadReq struct {
	Id int64
}

func (s *Server) HandleWishlistRead(w http.ResponseWriter, r *http.Request) (err error) {
	var req WishlistReadReq
	err = decodeWishlistReadReq(&req, r)
	if err != nil {
		return err
	}

	wl, err := s.Wishlists.Read(r.Context(), req.Id)
	if err != nil {
		return err
	}

	return Render(w, r, components.Wishlist(wl))
}

func decodeWishlistReadReq(req *WishlistReadReq, r *http.Request) (err error) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return &DecodingError{Field: "id", Value: idStr, Expectation: "be an integer"}
	}

	req.Id = int64(id)
	return nil

}

func (s *Server) HandleWishlistUpdate(w http.ResponseWriter, r *http.Request) (err error) {
	var req UpdateWishlistRequest
	err = decodeUpdateWishlistRequest(&req, r)
	if err != nil {
		return err
	}

	wl, err := s.Wishlists.Update(r.Context(), core.WishlistUpdateParams{
		Id: req.Id,
	})
	if err != nil {
		return err
	}

	return Render(w, r, components.Wishlist(wl))
}

type UpdateWishlistRequest struct {
	Id int64
}

func decodeUpdateWishlistRequest(req *UpdateWishlistRequest, r *http.Request) (err error) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return &DecodingError{Field: "id", Value: idStr, Expectation: "be an integer"}
	}

	req.Id = int64(id)
	return nil
}
