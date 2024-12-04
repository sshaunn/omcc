package private

import (
	"context"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"ohmycontrolcenter.tech/omcc/internal/common"
	"ohmycontrolcenter.tech/omcc/internal/domain/service"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/logger"
	"ohmycontrolcenter.tech/omcc/pkg/exception"
	"strconv"
	"strings"
	"sync"
)

type VerifyCommand struct {
	//log           logger.Logger
	bot *tele.Bot
	BaseCommand
	verifyService service.VerifyService
}

func NewVerifyCommand(bot *tele.Bot, log logger.Logger, verifyService service.VerifyService) *VerifyCommand {
	return &VerifyCommand{
		bot: bot,
		BaseCommand: BaseCommand{
			log:          log,
			validator:    &CommandValidator{2, 2, IsNumeric},
			errorHandler: exception.NewErrorHandler(log),
		},
		verifyService: verifyService,
	}
}

func (h *VerifyCommand) Handle(c tele.Context) error {
	uid, err := h.validator.validateUidInput(c, common.VerifyCommandName)
	if err != nil {
		return err
	}

	userInfo := h.buildUserInfoContext(c, uid)
	if err = h.sendProcessingMessage(c, common.ProcessingMessage); err != nil {
		return err
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "userInfo", userInfo)

	err = h.verifyService.HandleVerification(ctx, uid)
	return h.handleResponse(c, err, uid, userInfo)
}

func (h *VerifyCommand) handleResponse(c tele.Context, err error, args ...interface{}) error {
	h.logResponse(err, args)
	uid := args[0].(string)
	userInfo := args[1]
	if err != nil {
		return h.errorHandler.HandleServiceError(err, map[string]interface{}{
			"uid":      uid,
			"userInfo": userInfo,
		})
	}

	linkList, err := h.generateInviteLinks(h.bot)

	err = c.Send(common.SuccessVerifyReplyMessage)
	return h.concurrentlySendMessage(c, linkList)
}

func (h *VerifyCommand) generateInviteLinks(b *tele.Bot) ([]string, error) {
	groupIds := strings.Split(h.verifyService.Cfg.Group, ",")
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

func (h *VerifyCommand) concurrentlySendMessage(c tele.Context, messages []string) error {
	var wg sync.WaitGroup
	errs := make(chan error, len(messages))

	for _, msg := range messages {
		wg.Add(1)
		go func(message string) {
			defer wg.Done()
			if err := c.Send(message); err != nil {
				errs <- err
			}
		}(msg)
	}

	go func() {
		wg.Wait()
		close(errs)
	}()

	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
