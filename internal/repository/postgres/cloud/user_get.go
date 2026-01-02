package cloud

import (
	"database/sql"
	"fmt"

	cloudErrors "github.com/Trecer05/Swiftly/internal/errors/cloud"
	errors "github.com/Trecer05/Swiftly/internal/errors/file"
	models "github.com/Trecer05/Swiftly/internal/model/cloud"

	"github.com/google/uuid"
)

func (manager *Manager) GetUserFileByID(userId int, fileId uuid.UUID) (*models.File, error) {
	var file models.File
	row := manager.Conn.QueryRow(`SELECT 
				f.uuid,
				f.storage_path,
				f.original_filename,
				f.display_name,
				f.mime_type,
				f.size,
				f.visibility,
				f.created_by,
				f.owner_id,
				f.owner_type,
				f.uploaded_at,
				f.updated_at,
				f.hash,
				f.version
			FROM 
				files f
			WHERE 
			    owner_type = 'user' AND
				owner_id = $1 AND
				f.uuid = $2`, userId, fileId)
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

func (manager *Manager) GetShortUserFileByID(userId int, fileId uuid.UUID) (*models.FileShort, error) {
	var file models.FileShort
	row := manager.Conn.QueryRow(`SELECT 
				f.uuid,
				f.storage_path,
				f.original_filename,
				f.display_name,
				f.mime_type,
				f.size,
			FROM 
				files f
			WHERE 
			    owner_type = 'user' AND
				owner_id = $1 AND
				f.uuid = $2`, userId, fileId)
	err := row.Scan(
		&file.UUID,
		&file.StoragePath,
		&file.OriginalFilename,
		&file.DisplayName,
		&file.MimeType,
		&file.Size,
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

func (manager *Manager) GetUserFilesAndFolders(userId int, sort string) (*models.FilesAndFoldersResponse, error) {
	var response models.FilesAndFoldersResponse

	fileRows, err := manager.Conn.Query(fmt.Sprintf(`SELECT 
				f.uuid,
				f.storage_path,
				f.original_filename,
				f.display_name,
				f.mime_type,
				f.size,
				f.visibility,
				f.created_by,
				f.owner_id,
				f.owner_type,
				f.uploaded_at,
				f.updated_at,
				f.hash,
				f.version
			FROM 
				files f
			WHERE 
			    owner_type = 'user' AND
				owner_id = $1 AND
				folder_id IS NULL
			ORDER BY updated_at %s`, sort), userId)
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

		response.Files = append(response.Files, file)
	}

	folderRows, err := manager.Conn.Query(fmt.Sprintf(`SELECT 
				f.uuid,
				f.name,
				f.owner_id,
				f.owner_type,
				f.parent_folder_id,
				f.created_at,
				f.updated_at,
				f.visibility
			FROM 
				folders f
			WHERE
			    owner_type = 'user' AND
				owner_id = $1 AND
				parent_folder_id IS NULL
			ORDER BY updated_at %s`, sort), userId)
	if err != nil {
		return nil, err
	}
	defer folderRows.Close()

	for folderRows.Next() {
		var folder models.Folder

		err := folderRows.Scan(
			&folder.UUID,
			&folder.Name,
			&folder.OwnerID,
			&folder.OwnerType,
			&folder.ParentFolderID,
			&folder.CreatedAt,
			&folder.UpdatedAt,
			&folder.Visibility,
		)
		if err != nil {
			return nil, err
		}

		response.Folders = append(response.Folders, folder)
	}

	return &response, nil
}

func (manager *Manager) GetUserFilesAndFoldersByFolderID(userId int, folderID, sort string) (*models.FilesAndFoldersResponse, error) {
	var response models.FilesAndFoldersResponse

	fileRows, err := manager.Conn.Query(fmt.Sprintf(`SELECT 
				f.uuid,
				f.storage_path,
				f.original_filename,
				f.display_name,
				f.mime_type,
				f.size,
				f.visibility,
				f.created_by,
				f.owner_id,
				f.owner_type,
				f.uploaded_at,
				f.updated_at,
				f.hash,
				f.version
			FROM 
				files f
			WHERE 
			    owner_type = 'user' AND
				owner_id = $1 AND
				folder_id = $2
			ORDER BY updated_at %s`, sort), userId, folderID)
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

		response.Files = append(response.Files, file)
	}

	folderRows, err := manager.Conn.Query(fmt.Sprintf(`SELECT 
				f.uuid,
				f.name,
				f.owner_id,
				f.owner_type,
				f.parent_folder_id,
				f.created_at,
				f.updated_at,
				f.visibility
			FROM 
				folders f
			WHERE
			    owner_type = 'user' AND
				owner_id = $1 AND
				parent_folder_id = $2
			ORDER BY updated_at %s`, sort), userId, folderID)
	if err != nil {
		return nil, err
	}
	defer folderRows.Close()

	for folderRows.Next() {
		var folder models.Folder

		err := folderRows.Scan(
			&folder.UUID,
			&folder.Name,
			&folder.OwnerID,
			&folder.OwnerType,
			&folder.ParentFolderID,
			&folder.CreatedAt,
			&folder.UpdatedAt,
			&folder.Visibility,
		)
		if err != nil {
			return nil, err
		}

		response.Folders = append(response.Folders, folder)
	}

	return &response, nil
}

// type FileShare struct {
// 	UUID             uuid.UUID  `json:"uuid"`
// 	FolderID         *uuid.UUID `json:"folder_id,omitempty"`
// 	DisplayName      string     `json:"display_name"`
// 	MimeType         string     `json:"mime_type"`
// 	Size             int64      `json:"size"`
// 	CreatedAt        time.Time  `json:"created_at"`
// 	UpdatedAt      time.Time  `json:"updated_at"`
// }

// type FolderShare struct {
// 	UUID           uuid.UUID  `json:"uuid"`
// 	Name           string     `json:"name"`
// 	ParentFolderID *uuid.UUID `json:"parent_folder_id,omitempty"`
// 	CreatedAt      time.Time  `json:"created_at"`
// 	UpdatedAt      time.Time  `json:"updated_at"`
// }

func (manager *Manager) GetSharedFilesAndFoldersByFolderID(folderID, sort string) (*models.SharedFilesAndFoldersResponse, error) {
	var response models.SharedFilesAndFoldersResponse

	fileRows, err := manager.Conn.Query(fmt.Sprintf(`SELECT 
				f.uuid,
				f.display_name,
				f.mime_type,
				f.size,
				f.uploaded_at,
				f.updated_at,
			FROM 
				files f
			JOIN shared_access sa
				ON sa.file_id = f.uuid
			WHERE f.uuid = $1
				AND f.visibility = 'shared'
			ORDER BY updated_at %s`, sort), folderID)
	if err != nil {
		return nil, err
	}
	defer fileRows.Close()

	for fileRows.Next() {
		var file models.FileShare

		err := fileRows.Scan(
			&file.UUID,
			&file.DisplayName,
			&file.MimeType,
			&file.Size,
			&file.CreatedAt,
			&file.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		response.Files = append(response.Files, file)
	}

	folderRows, err := manager.Conn.Query(fmt.Sprintf(`SELECT 
				f.uuid,
				f.name,
				f.parent_folder_id,
				f.created_at,
				f.updated_at,
			FROM 
				folders f
			JOIN shared_access sa
				ON sa.folder_id = f.uuid
			WHERE f.uuid = $1
				AND f.visibility = 'shared'
			ORDER BY updated_at %s`, sort), folderID)
	if err != nil {
		return nil, err
	}
	defer folderRows.Close()

	for folderRows.Next() {
		var folder models.FolderShare

		err := folderRows.Scan(
			&folder.UUID,
			&folder.Name,
			&folder.ParentFolderID,
			&folder.CreatedAt,
			&folder.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		response.Folders = append(response.Folders, folder)
	}

	return &response, nil
}

func (manager *Manager) GetUserFilepathByID(userID int, fileID string) (string, error) {
	row := manager.Conn.QueryRow(`SELECT storage_path FROM files WHERE uuid = $1`, fileID)

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

func (manager *Manager) GetUserFolderpathByID(userID int, folderID string) (string, error) {
	row := manager.Conn.QueryRow(`SELECT storage_path FROM folders WHERE uuid = $1 AND owner_type = 'user' AND owner_id = $2`, folderID, userID)

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

func (manager *Manager) GetOriginalUserFilenameByID(userID int, fileID string) (string, error) {
	var originalFilename string

	if err := manager.Conn.QueryRow(`SELECT original_filename FROM files WHERE uuid = $1 AND created_by = $2 AND owner_type = 'user'`, fileID, userID).Scan(&originalFilename); err != nil {
		return "", err
	}

	return originalFilename, nil
}

func (manager *Manager) GetOriginalUserFilenameAndStoragePathByID(userID int, fileID string) (string, string, error) {
	var originalFilename, storagePath string

	if err := manager.Conn.QueryRow(`SELECT original_filename, storage_path FROM files WHERE uuid = $1 AND created_by = $2 AND owner_type = 'user'`, fileID, userID).Scan(&originalFilename, &storagePath); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return "", "", cloudErrors.ErrFileNotFound
		default:
			return "", "", err
		}
	}

	return originalFilename, storagePath, nil
}

func (manager *Manager) GetSharedFile(fileID string) (string, string, error) {
	var filepath, displayName string

	if err := manager.Conn.QueryRow(`
		SELECT f.storage_path, f.display_name
		FROM files f
		JOIN shared_access sa
			ON sa.file_id = f.uuid
		WHERE f.uuid = $1
			AND f.visibility = 'shared'
	'`, fileID).Scan(&filepath, &displayName); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return "", "", cloudErrors.ErrFileNotFound
		default:
			return "", "", err
		}
	}

	return filepath, displayName, nil
}
