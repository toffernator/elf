// TODO: In general for this package: Use transactions
// TODO: In general for this package: Don't SELECT * everywhere
// TODO: In gneral for stores: Make Creates only return the id
package sqlite

import (
	"context"
	"database/sql"
	"elf/internal/core"
	"elf/internal/store"
	"errors"

	"github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
	sqlite3 "modernc.org/sqlite/lib"
)

type Wishlist struct {
	Id      int64          `db:"id"`
	OwnerId int64          `db:"owner_id"`
	Name    string         `db:"name"`
	Image   sql.NullString `db:"image"`
}

type WishlistStore struct {
	db *sqlx.DB
}

func NewWishlistStore(db *sqlx.DB) *WishlistStore {
	return &WishlistStore{db: db}
}

// Create will not populate the the products field of the newly created wishlist.
// To read the products use the Read method following the Create method.
func (s *WishlistStore) Create(ctx context.Context, p core.WishlistCreateParams) (core.Wishlist, error) {
	res, err := s.db.Exec(`INSERT INTO wishlist (name, owner_id)
        VALUES ($1, $2)`, p.Name, p.OwnerId)

	if err != nil {
		switch e := err.(type) {
		case *sqlite.Error:
			switch e.Code() {
			case sqlite3.SQLITE_CONSTRAINT_FOREIGNKEY:
				// Assumes that the wishlist table has exactly one foreign key:
				// The id of the owner in the user table.
				return core.Wishlist{}, store.NewEntityDoesNotExistError("user", p.OwnerId)
			}
		default:
			return core.Wishlist{}, err
		}
	}

	id, err := res.LastInsertId()
	if err != nil {
		return core.Wishlist{}, err
	}

	return core.Wishlist{Id: id, OwnerId: p.OwnerId, Name: p.Name, Products: []core.Product{}}, nil
}

// Read will populate the Products field of the read wishlist.
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

// ReadBy will not populate the the products field read wishlists. To read the
// the products use the Read method on specific wishlists.
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

func (s *WishlistStore) IsOwnedBy(ctx context.Context, wishlistId int64, userId int64) (isOwnedBy bool, err error) {
	err = s.db.GetContext(ctx, &isOwnedBy, `IF EXISTS(SELECT COUNT(1) FROM wishlist
        WHERE id = $1 AND owner_id = $2)`, wishlistId, userId)
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
