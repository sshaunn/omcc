package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"ohmycontrolcenter.tech/omcc/internal/domain/model"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/logger"
)

type CustomerSocialBindingRepository struct {
	db  *gorm.DB
	log logger.Logger
}

func NewCustomerSocialRepository(db *gorm.DB, log logger.Logger) *CustomerSocialBindingRepository {
	return &CustomerSocialBindingRepository{db: db, log: log}
}

func (r *CustomerSocialBindingRepository) Create(ctx context.Context, tx *gorm.DB, binding *model.CustomerSocialBinding) (*model.CustomerSocialBinding, error) {

	db := tx
	if db == nil {
		db = r.db
	}

	if err := db.WithContext(ctx).Omit("Customer", "Platform").Create(binding).Error; err != nil {
		if IsUniqueViolation(err) {
			return nil, ErrSocialBindingExists
		}
		return nil, fmt.Errorf("failed to create social binding: %w", err)
	}
	return binding, nil
}
