package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"ohmycontrolcenter.tech/omcc/internal/domain/model"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/logger"
)

type CustomerRepository struct {
	db  *gorm.DB
	log logger.Logger
}

func NewCustomerRepository(db *gorm.DB, log logger.Logger) *CustomerRepository {
	return &CustomerRepository{db, log}
}

func (r *CustomerRepository) Create(ctx context.Context, tx *gorm.DB, customer *model.Customer) (*model.Customer, error) {
	db := tx
	if db == nil {
		db = r.db
	}
	if err := db.WithContext(ctx).Create(customer).Error; err != nil {
		if IsUniqueViolation(err) {
			return nil, ErrCustomerExists
		}
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}
	return customer, nil
}
