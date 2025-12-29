package cloud

import (
	models "github.com/Trecer05/Swiftly/internal/model/cloud"
)

// type Folder struct {
//     UUID           uuid.UUID      `json:"uuid"`
//     Name           string         `json:"name"`
//     CreatedBy      int            `json:"created_by"`
//     OwnerID        int            `json:"owner_id"`
//     OwnerType      OwnerType      `json:"owner_type"`
//     StoragePath    string         `json:"storage_path"`
//     Files          []File         `json:"files,omitempty"`
//     Visibility     VisibilityType `json:"visibility"`
//     ParentFolderID *uuid.UUID     `json:"parent_folder_id,omitempty"`
//     CreatedAt      time.Time      `json:"created_at"`
//     UpdatedAt      time.Time      `json:"updated_at"`
// }

// CREATE TABLE folders (
//     uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
//     name VARCHAR(255) NOT NULL,
//     created_by INTEGER NOT NULL,
//     owner_id INTEGER,  -- Владелец (пользователь ИЛИ команда)
//     owner_type VARCHAR(10) CHECK (owner_type IN ('user', 'team')),
//     storage_path VARCHAR(500) NOT NULL,
//     visibility file_visibility NOT NULL DEFAULT 'private',
//     parent_folder_id UUID REFERENCES folders(uuid) ON DELETE CASCADE,
//     created_at TIMESTAMP DEFAULT NOW(),
//     updated_at TIMESTAMP DEFAULT NOW(),
    
//     CONSTRAINT valid_owner CHECK (
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

func (manager *Manager) CreateUserFolder(req *models.Folder) (error) {
	err := manager.Conn.QueryRow(`
	INSERT INTO folders (name, created_by, owner_id, owner_type, storage_path, visibility, parent_folder_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING uuid, updated_at`, req.Name, req.CreatedBy, req.OwnerID, req.OwnerType, req.StoragePath, req.Visibility, req.ParentFolderID).Scan(&req.UUID, &req.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}
