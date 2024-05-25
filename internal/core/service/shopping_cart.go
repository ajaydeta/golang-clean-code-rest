package service

import (
	"context"
	"github.com/google/uuid"
	errors "github.com/rotisserie/eris"
	"synapsis-challenge/internal/core/domain"
	"synapsis-challenge/internal/core/port/outbound/registry"
	"synapsis-challenge/shared"
)

type ShoppingCartService struct {
	repositoryRegistry registry.RepositoryRegistry
}

func NewShoppingCartService(repositoryRegistry registry.RepositoryRegistry) *ShoppingCartService {
	return &ShoppingCartService{
		repositoryRegistry: repositoryRegistry,
	}
}

func (s *ShoppingCartService) Add(ctx context.Context, param *domain.ShoppingCart) (string, error) {
	var (
		err              error
		id               string
		shoppingCart     *domain.ShoppingCart
		productRepo      = s.repositoryRegistry.GetProductRepository()
		shoppingCartRepo = s.repositoryRegistry.GetShoppingCartRepository()
		customerId       = ctx.Value("customerId").(string)
	)

	_, err = productRepo.FindById(ctx, param.ProductID)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return id, err
		}
		return id, errors.Wrap(err, "Add.productRepo.FindById")
	}

	shoppingCart, err = shoppingCartRepo.FindByCustomerProductId(ctx, customerId, param.ProductID)
	if err != nil && !errors.Is(err, shared.ErrNotFound) {
		return id, errors.Wrap(err, "Add.shoppingCartRepo.FindByCustomerProductId")
	}

	if shoppingCart != nil {
		param.ID = shoppingCart.ID

		err = shoppingCartRepo.Update(ctx, param)
		if err != nil {
			return id, errors.Wrap(err, "Add.shoppingCartRepo.Update")
		}

		return param.ID, nil
	}

	id = uuid.NewString()
	param.ID = id
	param.CustomerID = customerId

	err = shoppingCartRepo.Add(ctx, param)
	if err != nil {
		return id, errors.Wrap(err, "Add.shoppingCartRepo.Add")
	}

	return id, nil
}

func (s *ShoppingCartService) FindAll(ctx context.Context, filter domain.Filter) ([]domain.ShoppingCart, error) {
	var (
		result           []domain.ShoppingCart
		err              error
		shoppingCartRepo = s.repositoryRegistry.GetShoppingCartRepository()
		customerId       = ctx.Value("customerId").(string)
	)

	result, err = shoppingCartRepo.FindAll(ctx, customerId, filter)
	if err != nil {
		return result, errors.Wrap(err, "error FindAll.shoppingCartRepo.FindAll")
	}

	return result, nil
}

func (s *ShoppingCartService) CountAll(ctx context.Context, filter domain.Filter) (int64, error) {
	var (
		count            int64
		err              error
		shoppingCartRepo = s.repositoryRegistry.GetShoppingCartRepository()
		customerId       = ctx.Value("customerId").(string)
	)

	count, err = shoppingCartRepo.CountAll(ctx, customerId, filter)
	if err != nil {
		return count, errors.Wrap(err, "error CountAll.shoppingCartRepo.CountAll")
	}

	return count, nil
}

func (s *ShoppingCartService) Delete(ctx context.Context, id string) error {
	var (
		err              error
		shoppingCartRepo = s.repositoryRegistry.GetShoppingCartRepository()
	)

	_, err = shoppingCartRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return shared.ErrNotFound
		}
		return errors.Wrap(err, "error Delete.shoppingCartRepo.FindById")
	}

	return shoppingCartRepo.DeleteById(ctx, id)
}
