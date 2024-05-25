package service

import (
	"context"
	"synapsis-challenge/internal/core/domain"
)

type ShoppingCartService interface {
	Add(ctx context.Context, shoppingCart *domain.ShoppingCart) (string, error)
}
