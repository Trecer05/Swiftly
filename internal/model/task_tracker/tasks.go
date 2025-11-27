package task_tracker

import (
	"time"
)

type Task struct {
	ID        int64     `json:"id"`
	AuthorID  int64     `json:"author_id"`
	DeveloperID int64   `json:"developer_id"`
	Label     string    `json:"label"`
	Title     string    `json:"title"`
	Description *string `json:"description"`
	Status    string    `json:"status"`
	Priority  Priority  `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

type PriorityType string

const (
	PriorityLow PriorityType = "low"
	PriorityMedium PriorityType = "medium"
	PriorityHigh PriorityType = "high"
)

type Priority struct {
	Type PriorityType `json:"type"`
	Title string `json:"title"`
	HexColor string `json:"hex_color"`
}
