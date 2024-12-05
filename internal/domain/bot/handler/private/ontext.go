package private

import (
	tele "gopkg.in/telebot.v3"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/logger"
)

type OnTextCommand struct {
	log logger.Logger
}

func NewOnTextCommand(log logger.Logger) OnTextCommand {
	return OnTextCommand{log: log}
}

func (h *OnTextCommand) Handle(c tele.Context) error {
	return c.Send("我不是聊天機器人 有活人不聊 你找我幹啥 我能跳舞嗎 輸入指令我才幹活 不然我會罵街的")
}
