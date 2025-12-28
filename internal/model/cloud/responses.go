package cloud

import (
	"time"

	"github.com/google/uuid"
)

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

type FilesAndFoldersResponse struct {
	Files   []File
	Folders []Folder
}

type FileUpdateResponse struct {
	UUID        	uuid.UUID
	UpdatedAt   	time.Time
	NewFilename     string
}
