package registry

import (
	ireg "synapsis-challenge/internal/core/port/inbound/registry"
	"synapsis-challenge/internal/core/port/inbound/service"
	oreg "synapsis-challenge/internal/core/port/outbound/registry"
	iservice "synapsis-challenge/internal/core/service"
)

type ServiceRegistry struct {
	customerSvc service.CustomerService
}

func (s ServiceRegistry) GetCustomerService() service.CustomerService {
	panic("implement me")
}

func NewServiceRegistry(reg oreg.RepositoryRegistry) ireg.ServiceRegistry {
	return &ServiceRegistry{
		customerSvc: iservice.NewAccountService(reg),
	}
}
