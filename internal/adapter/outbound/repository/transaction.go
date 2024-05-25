package repository

import (
	"context"
	"github.com/jinzhu/copier"
	errors "github.com/rotisserie/eris"
	"gorm.io/gorm"
	"synapsis-challenge/internal/core/domain"
	"synapsis-challenge/internal/core/port/outbound/repository"
	"synapsis-challenge/shared"
	"time"
)

type (
	TransactionRepository struct {
		db *gorm.DB
	}

	Transaction struct {
		ID                 string `gorm:"primaryKey"`
		CustomerID         string
		Subtotal           float64
		Discount           float64
		Total              float64
		CreatedAt          time.Time
		TransactionItem    []TransactionItem
		TransactionPayment *TransactionPayment
	}

	TransactionItem struct {
		ID            string `gorm:"primaryKey"`
		TransactionID string
		ProductID     string
		Notes         string
		Price         float64
		Qty           float64
		Total         float64
		CreatedAt     time.Time
		Product       *Product
	}

	TransactionPayment struct {
		ID            string `gorm:"primaryKey"`
		TransactionID string
		PaymentType   string
		Paid          int
		CreatedAt     time.Time
	}
)

func NewTransactionRepository(db *gorm.DB) repository.TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (t *TransactionRepository) CreateTransaction(ctx context.Context, data *domain.Transaction) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		e := tx.Create(data).Error
		if e != nil {
			return e
		}

		var productIds []string
		for _, item := range data.TransactionItem {
			productIds = append(productIds, item.ProductID)
		}

		e = tx.
			Where("customer_id = ?", data.CustomerID).
			Where("product_id IN (?)", productIds).
			Delete(&ShoppingCart{}).
			Error
		if e != nil {
			return e
		}

		return nil
	})
}

func (t *TransactionRepository) GetPaymentTransactionById(ctx context.Context, paymentId string) (*domain.TransactionPayment, error) {
	var (
		err    error
		model  = new(TransactionPayment)
		result domain.TransactionPayment
	)

	err = t.db.WithContext(ctx).
		Where("id = ?", paymentId).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrNotFound
		}
		return nil, errors.Wrap(err, "failed First")
	}

	copier.CopyWithOption(&result, model, copier.Option{DeepCopy: true})

	return &result, nil
}

func (t *TransactionRepository) PayTransaction(ctx context.Context, paymentId string) error {
	return t.db.WithContext(ctx).
		Model(&TransactionPayment{}).
		Where("id = ?", paymentId).
		Updates(map[string]any{"paid": 1}).
		Error
}

func (t *TransactionRepository) FindById(ctx context.Context, id string) (*domain.Transaction, error) {
	var (
		err    error
		model  = new(Transaction)
		result domain.Transaction
	)

	err = t.db.WithContext(ctx).
		Where("id = ?", id).
		Preload("TransactionItem.Product").
		Preload("TransactionPayment").
		First(&model).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrNotFound
		}
		return nil, errors.Wrap(err, "failed First")
	}

	copier.CopyWithOption(&result, model, copier.Option{DeepCopy: true})

	return &result, nil
}
