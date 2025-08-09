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

func (manager *Manager) GetGroupUsers(groupId int) ([]models.User, error) {
	var users []models.User

	rows, err := manager.Conn.Query(`SELECT 
			u.id, 
			u.name, 
			u.username, 
			u.description
		FROM 
			users u
		JOIN 
			group_users gu ON u.id = gu.user_id
		WHERE 
    gu.group_id = :group_id`, groupId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNoUsers
		} else {
			return nil, err
		}
	}

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Username, &user.Description); err != nil {
			return nil, err
		} else {
			users = append(users, user)
		}
	}

	return users, nil
}

func (manager *Manager) CreateChat(id1, id2 int) (int, error) {
	var chatID int

	err := manager.Conn.QueryRow(`WITH inserted_chat AS (
			INSERT INTO chats (name)
			SELECT 
				CASE 
					WHEN u1.id < u2.id THEN u1.name || ' и ' || u2.name
					ELSE u2.name || ' и ' || u1.name
				END AS chat_name
			FROM users u1, users u2
			WHERE u1.id = $1 AND u2.id = $2
			RETURNING id
		),
		inserted_chat_users AS (
			INSERT INTO chat_users (chat_id, user_id)
			SELECT ic.id, $1 FROM inserted_chat ic
			UNION ALL
			SELECT ic.id, $2 FROM inserted_chat ic
		)
		SELECT id FROM inserted_chat`, id1, id2).Scan(&chatID)
	if err != nil {
		return 0, err
	}

	return chatID, nil
}

func (manager *Manager) GetGroupInfo(groupId int) (models.Group, error) {}
