package exception

import "errors"

type ErrorType int

const (
	ErrInvalidFormat ErrorType = iota
	ErrServiceUnavailable
	ErrInternal
)

var (
	ErrSendingMessage = errors.New("failed to send message")
)

type CommandError struct {
	Message string
	Type    ErrorType
}

func (e *CommandError) Error() string {
	return e.Message
}
