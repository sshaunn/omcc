package private

import (
	tele "gopkg.in/telebot.v3"
	"ohmycontrolcenter.tech/omcc/internal/common"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
)

type HelpCommand struct {
	log logger.Logger
}

func NewHelpCommand(log logger.Logger) HelpCommand {
	return HelpCommand{log: log}
}

func (h *HelpCommand) Handle(c tele.Context) error {
	return c.Send(common.HelpMessage)
}
