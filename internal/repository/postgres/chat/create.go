package chat

import (
	"fmt"

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

	userIDs := make([]int, 0, len(group.Users))
	for _, user := range group.Users {
        userIDs = append(userIDs, user.ID)
    }

	_, err = tx.Exec(`
        INSERT INTO group_users (group_id, user_id)
        SELECT $1, unnest($2::bigint[])
        ON CONFLICT (group_id, user_id) DO NOTHING`,
        groupId, pq.Array(userIDs),
    )
    if err != nil {
        return 0, fmt.Errorf("failed to add users to group: %v", err)
    }

    if err := tx.Commit(); err != nil {
        return 0, errors.ErrFailedCommitTx
    }

    return groupId, nil
}