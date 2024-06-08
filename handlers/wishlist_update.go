package handlers

import (
	"elf/internal/config"
	"elf/internal/core"
	"net/http"
	"strconv"
)

func PatchWishlist(cfg *config.Config, srvcs *WishlistServices) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		req := NewUpdateWishlistRequst(r)

		_, err := srvcs.WishlistUpdater.AddProduct(req.Context(), req.WishlistId, req.Product)
		if err != nil {
			return err
		}

		return nil
	}
}

type UpdateWishlistRequest struct {
	*http.Request

	WishlistId int
	Product    core.Product
}

func NewUpdateWishlistRequst(r *http.Request) *UpdateWishlistRequest {
	return &UpdateWishlistRequest{Request: r}
}

func (r *UpdateWishlistRequest) Validate() error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	es := make(map[Field]FieldError, 0)

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		es["id"] = FieldError{
			Location: PATH_PARAM_LOCATION,
			Value:    idStr,
			Reason:   REASON_NOT_AN_INTEGER,
		}
	}
	name := r.PostFormValue("name")
	if name == "" {
		es["name"] = FieldError{
			Location: FORM_LOCATION,
			Value:    name,
			Reason:   REASON_REQUIRED,
		}
	}
	url := r.PostFormValue("url")
	if url == "" {
		es["url"] = FieldError{
			Location: FORM_LOCATION,
			Value:    url,
			Reason:   REASON_REQUIRED,
		}
	}
	priceStr := r.PostFormValue("price")
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
	currency := r.PostFormValue("currency")
	if currency == "" {
		es["currency"] = FieldError{
			Location: FORM_LOCATION,
			Value:    currency,
			Reason:   REASON_REQUIRED,
		}
	}

	if len(es) > 0 {
		return ValidationError(es)
	}

	r.WishlistId = id
	r.Product = core.Product{
		Name:     name,
		Url:      url,
		Price:    price,
		Currency: currency,
	}
	return nil
}
