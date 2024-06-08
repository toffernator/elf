package handlers

import (
	"elf/internal/config"
	"elf/middleware"
	components "elf/views/wishlist"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-playground/form"
)

var decoder *form.Decoder = form.NewDecoder()

func NewWishlist(cfg *config.Config, srvcs *WishlistServices) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		req := NewCreateWishlistRequest(r)
		values, err := req.Parse()
		if err != nil {
			return err
		}
		err = Parse(&req.Data, values)
		if err != nil {
			return err
		}

		err = Validate(req.Data)
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

func (r *CreateWishlistRequest) Parse() (values url.Values, err error) {
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
