package chat

import (
	"database/sql"
	"github.com/lib/pq"

	models "github.com/Trecer05/Swiftly/internal/model/chat"
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

func (manager *Manager) UpdateUserRole(teamID, ownerID, userID int, newRole string) error {
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
        "UPDATE users_projects SET role = $1 WHERE project_id = $2 AND user_id = $3",
        newRole, teamID, userID,
    )
    if err != nil {
        return err
    }

    if err := tx.Commit(); err != nil {
        return err
    }

    return nil
}

func (manager *Manager) CreateTeam(team *models.TeamCreate) (int, error) {
	tx, err := manager.Conn.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	
	var projectID int
	if err = tx.QueryRow(
		"INSERT INTO projects (name, description) VALUES ($1, $2) RETURNING id",
		team.Name, team.Description,
	).Scan(&projectID); err != nil {
		return 0, err
	}

	if len(team.Users) != 0 {
		users, err := manager.GetUsersFromUsernames(tx, team.Users)
		if err != nil {
			return 0, err
		}
		
		for _, user := range users {
			if err = tx.QueryRow(
				"INSERT INTO users_projects (project_id, user_id, role) VALUES ($1, $2, $3) RETURNING id",
				projectID, user.ID, user.Role,
			).Scan(&user.ID); err != nil {
				return 0, err
			}
		}
	}
	
	_, err = tx.Exec("INSERT INTO users_projects (project_id, user_id, role, is_admin) VALUES ($1, $2, $3, $4)", projectID, team.OwnerID, "creator", true)

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return projectID, nil
}

func (manager *Manager) GetUsersFromUsernames(tx *sql.Tx, oldUsers []models.TeamUser) ([]models.TeamUser, error) {
	var usernames []string
	for _, user := range oldUsers {
		usernames = append(usernames, user.Username)
	}
	
	query := `SELECT id, username FROM users WHERE username = ANY($1)`
	rows, err := tx.Query(query, pq.Array(usernames))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []models.TeamUser
	for rows.Next() {
		var user models.TeamUser
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
