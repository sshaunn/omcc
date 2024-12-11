package service

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"math/big"
	"ohmycontrolcenter.tech/omcc/internal/domain/service/exchange"
	"ohmycontrolcenter.tech/omcc/internal/domain/service/exchange/bitget"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/config"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/database"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/repository"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
	"ohmycontrolcenter.tech/omcc/util"
)

type VolumeService struct {
	client                 *exchange.Client
	db                     *gorm.DB
	customerTradingBinding repository.CustomerTradingBindingRepository
	tradingHistory         repository.TradingHistoryRepository
	log                    logger.Logger
}

func NewVolumeService(cfg *config.DatabaseConfig, client *exchange.Client, log logger.Logger) *VolumeService {
	db, _ := database.NewMySqlClient(cfg, log)
	return &VolumeService{
		client:                 client,
		db:                     db,
		customerTradingBinding: repository.NewCustomerTradingRepository(db, log),
		tradingHistory:         repository.NewTradingHistoryRepository(db, log),
		log:                    log,
	}
}

func (v *VolumeService) HandleVolumeCheck(ctx context.Context, uid string) (*big.Float, error) {
	return v.volumeCalculator(ctx, uid)
}

func (v *VolumeService) volumeCalculator(ctx context.Context, uid string) (*big.Float, error) {
	v.log.Info("Started volume telegram user uid",
		logger.String("uid", uid),
		logger.Any("userInfo", ctx.Value("userInfo")))

	response, err := v.client.GetCustomerVolumeList(ctx, uid)
	if err != nil {
		v.log.Error("failed to get customer info",
			logger.String("uid", uid),
			logger.Error(err),
		)
		return nil, repository.ErrServiceUnavailable
	}
	v.log.Info("Completed invoking bitget getCustomerVolumeList api by user uid",
		logger.String("uid", uid),
		logger.String("response", response),
		logger.Any("userInfo", ctx.Value("userInfo")))

	result, err := v.getValidResponse(response, uid)
	if err != nil {
		return nil, err
	}

	//go func() {
	//	if err := v.SaveTradingHistories(context.TODO(), uid, result); err != nil {
	//		v.log.Error("Failed to save trade histories",
	//			logger.String("uid", uid),
	//			logger.Error(err),
	//		)
	//	}
	//}()

	return sumMoneyDecimals(result)
}

func (v *VolumeService) getValidResponse(response string, uid string) ([]*bitget.CustomerVolume, error) {
	result, err := util.UnmarshalSafe[bitget.BaseResponse[[]*bitget.CustomerVolume]]([]byte(response))
	if err != nil {
		v.log.Error("failed to unmarshal response",
			logger.String("uid", uid),
			logger.Error(err),
			logger.String("response", response),
		)
		return nil, repository.ErrServiceUnavailable
	}

	if len(result.Data) == 0 {
		return nil, repository.ErrUIDNotFound
	}

	return result.Data, nil
}

func sumMoneyDecimals(volumeList []*bitget.CustomerVolume) (*big.Float, error) {
	sum := new(big.Float)

	for _, volume := range volumeList {
		// Create a new big.Rat for each float string
		f := new(big.Float)
		// Set the value of the big.Rat from the string float
		_, err := fmt.Sscan(volume.Volume, f)
		if err != nil {
			return nil, err
		}

		// Add the current big.Rat to the sum
		sum.Add(sum, f)
	}

	return sum, nil
}

//func (v *VolumeService) SaveTradingHistories(ctx context.Context, uid string, results []*bitget.CustomerVolume) error {
//	histories := make([]*model.TradingHistory, len(results))
//	binding, err := v.customerTradingBinding.FindTradingBindingByUid(ctx, v.db, uid)
//	if err != nil {
//		return err
//	}
//	for i, result := range results {
//		volume, _ := strconv.ParseFloat(result.Volume, 64)
//		date, _ := util.ToIsoTimeFormat(result.Time)
//		histories[i] = &model.TradingHistory{
//			BindingID:      binding.ID,
//			Volume:         volume,
//			TimePeriod:     common.DailyTrading,
//			TradingDate:    date,
//			TradingBinding: binding,
//		}
//	}
//
//	return database.WithTransaction(v.db, func(tx *gorm.DB) error {
//		return v.tradingHistory.CreateInBatches(ctx, tx, len(results), histories)
//	})
//}
