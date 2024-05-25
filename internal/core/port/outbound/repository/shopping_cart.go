package repository

import (
	"context"
	"synapsis-challenge/internal/core/domain"
)

type ShoppingCartRepository interface {
	Add(ctx context.Context, shoppingCart *domain.ShoppingCart) error
	Update(ctx context.Context, shoppingCart *domain.ShoppingCart) error
	FindByCustomerProductId(ctx context.Context, customerId, productId string) (*domain.ShoppingCart, error)
	FindById(ctx context.Context, id string) (*domain.ShoppingCart, error)
	FindAll(ctx context.Context, customerId string, filter domain.Filter) ([]domain.ShoppingCart, error)
	CountAll(ctx context.Context, customerId string, filter domain.Filter) (int64, error)
	DeleteById(ctx context.Context, id string) error
}
