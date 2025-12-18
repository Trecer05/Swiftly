package chat

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"

	errors "github.com/Trecer05/Swiftly/internal/errors/chat"
	models "github.com/Trecer05/Swiftly/internal/model/chat"
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

func (manager *Manager) EditTeam(team *models.TeamEdit) error {
	var realOwnerID int
	err := manager.Conn.QueryRow(
		"SELECT owner_id FROM projects WHERE id = $1",
		team.ID,
	).Scan(&realOwnerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrProjectNotFound
		}
		return err
	}

	if realOwnerID != team.OwnerID {
		return errors.ErrUserNotAOwner
	}

	setParts := []string{}
	args := []interface{}{}
	idx := 1

	if team.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name=$%d", idx))
		args = append(args, *team.Name)
		idx++
	}

	if team.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description=$%d", idx))
		args = append(args, *team.Description)
		idx++
	}

	if len(setParts) == 0 {
		return errors.ErrNoFieldsToUpdate
	}

	args = append(args, team.ID)

	query := fmt.Sprintf(`
        UPDATE projects
        SET %s
        WHERE id=$%d
    `, strings.Join(setParts, ", "), idx)

	res, err := manager.Conn.Exec(query, args...)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.ErrProjectNotFound
	}

	return nil
}

func (manager *Manager) DeleteTeam(userID, teamID int) error {
	var realOwnerID int
	err := manager.Conn.QueryRow(
		"SELECT owner_id FROM projects WHERE id = $1",
		teamID,
	).Scan(&realOwnerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrProjectNotFound
		}
		return err
	}

	if realOwnerID != userID {
		return errors.ErrUserNotAOwner
	}

	query := `DELETE FROM projects WHERE id = $1`
	res, err := manager.Conn.Exec(query, teamID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.ErrProjectNotFound
	}

	return nil
}

func (manager *Manager) GetTeamInfo(teamID int) (*models.TeamInfo, error) {
	var team models.TeamInfo
	err := manager.Conn.QueryRow(
		"SELECT id, name, description FROM projects WHERE id = $1",
		teamID,
	).Scan(&team.ID, &team.Name, &team.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrProjectNotFound
		}
		return nil, err
	}

	query := `
		SELECT id, name, avatar_url
		FROM users
		WHERE id IN (
			SELECT user_id FROM users_projects WHERE project_id = $1
		)
	`
	rows, err := manager.Conn.Query(query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.UserShort
		err := rows.Scan(&user.ID, &user.Name, &user.AvatarURL)
		if err != nil {
			return nil, err
		}
		team.Users = append(team.Users, user)
	}

	return &team, nil
}

func (manager *Manager) GetTeamApplications(teamID, ownerID int) ([]models.TeamApplication, error) {
	var realOwnerID int
	err := manager.Conn.QueryRow(
		"SELECT owner_id FROM projects WHERE id = $1",
		teamID,
	).Scan(&realOwnerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrProjectNotFound
		}
		return nil, err
	}

	if realOwnerID != ownerID {
		return nil, errors.ErrUserNotAOwner
	}

	var applications []models.TeamApplication
	query := `
		SELECT id, user_id, status, created_at
		FROM team_applications
		WHERE team_id = $1
	`
	rows, err := manager.Conn.Query(query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var application models.TeamApplication
		err := rows.Scan(&application.ID, &application.UserID, &application.Status, &application.CreatedAt)
		if err != nil {
			return nil, err
		}
		applications = append(applications, application)
	}

	return applications, nil
}

func (manager *Manager) UpdateTeamApplication(teamID, ownerID int, status models.TeamApplicationUpdate) error {
	var realOwnerID int
	err := manager.Conn.QueryRow(
		"SELECT owner_id FROM projects WHERE id = $1",
		teamID,
	).Scan(&realOwnerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrProjectNotFound
		}
		return err
	}

	if realOwnerID != ownerID {
		return errors.ErrUserNotAOwner
	}

	query := `
		UPDATE team_applications
		SET status = $1
		WHERE project_id = $2 and id = $3
		RETURNING user_id
	`
	var userID int
	err = manager.Conn.QueryRow(query, status.Status, teamID, status.ID).Scan(&userID)
	if err != nil {
		return err
	}

	switch status.Status {
	case models.TeamApplicationStatusAccepted:
		_, err = manager.Conn.Exec(
			"INSERT INTO users_projects (project_id, user_id) VALUES ($1, $2)",
			teamID, userID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (manager *Manager) CreateJoinCode(code string, req *models.CreateJoinCode) error {
	var realOwnerID int
	err := manager.Conn.QueryRow(
		"SELECT owner_id FROM projects WHERE id = $1",
		req.ProjectID,
	).Scan(&realOwnerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrProjectNotFound
		}
		return err
	}

	if realOwnerID != req.CreatorID {
		return errors.ErrUserNotAOwner
	}

	setParts := []string{}
	args := []string{}
	idx := 4

	if req.ExpiresAt != nil {
		setParts = append(setParts, "expires_at")
		args = append(args, fmt.Sprintf("$%d", idx))
		idx++
	}

	if req.IsSingleUse != nil {
		setParts = append(setParts, "is_single_use")
		args = append(args, fmt.Sprintf("$%d", idx))
	}

	query := fmt.Sprintf("INSERT INTO join_codes (code, project_id, creator_id, %s) VALUES ($1, $2, $3, %s)", strings.Join(setParts, ", "), strings.Join(args, ", "))

	_, err = manager.Conn.Exec(query, code, req.ProjectID, req.CreatorID)
	if err != nil {
		return err
	}

	return nil
}

func (manager *Manager) JoinTeam(userID int, joinCode string) error {
	var projectID int
	var expiresAt time.Time
	var isSingleUse bool
	var used bool
	err := manager.Conn.QueryRow(
		"SELECT project_id, expires_at, is_single_use, used FROM project_invites WHERE code = $1",
		joinCode,
	).Scan(&projectID, &expiresAt, &isSingleUse, &used)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrJoinCodeNotFound
		}
		return err
	}

	if isSingleUse && used {
		return errors.ErrJoinCodeAlreadyUsed
	}

	if time.Now().After(expiresAt) {
		return errors.ErrCodeExpired
	}

	if isSingleUse {
		_, err = manager.Conn.Exec("UPDATE project_invites SET used = true WHERE code = $1", joinCode)
		if err != nil {
			return err
		}
	}

	query := `
		INSERT INTO team_applications (project_id, user_id)
		VALUES ($1, $2)
		RETURNING id
	`
	var id int
	err = manager.Conn.QueryRow(query, projectID, userID).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

func (manager *Manager) IsUserInTeam(teamID int, userID int) (bool, error) {
	query := `
		SELECT 1
		FROM users_projects
		WHERE project_id = $1
		AND user_id = $2
		LIMIT 1;
	`
	err := manager.Conn.QueryRow(query, teamID, userID).Scan()
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil

}
