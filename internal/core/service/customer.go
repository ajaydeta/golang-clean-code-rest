package service

import (
	inservice "synapsis-challenge/internal/core/port/inbound/service"
	"synapsis-challenge/internal/core/port/outbound/registry"
)

type CustomerService struct {
	repositoryRegistry registry.RepositoryRegistry
}

func NewAccountService(repositoryRegistry registry.RepositoryRegistry) inservice.CustomerService {
	return &CustomerService{
		repositoryRegistry: repositoryRegistry,
	}
}

func (i *CustomerService) FindOne() {

}
