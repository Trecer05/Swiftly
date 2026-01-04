package cloud

import (
	"database/sql"
	"errors"

	cloudErrors "github.com/Trecer05/Swiftly/internal/errors/cloud"
	models "github.com/Trecer05/Swiftly/internal/model/cloud"
	"github.com/lib/pq"
)

func (manager *Manager) UpdateTeamFile(req *models.File) error {
	if err := manager.Conn.QueryRow(`
		UPDATE files SET original_filename = $1, display_name = $2, storage_path = $3, updated_at = NOW(), hash = $4, version = version + 1, size = $5, mime_type = $6
		WHERE uuid = $7 AND owner_id = $8 AND owner_type = 'team'
		RETURNING created_by, uploaded_at, updated_at, version
	`, req.OriginalFilename, req.DisplayName, req.StoragePath, req.Hash, req.Size, req.MimeType, req.UUID, req.OwnerID).Scan(&req.CreatedBy, &req.UploadedAt, &req.UpdatedAt, &req.Version); err != nil {
		return err
	}

	return nil
}

func (manager *Manager) ShareTeamFileByID(fileID string, teamID int, userID int) error {
	tx, err := manager.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // nolint: errcheck

	var created_by int
	if err := tx.QueryRow(`SELECT created_by FROM files WHERE uuid = $1 AND owner_type = 'team'`, fileID).Scan(&created_by); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return cloudErrors.ErrFileNotFound
		}
		return err
	}

	if created_by != userID {
		return cloudErrors.ErrNoPermissions
	}

	if _, err := tx.Exec(`
	    UPDATE files SET visibility = 'shared', updated_at = NOW()
		WHERE uuid = $1 AND created_by = $2 AND owner_type = 'team'
	`, fileID, userID); err != nil {
		return err
	}

	if _, err := tx.Exec(`
		INSERT INTO shared_access (file_id, shared_with_id, shared_with_type, shared_by)
			VALUES ($1, $2, 'team', $3)
	`, fileID, teamID, userID); err != nil {
		return err
	}

	return tx.Commit()
}

func (manager *Manager) ShareTeamFolderByID(folderID string, teamID int, userID int) error {
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
		  AND owner_type = 'team'
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
			$3,
			'team',
			$2
		FROM folders
		WHERE uuid = ANY($1)
	`, pq.Array(folderIDs), userID, teamID)
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
			$3,
			'team',
			$2
		FROM files f
		WHERE f.folder_id = ANY($1)
	`, pq.Array(folderIDs), userID, teamID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (manager *Manager) MoveTeamFileByID(fileID, folderID string, teamID int, storagePath string) error {
	if _, err := manager.Conn.Exec(`
		UPDATE files SET folder_id = $1, storage_path = $2, updated_at = NOW()
		WHERE uuid = $3 AND owner_id = $4 AND owner_type = 'team'
	`, folderID, storagePath, fileID, teamID); err != nil {
		return err
	}

	return nil
}

func (manager *Manager) MoveTeamFolderByID(folderID, newParentID string, teamID int, storagePath string) error {
	if _, err := manager.Conn.Exec(`
		UPDATE folders SET parent_folder_id = $1, storage_path = $2, updated_at = NOW()
		WHERE uuid = $3 AND owner_id = $4 AND owner_type = 'team'
	`, newParentID, storagePath, folderID, teamID); err != nil {
		return err
	}

	return nil
}
