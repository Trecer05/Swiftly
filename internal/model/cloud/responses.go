package cloud

import (
	"time"

	"github.com/google/uuid"
)

type FileResponse struct {
	UUID        uuid.UUID `json:"uuid"`
	FolderID    uuid.UUID `json:"folder_id"`
	DisplayName string    `json:"display_name"`
	MimeType    string    `json:"mime_type"`
	Size        int64     `json:"size"`
	Visibility  string    `json:"visibility"`
	CreatedBy   int       `json:"created_by"`
	OwnerID     int       `json:"owner_id"`
	OwnerType   string    `json:"owner_type"`
	UploadedAt  string    `json:"uploaded_at"`
	UpdatedAt   string    `json:"updated_at"`
	Hash        string    `json:"hash"`
	Version     int       `json:"version"`
}

type FilesAndFoldersResponse struct {
	Files   []File   `json:"files"`
	Folders []Folder `json:"folders"`
}

type SharedFilesAndFoldersResponse struct {
	Files   []FileShare
	Folders []FolderShare
}

type FileUpdateResponse struct {
	UUID        uuid.UUID `json:"uuid"`
	UpdatedAt   time.Time `json:"updated_at"`
	NewFilename string    `json:"new_filename"`
}

type ShareLinkResponse struct {
	Link string `json:"link"`
}
