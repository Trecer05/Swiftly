package chat

import (
	"database/sql"

	errors "github.com/Trecer05/Swiftly/internal/errors/chat"
)

func (manager *Manager) AddUserToTeamByUsername(teamID, ownerID int, username string) (int, error) {
    tx, err := manager.Conn.Begin()
    if err != nil {
        return 0, err
    }
    defer tx.Rollback()

    var isAdmin bool
    err = tx.QueryRow(
        "SELECT is_admin FROM users_projects WHERE project_id = $1 AND user_id = $2", 
        teamID, ownerID,
    ).Scan(&isAdmin)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return 0, errors.ErrNoUser
        }
        return 0, err
    }
    
    if !isAdmin {
        return 0, errors.ErrUserNotAOwner
    }

    var userID int
    err = tx.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&userID)
    if err != nil {
        if err == sql.ErrNoRows {
            return 0, errors.ErrNoUser
        }
        return 0, err
    }

    var exists bool
    err = tx.QueryRow(
        "SELECT EXISTS(SELECT 1 FROM users_projects WHERE project_id = $1 AND user_id = $2)",
        teamID, userID,
    ).Scan(&exists)
    if err != nil {
        return 0, err
    }
    
    if exists {
        return 0, errors.ErrUserAlreadyInTeam
    }

    _, err = tx.Exec(
        "INSERT INTO users_projects (project_id, user_id) VALUES ($1, $2)",
        teamID, userID,
    )
    if err != nil {
        return 0, err
    }

    if err := tx.Commit(); err != nil {
        return 0, err
    }

    return userID, nil
}

func (manager *Manager) RemoveUserFromTeamByID(teamID, ownerID, userID int) error {
    tx, err := manager.Conn.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    var isAdmin bool
    err = tx.QueryRow(
        "SELECT is_admin FROM users_projects WHERE project_id = $1 AND user_id = $2",
        teamID, ownerID,
    ).Scan(&isAdmin)
    if err != nil {
        if err == sql.ErrNoRows {
            return errors.ErrNoUser
        }
        return err
    }

    if !isAdmin {
        return errors.ErrUserNotAOwner
    }

    _, err = tx.Exec(
        "DELETE FROM users_projects WHERE project_id = $1 AND user_id = $2",
        teamID, userID,
    )
    if err != nil {
        return err
    }

    if err := tx.Commit(); err != nil {
        return err
    }

    return nil
}
