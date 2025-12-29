package cloud

import (
	"time"

	models "github.com/Trecer05/Swiftly/internal/model/cloud"
)

func (manager *Manager) UpdateFileFilenameByID(userID int, fileID, newOrigName, newFilename, newFilepath string) (time.Time, error) {
	var updatedAt time.Time
	
	if err := manager.Conn.QueryRow(`
		UPDATE files SET original_filename = $1, display_name = $2, storage_path = $3, updated_at = NOW()
		WHERE uuid = $4 AND created_by = $5
		RETURNING updated_at
	`, newOrigName, newFilename, newFilepath, fileID, userID).Scan(&updatedAt); err != nil {
		return time.Time{}, err
	}

	return updatedAt, nil
}

func (manager *Manager) UpdateFolderFoldernameByID(userID int, folderID, newFoldername, newFolderpath string) (time.Time, error) {
	var updatedAt time.Time

	if err := manager.Conn.QueryRow(`
		UPDATE folders SET name = $1, storage_path = $2, updated_at = NOW()
		WHERE uuid = $3 AND created_by = $4
		RETURNING updated_at
	`, newFoldername, newFolderpath, folderID, userID).Scan(&updatedAt); err != nil {
		return time.Time{}, err
	}

	return updatedAt, nil
}

func (manager *Manager) ShareUserFileByID(fileID string, userID int) error {
	tx, err := manager.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // nolint: errcheck

	if _, err := tx.Exec(`
	    UPDATE files SET visibility = 'shared' AND updated_at = NOW()
		WHERE uuid = $1 AND created_by = $2 AND owner_type = 'user'
	`, fileID, userID); err != nil {
		return err
	}

	if _, err := tx.Exec(`
		INSERT INTO shared_access (file_id, shared_with_id, shared_with_type, shared_by)
			VALUES ($1, $2, 'user', $2)
	`, fileID, userID); err != nil {
		return err
	}

	return tx.Commit()
}

func (manager *Manager) UpdateUserFile(req *models.File) error {
	if err := manager.Conn.QueryRow(`
		UPDATE files SET original_filename = $1, display_name = $2, storage_path = $3, updated_at = NOW(), hash = $4, version = version + 1, size = $5, mime_type = $6
		WHERE uuid = $7 AND created_by = $8
		RETURNING updated_at
	`, req.OriginalFilename, req.DisplayName, req.StoragePath, req.Hash, req.Size, req.MimeType, req.UUID, req.CreatedBy).Scan(&req.Hash); err != nil {
		return err
	}

	return nil
}
