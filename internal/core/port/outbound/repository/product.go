package repository

import (
	"context"
	"synapsis-challenge/internal/core/domain"
)

type ProductRepository interface {
	FindById(ctx context.Context, id string) (*domain.Product, error)
	FindAll(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, error)
	CountAll(ctx context.Context, filter domain.ProductFilter) (int64, error)
}
