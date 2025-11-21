package chat

import (
	"database/sql"

	errors "github.com/Trecer05/Swiftly/internal/errors/chat"
	pgErrors "github.com/Trecer05/Swiftly/internal/errors/postgres"
	models "github.com/Trecer05/Swiftly/internal/model/chat"
)

func (manager *Manager) GetUserInfo(userId int) (models.User, error) {
	var user models.User

	err := manager.Conn.QueryRow(`SELECT name, username, description, avatar_url FROM users WHERE id = $1`, userId).Scan(&user.Name, &user.Username, &user.Description, &user.AvatarUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.ErrNoUser
		}
		return user, err
	}

	return user, nil
}

func (manager *Manager) GetStartUserInfo(userId int) (models.StartUserInfo, error) {
	var user models.StartUserInfo
	user.ID = userId

	err := manager.Conn.QueryRow(`
		SELECT name, username, description, avatar_url 
		FROM users WHERE id = $1`, userId).Scan(
		&user.Name, &user.Username, &user.Description, &user.AvatarUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.ErrNoUser
		}
		return user, err
	}

	rows, err := manager.Conn.Query(`
		SELECT p.id, p.name, p.description, up.is_admin
		FROM projects p
		INNER JOIN users_projects up ON p.id = up.project_id
		WHERE up.user_id = $1
		ORDER BY p.name`, userId)
	if err != nil {
		return user, err
	}
	defer rows.Close()

	var projects []models.UserProject

	for rows.Next() {
		var project models.UserProject
		err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.IsAdmin)
		if err != nil {
			return user, err
		}

		userRows, err := manager.Conn.Query(`
			SELECT u.id, u.name, up.role, u.avatar_url
			FROM users u
			INNER JOIN users_projects up ON u.id = up.user_id
			WHERE up.project_id = $1
			ORDER BY u.name`, project.ID)
		if err != nil {
			return user, err
		}

		var users []models.UserShort
		for userRows.Next() {
			var userShort models.UserShort
			err := userRows.Scan(&userShort.ID, &userShort.Name, &userShort.Role, &userShort.AvatarURL)
			if err != nil {
				userRows.Close()
				return user, err
			}
			users = append(users, userShort)
		}
		userRows.Close()

		project.Users = users
		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		return user, err
	}

	user.Projects = projects
	return user, nil
}

func (manager *Manager) GetGroupUsers(groupId int) ([]models.User, error) {
	var users []models.User

	rows, err := manager.Conn.Query(`SELECT 
			u.id, 
			u.name, 
			u.username, 
			u.description,
			u.avatar_url
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
		if err := rows.Scan(&user.ID, &user.Name, &user.Username, &user.AvatarUrl, &user.Description); err != nil {
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

func (manager *Manager) GetGroupInfo(groupId int) (models.Group, error) {
	var group models.Group

	if err := manager.Conn.QueryRow(`SELECT name, description FROM groups WHERE id = $1`, groupId).Scan(&group.Name, &group.Description); err != nil {
		if err == sql.ErrNoRows {
			return group, errors.ErrNoGroupFound
		} else {
			return group, err
		}
	}

	return group, nil
}

func (manager *Manager) CreateUser(user models.RegisterUser, id int) (error) {
	err := manager.Conn.QueryRow(`INSERT INTO users (id, name, username, description)`, id, user.Name, user.Username, user.Description).Err()
	switch {
	case err == pgErrors.ErrUserIdDuplicate || err == pgErrors.ErrUserIdRusDuplicate:
		return pgErrors.ErrUserIdDuplicate
	case err != nil:
		return err
	}
	return nil
}

func (manager *Manager) EditProfileDescription(description string, userId int) error {
	if _, err := manager.Conn.Exec(`
		UPDATE users
		SET description = $1
		WHERE id = $2`, description, userId); err != nil {
		return err
	}

	return nil
}

func (manager *Manager) EditProfileName(name string, userId int) error {
	if _, err := manager.Conn.Exec(`
		UPDATE users
		SET name = $1
		WHERE id = $2`, name, userId); err != nil {
		return err
	}

	return nil
}
