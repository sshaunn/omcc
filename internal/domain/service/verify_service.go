package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"ohmycontrolcenter.tech/omcc/internal/common"
	"ohmycontrolcenter.tech/omcc/internal/domain/model"
	"ohmycontrolcenter.tech/omcc/internal/domain/service/exchange"
	"ohmycontrolcenter.tech/omcc/internal/domain/service/exchange/bitget"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/config"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/database"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/logger"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/repository"
	"ohmycontrolcenter.tech/omcc/util"
)

var (
	ErrInvalidUID         = errors.New("invalid UID format")
	ErrUIDNotFound        = errors.New("UID not found")
	ErrServiceUnavailable = errors.New("verification service unavailable")
	ErrDatabaseError      = errors.New("database query execution error")
)

type VerifyService struct {
	client                   *exchange.Client
	db                       *gorm.DB
	Cfg                      *config.TelegramConfig
	customerRepository       *repository.CustomerRepository
	socialBindingRepository  *repository.CustomerSocialBindingRepository
	tradingBindingRepository *repository.CustomerTradingBindingRepository
	log                      logger.Logger
}

func NewVerifyService(cfg *config.Config, client *exchange.Client, log logger.Logger) *VerifyService {
	db, _ := database.NewMySqlClient(&cfg.Database, log)
	customerRepo := repository.NewCustomerRepository(db, log)
	customerSocialRepo := repository.NewCustomerSocialRepository(db, log)
	customerTradingRepo := repository.NewCustomerTradingRepository(db, log)

	return &VerifyService{
		client:                   client,
		db:                       db,
		Cfg:                      &cfg.Telegram,
		customerRepository:       customerRepo,
		socialBindingRepository:  customerSocialRepo,
		tradingBindingRepository: customerTradingRepo,
		log:                      log,
	}
}

func (v *VerifyService) HandleVerification(ctx context.Context, uid string) error {

	result, err := v.getValidResultByUid(ctx, uid)
	if err != nil {
		return err
	}
	userInfo := ctx.Value("userInfo").(*common.UserInfo)
	customerId := uuid.New().String()
	customer := &model.Customer{
		Id: customerId,
	}
	socialBinding := buildSocialBinding(userInfo, customer)
	tradingBinding := buildTradingBinding(userInfo, customer, result)

	return database.WithTransaction(v.db, func(tx *gorm.DB) error {
		customerCreated, err := v.customerRepository.Create(ctx, tx, customer)
		if err != nil {
			return err
		}
		socialBinding.CustomerID = customerCreated.Id
		tradingBinding.CustomerID = customerCreated.Id

		if _, err = v.socialBindingRepository.Create(ctx, tx, socialBinding); err != nil {
			return err
		}
		if _, err = v.tradingBindingRepository.Create(ctx, tx, tradingBinding); err != nil {
			return err
		}
		return nil
	})
}

func (v *VerifyService) getValidResultByUid(ctx context.Context, uid string) (*bitget.CustomerInfo, error) {
	v.log.Info("Started verifying telegram user uid",
		logger.String("uid", uid),
		logger.Any("userInfo", ctx.Value("userInfo")))

	response, err := v.client.GetCustomerInfo(ctx, uid)
	if err != nil {
		v.log.Error("failed to get customer info",
			logger.String("uid", uid),
			logger.Error(err),
		)
		return nil, ErrServiceUnavailable
	}
	v.log.Info("Completed invoking bitget getCustomerList api by user uid",
		logger.String("uid", uid),
		logger.String("response", response),
		logger.Any("userInfo", ctx.Value("userInfo")))

	result, err := v.getValidResult(response, uid)
	if err != nil {
		return nil, err
	}
	v.log.Info("Completed verifying telegram user uid",
		logger.String("uid", uid),
		logger.Any("userInfo", ctx.Value("userInfo")))
	return result, nil
}

func (v *VerifyService) getValidResult(response string, uid string) (*bitget.CustomerInfo, error) {
	//var result bitget.BaseResponse[[]bitget.CustomerInfo]
	//if result, err := json.Unmarshal([]byte(response), &result); err != nil {

	result, err := util.UnmarshalSafe[bitget.BaseResponse[[]bitget.CustomerInfo]]([]byte(response))
	if err != nil {
		v.log.Error("failed to unmarshal response",
			logger.String("uid", uid),
			logger.Error(err),
			logger.String("response", response),
		)
		return nil, ErrServiceUnavailable
	}

	if len(result.Data) == 0 {
		return nil, ErrUIDNotFound
	}

	return &result.Data[0], nil
}

func buildTelegramSocialPlatform() *model.SocialPlatform {
	return &model.SocialPlatform{
		Id:       common.Telegram.Value(),
		Name:     common.Telegram.Name(),
		IsActive: true,
	}
}

func buildBitgetTradingPlatform() *model.TradingPlatform {
	return &model.TradingPlatform{
		Id:   common.Bitget.Value(),
		Name: common.Bitget.Name(),
	}
}

func buildSocialBinding(userInfo *common.UserInfo, customer *model.Customer) *model.CustomerSocialBinding {
	return &model.CustomerSocialBinding{
		CustomerID:   customer.Id,
		SocialID:     common.Telegram.Value(),
		UserID:       userInfo.UserId,
		Username:     userInfo.Username,
		Firstname:    userInfo.Firstname,
		Lastname:     userInfo.Lastname,
		IsActive:     true,
		MemberStatus: userInfo.MemberStatus,
		Status:       common.Normal,
		Customer:     customer,
		Platform:     buildTelegramSocialPlatform(),
	}
}

func buildTradingBinding(userInfo *common.UserInfo, customer *model.Customer, customerInfo *bitget.CustomerInfo) *model.CustomerTradingBinding {

	registerTime, _ := util.ToIsoTimeFormat(customerInfo.RegisterTime)

	return &model.CustomerTradingBinding{
		CustomerID:   customer.Id,
		TradingID:    common.Bitget.Value(),
		UID:          userInfo.UID,
		RegisterTime: registerTime,
		Customer:     customer,
		Platform:     buildBitgetTradingPlatform(),
	}
}
