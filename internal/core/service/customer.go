package service

import (
	"context"
	"github.com/google/uuid"
	errors "github.com/rotisserie/eris"
	"synapsis-challenge/internal/core/domain"
	inservice "synapsis-challenge/internal/core/port/inbound/service"
	"synapsis-challenge/internal/core/port/outbound/registry"
	"synapsis-challenge/shared"
)

type CustomerService struct {
	repositoryRegistry registry.RepositoryRegistry
}

func NewAccountService(repositoryRegistry registry.RepositoryRegistry) inservice.CustomerService {
	return &CustomerService{
		repositoryRegistry: repositoryRegistry,
	}
}

func (i *CustomerService) RegisterCustomer(ctx context.Context, customer *domain.Customer) (string, error) {
	var (
		id           string
		err          error
		dataCustomer *domain.Customer

		repo = i.repositoryRegistry.GetCustomerRepository()
	)

	dataCustomer, err = repo.FindByEmail(ctx, customer.Email)
	if err != nil && !errors.Is(err, shared.ErrNotFound) {
		return id, errors.Wrap(err, "RegisterCustomer.FindByEmail")
	}

	if dataCustomer != nil {
		return id, shared.ErrAlreadyExist
	}

	customer.Password, err = shared.EncryptPassword(customer.Password)
	if err != nil {
		return id, errors.Wrap(err, "RegisterCustomer.EncryptPassword")
	}

	customer.ID = uuid.NewString()
	id = customer.ID

	err = repo.Create(ctx, customer)
	if err != nil {
		return id, errors.Wrap(err, "RegisterCustomer.Create")
	}

	return id, nil
}

func (i *CustomerService) SignIn(ctx context.Context, customer *domain.Customer) error {
	return nil
}

func (i *CustomerService) SignOut(ctx context.Context, customer *domain.Customer) error {
	return nil
}
