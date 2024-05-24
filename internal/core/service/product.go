package service

import (
	"context"
	errors "github.com/rotisserie/eris"
	"synapsis-challenge/internal/core/domain"
	"synapsis-challenge/internal/core/port/outbound/registry"
	"synapsis-challenge/shared"
)

type ProductService struct {
	repositoryRegistry registry.RepositoryRegistry
}

func NewProductService(repositoryRegistry registry.RepositoryRegistry) *ProductService {
	return &ProductService{
		repositoryRegistry: repositoryRegistry,
	}
}

func (i *ProductService) FindAll(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, error) {
	var (
		result      []domain.Product
		err         error
		productRepo = i.repositoryRegistry.GetProductRepository()
	)

	result, err = productRepo.FindAll(ctx, filter)
	if err != nil {
		return result, errors.Wrap(err, "error productRepo.FindAll")
	}

	return result, nil
}

func (i *ProductService) CountAll(ctx context.Context, filter domain.ProductFilter) (int64, error) {
	var (
		count       int64
		err         error
		productRepo = i.repositoryRegistry.GetProductRepository()
	)

	count, err = productRepo.CountAll(ctx, filter)
	if err != nil {
		return count, errors.Wrap(err, "error productRepo.CountAll")
	}

	return count, nil
}

func (i *ProductService) FindId(ctx context.Context, id string) (*domain.Product, error) {
	var (
		result      *domain.Product
		err         error
		productRepo = i.repositoryRegistry.GetProductRepository()
	)

	result, err = productRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return nil, shared.ErrNotFound
		}

		return nil, errors.Wrap(err, "error productRepo.FindById")
	}

	return result, nil
}
