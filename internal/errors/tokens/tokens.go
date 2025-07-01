package tokens

import "errors"

var (
	ErrRefreshTokenExpired = errors.New("refresh token expired")
	ErrAccessTokenExpired = errors.New("access token expired")
)