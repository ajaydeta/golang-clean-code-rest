package repository

import (
	"context"
	"synapsis-challenge/internal/core/domain"
)

type ShoppingCartRepository interface {
	Add(ctx context.Context, shoppingCart *domain.ShoppingCart) error
	Update(ctx context.Context, shoppingCart *domain.ShoppingCart) error
	FindByCustomerProductId(ctx context.Context, customerId, productId string) (*domain.ShoppingCart, error)
}
