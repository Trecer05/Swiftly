package cloud

import models "github.com/Trecer05/Swiftly/internal/model/cloud"

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

	if _, err := tx.Exec(`
	    UPDATE files SET visibility = 'shared' AND updated_at = NOW()
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
