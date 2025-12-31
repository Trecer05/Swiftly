package cloud

import "github.com/google/uuid"

func (manager *Manager) DeleteTeamFileByID(teamID int, fileUUID uuid.UUID) (string, error) {
	var storagePath string

	err := manager.Conn.QueryRow(`DELETE FROM files WHERE uuid = $1 AND owner_type = 'team' AND owner_id = $2 RETURNING storage_path`, fileUUID, teamID).Scan(&storagePath)
	if err != nil {
		return "", err
	}
	return storagePath, nil
}

func (manager *Manager) DeleteTeamFolderByID(teamID int, folderUUID uuid.UUID) (string, error) {
	var storagePath string

	err := manager.Conn.QueryRow(`DELETE FROM folders WHERE uuid = $1 AND owner_type = 'team' AND owner_id = $2 RETURNING storage_path`, folderUUID, teamID).Scan(&storagePath)
	if err != nil {
		return "", err
	}
	return storagePath, nil
}
