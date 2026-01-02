package cloud

import (
	"database/sql"
	"fmt"

	"github.com/Trecer05/Swiftly/internal/config/logger"
	cloudErrors "github.com/Trecer05/Swiftly/internal/errors/cloud"
	errors "github.com/Trecer05/Swiftly/internal/errors/file"
	models "github.com/Trecer05/Swiftly/internal/model/cloud"

	"github.com/google/uuid"
)

func hasAccessToTeamFile(file *models.File, requestUserID int, isInTeam bool) error {
	switch file.OwnerType {
	case models.OwnerTypeUser:
		switch file.Visibility {
		case models.VisibilityPrivate:
			if file.OwnerID != requestUserID {
				return errors.ErrPermissionDenied
			} else {
				return nil
			}
		case models.VisibilityShared:
			if !isInTeam {
				return errors.ErrPermissionDenied
			} else {
				return nil
			}
		}
	case models.OwnerTypeTeam:
		switch file.Visibility {
		case models.VisibilityPublic:
			if !isInTeam {
				return errors.ErrPermissionDenied
			} else {
				return nil
			}
		case models.VisibilityPrivate:
			if file.OwnerID != requestUserID {
				return errors.ErrPermissionDenied
			} else {
				return nil
			}
		case models.VisibilityShared:
			return nil
		}
	}
	return errors.ErrPermissionDenied
}

func hasAccessToTeamFolder(folder *models.Folder, requestUserID int, isInTeam bool) error {
	switch folder.OwnerType {
	case models.OwnerTypeUser:
		switch folder.Visibility {
		case models.VisibilityPrivate:
			if folder.OwnerID != requestUserID {
				return errors.ErrPermissionDenied
			} else {
				return nil
			}
		case models.VisibilityShared:
			if !isInTeam {
				return errors.ErrPermissionDenied
			} else {
				return nil
			}
		}
	case models.OwnerTypeTeam:
		switch folder.Visibility {
		case models.VisibilityPublic:
			if !isInTeam {
				return errors.ErrPermissionDenied
			} else {
				return nil
			}
		case models.VisibilityPrivate:
			if folder.OwnerID != requestUserID {
				return errors.ErrPermissionDenied
			} else {
				return nil
			}
		case models.VisibilityShared:
			return nil
		}
	}
	return errors.ErrPermissionDenied
}

func (manager *Manager) GetTeamFileByID(teamId int, fileId uuid.UUID) (*models.File, error) {
	var file models.File
	row := manager.Conn.QueryRow(`SELECT 
				uuid,
				storage_path,
				original_filename,
				display_name,
				mime_type,
				size,
				visibility,
				created_by,
				owner_id,
				owner_type,
				uploaded_at,
				updated_at,
				hash,
				version
			FROM 
				files
			WHERE 
			    owner_type = 'team' AND
				owner_id = $1 AND
				uuid = $2`, teamId, fileId)
	err := row.Scan(
		&file.UUID,
		&file.StoragePath,
		&file.OriginalFilename,
		&file.DisplayName,
		&file.MimeType,
		&file.Size,
		&file.Visibility,
		&file.CreatedBy,
		&file.OwnerID,
		&file.OwnerType,
		&file.UploadedAt,
		&file.UpdatedAt,
		&file.Hash,
		&file.Version,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrFileNotFound
		} else {
			return nil, err
		}
	}
	return &file, nil
}

// CREATE TABLE folders (
//     uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
//     name VARCHAR(255) NOT NULL,
//     owner_id INTEGER,  -- Владелец (пользователь ИЛИ команда)
//     owner_type VARCHAR(10) CHECK (owner_type IN ('user', 'team')),
//     parent_folder_id UUID REFERENCES folders(uuid) ON DELETE CASCADE,
//     created_at TIMESTAMP DEFAULT NOW(),

//     CONSTRAINT valid_owner CHECK (
//         (owner_type = 'user' AND owner_id IS NOT NULL) OR
//         (owner_type = 'team' AND owner_id IS NOT NULL)
//     )
// );

