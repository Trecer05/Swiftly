package cloud

import (
	"time"

	"github.com/google/uuid"
)

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

type Folder struct {
	UUID       uuid.UUID
	Name 	 string
	OwnerID  int
	OwnerType OwnerType
	ParentFolderID   *uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
