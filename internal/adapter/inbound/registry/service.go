package registry

import (
	ireg "synapsis-challenge/internal/core/port/inbound/registry"
	"synapsis-challenge/internal/core/port/inbound/service"
	oreg "synapsis-challenge/internal/core/port/outbound/registry"
	iservice "synapsis-challenge/internal/core/service"
)

type ServiceRegistry struct {
	customerSvc     service.CustomerService
	productSvc      service.ProductService
	shoppingCartSvc service.ShoppingCartService
	transactionSvc  service.TransactionService
}

func NewServiceRegistry(reg oreg.RepositoryRegistry) ireg.ServiceRegistry {
	return &ServiceRegistry{
		customerSvc:     iservice.NewAccountService(reg),
		productSvc:      iservice.NewProductService(reg),
		shoppingCartSvc: iservice.NewShoppingCartService(reg),
		transactionSvc:  iservice.NewTransactionService(reg),
	}
}

func (s *ServiceRegistry) GetCustomerService() service.CustomerService {
	return s.customerSvc
}

func (s *ServiceRegistry) GetProductService() service.ProductService {
	return s.productSvc
}

func (s *ServiceRegistry) GetShoppingCartService() service.ShoppingCartService {
	return s.shoppingCartSvc
}

func (s *ServiceRegistry) GetTransactionService() service.TransactionService {
	return s.transactionSvc
}
