package kafka

type PasswordEdit struct {
	UserID int
	NewPassword string
	OldPassword string
}

type PhoneEdit struct {
	NewPhone string
	UserID int
}

type EmailEdit struct {
	NewEmail string
	UserID int
}
