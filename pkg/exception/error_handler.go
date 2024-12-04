package exception

import (
	"errors"
	"fmt"
	"ohmycontrolcenter.tech/omcc/internal/domain/service"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/logger"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/repository"
)

type ErrorHandler struct {
	log logger.Logger
}

func NewErrorHandler(log logger.Logger) *ErrorHandler {
	return &ErrorHandler{log: log}
}

// HandleServiceError error handler
func (h *ErrorHandler) HandleServiceError(err error, context map[string]interface{}) *CommandError {
	switch {
	case errors.Is(err, service.ErrUIDNotFound):
		return &CommandError{
			Message: fmt.Sprintf("❌ UID %s 不存在，请检查后重试。", context["uid"]),
			Type:    ErrInvalidFormat,
		}
	case errors.Is(err, service.ErrServiceUnavailable):
		return &CommandError{
			Message: "验证服务暂时不可用，请稍后重试",
			Type:    ErrServiceUnavailable,
		}
	case errors.Is(err, repository.ErrCustomerExists):
		return &CommandError{
			Message: "该用户已存在",
			Type:    ErrInvalidFormat,
		}
	case errors.Is(err, repository.ErrSocialBindingExists):
		return &CommandError{
			Message: "该社交账号已被绑定",
			Type:    ErrInvalidFormat,
		}
	case errors.Is(err, repository.ErrTradingBindingExists):
		return &CommandError{
			Message: "UID已被绑定 请联系",
			Type:    ErrInvalidFormat,
		}
	default:
		h.log.Error("unexpected error during operation",
			logger.Error(err),
			logger.Any("context", context),
		)
		return &CommandError{
			Message: "操作过程中发生错误，请稍后重试",
			Type:    ErrInternal,
		}
	}
}
