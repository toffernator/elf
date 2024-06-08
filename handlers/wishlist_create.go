package handlers

import (
	"elf/internal/config"
	"elf/middleware"
	components "elf/views/wishlist"
	"net/http"

	"github.com/go-playground/form"
	"github.com/go-playground/validator/v10"
)

var decoder *form.Decoder = form.NewDecoder()

func NewWishlist(cfg *config.Config, srvcs *WishlistServices) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		req := NewCreateWishlistRequest(r)
		err := req.Init()
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
		WishlistId string `validate:"required" form:"id"`
		Name       string `validate:"required" form:"name"`
		OwnerId    int    `validate:"required" form:"ownerId"`
	}
}

func NewCreateWishlistRequest(r *http.Request) *CreateWishlistRequest {
	return &CreateWishlistRequest{R: r}
}

func (r *CreateWishlistRequest) Init() error {
	if err := r.parse(); err != nil {
		return err
	}
	if err := r.mold(); err != nil {
		return err
	}
	if err := r.validate(); err != nil {
		return err
	}
	return nil
}

func (r *CreateWishlistRequest) parse() error {
	if err := r.R.ParseForm(); err != nil {
		return err
	}
	if err := decoder.Decode(&r.Data, r.R.Form); err != nil {
		if es, ok := err.(form.DecodeErrors); ok {
			return DecoderErrors(es)
		}

		return err
	}

	owner, err := middleware.GetUser(r.R.Context())
	if err != nil {
		return err
	}
	r.Data.OwnerId = owner.Id
	return nil
}

func (r *CreateWishlistRequest) mold() error {
	return nil
}

func (r *CreateWishlistRequest) validate() error {
	es := make(map[Field]FieldError, 0)
	if err := validate.StructCtx(r.R.Context(), r.Data); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			es[Field(err.Field())] = FieldError{
				Location: "",
				Value:    err.Param(),
				Reason:   err.Error(),
			}
		}
		if len(es) > 0 {
			return ValidationErrors(es)
		}
	}
	return nil
}
