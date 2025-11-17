package service

import "errors"

var (
	ErrPhoneExists          = errors.New("phone already exists")
	ErrEmailExists          = errors.New("email already exists")
	ErrUserExists           = errors.New("user already exists")
	ErrPhoneOrEmailRequired = errors.New("phone or email required")
	ErrPasswordRequired     = errors.New("password required")

	// Логин ошибки
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrInvalidLoginFormat = errors.New("invalid login format")
	ErrUserNotConfirmed   = errors.New("user not confirmed")
	ErrAccountLocked      = errors.New("account locked")
)
