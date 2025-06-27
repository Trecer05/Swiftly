package auth

import (
	"errors"
)

var (
	ErrNoUser = errors.New("user not found")
	ErrInvalidPassword = errors.New("bad password")
)