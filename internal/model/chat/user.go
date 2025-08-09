package chat

type User struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Username string `json:"username"`
	Description string `json:"description"`
}

type RegisterUser struct {
	Username string `json:"username"`
	Name string `json:"name"`
	Description string `json:"description"`
}
