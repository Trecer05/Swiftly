package cloud

import "github.com/google/uuid"

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

type OwnerType string

const (
	OwnerTypeUser OwnerType = "user"
	OwnerTypeTeam OwnerType = "team"
)

type VisibilityType string

const (
	VisibilityPrivate VisibilityType = "private"
	VisibilityPublic  VisibilityType = "public"
	VisibilityShared  VisibilityType = "shared"
)

type File struct {
	UUID             uuid.UUID
	FolderID         uuid.UUID
	StoragePath      string
	OriginalFilename string
	DisplayName      string
	MimeType         string
	Size             int64
	Visibility       VisibilityType
	CreatedBy        int
	OwnerID          int
	OwnerType        OwnerType
	UploadedAt       string
	UpdatedAt        string
	Hash             string
	Version          int
}

type FileShort struct {
	UUID             uuid.UUID
	StoragePath      string
	OriginalFilename string
	DisplayName      string
	MimeType         string
	Size             int64
}

type FileResponse struct {
	UUID        uuid.UUID
	FolderID    uuid.UUID
	DisplayName string
	MimeType    string
	Size        int64
	Visibility  string
	CreatedBy   int
	OwnerID     int
	OwnerType   string
	UploadedAt  string
	UpdatedAt   string
	Hash        string
	Version     int
}
