package private

import (
	"context"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"ohmycontrolcenter.tech/omcc/internal/common"
	"ohmycontrolcenter.tech/omcc/internal/domain/service"
	"ohmycontrolcenter.tech/omcc/pkg/exception"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
	"strconv"
	"strings"
)

type AccountCommand struct {
	log logger.Logger
	bot *tele.Bot
	BaseCommand
	accountService service.AccountCommandService
}

func NewAccountCommand(bot *tele.Bot, log logger.Logger, accountService service.AccountCommandService) *AccountCommand {
	return &AccountCommand{
		log: log,
		bot: bot,
		BaseCommand: BaseCommand{
			log:          log,
			validator:    &CommandValidator{2, 2, IsNumeric},
			errorHandler: exception.NewErrorHandler(log),
		},
		accountService: accountService,
	}
}

func (a *AccountCommand) Handle(c tele.Context) error {
	uid, err := a.validator.validateUidInput(c, common.AccountCommandName)
	if err != nil {
		return err
	}
	userInfo := a.BaseCommand.buildUserInfoContext(c, uid, common.Member)
	err = a.accountService.HandleUpdateCommandService(context.TODO(), uid, userInfo)
	return a.handleResponse(c, err, uid)
}

func (a *AccountCommand) handleResponse(c tele.Context, err error, args ...interface{}) error {
	a.logResponse(err, args)
	uid := args[0].(string)
	if err != nil {
		return a.errorHandler.HandleServiceError(err, map[string]interface{}{
			"uid": uid,
		})
	}
	links, _ := a.generateInviteLinks(a.bot)
	return a.BaseCommand.sendMultipleMessage(c, links)
}

func (a *AccountCommand) generateInviteLinks(b *tele.Bot) ([]string, error) {
	groupIds := strings.Split(a.accountService.Cfg.Telegram.Group, ",")
	var chatList []string
	for _, groupId := range groupIds {
		if parsedId, err := strconv.ParseInt(groupId, 10, 64); err == nil {
			chat, _ := b.ChatByID(parsedId)
			link, _ := b.CreateInviteLink(chat, &tele.ChatInviteLink{
				MemberLimit: 1,
			})
			linkStr := link.InviteLink
			chatList = append(chatList, linkStr)
		} else {
			return nil, fmt.Errorf("error occurred on generating invite links: %v", err)
		}
	}
	return chatList, nil
}
