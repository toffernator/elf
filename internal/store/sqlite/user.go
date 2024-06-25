package sqlite

import (
	"context"
	"elf/internal/core"

	"github.com/jmoiron/sqlx"
)

type User struct {
	Id   int64  `db:"id"`
	Name string `db:"name"`
}

type UserStore struct {
	db *sqlx.DB
}

func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{db}
}

func (s *UserStore) Create(ctx context.Context, p core.UserCreateParams) (u core.User, err error) {
	res, err := s.db.Exec(`INSERT INTO user (name) VALUES ($1)`, p.Name)
	if err != nil {
		return u, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return u, err
	}
	u = core.User{
		Id:   id,
		Name: p.Name,
	}
	return u, nil
}

func (s *UserStore) Read(ctx context.Context, id int64) (c core.User, err error) {
	var u User
	s.db.GetContext(ctx, &u, `SELECT * FROM user WHERE id = $1`, id)

	c = core.User{
		Id:   u.Id,
		Name: u.Name,
	}

	return
}
