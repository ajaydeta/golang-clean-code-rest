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
		Product    *Product
	}

	ShoppingCarts []ShoppingCart
)

func NewShoppingCartRepository(db *gorm.DB) repository.ShoppingCartRepository {
	return &ShoppingCartRepository{
		db: db,
	}
}

func (l *ShoppingCarts) ToDomain() []domain.ShoppingCart {
	if l == nil {
		return make([]domain.ShoppingCart, 0)
	}

	res := make([]domain.ShoppingCart, len(*l))
	for i, v := range *l {
		copier.CopyWithOption(&res[i], v, copier.Option{DeepCopy: true})
	}

	return res
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

func (s *ShoppingCartRepository) FindById(ctx context.Context, id string) (*domain.ShoppingCart, error) {
	var (
		err    error
		model  = new(ShoppingCart)
		result domain.ShoppingCart
	)

	err = s.db.
		WithContext(ctx).
		Where("id = ?", id).
		Preload("Product").
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

func (s *ShoppingCartRepository) FindByIds(ctx context.Context, ids []string) ([]domain.ShoppingCart, error) {
	var (
		err    error
		model  = new(ShoppingCarts)
		result []domain.ShoppingCart
	)

	err = s.db.
		WithContext(ctx).
		Where("id IN (?)", ids).
		Preload("Product").
		Find(model).
		Error
	if err != nil {
		return nil, errors.Wrap(err, "failed Find")
	}

	result = model.ToDomain()
	return result, nil
}

func (s *ShoppingCartRepository) FindAll(ctx context.Context, customerId string, filter domain.Filter) ([]domain.ShoppingCart, error) {
	var (
		err    error
		model  = new(ShoppingCarts)
		result []domain.ShoppingCart
	)

	err = s.getQueryList(customerId, filter, true).
		WithContext(ctx).
		Select("distinct sp.*").
		Preload("Product").
		Find(model).
		Error
	if err != nil {
		return nil, errors.Wrap(err, "failed Find")
	}

	result = model.ToDomain()
	return result, nil
}

func (s *ShoppingCartRepository) CountAll(ctx context.Context, customerId string, filter domain.Filter) (int64, error) {
	var (
		err   error
		count int64
	)

	err = s.getQueryList(customerId, filter, false).
		WithContext(ctx).
		Select("count(distinct sp.id) as count").
		Count(&count).
		Error
	if err != nil {
		return 0, errors.Wrap(err, "failed Count")
	}

	return count, nil
}

func (s *ShoppingCartRepository) DeleteById(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&ShoppingCart{}).
		Error
}

func (s *ShoppingCartRepository) getQueryList(customerId string, f domain.Filter, withLimit bool) *gorm.DB {
	db := s.db.
		Table("shopping_cart sp").
		Where("sp.customer_id = ?", customerId)

	if f.HasSearch() {
		db = db.
			Joins("JOIN product p ON p.id = sp.product_id").
			Where("p.name LIKE ?", f.Search+"%")
	}

	if withLimit {
		db = f.GetPagination(db).Order("sp.created_at desc")
	}

	return db

}
