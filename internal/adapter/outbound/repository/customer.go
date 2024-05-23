package repository

import (
	"gorm.io/gorm"
	"synapsis-challenge/internal/core/port/outbound/repository"
	"time"
)

type (
	CustomerRepository struct {
		db *gorm.DB
	}

	Customer struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		Password  string    `json:"password"`
		CreatedAt time.Time `json:"created_at"`
	}
)

func NewCustomerRepository(db *gorm.DB) repository.CustomerRepository {
	return &CustomerRepository{
		db: db,
	}
}

func (c *CustomerRepository) FindByID() {
	panic("implement me")
}
