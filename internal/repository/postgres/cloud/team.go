package cloud

import (
	"database/sql"

	errors "github.com/Trecer05/Swiftly/internal/errors/file"
	models "github.com/Trecer05/Swiftly/internal/model/cloud"

	"github.com/google/uuid"
)

// CREATE TABLE files (
//     uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
//     folder_id UUID REFERENCES folders(uuid) ON DELETE CASCADE,
//     storage_path VARCHAR(500) NOT NULL,
//     original_filename VARCHAR(500) NOT NULL,
//     display_name VARCHAR(500),
//     mime_type VARCHAR(255),
//     size BIGINT NOT NULL DEFAULT 0,
//     visibility file_visibility NOT NULL DEFAULT 'private',

//     created_by INTEGER NOT NULL,
//     owner_id INTEGER,
//     owner_type VARCHAR(10) CHECK (owner_type IN ('user', 'team')),

//     uploaded_at TIMESTAMP DEFAULT NOW(),
//     updated_at TIMESTAMP DEFAULT NOW(),
//     hash VARCHAR(64),
//     version INTEGER DEFAULT 1,

//     CONSTRAINT valid_file_owner CHECK (
//         (owner_type = 'user' AND owner_id IS NOT NULL) OR
//         (owner_type = 'team' AND owner_id IS NOT NULL)
//     )
// );

// type FileShort struct {
// 	UUID             uuid.UUID
// 	StoragePath      string
// 	OriginalFilename string
// 	DisplayName      string
// 	MimeType         string
// 	Size             int64
// }

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
