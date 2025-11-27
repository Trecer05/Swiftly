package chat

import (
	"time"
)

type GroupCreate struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Users []Client `json:"users"`
	OwnerID int
}

type Users struct {
	Users []Client `json:"users"`
}

type ChatCreate struct {
	UserSender string `json:"user_sender"`
	UserReceiver string `json:"user_receiver"`
}

type ProfileEdit struct {
	Name *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type GroupEdit struct {
	Name *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type PasswordEdit struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type PhoneEdit struct {
	NewPhone string `json:"new_phone"`
}

type EmailEdit struct {
	NewEmail string `json:"new_email"`
}

type UserRoleEdit struct {
	NewRole string `json:"new_role"`
}

type TeamUser struct {
	ID int
	Username string `json:"username"`
	Role string `json:"role"`
}

type TeamCreate struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Users []TeamUser `json:"users,omitempty"`
	OwnerID int
}

type TeamEdit struct {
	ID int
	OwnerID int
	Name *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type ApplicationStatus string

const (
	TeamApplicationStatusPending ApplicationStatus = "pending"
	TeamApplicationStatusAccepted ApplicationStatus = "accepted"
	TeamApplicationStatusRejected ApplicationStatus = "rejected"
)

type TeamApplicationUpdate struct {
	ID int
	Status ApplicationStatus `json:"status"`
}

type CreateJoinCode struct {
	ProjectID int
	CreatorID int
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	IsSingleUse *bool `json:"is_single_use,omitempty"`
}
