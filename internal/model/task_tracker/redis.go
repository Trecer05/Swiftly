package task_tracker

import (
	"encoding/json"
	"time"
)

type Types string

const (
	TaskCreateType Types = "create"
	TaskDeleteType Types = "delete"
	TaskTitleUpdateType Types = "title_update"
	TaskDescriptionUpdateType Types = "description_update"
	TaskDeadlineUpdateType Types = "deadline_update"
	TaskStatusUpdateType Types = "status_update"
	TaskDeveloperUpdateType Types = "developer_update"
	TaskColumnUpdateType Types = "column_update"
)

type Envelope struct {
	TeamID int `json:"team_id"`
	Data   json.RawMessage `json:"data"`
	Type   Types `json:"type"`
}

type TaskDeleteMessage struct {
	ID int `json:"id"`
	ColumnID int `json:"column_id"`
	ColumnPosition int `json:"column_position"`
	Message string `json:"message"`
}

type TaskCreateMessage struct {
	ID int `json:"id"`
	Title string `json:"title"`
	ColumnID int `json:"column_id"`
	ColumnPosition int `json:"column_position"`
	Deadline *time.Time `json:"deadline,omitempty"`
	Message string `json:"message"`
}

type TaskTitleUpdateMessage struct {
	ID int `json:"id"`
	Title string `json:"title"`
	ColumnID int `json:"column_id"`
	ColumnPosition int `json:"column_position"`
	Message string `json:"message"`
}

type TaskDeadlineUpdateMessage struct {
	ID int `json:"id"`
	Deadline *time.Time `json:"deadline"`
	ColumnID int `json:"column_id"`
	ColumnPosition int `json:"column_position"`
	Message string `json:"message"`
}

type TaskStatusUpdateMessage struct {
	ID int `json:"id"`
	Status string `json:"status"`
	ColumnID int `json:"column_id"`
	ColumnPosition int `json:"column_position"`
	Message string `json:"message"`
}

type TaskDeveloperUpdateMessage struct {
	ID int `json:"id"`
	DeveloperID int `json:"developer_id"`
	ColumnID int `json:"column_id"`
	ColumnPosition int `json:"column_position"`
	Message string `json:"message"`
}

type TaskColumnUpdateMessage struct {
	ID int `json:"id"`
	ColumnID int `json:"column_id"`
	ColumnPosition int `json:"column_position"`
	Message string `json:"message"`
}

type TaskDescriptionUpdateMessage struct {
	ID int `json:"id"`
	Description string `json:"description"`
	ColumnID int `json:"column_id"`
	ColumnPosition int `json:"column_position"`
	Message string `json:"message"`
}

type NotificationType string

type NotificationMessage struct {}

// TODO: сделать
type Notifications struct {
	Type NotificationType `json:"type"`
	Message NotificationMessage `json:"message"`
}
