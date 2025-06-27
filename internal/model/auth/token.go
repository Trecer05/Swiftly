package auth

import "time"

type RefreshToken struct {
	ID        int64 `json:"id"`
	Token     string `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}