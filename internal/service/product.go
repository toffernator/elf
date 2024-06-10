package service

import (
	"context"
	"elf/internal/core"
)

type ProductStore interface {
	Create(ctx context.Context, p core.ProductCreateParams) (core.Product, error)
}

// TODO: Mock with counterfeiter
type ProductService struct {
	store ProductStore
}

func NewProductService(u ProductStore) *ProductService {
	return &ProductService{store: u}
}

func (s *ProductService) Create(ctx context.Context, p core.ProductCreateParams) (prod core.Product, err error) {
	err = p.Validate()
	if err != nil {
		return prod, err
	}

	prod, err = s.store.Create(ctx, p)
	if err != nil {
		return prod, err
	}

	return prod, nil
}
