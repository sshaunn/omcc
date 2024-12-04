package handler

import (
	tele "gopkg.in/telebot.v3"
)

type CommandHandler interface {
	Handle(c tele.Context) error
}
