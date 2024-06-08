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
	return func(w http.ResponseWriter, r *http.Request) error {
		// req := NewReadWishlistRequest(r)
		req := NewApiRequest(r, func(r *http.Request) (vs map[string]interface{}, err error) {
			vs = make(map[string]interface{}, 1)
			vs["id"] = r.PathValue("id")
			return vs, nil
		})
		/*err := req.Init()
		if err != nil {
			return err
		}*/
		err := req.Validate(map[string]interface{}{
			"id": "required,number",
		})
		if err != nil {
			return err
		}

		wl, err := srvcs.WishlistReader.ReadById(req.R().Context(), 1)
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
		WishlistId int
	}
}

func NewReadWishlistRequest(r *http.Request) *ReadWishlistRequest {
	return &ReadWishlistRequest{R: r}
}

func (r *ReadWishlistRequest) Init() error {
	if err := r.validate(); err != nil {
		return err
	}
	if err := r.parse(); err != nil {
		return err
	}
	return nil
}

func (r *ReadWishlistRequest) parse() error {
	id, err := strconv.Atoi(r.R.PathValue("id"))
	if err != nil {
		return err
	}
	r.Data.WishlistId = id
	return nil
}

func (r *ReadWishlistRequest) validate() error {
	idStr := r.R.PathValue("id")

	es := make(map[Field]FieldError)
	if err := validate.VarCtx(r.R.Context(), idStr, "required,number"); err != nil {
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
