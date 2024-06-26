package service

import (
	"context"
	"elf/internal/core"

	"errors"
)

type WishlistStore interface {
	Create(ctx context.Context, p core.WishlistCreateParams) (core.Wishlist, error)
	Read(ctx context.Context, id int64) (core.Wishlist, error)
	ReadBy(ctx context.Context, p core.WishlistReadByParams) ([]core.Wishlist, error)
	IsOwnedBy(ctx context.Context, wishlistId int64, userId int64) (bool, error)
	Update(ctx context.Context, p core.WishlistUpdateParams) (core.Wishlist, error)
}

type WishlistService struct {
	wishlistStore WishlistStore
	productStore  ProductStore
}

func NewWishlistService(s WishlistStore) *WishlistService {
	return &WishlistService{wishlistStore: s}
}

func (w *WishlistService) Create(ctx context.Context, p core.WishlistCreateParams) (wl core.Wishlist, err error) {
	err = p.Validate()
	if err != nil {
		return wl, err
	}

	wl, err = w.wishlistStore.Create(ctx, p)
	if err != nil {
		return wl, err
	}

	return wl, nil
}

func (w *WishlistService) Read(ctx context.Context, id int64) (wl core.Wishlist, err error) {
	wl, err = w.wishlistStore.Read(ctx, id)
	if err != nil {
		return wl, err
	}

	return wl, nil
}

func (w *WishlistService) ReadBy(ctx context.Context, p core.WishlistReadByParams) (wls []core.Wishlist, err error) {
	err = p.Validate()
	if err != nil {
		return wls, err
	}

	wls, err = w.wishlistStore.ReadBy(ctx, p)
	if err != nil {
		return wls, err
	}

	return wls, nil
}

func (w *WishlistService) Update(ctx context.Context, p core.WishlistUpdateParams) (wl core.Wishlist, err error) {
	err = p.Validate()
	if err != nil {
		return wl, err
	}

	wl, err = w.wishlistStore.Update(ctx, p)
	if err != nil {
		return wl, err
	}

	return wl, nil
}

func (w *WishlistService) AddProduct(ctx context.Context, wishlistId int64, userId int64, p core.ProductCreateParams) (err error) {
	err = p.Validate()
	if err != nil {
		return err
	}

	doesOwnWishlist, err := w.wishlistStore.IsOwnedBy(ctx, wishlistId, userId)
	if err != nil {
		return err
	}
	if !doesOwnWishlist {
		// TODO: Better Unauthorized error
		return errors.New("You do not own this wishlist")
	}

	_, err = w.productStore.Create(ctx, p)
	if err != nil {
		return err
	}

	return nil
}
