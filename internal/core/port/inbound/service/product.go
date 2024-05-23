package service

import (
	"context"
	"synapsis-challenge/internal/core/domain"
)

type ProductService interface {
	FindAll(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, error)
	CountAll(ctx context.Context, filter domain.ProductFilter) (int64, error)
	FindId(ctx context.Context, id string) (*domain.Product, error)
}
