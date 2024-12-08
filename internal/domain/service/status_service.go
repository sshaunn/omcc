package service

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"ohmycontrolcenter.tech/omcc/internal/common"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/config"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/database"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/repository"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
)

type StatusService struct {
	log                logger.Logger
	db                 *gorm.DB
	tradingBindingRepo repository.CustomerTradingBindingRepository
}

func NewStatusService(cfg *config.Config, log logger.Logger) *StatusService {
	db, _ := database.NewMySqlClient(&cfg.Database, log)
	tradingBindingRepo := repository.NewCustomerTradingRepository(db, log)
	return &StatusService{log, db, tradingBindingRepo}
}

func (cs *StatusService) Check(ctx context.Context, uid string) (common.MemberStatus, error) {
	status, err := cs.tradingBindingRepo.CheckMemberStatus(ctx, cs.db, uid)
	if err != nil {
		return common.Unknown, fmt.Errorf("checking member status failed with uid=%s, error=%w", uid, err)
	}

	return status, nil
}
