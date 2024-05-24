package service

import (
	"context"
	"synapsis-challenge/internal/core/domain"
)

type CustomerService interface {
	RegisterCustomer(ctx context.Context, customer *domain.Customer) (string, error)
	SignIn(ctx context.Context, customer *domain.Customer) (*domain.SignIn, error)
	VerifyToken(token string) error
	SignOut(ctx context.Context, customer *domain.Customer) error
}
