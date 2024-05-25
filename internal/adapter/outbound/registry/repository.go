package registry

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	outrepo "synapsis-challenge/internal/adapter/outbound/repository"
	"synapsis-challenge/internal/core/port/outbound/registry"
	"synapsis-challenge/internal/core/port/outbound/repository"
)

type RepositoryRegistry struct {
	redisRepo        repository.RedisRepository
	customerRepo     repository.CustomerRepository
	productRepo      repository.ProductRepository
	shoppingCardRepo repository.ShoppingCartRepository
	transactionRepo  repository.TransactionRepository
}

func NewRepositoryRegistry(rdb *redis.Client, db *gorm.DB) registry.RepositoryRegistry {
	redisRepo := outrepo.NewRedisRepository(rdb)

	return &RepositoryRegistry{
		redisRepo:        redisRepo,
		customerRepo:     outrepo.NewCustomerRepository(db, redisRepo),
		productRepo:      outrepo.NewProductRepository(db),
		shoppingCardRepo: outrepo.NewShoppingCartRepository(db),
		transactionRepo:  outrepo.NewTransactionRepository(db),
	}
}

func (r *RepositoryRegistry) GetRedisRepository() repository.RedisRepository {
	return r.redisRepo
}

func (r *RepositoryRegistry) GetCustomerRepository() repository.CustomerRepository {
	return r.customerRepo
}

func (r *RepositoryRegistry) GetProductRepository() repository.ProductRepository {
	return r.productRepo
}

func (r *RepositoryRegistry) GetShoppingCartRepository() repository.ShoppingCartRepository {
	return r.shoppingCardRepo
}

func (r *RepositoryRegistry) GetTransactionRepository() repository.TransactionRepository {
	return r.transactionRepo
}
