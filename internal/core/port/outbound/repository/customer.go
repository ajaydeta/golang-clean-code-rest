package repository

import (
	"context"
	"synapsis-challenge/internal/core/domain"
)

type CustomerRepository interface {
	FindByID(ctx context.Context, id string) (*domain.Customer, error)
	FindByEmail(ctx context.Context, email string) (*domain.Customer, error)
	Create(ctx context.Context, customer *domain.Customer) error
}
