package chat

import (
	"database/sql"

	errors "github.com/Trecer05/Swiftly/internal/errors/chat"
	models "github.com/Trecer05/Swiftly/internal/model/chat"
)

func (manager *Manager) GetUserInfo(userId int) (models.User, error) {
	var user models.User

	err := manager.Conn.QueryRow(`SELECT name, username, description FROM users WHERE id = $1`, userId).Scan(&user.Name, &user.Username, &user.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.ErrNoUser
		}
		return user, err
	}

	return user, nil
}