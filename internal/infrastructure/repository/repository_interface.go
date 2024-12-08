package repository

import (
	"context"
	"gorm.io/gorm"
	"ohmycontrolcenter.tech/omcc/internal/common"
	"ohmycontrolcenter.tech/omcc/internal/domain/model"
)

type CustomerRepository interface {
	Create(ctx context.Context, tx *gorm.DB, customer *model.Customer) (*model.Customer, error)
}

type CustomerSocialBindingRepository interface {
	Create(ctx context.Context, tx *gorm.DB, binding *model.CustomerSocialBinding) (*model.CustomerSocialBinding, error)
	UpdateUserByUid(ctx context.Context, tx *gorm.DB, uid string, userInfo map[string]interface{}) error
	FindStatusByUid(ctx context.Context, tx *gorm.DB, uid string) (bool, error)
}

type CustomerTradingBindingRepository interface {
	Create(ctx context.Context, tx *gorm.DB, binding *model.CustomerTradingBinding) (*model.CustomerTradingBinding, error)
	CheckMemberStatus(ctx context.Context, tx *gorm.DB, uid string) (common.MemberStatus, error)
	GetTradingBindingById(ctx context.Context, tx *gorm.DB, uid string) (*model.CustomerTradingBinding, error)
}

type TradingHistoryRepository interface {
	Create(ctx context.Context, tx *gorm.DB, tradingHistory *model.TradingHistory) error
	CreateInBatches(ctx context.Context, tx *gorm.DB, batchSize int, tradingHistories []*model.TradingHistory) error
}
