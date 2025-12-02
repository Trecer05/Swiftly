package task_tracker

import (
	"time"
)

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description *string `json:"description,omitempty"`
	CreatorID   int
	TeamID      int
	ColumnID    int
	PositionInColumn int `json:"position_in_column"`
	Deadline    *time.Time `json:"deadline,omitempty"`
	Label       *string `json:"label,omitempty"`
	ExecutorUsername string `json:"executor_username,omitempty"`
}

type TaskTitleUpdateRequest struct {
	ID int `json:"id"`
	ColumnID int `json:"column_id"`
	PositionInColumn int `json:"position_in_column"`
	Title string `json:"title"`
}

type TaskDeadlineUpdateRequest struct {
	ID int `json:"id"`
	ColumnID int `json:"column_id"`
	PositionInColumn int `json:"position_in_column"`
	Deadline *time.Time `json:"deadline,omitempty"`
}

type TaskColumnUpdateRequest struct {
	ID int `json:"id"`
	ColumnID int `json:"column_id"`
	PositionInColumn int `json:"position_in_column"`
}

type TaskDescriptionUpdateRequest struct {
	ID int `json:"id"`
	ColumnID int `json:"column_id"`
	PositionInColumn int `json:"position_in_column"`
	Description string `json:"description"`
}

type TaskDeveloperUpdateRequest struct {
	ID int `json:"id"`
	ColumnID int `json:"column_id"`
	PositionInColumn int `json:"position_in_column"`
	DeveloperID int `json:"developer_id"`
}

type TaskStatusUpdateRequest struct {
	ID int `json:"id"`
	ColumnID int `json:"column_id"`
	PositionInColumn int `json:"position_in_column"`
	Status string `json:"status"`
}
