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
	ProductRepository struct {
		db    *gorm.DB
		cache repository.RedisRepository
	}

	Product struct {
		ID              string `gorm:"primaryKey"`
		Name            string
		Price           float64
		CreatedAt       time.Time
		ProductCategory []ProductCategory
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

const (
	productPrefixCacheKey         = "cacheProduct:"
	productPrefixCacheKeyFindById = productPrefixCacheKey + "FindById:"
	productPrefixCacheKeyFindAll  = productPrefixCacheKey + "FindAll:"
	productPrefixCacheKeyCountAll = productPrefixCacheKey + "CountAll:"
)

func NewProductRepository(db *gorm.DB, cache repository.RedisRepository) repository.ProductRepository {
	return &ProductRepository{
		db:    db,
		cache: cache,
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

	cacheKey := shared.GetCacheKey(productPrefixCacheKeyFindById, id)
	if err = i.cache.GetObject(cacheKey, &result); err == nil {
		return &result, nil
	}

	err = i.db.
		WithContext(ctx).
		Where("id = ?", id).
		Preload("ProductCategory.Category").
		First(model).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrNotFound
		}
		return nil, errors.Wrap(err, "failed First")
	}

	copier.CopyWithOption(&result, model, copier.Option{DeepCopy: true})

	go i.cache.SetObj(cacheKey, result, shared.CacheTtl)

	return &result, nil
}

func (i *ProductRepository) FindAll(ctx context.Context, f domain.ProductFilter) ([]domain.Product, error) {
	var (
		err    error
		model  = new(Products)
		result []domain.Product
	)

	cacheKey := shared.GetCacheKey(productPrefixCacheKeyFindAll, f)
	if err = i.cache.GetObject(cacheKey, &result); err == nil {
		return result, nil
	}

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

	go i.cache.SetObj(cacheKey, result, shared.CacheTtl)

	return result, nil
}

func (i *ProductRepository) CountAll(ctx context.Context, f domain.ProductFilter) (int64, error) {
	var (
		err   error
		count int64
	)

	cacheKey := shared.GetCacheKey(productPrefixCacheKeyCountAll, f)
	if count, err = i.cache.GetInt64(cacheKey); err == nil {
		return count, nil
	}

	err = i.getQueryList(f, false).
		WithContext(ctx).
		Select("count(distinct p.id) as count").
		Count(&count).
		Error
	if err != nil {
		return 0, errors.Wrap(err, "failed Count")
	}

	go i.cache.SetInt64(cacheKey, count, shared.CacheTtl)

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
