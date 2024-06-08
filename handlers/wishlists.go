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

func PatchWishlist(cfg *config.Config, srvcs *WishlistServices) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		req := NewUpdateWishlistRequest(r)

		_, err := srvcs.WishlistUpdater.AddProduct(req.Context(), req.Id, req.Product)
		if err != nil {
			return err
		}

		return nil
	}
}

type UpdateWishlistRequest struct {
	*http.Request

	Id      int
	Product core.Product
}

func NewUpdateWishlistRequest(r *http.Request) *UpdateWishlistRequest {
	return &UpdateWishlistRequest{
		Request: r,
	}
}

func (r *UpdateWishlistRequest) Validate() error {
	errors := make(map[string]string, 0)

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errors["id"] = err.Error()
	}
	if len(errors) > 0 {
		return MalformedUrl(errors)
	}
	r.Id = id

	err = r.ParseForm()
	if err != nil {
		return err
	}
	errors = make(map[string]string, 0)

	name := r.PostFormValue("name")
	if name == "" {
		errors["name"] = "'name' is required"
	}
	url := r.PostFormValue("url")
	if url == "" {
		errors["url"] = "'url' is required"
	}
	priceStr := r.PostFormValue("price")
	if priceStr == "" {
		errors["price"] = "'price' is required"
	}
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		errors["price"] = "'price' must be an integer"
	}
	currency := r.PostFormValue("currency")
	if currency == "" {
		errors["currency"] = "'currency' is required"
	}

	p := core.Product{
		Name:     name,
		Url:      url,
		Price:    price,
		Currency: currency,
	}

	r.Product = p
	return nil
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
	WishlistUpdater WishlistUpdater
}

type WishlistReader interface {
	ReadById(ctx context.Context, id int) (core.Wishlist, error)
	ReadByOwner(ctx context.Context, id int) (ws []core.Wishlist, err error)
}

type WishlistUpdater interface {
	AddProduct(ctx context.Context, id int, p core.Product) (core.Wishlist, error)
}

type WishlistCreator interface {
	Create(name string, ownerId int, products ...core.Product) (core.Wishlist, error)
}
