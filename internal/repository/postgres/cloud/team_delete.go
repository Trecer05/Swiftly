package cloud

import (
	"database/sql"

	cloudErrors "github.com/Trecer05/Swiftly/internal/errors/cloud"
	"github.com/google/uuid"
)

func (manager *Manager) DeleteTeamFileByID(teamID int, requestUserID int, fileUUID uuid.UUID) (string, error) {
	var storagePath string
	var created_by int

	err := manager.Conn.QueryRow(`SELECT storage_path, created_by FROM files WHERE uuid = $1 AND owner_type = 'team' AND owner_id = $2`, fileUUID, teamID).Scan(&storagePath, &created_by)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", cloudErrors.ErrFileNotFound
		}
		return "", err
	}
	if created_by != requestUserID {
		return "", cloudErrors.ErrNoPermissions
	}

	err = manager.Conn.QueryRow(`DELETE FROM files WHERE uuid = $1 AND owner_type = 'team' AND owner_id = $2 AND created_by = $3 RETURNING storage_path`, fileUUID, teamID, requestUserID).Scan(&storagePath)
	if err != nil {
		return "", err
	}
	return storagePath, nil
}

func (manager *Manager) DeleteTeamFolderByID(teamID int, requestUserID int, folderUUID uuid.UUID) (string, error) {
	var storagePath string
	var created_by int

	err := manager.Conn.QueryRow(`SELECT storage_path, created_by FROM folders WHERE uuid = $1 AND owner_type = 'team' AND owner_id = $2`, folderUUID, teamID).Scan(&storagePath, &created_by)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", cloudErrors.ErrFileNotFound
		}
		return "", err
	}
	if created_by != requestUserID {
		return "", cloudErrors.ErrNoPermissions
	}

	err = manager.Conn.QueryRow(`DELETE FROM folders WHERE uuid = $1 AND owner_type = 'team' AND owner_id = $2 AND created_by = $3 RETURNING storage_path`, folderUUID, teamID, requestUserID).Scan(&storagePath)
	if err != nil {
		return "", err
	}
	return storagePath, nil
}
