package registry

import (
	outrepo "synapsis-challenge/internal/adapter/outbound/repository"
	"synapsis-challenge/internal/core/port/outbound/registry"
	"synapsis-challenge/internal/core/port/outbound/repository"
)

type RepositoryRegistry struct {
	customerRepo repository.CustomerRepository
}

func NewRepositoryRegistry() registry.RepositoryRegistry {
	return &RepositoryRegistry{
		customerRepo: outrepo.NewCustomerRepository(),
	}
}

func (r RepositoryRegistry) GetCustomerRepository() repository.CustomerRepository {
	panic("implement me")
}
