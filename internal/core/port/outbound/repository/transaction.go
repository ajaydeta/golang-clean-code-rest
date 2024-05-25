package repository

import (
	"context"
	"synapsis-challenge/internal/core/domain"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, data *domain.Transaction) error
}
