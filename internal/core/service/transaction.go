package service

import (
	"context"
	"github.com/google/uuid"
	errors "github.com/rotisserie/eris"
	"synapsis-challenge/internal/core/domain"
	"synapsis-challenge/internal/core/port/outbound/registry"
)

type TransactionService struct {
	repositoryRegistry registry.RepositoryRegistry
}

func NewTransactionService(repositoryRegistry registry.RepositoryRegistry) *TransactionService {
	return &TransactionService{
		repositoryRegistry: repositoryRegistry,
	}
}

func (t *TransactionService) CreateTransaction(ctx context.Context, data *domain.TransactionCreateRequest) (*domain.Transaction, error) {
	var (
		err              error
		id               string
		totalTrx         float64
		shoppingCarts    []domain.ShoppingCart
		transactionItems []domain.TransactionItem
		transactionRepo  = t.repositoryRegistry.GetTransactionRepository()
		shoppingCartRepo = t.repositoryRegistry.GetShoppingCartRepository()
		customerId       = ctx.Value("customerId").(string)
	)

	shoppingCarts, err = shoppingCartRepo.FindByIds(ctx, data.ShoppingCartIDs)
	if err != nil {
		return nil, errors.Wrap(err, "CreateTransaction.shoppingCartRepo.FindByIds")
	}

	id = uuid.NewString()

	for _, cart := range shoppingCarts {
		total := cart.Product.Price * cart.Qty
		totalTrx += total

		transactionItems = append(transactionItems, domain.TransactionItem{
			ID:            uuid.NewString(),
			TransactionID: id,
			ProductID:     cart.ProductID,
			Notes:         cart.Notes,
			Price:         cart.Product.Price,
			Qty:           cart.Qty,
			Total:         total,
		})
	}

	insert := &domain.Transaction{
		ID:              id,
		CustomerID:      customerId,
		Subtotal:        totalTrx,
		Discount:        data.Discount,
		Total:           totalTrx - data.Discount,
		TransactionItem: transactionItems,
		TransactionPayment: &domain.TransactionPayment{
			ID:            uuid.NewString(),
			TransactionID: id,
			PaymentType:   data.PaymentType,
		},
	}

	err = transactionRepo.CreateTransaction(ctx, insert)
	if err != nil {
		return nil, errors.Wrap(err, "CreateTransaction.transactionRepo.CreateTransaction")
	}

	return insert, nil
}
