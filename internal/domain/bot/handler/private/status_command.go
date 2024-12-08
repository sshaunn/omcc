package private

import (
	"context"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"ohmycontrolcenter.tech/omcc/internal/common"
	"ohmycontrolcenter.tech/omcc/internal/domain/service"
	"ohmycontrolcenter.tech/omcc/pkg/exception"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
)

type StatusCommand struct {
	log logger.Logger
	BaseCommand
	statusService service.StatusService
}

func NewCheckCommand(log logger.Logger, checkService service.StatusService) *StatusCommand {
	return &StatusCommand{
		log: log,
		BaseCommand: BaseCommand{
			log:          log,
			validator:    &CommandValidator{2, 2, IsNumeric},
			errorHandler: exception.NewErrorHandler(log),
		},
		statusService: checkService,
	}
}

func (cc *StatusCommand) Handle(c tele.Context) error {
	uid, err := cc.validateUidInput(c, common.StatusCommandName)
	if err != nil {
		return err
	}
	status, err := cc.statusService.Check(context.TODO(), uid)
	return cc.handleResponse(c, err, uid, status)
}

func (cc *StatusCommand) handleResponse(c tele.Context, err error, args ...interface{}) error {
	cc.logResponse(err, args)
	uid := args[0].(string)
	memberStatus := args[1].(common.MemberStatus)
	if err != nil {
		return cc.errorHandler.HandleServiceError(err, map[string]interface{}{
			"uid": uid,
		})
	}
	return c.Send(fmt.Sprintf(common.MemberStatusReplyMessage, uid, memberStatus.Value()))
}
