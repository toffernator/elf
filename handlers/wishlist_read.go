package handlers

import (
	"elf/internal/config"
	"elf/views/wishlist"
	components "elf/views/wishlist"
	"net/http"
	"strconv"
)

func GetWishlist(cfg *config.Config, srvcs *WishlistServices) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		req := NewReadWishlistRequest(r)
		err := req.Validate()
		if err != nil {
			return err
		}

		wl, err := srvcs.WishlistReader.ReadById(req.Context(), req.WishlistId)
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
				return Render(w, r, wishlist.Modal(req.WishlistId))
			default:
				return ApiError{StatusCode: http.StatusUnprocessableEntity, Msg: "unsupported modal"}
			}
		}

		wl, err := srvcs.WishlistReader.ReadById(req.Context(), req.WishlistId)
		if err != nil {
			return err
		}

		return Render(w, r, wishlist.Page(wl))
	}
}

type ReadWishlistRequest struct {
	*http.Request
	WishlistId int
}

func NewReadWishlistRequest(r *http.Request) *ReadWishlistRequest {
	return &ReadWishlistRequest{Request: r}
}

func (r *ReadWishlistRequest) Validate() error {
	errors := make(map[Field]FieldError, 0)

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errors["id"] = FieldError{
			Location: PATH_PARAM_LOCATION,
			Value:    idStr,
			Reason:   REASON_NOT_AN_INTEGER,
		}
	}
	if len(errors) > 0 {
		return ValidationError(errors)
	}

	r.WishlistId = id
	return nil
}
