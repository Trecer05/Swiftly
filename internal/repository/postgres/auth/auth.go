package auth

import (
	"database/sql"

	errors "github.com/Trecer05/Swiftly/internal/errors/auth"
	model "github.com/Trecer05/Swiftly/internal/model/auth"

	"golang.org/x/crypto/bcrypt"
)

func (manager *Manager) Login(user *model.User) error {
	var passwordHash string

	if err := manager.Conn.QueryRow("SELECT password_hash, id FROM users WHERE email = $1 OR number = $2", user.Email, user.Number).Scan(&passwordHash, &user.ID); err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrNoUser
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(user.Password)); err != nil {
		return errors.ErrInvalidPassword
	}
	return nil
}