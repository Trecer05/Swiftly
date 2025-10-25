package chat

import (
	"database/sql"
	"fmt"

	chatErrors "github.com/Trecer05/Swiftly/internal/errors/chat"
	errors "github.com/Trecer05/Swiftly/internal/errors/postgres"
	models "github.com/Trecer05/Swiftly/internal/model/chat"

	"github.com/lib/pq"
)

func (manager *Manager) CreateGroup(group models.GroupCreate) (int, error) {
	var groupId int

	tx, err := manager.Conn.Begin()
    if err != nil {
        return 0, errors.ErrFailedBeginTx
    }
    defer tx.Rollback()

	err = tx.QueryRow(`
		INSERT INTO groups (name, description, owner_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`).Scan(&groupId)
	if err != nil {
		return 0, fmt.Errorf("failed to create group: %v", err)
	}

	userIDs := make(map[int]struct{})
	for _, user := range group.Users {
		userIDs[user.ID] = struct{}{}
	}
	userIDs[group.OwnerID] = struct{}{}

	idSlice := make([]int, 0, len(userIDs))
	for id := range userIDs {
		idSlice = append(idSlice, id)
	}

	_, err = tx.Exec(`
        INSERT INTO group_users (group_id, user_id)
        SELECT $1, unnest($2::bigint[])
        ON CONFLICT (group_id, user_id) DO NOTHING`,
        groupId, pq.Array(idSlice),
    )
    if err != nil {
        return 0, fmt.Errorf("failed to add users to group: %v", err)
    }

    if err := tx.Commit(); err != nil {
        return 0, errors.ErrFailedCommitTx
    }

    return groupId, nil
}

func (manager *Manager) DeleteGroup(groupId int) error {
	_, err := manager.Conn.Exec("DELETE FROM groups WHERE id = $1", groupId)
	return err
}

func (manager *Manager) DeleteChat(chatId int) error {
    _, err := manager.Conn.Exec("DELETE FROM chats WHERE id = $1", chatId)
	return err
}

func (manager *Manager) DeleteUsersFromGroup(users models.Users, groupId int) error {
	tx, err := manager.Conn.Begin()
    if err != nil {
        return errors.ErrFailedBeginTx
    }
    defer tx.Rollback()

	idSlice := make([]int, 0, len(users.Users))
	for _, user := range users.Users {
		idSlice = append(idSlice, user.ID)
	}

	var count int
    err = tx.QueryRow(`
        SELECT COUNT(*) 
        FROM unnest($1::bigint[]) AS user_id 
        WHERE NOT EXISTS (
            SELECT 1 FROM group_users 
            WHERE group_id = $2 AND group_users.user_id = user_id
        )`,
        pq.Array(idSlice), groupId,
    ).Scan(&count)
    if err != nil {
        return err
    }
    if count > 0 {
        return chatErrors.ErrUserNotInGroup
    }

	result, err := tx.Exec(`
        DELETE FROM group_users 
        WHERE group_id = $1 AND user_id = ANY($2)`,
        groupId, pq.Array(idSlice),
    )
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if rowsAffected != int64(len(idSlice)) {
        return chatErrors.ErrUserNotInGroup
    }

	if err := tx.Commit(); err != nil {
        return errors.ErrFailedCommitTx
    }

	return nil
}

func (manager *Manager) AddUsersToGroup(users models.Users, groupId int) error {
	tx, err := manager.Conn.Begin()
    if err != nil {
        return errors.ErrFailedBeginTx
    }
    defer tx.Rollback()

	idSlice := make([]int, 0, len(users.Users))
	for _, user := range users.Users {
		idSlice = append(idSlice, user.ID)
	}

	var count int
    err = tx.QueryRow(`
        SELECT COUNT(*) 
        FROM unnest($1::bigint[]) AS user_id 
        WHERE EXISTS (
            SELECT 1 FROM group_users 
            WHERE group_id = $2 AND group_users.user_id = user_id
        )`,
        pq.Array(idSlice), groupId,
    ).Scan(&count)
    if err != nil {
        return err
    }
    if count > 0 {
        return chatErrors.ErrUserAlreadyInGroup
    }

	_, err = tx.Exec(`
		INSERT INTO group_users (group_id, user_id)
		SELECT $1, unnest($2::bigint[])
		ON CONFLICT (group_id, user_id) DO NOTHING`,
		groupId, pq.Array(idSlice),
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return chatErrors.ErrUserAlreadyInGroup
		} else {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
        return errors.ErrFailedCommitTx
    }

	return nil
}

func (manager *Manager) ExitGroup(userId, groupId int) error {
    _, err := manager.Conn.Exec(`
        DELETE FROM group_users
        WHERE user_id = $1 and group_id = $2`,
        userId, groupId)
    if err != nil {
        return err
    }

    return nil
}
