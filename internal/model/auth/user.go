package auth

type User struct {
	ID int `json:"id,omitempty"`
	Email string `json:"email"`
	Number string `json:"number"`
	Password string `json:"password"`
}