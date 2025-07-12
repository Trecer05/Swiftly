package chat

type GroupCreate struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Users []Client `json:"users"`
	OwnerID int
}

type Users struct {
	Users []Client `json:"users"`
}
