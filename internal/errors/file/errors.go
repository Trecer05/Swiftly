package file

import "errors"

var (
	ErrFileNotFound   = errors.New("file not found")
	ErrFilesNotFound  = errors.New("files not found")
	ErrFolderNotFound = errors.New("folder not found")

	ErrPermissionDenied = errors.New("permission denied")
)
