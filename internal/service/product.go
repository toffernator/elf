package service

import (
	"context"
	"elf/internal/core"
	"log/slog"
)

type ProductStore interface {
	Create(ctx context.Context, p core.ProductCreateParams) (int64, error)
}

type ProductService struct {
	store ProductStore
}

func NewProductService(u ProductStore) *ProductService {
	return &ProductService{store: u}
}

func (s *ProductService) Create(ctx context.Context, p core.ProductCreateParams) (id int64, err error) {
	slog.Info("ProductService.Create is called", "p", p)
	err = p.Validate()
	if err != nil {
		return
	}

	id, err = s.store.Create(ctx, p)
	if err != nil {
		return
	}

	return
}
