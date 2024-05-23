package registry

import (
	"gorm.io/gorm"
	outrepo "synapsis-challenge/internal/adapter/outbound/repository"
	"synapsis-challenge/internal/core/port/outbound/registry"
	"synapsis-challenge/internal/core/port/outbound/repository"
)

type RepositoryRegistry struct {
	customerRepo repository.CustomerRepository
	productRepo  repository.ProductRepository
}

func NewRepositoryRegistry(db *gorm.DB) registry.RepositoryRegistry {
	return &RepositoryRegistry{
		customerRepo: outrepo.NewCustomerRepository(db),
		productRepo:  outrepo.NewProductRepository(db),
	}
}

func (r *RepositoryRegistry) GetCustomerRepository() repository.CustomerRepository {
	return r.customerRepo
}

func (r *RepositoryRegistry) GetProductRepository() repository.ProductRepository {
	return r.productRepo
}
