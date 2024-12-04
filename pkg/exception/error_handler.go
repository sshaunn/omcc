package exception

import (
	"errors"
	"fmt"
	"ohmycontrolcenter.tech/omcc/internal/common"
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
			Message: fmt.Sprintf(common.InvalidUidVerifyReplyMessage),
			Type:    ErrInvalidFormat,
		}
	case errors.Is(err, service.ErrServiceUnavailable):
		return &CommandError{
			Message: common.ServerErrorMessage,
			Type:    ErrServiceUnavailable,
		}
	case errors.Is(err, repository.ErrCustomerExists):
		return &CommandError{
			Message: "Unknown Error",
			Type:    ErrInvalidFormat,
		}
	case errors.Is(err, repository.ErrSocialBindingExists):
		return &CommandError{
			Message: common.ExistsSocialUserIdVerifyReplyMessage,
			Type:    ErrInvalidFormat,
		}
	case errors.Is(err, repository.ErrTradingBindingExists):
		return &CommandError{
			Message: common.ExistsUidVerifyReplyMessage,
			Type:    ErrInvalidFormat,
		}
	default:
		h.log.Error("unexpected error during operation",
			logger.Error(err),
			logger.Any("context", context),
		)
		return &CommandError{
			Message: common.InternalServerErrorMessage,
			Type:    ErrInternal,
		}
	}
}
