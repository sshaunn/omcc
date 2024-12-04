package handler

import (
	tele "gopkg.in/telebot.v3"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/logger"
)

type StartCommand struct {
	log logger.Logger
}

func NewStartCommand(log logger.Logger) StartCommand {
	return StartCommand{log: log}
}

func (h *StartCommand) Handle(c tele.Context) error {
	return c.Send("欢迎使用！输入 /help 查看可用命令。")
}
