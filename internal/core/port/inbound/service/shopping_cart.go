package service

import (
	"context"
	"synapsis-challenge/internal/core/domain"
)

type ShoppingCartService interface {
	Add(ctx context.Context, shoppingCart *domain.ShoppingCart) (string, error)
	FindAll(ctx context.Context, filter domain.Filter) ([]domain.ShoppingCart, error)
	CountAll(ctx context.Context, filter domain.Filter) (int64, error)
}
