package rest

import (
	"elf/internal/core"
	components "elf/internal/rest/views/wishlist"
	restcontext "elf/internal/rest_context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

func (s *Server) HandleWishlistNew(w http.ResponseWriter, r *http.Request) (err error) {
	return Render(w, r, components.NewWishlist())
}

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

	_, err = s.Wishlists.Create(r.Context(), core.WishlistCreateParams{
		OwnerId: req.OwnerId,
		Name:    req.Name,
		Image:   req.Image,
	})
	if err != nil {
		return err
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func decodeWishlistCreateReq(req *WishlistCreateReq, r *http.Request) (err error) {
	err = Decode(&req, r)
	if err != nil {
		return err
	}

	u, err := restcontext.GetUser(r.Context())
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

func (s *Server) HandleWishlistAddProduct(w http.ResponseWriter, r *http.Request) (err error) {
	var req WishlistAddProductReq
	err = decodeWishlistAddProductReq(&req, r)
	if err != nil {
		return err
	}
	slog.Info("HandleProductCreate is called", "args", req)

	_, err = s.Products.Create(r.Context(), core.ProductCreateParams{
		BelongsToId: req.BelongsToId,
		Name:        req.Name,
		Url:         req.Url,
		Price:       req.Price,
		// TODO: Parse the string to a currency integer (in decodeWishlistAddProductReq, with a helper function)
		Currency: core.CurrencyEur,
	})
	if err != nil {
		return err
	}

	http.Redirect(w, r, fmt.Sprintf("/wishlist/%d", req.BelongsToId), http.StatusSeeOther)
	return
}

type WishlistAddProductReq struct {
	BelongsToId int64
	Name        string `form:"name"`
	Url         string `form:"url"`
	Price       int    `form:"price"`
	Currency    string `form:"currency"`
}

func decodeWishlistAddProductReq(req *WishlistAddProductReq, r *http.Request) (err error) {
	err = Decode(&req, r)
	if err != nil {
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return &DecodingError{Field: "id", Value: idStr, Expectation: "be an integer"}
	}

	req.BelongsToId = int64(id)
	return
}
