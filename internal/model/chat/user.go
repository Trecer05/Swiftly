package chat

type User struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Username string `json:"username"`
	Description string `json:"description"`
	AvatarUrl string `json:"avatar_url"`
}

type RegisterUser struct {
	Username string `json:"username"`
	Name string `json:"name"`
	Description string `json:"description"`
}

type StartUserInfo struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Username string `json:"username"`
	Description string `json:"description"`
	AvatarUrl string `json:"avatar_url"`
	Projects []UserProject `json:"projects"`
}
