package postgres

import (
	"errors"
)

var (
	ErrFailedBeginTx = errors.New("failed to begin transaction")
	ErrFailedCommitTx = errors.New("failed to commit transaction")
)