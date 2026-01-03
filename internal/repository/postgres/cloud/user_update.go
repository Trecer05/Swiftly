package cloud

import (
	"database/sql"
	"time"

	cloudErrors "github.com/Trecer05/Swiftly/internal/errors/cloud"
	models "github.com/Trecer05/Swiftly/internal/model/cloud"
	"github.com/lib/pq"
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

func (manager *Manager) ShareUserFolderByID(folderID string, userID int) error {
	tx, err := manager.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // nolint: errcheck

	// проверяем, что пользователь владелец папки
	var ownerID int
	err = tx.QueryRow(`
		SELECT created_by
		FROM folders
		WHERE uuid = $1
		  AND owner_type = 'user'
	`, folderID).Scan(&ownerID)

	if err != nil {
		if err == sql.ErrNoRows {
			return cloudErrors.ErrFolderNotFound
		}
		return err
	}

	if ownerID != userID {
		return cloudErrors.ErrNoPermissions
	}

	rows, err := tx.Query(`
		WITH RECURSIVE folder_tree AS (
			SELECT uuid
			FROM folders
			WHERE uuid = $1

			UNION ALL

			SELECT f.uuid
			FROM folders f
			INNER JOIN folder_tree ft ON f.parent_folder_id = ft.uuid
		)
		SELECT uuid FROM folder_tree
	`, folderID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var folderIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return err
		}
		folderIDs = append(folderIDs, id)
	}

	if len(folderIDs) == 0 {
		return cloudErrors.ErrFolderNotFound
	}

	_, err = tx.Exec(`
		UPDATE folders
		SET visibility = 'shared', updated_at = NOW()
		WHERE uuid = ANY($1)
	`, pq.Array(folderIDs))
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE files
		SET visibility = 'shared', updated_at = NOW()
		WHERE folder_id = ANY($1)
	`, pq.Array(folderIDs))
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO shared_access (
			folder_id,
			shared_with_id,
			shared_with_type,
			shared_by
		)
		SELECT
			uuid,
			created_by,
			'user',
			$2
		FROM folders
		WHERE uuid = ANY($1)
	`, pq.Array(folderIDs), userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO shared_access (
			file_id,
			shared_with_id,
			shared_with_type,
			shared_by
		)
		SELECT
			f.uuid,
			f.created_by,
			'user',
			$2
		FROM files f
		WHERE f.folder_id = ANY($1)
	`, pq.Array(folderIDs), userID)
	if err != nil {
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

func (manager *Manager) MoveUserFileByID(fileID, folderID string, userID int, storagePath string) error {
	if _, err := manager.Conn.Exec(`
		UPDATE files SET folder_id = $1, storage_path = $2, updated_at = NOW()
		WHERE uuid = $3 AND created_by = $4 AND owner_type = 'user'
	`, folderID, storagePath, fileID, userID); err != nil {
		return err
	}

	return nil
}

func (manager *Manager) MoveUserFolderByID(folderID, newParentID string, userID int, storagePath string) error {
	if _, err := manager.Conn.Exec(`
		UPDATE folders SET parent_folder_id = $1, storage_path = $2, updated_at = NOW()
		WHERE uuid = $3 AND created_by = $4 AND owner_type = 'user'
	`, newParentID, storagePath, folderID, userID); err != nil {
		return err
	}

	return nil
}
