package store

import (
	"elf/internal/core"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type User struct {
	db *sqlx.DB
}

func NewUser(db *sqlx.DB) *User {
	return &User{db}
}

func (s *User) Create(sub string, name string) (u core.User, err error) {
	slog.Info("Called with", "sub", sub, "name", name)

	u = core.User{
		Sub:  sub,
		Name: name,
	}

	res, err := s.db.NamedExec(`INSERT INTO user (sub, name)
  VALUES (:sub, :name)`, u)
	if err != nil {
		return u, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return u, err
	}
	u.Id = int(id)

	return u, nil
}
