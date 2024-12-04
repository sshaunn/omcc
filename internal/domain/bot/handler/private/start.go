package private

import (
	tele "gopkg.in/telebot.v3"
	"ohmycontrolcenter.tech/omcc/internal/common"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/logger"
)

type StartCommand struct {
	log logger.Logger
}

func NewStartCommand(log logger.Logger) StartCommand {
	return StartCommand{log: log}
}

func (h *StartCommand) Handle(c tele.Context) error {
	return c.Send(common.WelcomeMessage)
}
