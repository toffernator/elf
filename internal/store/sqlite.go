package store

import (
	"elf/internal/config"

	"github.com/jmoiron/sqlx"

	_ "github.com/glebarez/go-sqlite"
)

type SqliteStore struct {
	db *sqlx.DB
}

func NewSqliteStore(c config.Db) (*SqliteStore, error) {
	db, err := sqlx.Connect("sqlite", c.Name)
	if err != nil {
		return nil, err
	}

	return &SqliteStore{db}, nil
}
