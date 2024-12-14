package exchange

import (
	"context"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/config"
	"ohmycontrolcenter.tech/omcc/pkg/client"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
	"ohmycontrolcenter.tech/omcc/util"
)

type Client struct {
	BitgetApiClient *client.BitgetClient
	config          *config.BitgetConfig
	log             logger.Logger
}

func NewBitgetClient(config *config.BitgetConfig, log logger.Logger) *Client {
	c := client.NewBitgetClient(config, log)
	return &Client{
		BitgetApiClient: c,
		config:          config,
		log:             log,
	}
}

func (b *Client) GetCustomerInfo(ctx context.Context, uid string) (string, error) {
	params := map[string]string{
		"uid":      uid,
		"pageNo":   "1",
		"pageSize": "100",
	}

	response, err := b.BitgetApiClient.Post(ctx, b.config.CustomerList, params)
	if err != nil {
		return "", err
	}
	return string(response), nil
}

func (b *Client) GetCustomerVolumeList(ctx context.Context, uid string) (string, error) {
	currentDate := util.CurrentDateToEpoch()
	firstDateInCurrentMonth := util.FirstDateInCurrentMonthToEpoch()
	params := map[string]string{
		"uid":       uid,
		"startTime": firstDateInCurrentMonth,
		"endTime":   currentDate,
		"pageNo":    "1",
		"pageSize":  "100",
	}
	b.log.Info("Started invoking Bitget customerVolumeList endpoint",
		logger.String("endpoint", b.config.CustomerList),
		logger.Any("params", params))

	response, err := b.BitgetApiClient.Post(ctx, b.config.CustomerTradeVolume, params)
	if err != nil {
		return "", err
	}
	return string(response), nil
}
