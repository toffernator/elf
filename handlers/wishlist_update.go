package handlers

import (
	"elf/internal/config"
	"net/http"
	"strconv"
)

func PatchWishlist(cfg *config.Config, srvcs *WishlistServices) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		req := NewUpdateWishlistRequst(r)
		req.Init()

		// _, err := srvcs.WishlistUpdater.AddProduct(req.R.Context(), req.Data.WishlistId, req.Data.Product)
		// if err != nil {
		//return err
		//}

		return nil
	}
}

type UpdateWishlistRequest struct {
	R *http.Request

	Data struct {
		WishlistId int `validate:"required"`
		Product    struct {
			Name     string  `validate:"required"`
			Url      string  `validate:"required"`
			Price    float32 `validate:"number"`
			Currency string
		}
	}
}

func NewUpdateWishlistRequst(r *http.Request) *UpdateWishlistRequest {
	return &UpdateWishlistRequest{R: r}
}

func (r *UpdateWishlistRequest) Init() error {
	if err := r.Validate(); err != nil {
		return err
	}

	return nil
}

func (r *UpdateWishlistRequest) Validate() error {
	err := r.R.ParseForm()
	if err != nil {
		return err
	}

	es := make(map[Field]FieldError, 0)

	idStr := r.R.PathValue("id")
	// id, err := strconv.Atoi(idStr)
	if err != nil {
		es["id"] = FieldError{
			Location: PATH_PARAM_LOCATION,
			Value:    idStr,
			Reason:   REASON_NOT_AN_INTEGER,
		}
	}
	name := r.R.PostFormValue("name")
	if name == "" {
		es["name"] = FieldError{
			Location: FORM_LOCATION,
			Value:    name,
			Reason:   REASON_REQUIRED,
		}
	}
	url := r.R.PostFormValue("url")
	if url == "" {
		es["url"] = FieldError{
			Location: FORM_LOCATION,
			Value:    url,
			Reason:   REASON_REQUIRED,
		}
	}
	priceStr := r.R.PostFormValue("price")
	if priceStr == "" {
		es["price"] = FieldError{
			Location: FORM_LOCATION,
			Value:    priceStr,
			Reason:   REASON_REQUIRED,
		}
	}
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		es["price"] = FieldError{
			Location: FORM_LOCATION,
			Value:    price,
			Reason:   REASON_NOT_AN_INTEGER,
		}
	}
	currency := r.R.PostFormValue("currency")
	if currency == "" {
		es["currency"] = FieldError{
			Location: FORM_LOCATION,
			Value:    currency,
			Reason:   REASON_REQUIRED,
		}
	}

	if len(es) > 0 {
		return ValidationErrors(es)
	}

	return nil
}
