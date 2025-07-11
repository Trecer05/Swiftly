package auth

import (
	"errors"
)

var (
	ErrNoUser = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrEmailExists = errors.New("email already exists")
	ErrNumberExists = errors.New("number already exists")

	ErrUnauthorized = errors.New("unauthorized")
	ErrGroupForbidden = errors.New("you are not allowed to delete this group")
)
