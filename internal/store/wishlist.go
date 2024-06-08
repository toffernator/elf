package store

import (
	"context"
	"elf/internal/core"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Wishlist struct {
	db *sqlx.DB
}

func NewWishlist(db *sqlx.DB) *Wishlist {
	return &Wishlist{db: db}
}

func (s *Wishlist) Create(name string, ownerId int, products ...core.Product) (core.Wishlist, error) {
	wl := core.Wishlist{
		Name:    name,
		OwnerId: ownerId,
	}

	res, err := s.db.NamedExec(`INSERT INTO wishlist (name, owner_id)
        VALUES (:name, :owner_id)`, wl)
	if err != nil {
		return core.Wishlist{}, fmt.Errorf("Wishlist create error: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return core.Wishlist{}, fmt.Errorf("Wishlist last insert id error: %w", err)
	}

	wl.Id = int(id)

	return wl, nil
}

func (s *Wishlist) ReadById(ctx context.Context, id int) (w core.Wishlist, err error) {
	err = s.db.GetContext(ctx, &w, "SELECT * FROM wishlist WHERE id = $1", id)
	return
}

func (s *Wishlist) ReadByOwner(ctx context.Context, id int) (ws []core.Wishlist, err error) {
	err = s.db.SelectContext(ctx, &ws, "SELECT * FROM wishlist WHERE owner_id = $1", id)
	return
}

func (s *Wishlist) AddProduct(ctx context.Context, id int, p core.Product) (w core.Wishlist, err error) {
	return
}
