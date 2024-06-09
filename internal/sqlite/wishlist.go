package sqlite

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

func (s *Wishlist) Create(ctx context.Context, p core.WishlistCreateParams) (core.Wishlist, error) {
	res, err := s.db.NamedExec(`INSERT INTO wishlist (name, owner_id)
        VALUES (:Name, :OwnerId)`, p)
	if err != nil {
		return core.Wishlist{}, fmt.Errorf("Wishlist create error: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return core.Wishlist{}, fmt.Errorf("Wishlist last insert id error: %w", err)
	}

	return core.Wishlist{Id: id, OwnerId: p.OwnerId, Name: p.Name}, nil
}

func (s *Wishlist) ReadById(ctx context.Context, id int) (w core.Wishlist, err error) {
	err = s.db.GetContext(ctx, &w, "SELECT * FROM wishlist WHERE id = $1", id)
	products, err := s.readProducts(ctx, id)
	w.Products = products
	return
}

func (s *Wishlist) ReadByOwner(ctx context.Context, id int) (ws []core.Wishlist, err error) {
	err = s.db.SelectContext(ctx, &ws, "SELECT * FROM wishlist WHERE owner_id = $1", id)
	return
}

func (s *Wishlist) readProducts(ctx context.Context, id int) (ps []core.Product, err error) {
	err = s.db.SelectContext(ctx, &ps, `SELECT * FROM product WHERE belongs_to_id = $1`, id)
	return
}

func (s *Wishlist) AddProduct(ctx context.Context, id int, p core.Product) (w core.Wishlist, err error) {
	_, err = s.db.NamedExecContext(ctx, `INSERT INTO product (name, url, price, currency, belongs_to_id)`, p)
	if err != nil {
		return w, err
	}

	return w, err
}
