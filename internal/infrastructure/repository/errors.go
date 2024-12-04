package repository

import (
	"errors"
	"strings"
)

var (
	ErrCustomerExists       = errors.New("customer already exists")
	ErrSocialBindingExists  = errors.New("social binding already exists")
	ErrTradingBindingExists = errors.New("trading binding already exists")
	ErrGeneralDatabaseError = errors.New("general database error")
)

func IsUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "UNIQUE constraint failed: ") ||
		strings.Contains(err.Error(), "1062")
}
