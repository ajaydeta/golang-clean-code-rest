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
	ShoppingCartRepository struct {
		db *gorm.DB
	}

	ShoppingCart struct {
		ID         string `gorm:"primaryKey"`
		CustomerID string
		ProductID  string
		Notes      string
		Qty        float64
		CreatedAt  time.Time
	}
)

func NewShoppingCartRepository(db *gorm.DB) repository.ShoppingCartRepository {
	return &ShoppingCartRepository{
		db: db,
	}
}

func (s *ShoppingCartRepository) Add(ctx context.Context, shoppingCart *domain.ShoppingCart) error {
	return s.db.WithContext(ctx).Create(shoppingCart).Error
}

func (s *ShoppingCartRepository) Update(ctx context.Context, shoppingCart *domain.ShoppingCart) error {
	return s.db.WithContext(ctx).
		Model(&ShoppingCart{}).
		Where("id = ?", shoppingCart.ID).
		Updates(map[string]any{
			"notes": shoppingCart.Notes,
			"qty":   shoppingCart.Qty,
		}).
		Error
}

func (s *ShoppingCartRepository) FindByCustomerProductId(ctx context.Context, customerId, productId string) (*domain.ShoppingCart, error) {
	var (
		err    error
		model  = new(ShoppingCart)
		result domain.ShoppingCart
	)

	err = s.db.
		WithContext(ctx).
		Where("customer_id = ?", customerId).
		Where("product_id = ?", productId).
		First(model).
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
