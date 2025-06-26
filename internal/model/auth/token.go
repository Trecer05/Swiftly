package auth

import "time"

type RefreshToken struct {
	Token     string `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}