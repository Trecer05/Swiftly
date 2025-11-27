package chat

import (
	"time"
)

type UserProject struct {
	ID        int    `json:"id"`
	IsAdmin   bool   `json:"is_admin"`
	Name      string `json:"name"`
	Description string `json:"description"`
	Users     []UserShort `json:"users"`
	Tasks []UserTask `json:"tasks"`
}

type UserShort struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	AvatarURL string `json:"avatar_url"`
}

type UserTask struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Description string `json:"description"`
	Priority    Priority `json:"priority"`
	Label     string `json:"label"`
	EndTime   time.Time `json:"end_time"`
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

type TeamInfo struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Description string `json:"description"`
	Users     []UserShort `json:"users"`
}

type TeamApplication struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Status    string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
