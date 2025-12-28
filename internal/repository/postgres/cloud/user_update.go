package cloud

import "time"

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
