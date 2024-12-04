package exception

type ErrorType int

const (
	ErrInvalidFormat ErrorType = iota
	ErrServiceUnavailable
	ErrInternal
)

type CommandError struct {
	Message string
	Type    ErrorType
}

func (e *CommandError) Error() string {
	return e.Message
}
