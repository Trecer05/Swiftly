package auth

import (
	"database/sql"
	"errors"

	"github.com/Trecer05/Swiftly/internal/errors/auth"
	pgErrors "github.com/Trecer05/Swiftly/internal/errors/postgres"
	model "github.com/Trecer05/Swiftly/internal/model/auth"

	"golang.org/x/crypto/bcrypt"
)

func (manager *Manager) Login(user *model.User) error {
	var passwordHash string

	if err := manager.Conn.QueryRow("SELECT password_hash, id FROM users WHERE email = $1 OR number = $2", user.Email, user.Number).Scan(&passwordHash, &user.ID); err != nil {
		if err == sql.ErrNoRows {
			return auth.ErrNoUser
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(user.Password)); err != nil {
		return auth.ErrInvalidPassword
	}
	return nil
}

func (manager *Manager) Register(user *model.User) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := manager.Conn.QueryRow("INSERT INTO users (email, number, password_hash) VALUES ($1, $2, $3) RETURNING id", user.Email, user.Number, string(passwordHash)).Scan(&user.ID); err != nil {
		if errors.Is(err, pgErrors.ErrUsersEmailDuplicate) {
			return auth.ErrEmailExists
		} else if errors.Is(err, pgErrors.ErrUsersNumberDuplicate) {
			return auth.ErrNumberExists
		}

		return err
	}
	return nil
}
