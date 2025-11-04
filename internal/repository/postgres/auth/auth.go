package auth

import (
	"database/sql"

	"github.com/Trecer05/Swiftly/internal/errors/auth"
	pgErrors "github.com/Trecer05/Swiftly/internal/errors/postgres"
	model "github.com/Trecer05/Swiftly/internal/model/auth"
	"github.com/Trecer05/Swiftly/internal/model/kafka"

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

func (manager *Manager) ValidatePasswordByID(id int, password string) error {
	var passwordHash string

	if err := manager.Conn.QueryRow("SELECT password_hash FROM users WHERE id = $1", id).Scan(&passwordHash); err != nil {
		if err == sql.ErrNoRows {
			return auth.ErrNoUser
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return auth.ErrInvalidPassword
	}
	return nil
}
 
func (manager *Manager) EditPassword(req *kafka.PasswordEdit) error {
	if err := manager.Conn.QueryRow("UPDATE users SET password_hash = $1 WHERE id = $2", req.NewPassword, req.UserID).Err(); err != nil {
		if err == sql.ErrNoRows {
			return auth.ErrNoUser
		}
		return err
	}

	return nil
}

func (manager *Manager) Register(user *model.User) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := manager.Conn.QueryRow("INSERT INTO users (email, number, password_hash) VALUES ($1, $2, $3) RETURNING id", user.Email, user.Number, string(passwordHash)).Scan(&user.ID); err != nil {
		if pgErrors.IsPgError(err, pgErrors.ErrUsersEmailDuplicate) || pgErrors.IsPgError(err, pgErrors.ErrUsersEmailRusDuplicate){
			return auth.ErrEmailExists
		} else if pgErrors.IsPgError(err, pgErrors.ErrUsersNumberDuplicate) || pgErrors.IsPgError(err, pgErrors.ErrUsersNumberRusDuplicate){
			return auth.ErrNumberExists
		}

		return err
	}
	return nil
}

func (manager *Manager) EditPhone(req *kafka.PhoneEdit) error {
	if err := manager.Conn.QueryRow("UPDATE users SET number = $1 WHERE id = $2", req.NewPhone, req.UserID).Err(); err != nil {
		if err == sql.ErrNoRows {
			return auth.ErrNoUser
		}
		return err
	}

	return nil
}

func (manager *Manager) EditEmail(req *kafka.EmailEdit) error {
	if err := manager.Conn.QueryRow("UPDATE users SET email = $1 WHERE id = $2", req.NewEmail, req.UserID).Err(); err != nil {
		if err == sql.ErrNoRows {
			return auth.ErrNoUser
		}
		return err
	}

	return nil
}
