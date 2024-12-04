package handler

import (
	tele "gopkg.in/telebot.v3"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/logger"
)

type HelpCommand struct {
	log logger.Logger
}

func NewHelpCommand(log logger.Logger) HelpCommand {
	return HelpCommand{log: log}
}

func (h *HelpCommand) Handle(c tele.Context) error {
	help := `可用命令：
			/start - 开始使用
			/help - 显示帮助信息
			/verify - 开始验证流程`
	return c.Send(help)
}
