package auth

import (
	"time"

	"github.com/Trecer05/Swiftly/internal/model/auth"
)

func (manager *Manager) SaveRefreshToken(token string, userId int) error {
	expiredAt := time.Now().Add(time.Hour * 24 * 7).Unix()

	_, err := manager.Conn.Exec(`INSERT INTO user_tokens (user_id, refresh, expired_at) VALUES ($1, $2, $3)`, userId, token, expiredAt)
	return err
}

func (manager *Manager) GetRefreshToken(userId int) (auth.RefreshToken, error) {
	var token auth.RefreshToken

	if err := manager.Conn.QueryRow(`SELECT refresh, expired_at FROM user_tokens WHERE user_id = $1`, userId).Scan(&token.Token, &token.ExpiredAt); err != nil {
		return auth.RefreshToken{}, err
	}

	return token, nil
}

func (manager *Manager) DeleteRefreshToken(userId int) error {
	if _, err := manager.Conn.Exec(`DELETE FROM user_tokens WHERE user_id = $1`, userId); err != nil {
		return err
	}

	return nil
}

func (manager *Manager) UpdateRefreshToken(userId int, token string) error {
	if _, err := manager.Conn.Exec(`UPDATE user_tokens SET expired_at = $1, refresh = $2 WHERE user_id = $3`, time.Now().Add(time.Hour * 24 * 7).Unix(), token, userId); err != nil {
		return err
	}

	return nil
}