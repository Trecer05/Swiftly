package auth

type User struct {
	ID int `json:"id,omitempty"`
	Email *string `json:"email,omitempty"`
	Number string `json:"number"`
	Password string `json:"password"`
}