package service

import (
	"context"
	"fmt"
	"math/big"
	"ohmycontrolcenter.tech/omcc/internal/domain/service/exchange"
	"ohmycontrolcenter.tech/omcc/internal/domain/service/exchange/bitget"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/logger"
	"ohmycontrolcenter.tech/omcc/util"
)

type VolumeService struct {
	client *exchange.Client
	log    logger.Logger
}

func NewVolumeService(client *exchange.Client, log logger.Logger) *VolumeService {
	return &VolumeService{
		client: client,
		log:    log,
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
		return nil, ErrServiceUnavailable
	}
	v.log.Info("Completed invoking bitget getCustomerVolumeList api by user uid",
		logger.String("uid", uid),
		logger.String("response", response),
		logger.Any("userInfo", ctx.Value("userInfo")))

	result, err := v.getValidResponse(response, uid)
	if err != nil {
		return nil, err
	}

	return sumMoneyDecimals(result)
}

func (v *VolumeService) getValidResponse(response string, uid string) ([]bitget.CustomerVolume, error) {
	result, err := util.UnmarshalSafe[bitget.BaseResponse[[]bitget.CustomerVolume]]([]byte(response))
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

	return result.Data, nil
}

func sumMoneyDecimals(volumeList []bitget.CustomerVolume) (*big.Float, error) {
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
