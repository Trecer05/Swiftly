package chat

import (
	"context"
	"sync"

	models "github.com/Trecer05/Swiftly/internal/model/chat"
	"github.com/gorilla/websocket"

	"github.com/redis/go-redis/v9"
)

type Manager struct {
	WG   *sync.WaitGroup
	MU   *sync.RWMutex
	RDB  *redis.Client
	Sessions map[models.SessionKey]map[int]*websocket.Conn
	Calls map[models.CallsKey]*models.Room
	SubscribedChats map[models.SessionKey]bool
}

var ctx = context.Background()

func NewChatManager(addr string) *Manager {
	return &Manager{
		WG:   &sync.WaitGroup{},
		MU:   &sync.RWMutex{},
		RDB:  redis.NewClient(&redis.Options{Addr: addr}),
		Sessions: make(map[models.SessionKey]map[int]*websocket.Conn),
		Calls: make(map[models.CallsKey]*models.Room),
		SubscribedChats: make(map[models.SessionKey]bool),
	}
}

func (m *Manager) Close() {
	if m.RDB != nil {
		m.RDB.Close()
	}
}
