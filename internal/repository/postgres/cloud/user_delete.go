package cloud

import (
	"database/sql"

	errorCloudTypes "github.com/Trecer05/Swiftly/internal/errors/cloud"
)

func (manager *Manager) DeleteUserFileByID(userID int, fileID string) (string, error) {
	var filepath string

	err := manager.Conn.QueryRow("DELETE FROM files WHERE uuid = $1 AND created_by = $2 AND owner_type = 'user' RETURNING storage_path", fileID, userID).Scan(&filepath)
	switch {
	case err == sql.ErrNoRows:
		return "", errorCloudTypes.ErrFileNotFound
	case err != nil:
		return "", err
	}

	return filepath, nil
}

func (manager *Manager) DeleteUserFolderByID(userID int, folderID string) (string, error) {
	var folderpath string

	if err := manager.Conn.QueryRow("DELETE FROM folders WHERE uuid = $1 AND created_by = $2 AND owner_type = 'user' RETURNING storage_path", folderID, userID).Scan(&folderpath); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return "", errorCloudTypes.ErrFolderNotFound
		default:
			return "", err
		}
	}

	return folderpath, nil
}
