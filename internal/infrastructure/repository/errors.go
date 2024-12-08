package repository

import (
	"errors"
	"strings"
)

var (
	ErrCustomerExists       = errors.New("customer already exists")
	ErrSocialBindingExists  = errors.New("social binding already exists")
	ErrTradingBindingExists = errors.New("trading binding already exists")
	ErrRecordNotFound       = errors.New("record not found")
	ErrGeneralDatabaseError = errors.New("general database error")
)

var (
	ErrInvalidUID                = errors.New("invalid UID format")
	ErrUIDNotFound               = errors.New("UID not found")
	ErrServiceUnavailable        = errors.New("verification service unavailable")
	ErrDatabaseError             = errors.New("database query execution error")
	ErrDuplicatedSocialUserError = errors.New("duplicated social user")
)

func IsUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "UNIQUE constraint failed: ") ||
		strings.Contains(err.Error(), "1062")
}

func isRecordNotFound(err error) bool {
	return strings.Contains(err.Error(), "record not found")
}
