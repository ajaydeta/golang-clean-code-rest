package registry

import "synapsis-challenge/internal/core/port/outbound/repository"

type RepositoryRegistry interface {
	GetCustomerRepository() repository.CustomerRepository
}
