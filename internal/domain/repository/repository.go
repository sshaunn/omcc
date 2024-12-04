package repository

import (
	"context"
	"gorm.io/gorm"
	"ohmycontrolcenter.tech/omcc/internal/domain/model"
)

type CustomerRepository interface {
	Create(ctx context.Context, tx *gorm.DB, customer *model.Customer) (*model.Customer, error)
}

type CustomerSocialBindingRepository interface {
	Create(ctx context.Context, tx *gorm.DB, binding *model.CustomerSocialBinding) (*model.CustomerSocialBinding, error)
	//FindByUserId(ctx context.Context, userId string) (*model.CustomerSocialBinding, error)
	//UpdateStatus(ctx context.Context, userId string, status string) error
	//DeActivate(ctx context.Context, userId string) error
}

type CustomerTradingBindingRepository interface {
	Create(ctx context.Context, tx *gorm.DB, binding *model.CustomerTradingBinding) (*model.CustomerTradingBinding, error)
	//FindByUid(ctx context.Context, uid string) (*model.CustomerTradingBinding, error)
}

type TradingHistoryRepository interface {
	Create(ctx context.Context, tradingHistory *model.TradingHistory) (*model.TradingHistory, error)
	//UpdateVolume(ctx context.Context, bindingId int64, volume float64) error
}
