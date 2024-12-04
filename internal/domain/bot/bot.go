package bot

import (
	"context"
	tele "gopkg.in/telebot.v3"
)

type Bot interface {
	// Start bot
	Start(ctx context.Context) error
	// Stop bot
	Stop()
	// SendMessage send telegram message
	SendMessage(ctx context.Context, chatID int64, message string) error
}

type CommandHandler interface {
	Handle(c tele.Context) error
}
