package cloud

import (
	"database/sql"

	errors "github.com/Trecer05/Swiftly/internal/errors/file"
	models "github.com/Trecer05/Swiftly/internal/model/cloud"

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
