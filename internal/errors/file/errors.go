package file

import "errors"

var (
	ErrFileNotFound = errors.New("file not found")
	ErrFilesNotFound = errors.New("files not found")
)