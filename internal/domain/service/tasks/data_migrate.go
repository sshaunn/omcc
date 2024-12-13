package tasks

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm"
	"log"
	"ohmycontrolcenter.tech/omcc/internal/common"
	"ohmycontrolcenter.tech/omcc/internal/domain/model"
	"ohmycontrolcenter.tech/omcc/internal/domain/service/exchange/bitget"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/config"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/database"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/repository"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
	"ohmycontrolcenter.tech/omcc/util"
	"time"
)

type MigrateService struct {
	pgDB                     *gorm.DB
	myDB                     *gorm.DB
	customerRepository       repository.CustomerRepository
	socialBindingRepository  repository.CustomerSocialBindingRepository
	tradingBindingRepository repository.CustomerTradingBindingRepository
}

func NewMigrateService(cfg *config.Config, log logger.Logger) *MigrateService {
	pgDB, _ := database.NewPostgresClient(log)
	myDB, _ := database.NewMySqlClient(&cfg.Database, log)
	customerRepo := repository.NewCustomerRepository(myDB, log)
	customerSocialRepo := repository.NewCustomerSocialRepository(myDB, log)
	customerTradingRepo := repository.NewCustomerTradingRepository(myDB, log)
	return &MigrateService{
		pgDB:                     pgDB,
		myDB:                     myDB,
		customerRepository:       customerRepo,
		socialBindingRepository:  customerSocialRepo,
		tradingBindingRepository: customerTradingRepo,
	}
}

type customerBinding struct {
	Tgid         string    `json:"tgid"`
	Uid          string    `json:"uid"`
	RegisterTime time.Time `json:"register_time"`
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
}

func (m *MigrateService) Migrate() error {

	list, _ := GetOldDatabaseRecords(m.pgDB)
	return database.WithTransaction(m.myDB, func(tx *gorm.DB) error {
		for _, binding := range list {
			customerId := uuid.New().String()
			customer := &model.Customer{
				Id: customerId,
			}
			userInfo := &common.UserInfo{
				UID:            binding.Uid,
				UserId:         binding.Tgid,
				Firstname:      binding.Firstname,
				Lastname:       binding.Lastname,
				MemberStatus:   common.Member,
				SocialPlatform: common.Telegram,
			}
			epochStr := fmt.Sprintf("%d", binding.RegisterTime.UnixMilli())
			result := &bitget.CustomerInfo{
				Uid:          binding.Uid,
				RegisterTime: epochStr,
			}
			log.Printf("customerInfo=%s", result)
			socialBinding := buildSocialBinding(userInfo, customer)
			tradingBinding := buildTradingBinding(userInfo, customer, result)

			customerCreated, err := m.customerRepository.Create(context.TODO(), tx, customer)
			if err != nil {
				return err
			}
			socialBinding.CustomerID = customerCreated.Id
			tradingBinding.CustomerID = customerCreated.Id

			if _, err = m.socialBindingRepository.Create(context.TODO(), tx, socialBinding); err != nil {
				return err
			}
			if _, err = m.tradingBindingRepository.Create(context.TODO(), tx, tradingBinding); err != nil {
				return err
			}

		}

		return nil
	})
}

func GetOldDatabaseRecords(pgDB *gorm.DB) ([]customerBinding, error) {
	var whitelistedCustomers []customerBinding
	err := pgDB.Table("customers_backup").
		Select("tgid, uid, register_time, firstname, lastname").
		Where("tgid IS NOT NULL AND uid IS NOT NULL").
		Find(&whitelistedCustomers).Error
	if err != nil {
		return nil, err
	}
	return whitelistedCustomers, nil
}

func GetTelegramUser(bot *tele.Bot, groupId, userID int64) (*tele.User, error) {
	// 创建 chat 和 user recipients
	chat := &tele.Chat{ID: groupId}
	user := &tele.User{ID: userID}

	// 获取成员信息
	member, err := bot.ChatMemberOf(chat, user)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	return member.User, nil
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

	registerTime, _ := util.ToIsoTimeStringFormat(customerInfo.RegisterTime)

	return &model.CustomerTradingBinding{
		CustomerID:   customer.Id,
		TradingID:    common.Bitget.Value(),
		UID:          userInfo.UID,
		RegisterTime: registerTime,
		Customer:     customer,
		Platform:     buildBitgetTradingPlatform(),
	}
}
