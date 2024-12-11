package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"ohmycontrolcenter.tech/omcc/internal/domain/model"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
)

type SocialPlatformRepositoryImpl struct {
	db  *gorm.DB
	log logger.Logger
}

func (s SocialPlatformRepositoryImpl) FindById(ctx context.Context, tx *gorm.DB, id string) (*model.SocialPlatform, error) {
	db := tx
	if db == nil {
		db = s.db
	}
	var platform *model.SocialPlatform
	result := db.WithContext(ctx).Where("id = ?", id).First(&platform)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find social platform: %w", result.Error)
	}
	return platform, nil
}

func NewSocialPlatformRepository(db *gorm.DB, log logger.Logger) SocialPlatformRepository {
	return SocialPlatformRepositoryImpl{db: db, log: log}
}