func (manager *Manager) GetTeamFilesFromFolder(teamId int, requestUserID int, folderId uuid.UUID) ([]models.File, error) {
	rows, err := manager.Conn.Query(`SELECT
				uuid,
				storage_path,
				original_filename,
				display_name,
				mime_type,
				size,
				visibility,
				created_by,
				owner_id,
				owner_type,
				uploaded_at,
				updated_at,
				hash,
				version
			FROM
				files
			WHERE
			    owner_type = 'team' AND
				owner_id = $1 AND
				folder_id = $2`, teamId, folderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var files []models.File
	for rows.Next() {
		var file models.File
		err := rows.Scan(
			&file.UUID,
			&file.StoragePath,
			&file.OriginalFilename,
			&file.DisplayName,
			&file.MimeType,
			&file.Size,
			&file.Visibility,
			&file.CreatedBy,
			&file.OwnerID,
			&file.OwnerType,
			&file.UploadedAt,
			&file.UpdatedAt,
			&file.Hash,
			&file.Version,
		)
		if err != nil {
			logger.Logger.Error("Error parsing data from query", err.Error())
			return nil, err
		}
		// Если получение папки организации уже подразумевает, что пользователь в этой организации,
		// то передаем isInTeam = true в функцию проверки доступа к файлу
		if hasAccessToTeamFile(&file, requestUserID, true) == nil {
			files = append(files, file)
		}
	}
	return files, nil
}

func (manager *Manager) GetTeamFolderByTeamID(teamId int, requestUserID int, folderId uuid.UUID) (*models.FilesAndFoldersResponse, error) {
	// var folder models.Folder
	var response models.FilesAndFoldersResponse

	response.Files = make([]models.File, 0)
	response.Folders = make([]models.Folder, 0)

	// Get main folder data
	folderRows, err := manager.Conn.Query(`SELECT 
				uuid,
				name,
				created_by,
				owner_id,
				owner_type,
				visibility,
				parent_folder_id,
				created_at
			FROM 
				folders
			WHERE
			    owner_type = 'team' AND
				owner_id = $1 AND
				parent_folder_id = $2`, teamId, folderId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrFolderNotFound
		} else {
			return nil, err
		}
	}

	defer folderRows.Close()
	for folderRows.Next() {
		var folder models.Folder
		var parent sql.NullString
		err := folderRows.Scan(
			&folder.UUID,
			&folder.Name,
			&folder.CreatedBy,
			&folder.OwnerID,
			&folder.OwnerType,
			&folder.Visibility,
			&parent,
			&folder.CreatedAt,
		)
		if err != nil {
			continue
		}
		if parent.Valid {
			parentUUID, err := uuid.Parse(parent.String)
			if err != nil {
				return nil, err
			}
			folder.ParentFolderID = &parentUUID
		}
		if hasAccessToTeamFolder(&folder, requestUserID, true) == nil {
			response.Folders = append(response.Folders, folder)
		}
	}

	response.Files, err = manager.GetTeamFilesFromFolder(teamId, requestUserID, folderId)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (manager *Manager) GetTeamFilesAndFolders(teamID int, requestUserID int, sort string) (*models.FilesAndFoldersResponse, error) {
	var response models.FilesAndFoldersResponse

	response.Folders = make([]models.Folder, 0)
	response.Files = make([]models.File, 0)

	// Получаем папки
	rows, err := manager.Conn.Query(fmt.Sprintf(`
		SELECT 
			uuid,
			name,
			created_by,
			owner_id,
			owner_type,
			visibility,
			parent_folder_id,
			created_at
			FROM 
				folders
				WHERE
				    owner_type = 'team' AND
					owner_id = $1 AND
					parent_folder_id IS NULL
			ORDER BY updated_at %s`, sort), teamID)
	if err != nil {
		logger.Logger.Error("Failed to get folder data from database ", err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var folder models.Folder
		err := rows.Scan(
			&folder.UUID,
			&folder.Name,
			&folder.CreatedBy,
			&folder.OwnerID,
			&folder.OwnerType,
			&folder.Visibility,
			&folder.ParentFolderID,
			&folder.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if hasAccessToTeamFolder(&folder, requestUserID, true) == nil {
			response.Folders = append(response.Folders, folder)
		}
	}

	// Получаем файлы
	fileRows, err := manager.Conn.Query(fmt.Sprintf(`SELECT
				uuid,
				storage_path,
				original_filename,
				display_name,
				mime_type,
				size,
				visibility,
				created_by,
				owner_id,
				owner_type,
				uploaded_at,
				updated_at,
				hash,
				version
			FROM
				files
			WHERE
			    owner_type = 'team' AND
				owner_id = $1 AND
				folder_id is null
			ORDER BY updated_at %s`, sort), teamID)

	if err != nil {
		return nil, err
	}
	defer fileRows.Close()

	for fileRows.Next() {
		var file models.File
		err := fileRows.Scan(
			&file.UUID,
			&file.StoragePath,
			&file.OriginalFilename,
			&file.DisplayName,
			&file.MimeType,
			&file.Size,
			&file.Visibility,
			&file.CreatedBy,
			&file.OwnerID,
			&file.OwnerType,
			&file.UploadedAt,
			&file.UpdatedAt,
			&file.Hash,
			&file.Version,
		)
		if err != nil {
			return nil, err
		}
		if hasAccessToTeamFile(&file, requestUserID, true) == nil {
			response.Files = append(response.Files, file)
		}
	}
	return &response, nil
}

func (manager *Manager) GetTeamFolderpathByID(teamID int, folderID string) (string, error) {
	row := manager.Conn.QueryRow(`SELECT storage_path FROM folders WHERE uuid = $1 AND owner_type = 'team' AND owner_id = $2`, folderID, teamID)

	var storagePath string
	if err := row.Scan(&storagePath); err != nil {
		if err == sql.ErrNoRows {
			return "", cloudErrors.ErrFolderNotFound
		} else {
			return "", err
		}
	}

	return storagePath, nil
}

func (manager *Manager) GetTeamFilepathByID(teamID int, fileID string) (string, error) {
	row := manager.Conn.QueryRow(`SELECT storage_path FROM files WHERE uuid = $1 AND owner_type = 'team' AND owner_id = $2`, fileID, teamID)

	var storagePath string
	if err := row.Scan(&storagePath); err != nil {
		if err == sql.ErrNoRows {
			return "", cloudErrors.ErrFileNotFound
		} else {
			return "", err
		}
	}

	return storagePath, nil
}

func (manager *Manager) GetOriginalTeamFilenameByID(teamID int, fileID string) (string, error) {
	var originalFilename string

	if err := manager.Conn.QueryRow(`SELECT original_filename FROM files WHERE uuid = $1 AND owner_id = $2 AND owner_type = 'team'`, fileID, teamID).Scan(&originalFilename); err != nil {
		return "", err
	}

	return originalFilename, nil
}

func (manager *Manager) GetOriginalTeamFilenameAndStoragePathByID(teamID int, fileID string) (string, string, error) {
	var originalFilename, storagePath string

	if err := manager.Conn.QueryRow(`SELECT original_filename, storage_path FROM files WHERE uuid = $1 AND owner_id = $2 AND owner_type = 'team'`, fileID, teamID).Scan(&originalFilename, &storagePath); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return "", "", cloudErrors.ErrFileNotFound
		default:
			return "", "", err
		}
	}

	return originalFilename, storagePath, nil
}
