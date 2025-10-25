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
