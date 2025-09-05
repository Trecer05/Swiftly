package chat

import (
	"time"
)

type Message struct {
	ID     int    `json:"id"`
	ChatID int    `json:"chat_id"`
	Type   MessageType `json:"type"`
	Read   bool `json:"read"`
	Text   string    `json:"text"`
	Author Client    `json:"author"`
	Time   time.Time `json:"time"`
}

type ChatRoom struct {
	ID int `json:"id"`
	Name string `json:"name"`
	LastMessage *Message `json:"last_message"`
	Type ChatType `json:"type"`
}

type Client struct {
	ID    int	`json:"id"`
	Name  string `json:"name"`
}

type Group struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
}
