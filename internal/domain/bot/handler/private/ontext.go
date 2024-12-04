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
	return c.Send("我不是聊天机器人 有活人不聊 你找我干啥 我能跳舞吗 输入指令我才干活 不然我骂街了")
}
