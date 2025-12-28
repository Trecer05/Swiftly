package cloud

import "errors"

var (
	ErrEmptyMetadata = errors.New("empty metadata")
	ErrFileTooLarge = errors.New("file too large")
	ErrFileNotFound = errors.New("file not found")
	ErrFolderNotFound = errors.New("folder not found")
)
