package service

import (
	"context"
	"ohmycontrolcenter.tech/omcc/internal/domain/service/exchange"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/logger"
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

func (v *VolumeService) HandleVolumeCheck(ctx context.Context, uid string) error {
	return nil
}
