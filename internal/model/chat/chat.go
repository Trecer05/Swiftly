package chat

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	ID     string    `json:"id"`
	Text   string    `json:"text"`
	Author Client    `json:"author"`
	Time   time.Time `json:"time"`
}

type Chat struct {
	Lock sync.Mutex
	WG  sync.WaitGroup
	Broadcaster chan Message
	Clients   map[*Client]bool
	ErrCh    chan error
}

type Client struct {
	ID    int	`json:"id"`
	Name  string `json:"name"`
	Conn  *websocket.Conn
}