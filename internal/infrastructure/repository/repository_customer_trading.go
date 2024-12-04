package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"ohmycontrolcenter.tech/omcc/internal/domain/model"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/logger"
)

type CustomerTradingBindingRepository struct {
	db  *gorm.DB
	log logger.Logger
}

func NewCustomerTradingRepository(db *gorm.DB, log logger.Logger) *CustomerTradingBindingRepository {
	return &CustomerTradingBindingRepository{db: db, log: log}
}

func (r *CustomerTradingBindingRepository) Create(ctx context.Context, tx *gorm.DB, binding *model.CustomerTradingBinding) (*model.CustomerTradingBinding, error) {
	db := tx
	if db == nil {
		db = r.db
	}

	if err := db.WithContext(ctx).Omit("Customer", "Platform").Create(binding).Error; err != nil {
		if IsUniqueViolation(err) {
			return nil, ErrTradingBindingExists
		}
		return nil, fmt.Errorf("failed to create trading binding: %w", err)
	}
	return binding, nil
}
