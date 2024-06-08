package handlers

import (
	"elf/internal/config"
	"elf/views/wishlist"
	components "elf/views/wishlist"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func GetWishlist(cfg *config.Config, srvcs *WishlistServices) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		req := NewReadWishlistRequest(r)
		err = req.Parse()
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
		err := req.Init()
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
		WishlistId int `validate:"required"`
	}
}

func NewReadWishlistRequest(r *http.Request) *ReadWishlistRequest {
	return &ReadWishlistRequest{R: r}
}

func (r *ReadWishlistRequest) Init() error {
	if err := r.Parse(); err != nil {
		return err
	}
	if err := r.Validate(); err != nil {
		return err
	}
	return nil
}

func (r *ReadWishlistRequest) Parse() error {
	idStr := r.R.PathValue("id")
	err := validate.Var(idStr, "number")
	if err != nil {
		if es, ok := err.(validator.ValidationErrors); ok {
			return ValidationErrors2(es)
		}
		return err
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	r.Data.WishlistId = id
	return nil
}

func (r *ReadWishlistRequest) Validate() error {
	err := validate.Struct(r.Data)
	if err != nil {
		if es, ok := err.(validator.ValidationErrors); ok {
			return ValidationErrors2(es)
		}
		return err
	}

	return nil
}
