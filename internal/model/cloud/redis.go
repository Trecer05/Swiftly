package cloud

import (
	"encoding/json"
)

type MessageType string

const (
	FileCreateType = "file_create"
	FileDeleteType = "file_delete"
	FileUpdateType = "file_update"
	FileNameUpdateType = "file_name_update"
	FolderCreateType = "folder_create"
	FolderDeleteType = "folder_delete"
	FolderNameUpdateType = "folder_name_update"
	FolderMoveType = "folder_move"
	FileMoveType = "file_move"
)

type Envelope struct {
	TeamID int `json:"team_id"`
	Data   json.RawMessage `json:"data"`
	Type   MessageType `json:"type"`
}