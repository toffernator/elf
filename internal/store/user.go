package store

import (
	"elf/internal/auth"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type User struct {
	db *sqlx.DB
}

func NewUser(db *sqlx.DB) *User {
	return &User{db}
}

func (s *User) Create(sub string, name string, email string) (auth.AuthenticatedUser, error) {
	slog.Info("Creating a user", "sub", sub, "name", name, "email", email)

	u := auth.AuthenticatedUser{
		Profile: auth.Profile{
			Sub:   sub,
			Name:  name,
			Email: email,
		},
	}

	res, err := s.db.NamedExec(`INSERT INTO user (sub, name, email)
  VALUES (:sub, :name, :email)`, u)
	if err != nil {
		return auth.AuthenticatedUser{}, fmt.Errorf("User create error: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return auth.AuthenticatedUser{}, fmt.Errorf("User last insert id error: %w", err)
	}
	u.User.Id = int(id)

	return u, nil
}
