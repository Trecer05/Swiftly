package cloud

import "github.com/google/uuid"

type CreateFileRequest struct {
	FolderID    *uuid.UUID `json:"folder_id,omitempty"` // Можно добавить в корень хранилища
	DisplayName string    `json:"display_name"`
	Visibility  VisibilityType    `json:"visibility"`
	OwnerType   OwnerType    `json:"owner_type"`
}

type FilenameUpdateRequest struct {
	NewFilename string `json:"new_filename"`
}

type FoldernameUpdateRequest struct {
	NewFoldername string `json:"new_foldername"`
}

type CreateFolderRequest struct {
	ParentID    *uuid.UUID `json:"parent_id,omitempty"` // Можно добавить в корень хранилища
	DisplayName string    `json:"display_name"`
	Visibility  VisibilityType    `json:"visibility"`
	OwnerType   OwnerType    `json:"owner_type"`
}

type UpdateFileRequest struct {
	ParentID    *uuid.UUID `json:"parent_id,omitempty"` // Можно добавить в корень хранилища
	DisplayName string    `json:"display_name"`
	Visibility  VisibilityType    `json:"visibility"`
	OwnerType   OwnerType    `json:"owner_type"`
}

type MoveUserFileRequest struct {
	NewFolderID *uuid.UUID `json:"new_folder_id,omitempty"` // Можно добавить в корень хранилища
}

type MoveUserFolderRequest struct {
	NewFolderID *uuid.UUID `json:"new_folder_id,omitempty"` // Можно добавить в корень хранилища
	FolderName  string    `json:"folder_name,omitempty"`
}
