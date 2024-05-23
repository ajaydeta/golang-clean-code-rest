package repository

import (
	"context"
	"github.com/jinzhu/copier"
	errors "github.com/rotisserie/eris"
	"gorm.io/gorm"
	"synapsis-challenge/internal/core/domain"
	"synapsis-challenge/internal/core/port/outbound/repository"
	"time"
)

type (
	ProductRepository struct {
		db *gorm.DB
	}

	Product struct {
		ID              string            `json:"id" gorm:"primaryKey"`
		Name            string            `json:"name"`
		Price           float64           `json:"price"`
		CreatedAt       time.Time         `json:"created_at"`
		ProductCategory []ProductCategory `json:"product_category"`
	}

	ProductCategory struct {
		ProductId  string    `json:"product_id"`
		CategoryId string    `json:"category_id"`
		Category   *Category `json:"category" `
	}

	Category struct {
		ID        string    `json:"id" gorm:"primaryKey"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
	}

	Products []Product
)

func NewProductRepository(db *gorm.DB) repository.ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (l *Products) ToDomain() []domain.Product {
	if l == nil {
		return make([]domain.Product, 0)
	}

	res := make([]domain.Product, len(*l))
	for i, v := range *l {
		copier.CopyWithOption(&res[i], v, copier.Option{DeepCopy: true})
	}

	return res
}

func (i *ProductRepository) FindById(ctx context.Context, id string) (*domain.Product, error) {
	var (
		err    error
		model  = new(Product)
		result domain.Product
	)

	err = i.db.
		WithContext(ctx).
		Where("id = ?", id).
		Preload("ProductCategory.Category").
		First(model).
		Error
	if err != nil {
		return nil, errors.Wrap(err, "failed Find")
	}

	copier.CopyWithOption(&result, model, copier.Option{DeepCopy: true})

	return &result, nil
}

func (i *ProductRepository) FindAll(ctx context.Context, f domain.ProductFilter) ([]domain.Product, error) {
	var (
		err    error
		model  = new(Products)
		result []domain.Product
	)

	err = i.getQueryList(f, true).
		WithContext(ctx).
		Select("distinct p.*").
		Preload("ProductCategory.Category").
		Find(model).
		Error
	if err != nil {
		return nil, errors.Wrap(err, "failed Find")
	}

	result = model.ToDomain()
	return result, nil
}

func (i *ProductRepository) CountAll(ctx context.Context, f domain.ProductFilter) (int64, error) {
	var (
		err   error
		count int64
	)

	err = i.getQueryList(f, false).
		WithContext(ctx).
		Select("count(distinct p.id) as count").
		Count(&count).
		Error
	if err != nil {
		return 0, errors.Wrap(err, "failed Count")
	}

	return count, nil
}

func (i *ProductRepository) getQueryList(f domain.ProductFilter, withLimit bool) *gorm.DB {
	var (
		orderField = map[string]string{
			"price": "p.price",
			"name":  "p.name",
		}
	)

	db := i.db.Table("product p")

	if len(f.CategoryID) > 0 {
		db = db.
			Joins("JOIN product_category pc ON p.id = pc.product_id").
			Where("pc.category_id IN (?)", f.CategoryID)
	}

	if f.HasSearch() {
		db = db.Where("name LIKE ?", f.Search+"%")
	}

	if withLimit {
		db = f.GetSortAndPaginationWithDefaultQuery(db, "p.name ASC", orderField)
	}

	return db

}
