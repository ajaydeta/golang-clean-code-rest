package registry

import (
	ireg "synapsis-challenge/internal/core/port/inbound/registry"
	"synapsis-challenge/internal/core/port/inbound/service"
	oreg "synapsis-challenge/internal/core/port/outbound/registry"
	iservice "synapsis-challenge/internal/core/service"
)

type ServiceRegistry struct {
	customerSvc service.CustomerService
	productSvc  service.ProductService
}

func NewServiceRegistry(reg oreg.RepositoryRegistry) ireg.ServiceRegistry {
	return &ServiceRegistry{
		customerSvc: iservice.NewAccountService(reg),
		productSvc:  iservice.NewProductService(reg),
	}
}

func (s *ServiceRegistry) GetCustomerService() service.CustomerService {
	return s.customerSvc
}

func (s *ServiceRegistry) GetProductService() service.ProductService {
	return s.productSvc
}
