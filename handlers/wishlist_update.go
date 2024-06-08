package handlers

import (
	"elf/internal/config"
	"elf/internal/core"
	"net/http"
	"net/url"
)

func PatchWishlist(cfg *config.Config, srvcs *WishlistServices) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		req := NewUpdateWishlistRequst(r)
		err := req.Validate()
		if err != nil {
			return err
		}

		_, err = srvcs.WishlistUpdater.AddProduct(req.R.Context(), req.Data.WishlistId, core.Product{
			Name:     req.Data.Product.Name,
			Url:      req.Data.Product.Url,
			Price:    int(req.Data.Product.Price),
			Currency: req.Data.Product.Currency,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

type UpdateWishlistRequest struct {
	R *http.Request

	Data struct {
		WishlistId int `validate:"required"`
		Product    struct {
			Name     string  `validate:"required" form:"name"`
			Url      string  `validate:"required" form:"url"`
			Price    float32 `form:"price"`
			Currency string  `form:"currency"`
		}
	}
}

func NewUpdateWishlistRequst(r *http.Request) *UpdateWishlistRequest {
	return &UpdateWishlistRequest{R: r}
}

func (r *UpdateWishlistRequest) Validate() error {
	values, err := r.parse()
	if err != nil {
		return nil
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

func (r *UpdateWishlistRequest) parse() (values url.Values, err error) {
	err = r.R.ParseForm()
	if err != nil {
		return values, err
	}

	values = r.R.Form
	return values, nil
}
