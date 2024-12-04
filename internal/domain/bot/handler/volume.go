package handler

import (
	tele "gopkg.in/telebot.v3"
	"ohmycontrolcenter.tech/omcc/internal/domain/service"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/logger"
	"ohmycontrolcenter.tech/omcc/pkg/exception"
)

type VolumeCommand struct {
	log           logger.Logger
	volumeService service.VolumeService
	errorHandler  *exception.ErrorHandler
}

func NewVolumeCommand(log logger.Logger, volumeService service.VolumeService) *VolumeCommand {
	return &VolumeCommand{
		log:           log,
		volumeService: volumeService,
		errorHandler:  exception.NewErrorHandler(log),
	}
}

func (v *VolumeCommand) Handle(c tele.Context) error {
	return nil
}
