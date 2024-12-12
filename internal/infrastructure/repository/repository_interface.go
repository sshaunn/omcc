package repository

import (
	"context"
	"gorm.io/gorm"
	"ohmycontrolcenter.tech/omcc/internal/common"
	"ohmycontrolcenter.tech/omcc/internal/domain/model"
)

type CustomerRepository interface {
	Create(ctx context.Context, tx *gorm.DB, customer *model.Customer) (*model.Customer, error)
	FindById(ctx context.Context, tx *gorm.DB, id string) (*model.Customer, error)
	FindAllCustomers(ctx context.Context, tx *gorm.DB, page, limit int) ([]*model.CustomerWithBindings, int64, error)
	DeleteCustomer(ctx context.Context, tx *gorm.DB, ids []string) ([]string, error)
}

type CustomerSocialBindingRepository interface {
	Create(ctx context.Context, tx *gorm.DB, binding *model.CustomerSocialBinding) (*model.CustomerSocialBinding, error)
	UpdateUserByUid(ctx context.Context, tx *gorm.DB, uid string, userInfo map[string]interface{}) error
	FindStatusByUid(ctx context.Context, tx *gorm.DB, uid string) (bool, error)
	FindSocialBindingByCustomerId(ctx context.Context, tx *gorm.DB, customerId string) (*model.CustomerSocialBinding, error)
	UpdateCustomerStatus(ctx context.Context, tx *gorm.DB, customerID string, socialID string, status string, memberStatus common.MemberStatus) error
}

type CustomerTradingBindingRepository interface {
	Create(ctx context.Context, tx *gorm.DB, binding *model.CustomerTradingBinding) (*model.CustomerTradingBinding, error)
	CheckMemberStatus(ctx context.Context, tx *gorm.DB, uid string) (common.MemberStatus, error)
	FindTradingBindingByUid(ctx context.Context, tx *gorm.DB, uid string) (*model.CustomerInfoResponse, error)
}

type TradingHistoryRepository interface {
	Create(ctx context.Context, tx *gorm.DB, tradingHistory *model.TradingHistory) error
	CreateInBatches(ctx context.Context, tx *gorm.DB, batchSize int, tradingHistories []*model.TradingHistory) error
}

type TradingPlatformRepository interface {
	FindById(ctx context.Context, tx *gorm.DB, id string) (*model.TradingPlatform, error)
}

type SocialPlatformRepository interface {
	FindById(ctx context.Context, tx *gorm.DB, id string) (*model.SocialPlatform, error)
}
