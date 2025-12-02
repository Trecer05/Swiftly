package task_tracker

import (
	"context"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/redis/go-redis/v9"
)

type WebSocketManager struct {
	WG   *sync.WaitGroup
	MU   *sync.RWMutex
	RDB  *redis.Client
	Sessions map[int]map[int]*websocket.Conn
	SubscribedSessions map[int]bool
}

var ctx = context.Background()

func NewTaskManager(addr string) *WebSocketManager {
	return &WebSocketManager{
		WG:   &sync.WaitGroup{},
		MU:   &sync.RWMutex{},
		RDB:  redis.NewClient(&redis.Options{Addr: addr}),
		Sessions: make(map[int]map[int]*websocket.Conn),
		SubscribedSessions: make(map[int]bool),
	}
}

func (m *WebSocketManager) Close() {
	m.MU.Lock()
	defer m.MU.Unlock()
	
	for _, team := range m.Sessions {
		for _, conn := range team {
			conn.WriteMessage(websocket.CloseMessage, []byte{})
			conn.Close()
		}
	}
	
	m.Sessions = make(map[int]map[int]*websocket.Conn)
	
	if m.RDB != nil {
		m.RDB.Close()
	}
}
