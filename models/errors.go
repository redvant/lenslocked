package models

import "errors"

var (
	ErrEmailTaken  = errors.New("models: email address is already in use")
	ErrBadPassword = errors.New("models: incorrect password")

	ErrInvalidPwResetToken = errors.New("models: invalid password reset token")
	ErrExpiredPwResetToken = errors.New("models: expired password reset token")

	ErrNotFound = errors.New("models: resource could not be found")
)
