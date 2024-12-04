package handler

import (
	"context"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"ohmycontrolcenter.tech/omcc/internal/common"
	"ohmycontrolcenter.tech/omcc/internal/domain/service"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/logger"
	"ohmycontrolcenter.tech/omcc/pkg/exception"
	"regexp"
	"strconv"
	"strings"
)

type VerifyCommand struct {
	log           logger.Logger
	verifyService service.VerifyService
	errorHandler  *exception.ErrorHandler
}

func NewVerifyCommand(log logger.Logger, verifyService service.VerifyService) *VerifyCommand {
	return &VerifyCommand{
		log:           log,
		verifyService: verifyService,
		errorHandler:  exception.NewErrorHandler(log),
	}
}

func (h *VerifyCommand) Handle(c tele.Context) error {
	uid, err := h.validateInput(c)
	if err != nil {
		return err
	}

	userInfo := h.buildUserInfoContext(c, uid)
	if err = h.sendProcessingMessage(c); err != nil {
		return err
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "userInfo", userInfo)

	if err = h.verifyService.HandleVerification(ctx, uid); err != nil {
		return h.errorHandler.HandleServiceError(err, map[string]interface{}{
			"uid":      uid,
			"userInfo": userInfo,
		})
	}

	return h.sendResponse(c, err, uid, userInfo)
}

// validateInput 验证输入参数
func (h *VerifyCommand) validateInput(c tele.Context) (string, error) {
	args := strings.Fields(c.Text())
	if len(args) != 2 {
		return "", &exception.CommandError{
			Message: "请使用正确的格式：/verify <uid>\n示例：/verify 123456",
			Type:    exception.ErrInvalidFormat,
		}
	}

	uid := args[1]
	if !h.isValidUID(uid) {
		return "", &exception.CommandError{
			Message: "无效的 UID 格式，UID 必须为数字\n示例：/verify 123456",
			Type:    exception.ErrInvalidFormat,
		}
	}

	return uid, nil
}

// buildUserInfoContext 创建用户信息
func (h *VerifyCommand) buildUserInfoContext(c tele.Context, uid string) *common.UserInfo {
	return &common.UserInfo{
		UID:            uid,
		UserId:         strconv.FormatInt(c.Chat().ID, 10),
		Firstname:      c.Chat().FirstName,
		Lastname:       c.Chat().LastName,
		Username:       c.Chat().Username,
		SocialPlatform: common.Telegram,
	}
}

// sendProcessingMessage 发送处理中消息
func (h *VerifyCommand) sendProcessingMessage(c tele.Context) error {
	if err := c.Send("正在验证 UID，请稍候..."); err != nil {
		return &exception.CommandError{
			Message: "验证服务暂时不可用，请稍后重试",
			Type:    exception.ErrServiceUnavailable,
		}
	}

	return nil
}

func (h *VerifyCommand) isValidUID(uid string) bool {
	return regexp.MustCompile(`^\d+$`).MatchString(uid)
}

// sendResponse 发送响应
func (h *VerifyCommand) sendResponse(c tele.Context, err error, uid string, userInfo *common.UserInfo) error {
	if err != nil {
		h.log.Info("uid verification failed",
			logger.String("uid", uid),
			logger.Any("userInfo", userInfo),
		)
		return c.Send(fmt.Sprintf("❌ UID %s 验证失败，该 UID 无效。", uid))
	}

	h.log.Info("processing verify command success",
		logger.String("uid", uid))
	return c.Send(fmt.Sprintf("✅ UID %s 验证成功！", uid))

}
