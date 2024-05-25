package repository

import (
	"context"
	"gorm.io/gorm"
	"synapsis-challenge/internal/core/domain"
	"synapsis-challenge/internal/core/port/outbound/repository"
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
