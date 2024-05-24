package registry

import "synapsis-challenge/internal/core/port/outbound/repository"

type RepositoryRegistry interface {
	GetCustomerRepository() repository.CustomerRepository
	GetProductRepository() repository.ProductRepository
	GetRedisRepository() repository.RedisRepository
}
