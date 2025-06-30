package auth

import (
	"errors"
)

var (
	ErrNoUser = errors.New("user not found")
	ErrInvalidPassword = errors.New("bad password")
	ErrEmailExists = errors.New("email already exists")
	ErrNumberExists = errors.New("number already exists")

	ErrUnauthorized = errors.New("unauthorized")
)