package service

import (
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"ohmycontrolcenter.tech/omcc/internal/common"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/config"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/database"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/repository"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
)

type AccountCommandService struct {
	db                        *gorm.DB
	Cfg                       *config.Config
	customerSocialBindingRepo repository.CustomerSocialBindingRepository
}

func NewAccountService(cfg *config.Config, log logger.Logger) *AccountCommandService {
	db, _ := database.NewMySqlClient(&cfg.Database, log)
	return &AccountCommandService{
		db:                        db,
		Cfg:                       cfg,
		customerSocialBindingRepo: repository.NewCustomerSocialRepository(db, log),
	}
}

func (u *AccountCommandService) HandleUpdateCommandService(ctx context.Context, uid string, userInfo *common.UserInfo) error {
	active, _ := u.customerSocialBindingRepo.FindStatusByUid(ctx, u.db, uid)
	if !active {
		return repository.ErrInvalidUID
	}
	update := map[string]interface{}{
		"user_id":   userInfo.UserId,
		"username":  userInfo.Username,
		"firstname": userInfo.Firstname,
		"lastname":  userInfo.Lastname,
	}
	return u.customerSocialBindingRepo.UpdateUserByUid(ctx, u.db, userInfo.UID, update)
}
