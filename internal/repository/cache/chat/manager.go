package chat

import (
	"context"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type Manager struct {
	WG   *sync.WaitGroup
	MU   *sync.RWMutex
	RDB  *redis.Client
	Sessions map[int]map[int]*websocket.Conn
	SubscribedChats map[int]bool
}

var ctx = context.Background()

func NewChatManager(addr string) *Manager {
	return &Manager{
		WG:   &sync.WaitGroup{},
		MU:   &sync.RWMutex{},
		RDB:  redis.NewClient(&redis.Options{Addr: addr}),
		Sessions: make(map[int]map[int]*websocket.Conn),
		SubscribedChats: make(map[int]bool),
	}
}

func (m *Manager) Close() {
	if m.RDB != nil {
		m.RDB.Close()
	}
}