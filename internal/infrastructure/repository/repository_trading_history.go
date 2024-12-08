package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"ohmycontrolcenter.tech/omcc/internal/domain/model"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
)

type TradingHistoryRepositoryImpl struct {
	db  *gorm.DB
	log logger.Logger
}

func NewTradingHistoryRepository(db *gorm.DB, log logger.Logger) TradingHistoryRepository {
	return &TradingHistoryRepositoryImpl{
		db:  db,
		log: log,
	}
}

func (t *TradingHistoryRepositoryImpl) CreateInBatches(ctx context.Context, tx *gorm.DB, batchSize int, tradingHistories []*model.TradingHistory) error {
	db := tx
	if db == nil {
		db = t.db
	}
	if err := db.WithContext(ctx).Omit("TradingBinding").CreateInBatches(tradingHistories, batchSize).Error; err != nil {
		return fmt.Errorf("failed to batch create trading histories: %w", err)
	}
	return nil
}

func (t *TradingHistoryRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, tradingHistory *model.TradingHistory) error {
	db := tx
	if db == nil {
		db = t.db
	}
	if err := db.WithContext(ctx).Omit("TradingBinding").Create(tradingHistory).Error; err != nil {
		return fmt.Errorf("failed to create trading history: %w", err)
	}
	return nil
}
