package bot

import (
	"context"
	tele "gopkg.in/telebot.v3"
)

type Bot interface {
	// Start 启动机器人
	Start(ctx context.Context) error
	// Stop 停止机器人
	Stop()
	// SendMessage 发送消息
	SendMessage(ctx context.Context, chatID int64, message string) error
}

type CommandHandler interface {
	Handle(c tele.Context) error
}
