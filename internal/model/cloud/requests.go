package cloud

import "github.com/google/uuid"

type CreateFileRequest struct {
	FolderID    *uuid.UUID `json:"folder_id,omitempty"` // Можно добавить в корень хранилища
	DisplayName string    `json:"display_name"`
	Visibility  VisibilityType    `json:"visibility"`
	OwnerType   OwnerType    `json:"owner_type"`
}
