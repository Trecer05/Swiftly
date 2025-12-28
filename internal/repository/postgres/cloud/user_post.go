package cloud

import (
	models "github.com/Trecer05/Swiftly/internal/model/cloud"
)

// dbReq := models.File{
// 	FolderID: req.FolderID,
// 	OriginalFilename: origFilename,
// 	DisplayName: req.DisplayName,
// 	StoragePath: storagePath,
// 	CreatedBy: userID,
// 	OwnerID: userID,
// 	OwnerType: models.OwnerTypeUser,
// 	Hash: hash,
// 	MimeType: mimeType,
// 	Visibility: req.Visibility,
// 	Size: header.Size,
// }


// CREATE TABLE files (
//     uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
//     folder_id UUID REFERENCES folders(uuid) ON DELETE CASCADE,
//     storage_path VARCHAR(500) NOT NULL,
//     original_filename VARCHAR(500) NOT NULL,
//     display_name VARCHAR(500),
//     mime_type VARCHAR(255),
//     size BIGINT NOT NULL DEFAULT 0,
//     visibility file_visibility NOT NULL DEFAULT 'private',
    
//     created_by INTEGER NOT NULL, -- Кто загрузил файл
//     owner_id INTEGER, -- Владелец (пользователь ИЛИ команда)
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

func (manager *Manager) CreateUserFile(req *models.File) (error) {
	err := manager.Conn.QueryRow(`
	INSERT INTO files (folder_id, storage_path, original_filename, display_name, mime_type, size, visibility, created_by, owner_id, owner_type, hash)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING uuid, updated_at, uploaded_at`, req.FolderID, req.StoragePath, req.OriginalFilename, req.DisplayName, req.MimeType, req.Size, req.Visibility, req.CreatedBy, req.OwnerID, req.OwnerType, req.Hash).Scan(&req.UUID, &req.UpdatedAt, &req.UploadedAt)
	if err != nil {
		return err
	}

	return nil
}
