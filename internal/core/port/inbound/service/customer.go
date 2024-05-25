package service

import (
	"context"
	"synapsis-challenge/internal/core/domain"
)

type CustomerService interface {
	RegisterCustomer(ctx context.Context, customer *domain.Customer) (string, error)
	SignIn(ctx context.Context, customer *domain.Customer) (*domain.SignIn, error)
	VerifyToken(token string) (string, error)
	SignOut(ctx context.Context, customer *domain.Customer) error
	RefreshToken(ctx context.Context, token string) (*domain.SignIn, error)
}
