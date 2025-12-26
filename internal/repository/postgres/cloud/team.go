package cloud

import (
	"database/sql"

	errors "github.com/Trecer05/Swiftly/internal/errors/file"
	models "github.com/Trecer05/Swiftly/internal/model/cloud"
	"github.com/Trecer05/Swiftly/internal/service/cloud"

	"github.com/google/uuid"
)

func (manager *Manager) GetTeamFileByID(teamId int, fileId uuid.UUID) (*models.File, error) {
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
			    owner_type = 'team' AND
				owner_id = $1 AND
				f.uuid = $2`, teamId, fileId)
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
			return nil, err
		}
		// Если получение папки организации уже подразумевает, что пользователь в этой организации,
		// то передаем isInTeam = true в функцию проверки доступа к файлу
		if cloud.HasAccessToTeamFile(&file, requestUserID, true) != nil {
			files = append(files, file)
		}
	}
	return files, nil
}

func (manager *Manager) GetTeamFolderByTeamID(teamId int, folderId uuid.UUID) (*models.Folder, error) {
	var folder models.Folder
	var parent sql.NullString

	err := manager.Conn.QueryRow(`SELECT 
				uuid,
				name,
				owner_id,
				owner_type,
				parent_folder_id,
				created_at
			FROM 
				folders
			WHERE
			    owner_type = 'team' AND
				owner_id = $1 AND
				uuid = $2`, teamId, folderId).Scan(
		&folder.UUID,
		&folder.Name,
		&folder.OwnerID,
		&folder.OwnerType,
		&parent,
		&folder.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrFolderNotFound
		} else {
			return nil, err
		}
	}
	if parent.Valid {
		parentUUID, err := uuid.Parse(parent.String)
		if err != nil {
			return nil, err
		}
		folder.ParentFolderID = &parentUUID
	}
	folder.Files, err = manager.GetTeamFilesFromFolder(teamId, folder.OwnerID, folderId)
	if err != nil {
		return nil, err
	}
	return &folder, nil
}
