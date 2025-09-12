package global

import "errors"

var (
	ErrNoJsonData = errors.New("no json data found")

	ErrNoPhotos = errors.New("no photos found")
)