package handlers

import (
	"elf/internal/config"
	"elf/views/wishlist"
	components "elf/views/wishlist"
	"net/http"
	"net/url"
)

func GetWishlist(cfg *config.Config, srvcs *WishlistServices) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		req := NewReadWishlistRequest(r)
		err = req.Validate()
		if err != nil {
			return err
		}
		err = req.Validate()
		if err != nil {
			return err
		}

		wl, err := srvcs.WishlistReader.ReadById(req.R.Context(), 1)
		if err != nil {
			return err
		}

		return Render(w, r, components.Wishlist(wl))
	}
}

func GetWishlistPage(cfg *config.Config, srvcs *WishlistServices) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		req := NewReadWishlistRequest(r)
		err := req.Validate()
		if err != nil {
			return err
		}

		if isModalOpen := r.URL.Query().Has("openModal"); isModalOpen {
			switch r.URL.Query().Get("openModal") {
			case "addProduct":
				return Render(w, r, wishlist.Modal(req.Data.WishlistId))
			default:
				return ApiError{StatusCode: http.StatusUnprocessableEntity, Msg: "unsupported modal"}
			}
		}

		wl, err := srvcs.WishlistReader.ReadById(req.R.Context(), req.Data.WishlistId)
		if err != nil {
			return err
		}

		return Render(w, r, wishlist.Page(wl))
	}
}

type ReadWishlistRequest struct {
	R    *http.Request
	Data struct {
		WishlistId int `validate:"required" form:"id"`
	}
}

func NewReadWishlistRequest(r *http.Request) *ReadWishlistRequest {
	return &ReadWishlistRequest{R: r}
}

func (r *ReadWishlistRequest) Validate() error {
	values, err := r.parse()
	if err != nil {
		return err
	}
	err = Parse(&r.Data, values)
	if err != nil {
		return err
	}

	err = Validate(&r.Data)
	if err != nil {
		return err
	}

	return nil
}

func (r *ReadWishlistRequest) parse() (values url.Values, err error) {
	values = make(url.Values, 1)
	idStr := r.R.PathValue("id")
	values.Set("id", idStr)
	return values, nil
}
