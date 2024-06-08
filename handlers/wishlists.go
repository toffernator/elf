package handlers

import (
	"context"
	"elf/internal/config"
	"elf/internal/core"
	"elf/middleware"
	"elf/views/wishlist"
	components "elf/views/wishlist"
	"net/http"
	"strconv"
)

func NewWishlist(cfg *config.Config, srvcs *WishlistServices) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		r.ParseForm()

		owner, err := middleware.GetUser(r.Context())
		if err != nil {
			return err
		}

		wl, err := srvcs.WishlistCreator.Create(r.FormValue("name"), owner.Id)
		if err != nil {
			return err
		}

		return Render(w, r, components.Wishlist(wl))
	}
}

func GetWishlistPage(cfg *config.Config, srvcs *WishlistServices) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		readWishlistRequest := NewReadWishlistRequest(r)
		err := readWishlistRequest.Validate()
		if err != nil {
			return err
		}

		if isModalOpen := r.URL.Query().Has("openModal"); isModalOpen {
			switch r.URL.Query().Get("openModal") {
			case "addProduct":
				return Render(w, r, wishlist.Modal(readWishlistRequest.Id))
			default:
				return ApiError{StatusCode: http.StatusUnprocessableEntity, Msg: "unsupported modal"}
			}
		}

		wl, err := srvcs.WishlistReader.ReadById(readWishlistRequest.r.Context(), readWishlistRequest.Id)
		if err != nil {
			return err
		}

		return Render(w, r, wishlist.Page(wl))
	}
}

func GetWishlist(cfg *config.Config, srvcs *WishlistServices) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		readWishlistRequest := NewReadWishlistRequest(r)
		err := readWishlistRequest.Validate()
		if err != nil {
			return err
		}

		wl, err := srvcs.WishlistReader.ReadById(readWishlistRequest.r.Context(), readWishlistRequest.Id)
		if err != nil {
			return err
		}

		return Render(w, r, components.Wishlist(wl))
	}
}

type ReadWishlistRequest struct {
	Id int
	r  *http.Request
}

func NewReadWishlistRequest(r *http.Request) *ReadWishlistRequest {
	return &ReadWishlistRequest{r: r}
}

func (r *ReadWishlistRequest) Validate() error {
	errors := make(map[string]string, 0)

	idStr := r.r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errors["id"] = err.Error()
	}
	if len(errors) > 0 {
		return MalformedUrl(errors)
	}

	r.Id = id
	return nil
}

type WishlistServices struct {
	WishlistCreator WishlistCreator
	WishlistReader  WishlistReader
}

type WishlistReader interface {
	ReadById(ctx context.Context, id int) (core.Wishlist, error)
	ReadByOwner(ctx context.Context, id int) (ws []core.Wishlist, err error)
}

type WishlistCreator interface {
	Create(name string, ownerId int, products ...core.Product) (core.Wishlist, error)
}
