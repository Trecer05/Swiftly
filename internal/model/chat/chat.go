package chat

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	ID     int    `json:"id"`
	ChatID int    `json:"chat_id"`
	Text   string    `json:"text"`
	Author Client    `json:"author"`
	Time   time.Time `json:"time"`
}

// OnChat, OnLeave - проверка, в чате ли пользователь конкретный, для функционала прочитанных сообщений
type ChatRoom struct {
	ID int
	Name string
	Lock sync.RWMutex
	WG sync.WaitGroup
	Users map[*Client]bool
	Broadcast chan Message
	LastMessage *Message
	OnChat chan *Client
	OnLeave chan *Client
	ErrCh chan error
}

type Client struct {
	ID    int	`json:"id"`
	Name  string `json:"name"`
	Send  chan Message
	Conn  *websocket.Conn
}