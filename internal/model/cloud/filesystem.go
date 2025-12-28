package cloud

import (
	"time"

	"github.com/google/uuid"
)

const MaxUploadSize = 50 << 20

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

type Folder struct {
	UUID           uuid.UUID  `json:"uuid"`
	Name           string     `json:"name"`
	OwnerID        int        `json:"owner_id"`
	OwnerType      OwnerType  `json:"owner_type"`
	ChildFolders   []Folder   `json:"child_folders,omitempty"`
	Files          []File     `json:"files,omitempty"`
	ParentFolderID *uuid.UUID `json:"parent_folder_id,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type File struct {
	UUID             uuid.UUID      `json:"uuid"`
	FolderID         *uuid.UUID     `json:"folder_id,omitempty"`
	StoragePath      string         `json:"storage_path"`
	OriginalFilename string         `json:"original_filename"`
	DisplayName      string         `json:"display_name"`
	MimeType         string         `json:"mime_type"`
	Size             int64          `json:"size"`
	Visibility       VisibilityType `json:"visibility"`
	CreatedBy        int            `json:"created_by"`
	OwnerID          int            `json:"owner_id"`
	OwnerType        OwnerType      `json:"owner_type"`
	UploadedAt       string         `json:"uploaded_at"`
	UpdatedAt        string         `json:"updated_at"`
	Hash             string         `json:"hash"`
	Version          int            `json:"version"`
}

type FileShort struct {
	UUID             uuid.UUID
	StoragePath      string
	OriginalFilename string
	DisplayName      string
	MimeType         string
	Size             int64
}

// type Folder struct {
// 	UUID           uuid.UUID
// 	Name           string
// 	OwnerID        int
// 	OwnerType      OwnerType
// 	ParentFolderID *uuid.UUID
// 	CreatedAt      time.Time
// 	UpdatedAt      time.Time
// }
