package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/toffernator/elf/auth"
	"github.com/toffernator/elf/internal/core"
)

type WishlistCreator interface {
	Create(ownerId int, products ...core.Product) core.Wishlist
}

type WishlistReader interface {
	ReadAll() []core.Wishlist
	Read(id int) (core.Wishlist, error)
}

type WishlistCreatorReader interface {
	WishlistCreator
	WishlistReader
}

func GetWishlists(wishlists WishlistReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := json.Marshal(wishlists.ReadAll())
		if err != nil {
			slog.Error("", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(body)
	}
}

func GetWishlist(wishlists WishlistReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idPathPart := r.PathValue("id")
		id, err := strconv.Atoi(idPathPart)
		if err != nil {
			slog.Error("The 'id' path value is not a valid integer", "id", idPathPart)
			http.Error(w, "The 'id' path value is not a valid integer", http.StatusBadRequest)
			return
		}

		wishlist, err := wishlists.Read(id)
		if errors.Is(err, core.ErrWishlistDoesNotExist) {
			slog.Error("", err.Error(), slog.Int("id", id))
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("", "err", err.Error())
			http.Error(w, "We did an oopsie!", http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(wishlist)
		if err != nil {
			slog.Error("", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(body)
	}
}

type PostWishlistBody struct {
	OwnerId  int            `json:"ownerId"`
	Products []core.Product `json:"products"`
}

func PostWishlist(store sessions.Store, wishlists WishlistCreatorReader) (f http.HandlerFunc) {
	f = func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(AuthenticatedUserKey).(auth.AuthenticatedUser)
		if !ok {
			slog.Error("'user' is not of type 'auth.AuthenticatedUser'")
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("Reading the request's body failed", "err", err.Error())
			http.Error(w, "We did an oopsie!", http.StatusInternalServerError)
			return
		}

		var wishlist PostWishlistBody
		err = json.Unmarshal(requestBody, &wishlist)
		if err != nil {
			slog.Error("Unmarshalling the request's body into a 'PostWishlistBody' failed", "err", err.Error())
			http.Error(w, "We did an oopsie!", http.StatusInternalServerError)
			return
		}

		if user.User.Id != wishlist.OwnerId {
			slog.Error("The creation of a wishlist cannot be authorized because the user's Id does not match the wishlist's owner id", "user id", user.User.Id, "owner id", wishlist.OwnerId)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		createdWishlist := wishlists.Create(wishlist.OwnerId, wishlist.Products...)

		responseBody, err := json.Marshal(createdWishlist)
		if err != nil {
			slog.Error("", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(responseBody)
	}

	return EnsureAuthenticated(f, store)
}
