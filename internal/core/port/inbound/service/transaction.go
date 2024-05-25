package service

import (
	"context"
	"synapsis-challenge/internal/core/domain"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, data *domain.TransactionCreateRequest) (*domain.Transaction, error)
	PayoffTransaction(ctx context.Context, paymentId string) (*domain.TransactionPayment, error)
}
