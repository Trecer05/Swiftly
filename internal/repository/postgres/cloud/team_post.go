package cloud

import models "github.com/Trecer05/Swiftly/internal/model/cloud"

func (manager *Manager) CreateTeamFile(req *models.File) error {
	err := manager.Conn.QueryRow(`INSERT INTO public.files(
		storage_path, 
		original_filename, 
		display_name, 
		folder_id,
		mime_type, 
		size, 
		visibility, 
		created_by, 
		owner_id, 
		owner_type, 
		hash
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'team', $10) 
	RETURNING uuid, uploaded_at, updated_at`,
		req.StoragePath,
		req.OriginalFilename,
		req.DisplayName,
		req.FolderID,
		req.MimeType,
		req.Size,
		req.Visibility,
		req.CreatedBy,
		req.OwnerID,
		req.Hash).Scan(&req.UUID, &req.UploadedAt, &req.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

// CREATE TABLE folders (
//     uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
//     name VARCHAR(255) NOT NULL,
//     created_by INTEGER NOT NULL,
//     owner_id INTEGER,  -- Владелец (пользователь ИЛИ команда)
//     owner_type VARCHAR(10) CHECK (owner_type IN ('user', 'team')),
//     storage_path VARCHAR(500) NOT NULL,
//     parent_folder_id UUID REFERENCES folders(uuid) ON DELETE CASCADE,
//     created_at TIMESTAMP DEFAULT NOW(),
//     updated_at TIMESTAMP DEFAULT NOW(),

//     CONSTRAINT valid_owner CHECK (
//         (owner_type = 'user' AND owner_id IS NOT NULL) OR
//         (owner_type = 'team' AND owner_id IS NOT NULL)
//     )
// );

func (manager *Manager) CreateTeamFolder(req *models.Folder, storagePath string) error {
	err := manager.Conn.QueryRow(`INSERT INTO public.folders(
		name,
		created_by,
		owner_id,
		owner_type,
		storage_path,
		visibility,
		parent_folder_id
		)
	VALUES ($1, $2, $3, 'team', $4, $5, $6)
	RETURNING uuid, created_at, updated_at`,
		req.Name,
		req.CreatedBy,
		req.OwnerID,
		storagePath,
		req.Visibility,
		req.ParentFolderID).Scan(&req.UUID, &req.CreatedAt, &req.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}
