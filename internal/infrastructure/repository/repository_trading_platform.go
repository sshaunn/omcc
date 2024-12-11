package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"ohmycontrolcenter.tech/omcc/internal/domain/model"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
)

type TradingPlatformRepositoryImpl struct {
	db  *gorm.DB
	log logger.Logger
}

func (t TradingPlatformRepositoryImpl) FindById(ctx context.Context, tx *gorm.DB, id string) (*model.TradingPlatform, error) {
	db := tx
	if db == nil {
		db = t.db
	}
	var platform model.TradingPlatform
	result := db.WithContext(ctx).Where("id = ?", id).First(&platform)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find out the trading platform: %w", result.Error)
	}
	return &platform, nil
}

func NewTradingPlatformRepository(db *gorm.DB, log logger.Logger) TradingPlatformRepository {
	return TradingPlatformRepositoryImpl{db: db, log: log}
}
