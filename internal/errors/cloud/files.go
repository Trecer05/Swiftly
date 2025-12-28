package cloud

import "errors"

var (
	ErrEmptyMetadata = errors.New("empty metadata")
	ErrFileTooLarge = errors.New("file too large")
)
