package registry

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	outrepo "synapsis-challenge/internal/adapter/outbound/repository"
	"synapsis-challenge/internal/core/port/outbound/registry"
	"synapsis-challenge/internal/core/port/outbound/repository"
)

type RepositoryRegistry struct {
	customerRepo repository.CustomerRepository
	productRepo  repository.ProductRepository
}

func NewRepositoryRegistry(rdb *redis.Client, db *gorm.DB) registry.RepositoryRegistry {
	redisRepo := outrepo.NewRedisRepository(rdb)

	return &RepositoryRegistry{
		customerRepo: outrepo.NewCustomerRepository(db, redisRepo),
		productRepo:  outrepo.NewProductRepository(db),
	}
}

func (r *RepositoryRegistry) GetCustomerRepository() repository.CustomerRepository {
	return r.customerRepo
}

func (r *RepositoryRegistry) GetProductRepository() repository.ProductRepository {
	return r.productRepo
}
