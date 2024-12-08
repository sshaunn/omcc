package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"ohmycontrolcenter.tech/omcc/internal/common"
	"ohmycontrolcenter.tech/omcc/internal/domain/model"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
)

type CustomerTradingBindingRepositoryImpl struct {
	db  *gorm.DB
	log logger.Logger
}

func NewCustomerTradingRepository(db *gorm.DB, log logger.Logger) CustomerTradingBindingRepository {
	return &CustomerTradingBindingRepositoryImpl{db: db, log: log}
}

func (r *CustomerTradingBindingRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, binding *model.CustomerTradingBinding) (*model.CustomerTradingBinding, error) {
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

func (r *CustomerTradingBindingRepositoryImpl) CheckMemberStatus(ctx context.Context, tx *gorm.DB, uid string) (common.MemberStatus, error) {
	db := tx
	if db == nil {
		db = r.db
	}

	var socialBinding model.CustomerSocialBinding
	result := db.WithContext(ctx).
		Select("customer_social_bindings.*").
		Joins("JOIN customer_trading_bindings ON customer_social_bindings.customer_id = customer_trading_bindings.customer_id").
		Where("customer_trading_bindings.uid = ?", uid).
		First(&socialBinding)

	if result.Error != nil {
		if isRecordNotFound(result.Error) {
			return "", ErrRecordNotFound
		}
		return "", fmt.Errorf("failed to find customer social binding with uid=%s, error=%w", uid, result.Error)
	}

	return socialBinding.MemberStatus, nil
}

func (r *CustomerTradingBindingRepositoryImpl) GetTradingBindingById(ctx context.Context, tx *gorm.DB, uid string) (*model.CustomerTradingBinding, error) {
	db := tx
	if db == nil {
		db = r.db
	}
	var binding *model.CustomerTradingBinding
	if err := db.WithContext(ctx).First(&binding, "uid = ?", uid).Error; err != nil {
		return nil, fmt.Errorf("failed to find customer trading binding with uid=%s: %w", uid, err)
	}
	return binding, nil
}
