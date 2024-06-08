package handlers

import (
	"elf/internal/config"
	"elf/middleware"
	components "elf/views/wishlist"
	"net/http"
	"net/url"
	"strconv"
)

func NewWishlist(cfg *config.Config, srvcs *WishlistServices) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		req := NewCreateWishlistRequest(r)
		err := req.Validate()
		if err != nil {
			return err
		}

		wl, err := srvcs.WishlistCreator.Create(req.Data.Name, req.Data.OwnerId)
		if err != nil {
			return err
		}

		return Render(w, r, components.Wishlist(wl))
	}
}

type CreateWishlistRequest struct {
	R *http.Request

	Data struct {
		Name    string `validate:"required" form:"name"`
		OwnerId int    `validate:"required" form:"ownerId"`
	}
}

func NewCreateWishlistRequest(r *http.Request) *CreateWishlistRequest {
	return &CreateWishlistRequest{R: r}
}

func (r *CreateWishlistRequest) Validate() (err error) {
	values, err := r.parse()
	if err != nil {
		return err
	}
	err = Parse(&r.Data, values)
	if err != nil {
		return err
	}

	err = Validate(r.Data)
	if err != nil {
		return err
	}

	return nil
}

func (r *CreateWishlistRequest) parse() (values url.Values, err error) {
	err = r.R.ParseForm()
	if err != nil {
		return values, err
	}
	values = r.R.Form

	owner, err := middleware.GetUser(r.R.Context())
	if err != nil {
		return values, err
	}

	values.Set("ownerId", strconv.Itoa(owner.Id))
	return values, nil
}
