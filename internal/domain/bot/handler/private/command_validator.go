package private

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
	"ohmycontrolcenter.tech/omcc/internal/common"
	"ohmycontrolcenter.tech/omcc/pkg/exception"
	"strconv"
	"strings"
)

type CommandValidator struct {
	MinArgs     int
	MaxArgs     int
	ValidateArg func(string) bool
}

func (v *CommandValidator) ValidateGeneralCommand(text string, commandName string) ([]string, error) {
	args := strings.Fields(text)

	if len(args) < v.MinArgs || (v.MaxArgs > 0 && len(args) > v.MaxArgs) {
		return nil, &exception.CommandError{
			Message: fmt.Sprintf(common.InvalidCommandFormatMessage,
				commandName, commandName),
			Type: exception.ErrInvalidFormat,
		}
	}

	// if has args validator then execute func ValidateArg
	if v.ValidateArg != nil {
		for i := 1; i < len(args); i++ {
			if !v.ValidateArg(args[i]) {
				return nil, &exception.CommandError{
					Message: fmt.Sprintf(common.InvalidUIDFormatMessage,
						commandName),
					Type: exception.ErrInvalidFormat,
				}
			}
		}
	}

	return args[1:], nil
}

func IsNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// validateUidInput validate input args
func (v *CommandValidator) validateUidInput(c tele.Context, commandName string) (string, error) {
	args, err := v.ValidateGeneralCommand(c.Text(), commandName)
	if err != nil {
		return "", err
	}
	return args[0], nil
}
