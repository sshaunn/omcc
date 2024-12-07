package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"ohmycontrolcenter.tech/omcc/internal/domain/model"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
)

type CustomerSocialBindingRepositoryImpl struct {
	db  *gorm.DB
	log logger.Logger
}

func (r *CustomerSocialBindingRepositoryImpl) FindStatusByUid(ctx context.Context, tx *gorm.DB, uid string) (bool, error) {
	db := tx
	if db == nil {
		db = r.db
	}
	type Active struct {
		IsActive bool `gorm:"column:is_active"`
	}
	var value Active
	subQuery := db.Table("customer_trading_bindings").Select("customer_id").Where("uid = ?", uid)
	result := db.WithContext(ctx).
		Table("customer_social_bindings").
		Select("is_active").
		Where("customer_id IN (?)", subQuery).
		First(&value)
	if result.Error != nil {
		return false, fmt.Errorf("failed to find customer social bindings: %w", result.Error)
	}

	return value.IsActive, nil
}

func (r *CustomerSocialBindingRepositoryImpl) UpdateUserByUid(ctx context.Context, tx *gorm.DB, uid string, userInfo map[string]interface{}) error {
	db := tx
	if db == nil {
		db = r.db
	}
	subQuery := db.Table("customer_trading_bindings").Select("customer_id").Where("uid = ?", uid)
	result := db.WithContext(ctx).
		Table("customer_social_bindings").
		Where("customer_id IN (?)", subQuery).
		Updates(userInfo)

	if result.Error != nil {
		return fmt.Errorf("failed to find customer social bindings with uid=%s, error=%w", uid, result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrDuplicatedSocialUserError
	}
	return nil
}

func (r *CustomerSocialBindingRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, binding *model.CustomerSocialBinding) (*model.CustomerSocialBinding, error) {
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

func NewCustomerSocialRepository(db *gorm.DB, log logger.Logger) CustomerSocialBindingRepository {
	return &CustomerSocialBindingRepositoryImpl{db: db, log: log}
}
