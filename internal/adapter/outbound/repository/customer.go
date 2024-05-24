package repository

import (
	"context"
	"github.com/jinzhu/copier"
	errors "github.com/rotisserie/eris"
	"gorm.io/gorm"
	"synapsis-challenge/internal/core/domain"
	"synapsis-challenge/internal/core/port/outbound/repository"
	"synapsis-challenge/shared"
	"time"
)

type (
	CustomerRepository struct {
		redRepo repository.RedisRepository
		db      *gorm.DB
	}

	Customer struct {
		ID        string `gorm:"primaryKey"`
		Name      string
		Email     string
		Password  string
		CreatedAt time.Time
	}
)

func NewCustomerRepository(db *gorm.DB, redRepo repository.RedisRepository) repository.CustomerRepository {
	return &CustomerRepository{
		db:      db,
		redRepo: redRepo,
	}
}

func (c *CustomerRepository) FindByID(ctx context.Context, id string) (*domain.Customer, error) {
	var (
		err    error
		model  = new(Customer)
		result domain.Customer
	)

	err = c.db.
		WithContext(ctx).
		Where("id = ?", id).
		First(model).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrNotFound
		}
		return nil, errors.Wrap(err, "FindByID failed First")
	}

	copier.CopyWithOption(&result, model, copier.Option{DeepCopy: true})

	return &result, nil
}

func (c *CustomerRepository) FindByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	var (
		err    error
		model  = new(Customer)
		result domain.Customer
	)

	err = c.db.
		WithContext(ctx).
		Where("email = ?", email).
		First(model).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrNotFound
		}
		return nil, errors.Wrap(err, "FindByEmail failed First")
	}

	copier.CopyWithOption(&result, model, copier.Option{DeepCopy: true})

	return &result, nil
}

func (c *CustomerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	return c.db.WithContext(ctx).Create(customer).Error
}
