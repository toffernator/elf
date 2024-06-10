package sqlite

import (
	"context"
	"database/sql"
	"elf/internal/core"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Wishlist struct {
	Id      int64          `db:"id"`
	OwnerId int64          `db:"owner_id"`
	Name    string         `db:"name"`
	Image   sql.NullString `db:"image"`
}

type Product struct {
	Id    int64          `db:"id"`
	Name  string         `db:"name"`
	Url   sql.NullString `db:"url"`
	Price sql.NullInt64  `db:"price"`
	// TODO: Enfocrce enumeration of values
	Currency    sql.NullInt16 `db:"currency"`
	BelongsToId int64         `db:"belongs_to_id"`
}

// TODO: Mock with counterfeiter
type WishlistStore struct {
	db *sqlx.DB
}

func NewWishlistStore(db *sqlx.DB) *WishlistStore {
	return &WishlistStore{db: db}
}

func (s *WishlistStore) Create(ctx context.Context, p core.WishlistCreateParams) (core.Wishlist, error) {
	res, err := s.db.Exec(`INSERT INTO wishlist (name, owner_id)
        VALUES ($1, $2)`, p.Name, p.OwnerId)
	if err != nil {
		return core.Wishlist{}, fmt.Errorf("Wishlist create error: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return core.Wishlist{}, fmt.Errorf("Wishlist last insert id error: %w", err)
	}

	return core.Wishlist{Id: id, OwnerId: p.OwnerId, Name: p.Name}, nil
}

func (s *WishlistStore) Read(ctx context.Context, id int64) (c core.Wishlist, err error) {
	var w Wishlist
	err = s.db.GetContext(ctx, &w, "SELECT * FROM wishlist WHERE id = $1", id)
	if err != nil {
		return
	}

	ps, err := s.readProducts(ctx, id)
	if err != nil {
		return
	}

	products := make([]core.Product, len(ps))
	for i, p := range ps {
		products[i] = storeProductToCoreProduct(p)
	}

	c = core.Wishlist{
		Id:       w.Id,
		OwnerId:  w.OwnerId,
		Name:     w.Name,
		Products: products,
	}

	if w.Image.Valid {
		c.Image = w.Image.String
	} else {
		c.Image = ""
	}

	return
}

func (s *WishlistStore) ReadBy(ctx context.Context, p core.WishlistReadByParams) (cs []core.Wishlist, err error) {
	var ws []Wishlist

	err = s.db.SelectContext(ctx, &ws, "SELECT * FROM wishlist WHERE owner_id = $1", p.OwnerId)
	if err != nil {
		return
	}

	cs = make([]core.Wishlist, len(ws))
	for i, w := range ws {
		cs[i] = core.Wishlist{
			Id:      w.Id,
			OwnerId: w.OwnerId,
			Name:    w.Name,
			Image:   w.Image.String,
		}
	}

	return
}

func (s *WishlistStore) readProducts(ctx context.Context, id int64) (ps []Product, err error) {
	err = s.db.SelectContext(ctx, &ps, `SELECT * FROM product WHERE belongs_to_id = $1`, id)
	return
}

func storeProductToCoreProduct(p Product) (c core.Product) {
	c = core.Product{
		Id:   p.Id,
		Name: p.Name,
	}

	if p.Url.Valid {
		c.Url = p.Url.String
	} else {
		c.Url = ""
	}

	if p.Price.Valid {
		c.Price = int(p.Price.Int64)
	} else {
		c.Price = 0
	}

	if p.Currency.Valid {
		c.Currency = core.Currency(p.Currency.Int16)
	} else {
		c.Currency = core.CurrencyNone
	}

	return
}

func (s *WishlistStore) Update(ctx context.Context, p core.WishlistUpdateParams) (c core.Wishlist, err error) {
	return c, errors.New("sqlite.WishlistStore.Update is not yet implemented")
}

func (s *WishlistStore) AddProduct(ctx context.Context, id int, p core.Product) (w core.Wishlist, err error) {
	_, err = s.db.NamedExecContext(ctx, `INSERT INTO product (name, url, price, currency, belongs_to_id)`, p)
	if err != nil {
		return w, err
	}

	return w, err
}
