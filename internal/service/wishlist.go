package service

import (
	"context"
	"elf/internal/core"
)

type WishlistStore interface {
	Create(ctx context.Context, p core.WishlistCreateParams) (core.Wishlist, error)
	Read(ctx context.Context, id int64) (core.Wishlist, error)
	ReadBy(ctx context.Context, p core.WishlistReadByParams) ([]core.Wishlist, error)
	Update(ctx context.Context, p core.WishlistUpdateParams) (core.Wishlist, error)
}

// TODO: Mock with counterfeiter
type WishlistService struct {
	store WishlistStore
}

func NewWishlistService(s WishlistStore) *WishlistService {
	return &WishlistService{store: s}
}

func (w *WishlistService) Create(ctx context.Context, p core.WishlistCreateParams) (wl core.Wishlist, err error) {
	err = p.Validate()
	if err != nil {
		return wl, err
	}

	wl, err = w.store.Create(ctx, p)
	if err != nil {
		return wl, err
	}

	return wl, nil
}

func (w *WishlistService) Read(ctx context.Context, id int64) (wl core.Wishlist, err error) {
	wl, err = w.store.Read(ctx, id)
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

	wls, err = w.store.ReadBy(ctx, p)
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

	wl, err = w.store.Update(ctx, p)
	if err != nil {
		return wl, err
	}

	return wl, nil
}
