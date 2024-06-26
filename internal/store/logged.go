package store

import (
	"context"
	"elf/internal/core"
	"log/slog"
)

// TODO: not a big fan of having to redefine the interface here...
type UserStore interface {
	Create(ctx context.Context, p core.UserCreateParams) (core.User, error)
	Read(ctx context.Context, id int64) (core.User, error)
}

type LoggedUserStore struct {
	users  UserStore
	logger *slog.Logger
}

func NewLoggedUserStore(users UserStore, logger *slog.Logger) *LoggedUserStore {
	return &LoggedUserStore{users, logger}
}

func (s *LoggedUserStore) Create(ctx context.Context, p core.UserCreateParams) (core.User, error) {
	s.logger.Info("UserStore.Create is called.", "p", p)
	return s.users.Create(ctx, p)
}

func (s *LoggedUserStore) Read(ctx context.Context, id int64) (core.User, error) {
	s.logger.Info("UserStore.Read is called.", "id", id)
	return s.users.Read(ctx, id)
}

type WishlistStore interface {
	Create(ctx context.Context, p core.WishlistCreateParams) (core.Wishlist, error)
	Read(ctx context.Context, id int64) (core.Wishlist, error)
	ReadBy(ctx context.Context, p core.WishlistReadByParams) ([]core.Wishlist, error)
	IsOwnedBy(ctx context.Context, wishlistId int64, userId int64) (bool, error)
	Update(ctx context.Context, p core.WishlistUpdateParams) (core.Wishlist, error)
}

type LoggedWishlistStore struct {
	wishlists WishlistStore
	logger    *slog.Logger
}

func NewLoggedWishlistStore(wishlists WishlistStore, logger *slog.Logger) *LoggedWishlistStore {
	return &LoggedWishlistStore{wishlists, logger}
}

func (s *LoggedWishlistStore) Create(ctx context.Context, p core.WishlistCreateParams) (core.Wishlist, error) {
	s.logger.Info("WishlistStore.Create is called.", "p", p)
	return s.wishlists.Create(ctx, p)
}

func (s *LoggedWishlistStore) Read(ctx context.Context, id int64) (core.Wishlist, error) {
	s.logger.Info("WishlistStore.Read is called.", "id", id)
	return s.wishlists.Read(ctx, id)
}

func (s *LoggedWishlistStore) ReadBy(ctx context.Context, p core.WishlistReadByParams) ([]core.Wishlist, error) {
	go s.logger.Info("WishlistStore.ReadBy is called.", "p", p)
	return s.wishlists.ReadBy(ctx, p)
}

func (s *LoggedWishlistStore) IsOwnedBy(ctx context.Context, wishlistId int64, userId int64) (bool, error) {
	slog.Info("WishlistStore.IsOwnedBy is called.", "wishlistId", wishlistId, "userId", userId)
	return s.IsOwnedBy(ctx, wishlistId, userId)
}

func (s *LoggedWishlistStore) Update(ctx context.Context, p core.WishlistUpdateParams) (core.Wishlist, error) {
	s.logger.Info("WishlistStore.Update is called.", "p", p)
	return s.wishlists.Update(ctx, p)
}

type ProductStore interface {
	Create(ctx context.Context, p core.ProductCreateParams) (int64, error)
}

type LoggedProductStore struct {
	products ProductStore
	logger   *slog.Logger
}

func NewLoggedProductStore(products ProductStore, logger *slog.Logger) *LoggedProductStore {
	return &LoggedProductStore{products, logger}
}

func (s LoggedProductStore) Create(ctx context.Context, p core.ProductCreateParams) (int64, error) {
	s.logger.InfoContext(ctx, "ProductStore.Create is called.", "p", p)
	return s.products.Create(ctx, p)
}
