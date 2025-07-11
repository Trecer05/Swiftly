package chat

import (
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	ID     int    `json:"id"`
	ChatID int    `json:"chat_id"`
	Type   MessageType `json:"type"`
	Text   string    `json:"text"`
	Author Client    `json:"author"`
	Time   time.Time `json:"time"`
}

type ChatRoom struct {
	ID int
	Name string
	Users map[*Client]bool
	LastMessage *Message
	Type ChatType
}

type GroupCreate struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Users []Client `json:"users"`
	OwnerID int `json:"owner_id"`
}

type Client struct {
	ID    int	`json:"id"`
	Name  string `json:"name"`
	Send  chan Message
	Conn  *websocket.Conn
}