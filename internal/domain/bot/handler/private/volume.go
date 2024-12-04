package private

import (
	"context"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"math/big"
	"ohmycontrolcenter.tech/omcc/internal/common"
	"ohmycontrolcenter.tech/omcc/internal/domain/service"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/logger"
	"ohmycontrolcenter.tech/omcc/pkg/exception"
)

type VolumeCommand struct {
	log logger.Logger
	BaseCommand
	volumeService service.VolumeService
	//validator     *CommandValidator
	//errorHandler  *exception.ErrorHandler
}

func NewVolumeCommand(log logger.Logger, volumeService service.VolumeService) *VolumeCommand {
	return &VolumeCommand{
		log: log,
		BaseCommand: BaseCommand{
			log:          log,
			validator:    &CommandValidator{2, 2, IsNumeric},
			errorHandler: exception.NewErrorHandler(log),
		},
		volumeService: volumeService,
	}
}

func (v *VolumeCommand) Handle(c tele.Context) error {
	uid, err := v.validator.validateUidInput(c, common.VolumeCommandName)
	if err != nil {
		return err
	}
	volume, err := v.volumeService.HandleVolumeCheck(context.TODO(), uid)

	return v.handleResponse(c, err, uid, volume)
}

func (v *VolumeCommand) handleResponse(c tele.Context, err error, args ...interface{}) error {
	v.logResponse(err, args)
	uid := args[0].(string)
	volume := args[1].(*big.Float)
	if err != nil {
		return v.errorHandler.HandleServiceError(err, map[string]interface{}{
			"uid": uid,
		})
	}
	return c.Send(fmt.Sprintf(common.SuccessVolumeReplyMessage, volume))
}

//// sendResponse 发送响应
//func (h *VerifyCommand) sendResponse(c tele.Context, err error, uid string, userInfo *common.UserInfo) error {
//	if err != nil {
//		h.log.Info("uid verification failed",
//			logger.String("uid", uid),
//			logger.Any("userInfo", userInfo),
//		)
//		return c.Send(fmt.Sprintf("❌ UID %s 验证失败，该 UID 无效。", uid))
//	}
//
//	h.log.Info("processing verify command success",
//		logger.String("uid", uid))
//	return c.Send(fmt.Sprintf("✅ UID %s 验证成功！", uid))
//
//}
