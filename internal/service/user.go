package service

import (
	"context"
	"elf/internal/core"
)

type UserStore interface {
	Create(ctx context.Context, p core.UserCreateParams) (core.User, error)
	Read(ctx context.Context, id int64) (core.User, error)
}

type UserService struct {
	store UserStore
}

func NewUserService(u UserStore) *UserService {
	return &UserService{store: u}
}

func (u *UserService) Create(ctx context.Context, p core.UserCreateParams) (usr core.User, err error) {
	err = p.Validate()
	if err != nil {
		return usr, err
	}

	usr, err = u.store.Create(ctx, p)
	if err != nil {
		return usr, err
	}

	return usr, nil
}

func (u *UserService) Read(ctx context.Context, id int64) (usr core.User, err error) {
	usr, err = u.store.Read(ctx, id)
	if err != nil {
		return usr, err
	}

	return usr, nil
}
