package store

import (
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

func (s *Wishlist) Seed() {
	products := []core.Product{
		{Name: "iPad", Url: "www.example.com", Price: 100, Currency: "eur"},
		{Name: "Macbook", Url: "www.example.com", Price: 200, Currency: "eur"},
	}
	_ = []core.Wishlist{
		{Id: 1, Products: products[:], OwnerId: 0},
		{Id: 2, Products: products[:1], OwnerId: 1},
		{Id: 3, Products: products[0:], OwnerId: 0},
	}
	panic("store.Wishlist.Seed() is not implemented")
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
