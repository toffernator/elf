package sqlite

import (
	"context"
	"database/sql"
	"elf/internal/core"

	"github.com/jmoiron/sqlx"
)

type Product struct {
	Id    int64          `db:"id"`
	Name  string         `db:"name"`
	Url   sql.NullString `db:"url"`
	Price sql.NullInt64  `db:"price"`
	// TODO: Enforce enumeration of values
	Currency    sql.NullInt16 `db:"currency"`
	BelongsToId int64         `db:"belongs_to_id"`
}

type ProductStore struct {
	db *sqlx.DB
}

func NewProductStore(db *sqlx.DB) *ProductStore {
	return &ProductStore{db}
}

func (s ProductStore) Create(ctx context.Context, p core.ProductCreateParams) (id int64, err error) {
	res, err := s.db.ExecContext(ctx, `INSERT INTO product (name, url, price, currency, belongs_to_id)
        VALUES($1, $2, $3, $4, $5)`, p.Name, p.Url, p.Price, p.Currency, p.BelongsToId)
	if err != nil {
		return
	}

	id, err = res.LastInsertId()
	return
}
