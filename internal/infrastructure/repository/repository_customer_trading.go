package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"ohmycontrolcenter.tech/omcc/internal/common"
	"ohmycontrolcenter.tech/omcc/internal/domain/model"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
	"time"
)

type CustomerTradingBindingRepositoryImpl struct {
	db  *gorm.DB
	log logger.Logger
}

func NewCustomerTradingRepository(db *gorm.DB, log logger.Logger) CustomerTradingBindingRepository {
	return &CustomerTradingBindingRepositoryImpl{db: db, log: log}
}

func (r *CustomerTradingBindingRepositoryImpl) Create(
	ctx context.Context,
	tx *gorm.DB,
	binding *model.CustomerTradingBinding) (*model.CustomerTradingBinding, error) {
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

func (r *CustomerTradingBindingRepositoryImpl) CheckMemberStatus(
	ctx context.Context,
	tx *gorm.DB, uid string) (common.MemberStatus, error) {
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

func (r *CustomerTradingBindingRepositoryImpl) FindTradingBindingByUid(
	ctx context.Context,
	tx *gorm.DB, uid string) (*model.CustomerInfoResponse, error) {
	db := tx
	if db == nil {
		db = r.db
	}

	var result struct {
		// Customer
		CustomerID        string    `gorm:"column:customer_id"`
		CustomerUsername  string    `gorm:"column:customer_username"`
		CustomerCreatedAt time.Time `gorm:"column:customer_created_at"`

		// Social
		UserID          string    `gorm:"column:social_user_id"`
		SocialUsername  string    `gorm:"column:social_username"`
		Firstname       string    `gorm:"column:social_firstname"`
		Lastname        string    `gorm:"column:social_lastname"`
		IsActive        bool      `gorm:"column:social_is_active"`
		Status          string    `gorm:"column:social_status"`
		MemberStatus    string    `gorm:"column:social_member_status"`
		SocialPlatform  string    `gorm:"column:social_type"`
		SocialCreatedAt time.Time `gorm:"column:social_created_at"`

		// Trading
		TradingUID       string    `gorm:"column:trade_uid"`
		RegisterTime     time.Time `gorm:"column:trade_register_time"`
		TradingPlatform  string    `gorm:"column:trade_type"`
		TradingCreatedAt time.Time `gorm:"column:trade_created_at"`
	}

	dbResult := db.WithContext(ctx).Table("customer_trading_bindings as t").
		Select(`
           c.id as customer_id,
           c.username as customer_username,
           c.created_at as customer_created_at,
           s.user_id as social_user_id,
           s.username as social_username,
           s.firstname as social_firstname,
           s.lastname as social_lastname,
           s.is_active as social_is_active,
           s.status as social_status,
           s.member_status as social_member_status,
           s.created_at as social_created_at,
           sp.name as social_type,
           t.uid as trade_uid,
           t.register_time as trade_register_time,
           t.created_at as trade_created_at,
           tp.name as trade_type`).
		Joins("JOIN customers c ON t.customer_id = c.id").
		Joins(`JOIN customer_social_bindings s ON c.id = s.customer_id`).
		Joins(`JOIN social_platforms sp ON s.social_id = sp.id`).
		Joins(`JOIN trading_platforms tp ON t.trading_id = tp.id`).
		Where("t.uid = ?", uid).
		First(&result)

	//var binding *model.CustomerTradingBinding
	//if err := db.WithContext(ctx).First(&binding, "uid = ?", uid).Error; err != nil {
	//	return nil, fmt.Errorf("failed to find customer trading binding with uid=%s: %w", uid, err)
	//}
	if dbResult.Error != nil {
		return nil, fmt.Errorf("failed to find customer social binding with uid=%s, err=%w", uid, dbResult.Error)
	}
	return &model.CustomerInfoResponse{
		Customer: model.CustomerInfo{
			ID:        result.CustomerID,
			Username:  result.CustomerUsername,
			CreatedAt: result.CustomerCreatedAt,
		},
		SocialAccountInfo: model.CustomerSocialInfo{
			UserID:       result.UserID,
			Username:     result.SocialUsername,
			Firstname:    result.Firstname,
			Lastname:     result.Lastname,
			IsActive:     result.IsActive,
			Status:       result.Status,
			MemberStatus: result.MemberStatus,
			SocialType:   result.SocialPlatform,
			CreatedAt:    result.SocialCreatedAt,
		},
		TradingAccountInfo: model.CustomerTradingInfo{
			UID:          result.TradingUID,
			RegisterTime: result.RegisterTime.String(),
			Platform:     result.TradingPlatform,
			CreatedAt:    result.TradingCreatedAt,
		},
	}, nil
}
