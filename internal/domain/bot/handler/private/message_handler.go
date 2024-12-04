package private

import (
	tele "gopkg.in/telebot.v3"
	"ohmycontrolcenter.tech/omcc/internal/common"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/logger"
	"ohmycontrolcenter.tech/omcc/pkg/exception"
	"strconv"
)

type CommandHandler interface {
	Handle(c tele.Context) error
}

// TelegramCommandHandler extend base CommandHandler interface
type TelegramCommandHandler interface {
	CommandHandler
	validateUidInput(c tele.Context, command string) (string, error)
	buildUserInfoContext(c tele.Context, uid string) *common.UserInfo
	sendProcessingMessage(c tele.Context, text string) error
	handleResponse(c tele.Context, err error, args ...interface{}) error
}

type BaseCommand struct {
	log          logger.Logger
	validator    *CommandValidator
	errorHandler *exception.ErrorHandler
}

func (b *BaseCommand) logResponse(err error, args ...interface{}) {
	if err != nil {
		b.log.Info("command failed",
			logger.Error(err),
			logger.Any("args", args),
		)
	} else {
		b.log.Info("command success",
			logger.Any("args", args),
		)
	}
}

func (b *BaseCommand) validateUidInput(c tele.Context, command string) (string, error) {
	return b.validator.validateUidInput(c, command)
}

func (b *BaseCommand) buildUserInfoContext(c tele.Context, uid string) *common.UserInfo {
	return &common.UserInfo{
		UID:            uid,
		UserId:         strconv.FormatInt(c.Chat().ID, 10),
		Firstname:      c.Chat().FirstName,
		Lastname:       c.Chat().LastName,
		Username:       c.Chat().Username,
		SocialPlatform: common.Telegram,
	}
}

func (b *BaseCommand) sendProcessingMessage(c tele.Context, text string) error {
	if err := c.Send(text); err != nil {
		return &exception.CommandError{
			Message: common.ServerErrorMessage,
			Type:    exception.ErrServiceUnavailable,
		}
	}
	return nil
}

func (b *BaseCommand) generateInviteLinks(c tele.Context) error {

	return nil
}
