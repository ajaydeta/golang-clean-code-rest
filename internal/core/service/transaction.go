package service

import (
	"context"
	"github.com/google/uuid"
	errors "github.com/rotisserie/eris"
	"synapsis-challenge/internal/core/domain"
	"synapsis-challenge/internal/core/port/outbound/registry"
	"synapsis-challenge/shared"
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

func (t *TransactionService) PayoffTransaction(ctx context.Context, paymentId string) (*domain.TransactionPayment, error) {
	var (
		err             error
		paymentData     *domain.TransactionPayment
		transactionRepo = t.repositoryRegistry.GetTransactionRepository()
	)

	paymentData, err = transactionRepo.GetPaymentTransactionById(ctx, paymentId)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return nil, shared.ErrNotFound
		}
		return nil, errors.Wrap(err, "PayoffTransaction.transactionRepo.GetPaymentTransactionById")
	}

	if paymentData.Paid == 1 {
		return nil, shared.ErrAlreadyPaid
	}

	err = transactionRepo.PayTransaction(ctx, paymentId)
	if err != nil {
		return nil, errors.Wrap(err, "PayoffTransaction.transactionRepo.PayTransaction")
	}

	return paymentData, nil
}

func (t *TransactionService) FindId(ctx context.Context, id string) (*domain.Transaction, error) {
	var (
		result          *domain.Transaction
		err             error
		transactionRepo = t.repositoryRegistry.GetTransactionRepository()
	)

	result, err = transactionRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return nil, shared.ErrNotFound
		}

		return nil, errors.Wrap(err, "error transactionRepo.FindById")
	}

	return result, nil
}
