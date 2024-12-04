package exchange

import (
	"context"
	c "github.com/sshaunn/pkg/bitget-golang-sdk-api/pkg/client"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/config"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/logger"
	"ohmycontrolcenter.tech/omcc/util"
)

type Client struct {
	BitgetApiClient *c.BitgetApiClient
	config          *config.BitgetConfig
	log             logger.Logger
}

func NewBitgetClient(config *config.BitgetConfig, log logger.Logger) *Client {
	client := c.NewBitgetApiClient(config.ApiKey, config.SecretKey, config.Passphrase)
	return &Client{
		BitgetApiClient: client,
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

	b.log.Info("Started invoking Bitget customerList endpoint",
		logger.String("endpoint", b.config.CustomerList),
		logger.Any("params", params))

	return b.BitgetApiClient.Post(b.config.CustomerList, params)
}

func (b *Client) getCustomerVolumeList(ctx context.Context, uid string) (string, error) {
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

	return b.BitgetApiClient.Post(b.config.CustomerList, params)
}
