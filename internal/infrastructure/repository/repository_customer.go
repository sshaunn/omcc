package repository

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"ohmycontrolcenter.tech/omcc/internal/domain/model"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
)

type CustomerRepositoryImpl struct {
	db  *gorm.DB
	log logger.Logger
}

func (r *CustomerRepositoryImpl) DeleteCustomer(ctx context.Context, tx *gorm.DB, ids []string) ([]string, error) {
	db := tx
	if db == nil {
		db = r.db
	}
	var errorIds []string
	var lastErr error
	for _, id := range ids {
		result := db.WithContext(ctx).Model(&model.CustomerSocialBinding{}).
			Where("customer_id = ?", id).
			Updates(map[string]interface{}{
				"is_active": false,
			})

		if result.Error != nil {
			errorIds = append(errorIds, id)
			lastErr = result.Error
			r.log.Error("failed to update customer status",
				logger.String("customer_id", id),
				logger.Error(result.Error))
		}
	}
	if len(errorIds) > 0 {
		return errorIds, fmt.Errorf("failed to inActive customers %v with error: %v", errorIds, lastErr)
	}
	return ids, nil
}

func (r *CustomerRepositoryImpl) FindAllCustomers(ctx context.Context, tx *gorm.DB, page, limit int) ([]*model.CustomerWithBindings, int64, error) {
	db := tx
	if db == nil {
		db = r.db
	}

	var total int64

	if err := db.WithContext(ctx).Model(&model.Customer{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count customers: %w", err)
	}

	offset := (page - 1) * limit
	var results []*model.CustomerWithBindings

	err := db.WithContext(ctx).
		Table("customers c").
		Select(`
            c.id as c_id,
            c.username as c_username,
            c.created_at as c_created_at,
            c.updated_at as c_updated_at,
            s.user_id as s_user_id,
            s.username as s_username,
            s.firstname as s_firstname,
            s.lastname as s_lastname,
            s.is_active as s_is_active,
            s.status as s_status,
            s.member_status as s_member_status,
            s.created_at as s_created_at,
            t.uid as t_uid,
            t.register_time as t_register_time,
            t.created_at as t_created_at
        `).
		Joins("LEFT JOIN customer_social_bindings s ON c.id = s.customer_id").
		Joins("LEFT JOIN customer_trading_bindings t ON c.id = t.customer_id").
		Offset(offset).
		Limit(limit).
		Order("c.created_at DESC").
		Find(&results).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find customers with bindings: %w", err)
	}

	return results, total, nil

}

func (r *CustomerRepositoryImpl) FindById(ctx context.Context, tx *gorm.DB, id string) (*model.Customer, error) {
	db := tx
	if db == nil {
		db = r.db
	}

	var customer model.Customer
	result := db.WithContext(ctx).
		Where("id = ?", id).
		First(&customer)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("customer not found by uid=%s", id)
		}
		return nil, fmt.Errorf("failed to find customer: %w", result.Error)
	}
	return &customer, nil
}

func (r *CustomerRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, customer *model.Customer) (*model.Customer, error) {
	db := tx
	if db == nil {
		db = r.db
	}
	if err := db.WithContext(ctx).Create(customer).Error; err != nil {
		if IsUniqueViolation(err) {
			return nil, ErrCustomerExists
		}
		return nil, fmt.Errorf("failed to create customer error=%w", err)
	}
	return customer, nil
}

func NewCustomerRepository(db *gorm.DB, log logger.Logger) CustomerRepository {
	return &CustomerRepositoryImpl{db, log}
}
