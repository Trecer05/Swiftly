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
	Sessions map[models.SessionKey]map[int]*Calls
	SubscribedChats map[models.SessionKey]bool
}

type Calls struct {
	WS *websocket.Conn
	Room *models.Room
}

var ctx = context.Background()

func newRoom() *models.Room {
	var room models.Room

	room.Peers = make(map[string]*models.PeerState)
	room.Published = make(map[string]*models.PublishedTrack)

	return &room
}

func NewChatManager(addr string) *Manager {
	room := newRoom()

	var calls Calls

	calls.Room = room

	return &Manager{
		WG:   &sync.WaitGroup{},
		MU:   &sync.RWMutex{},
		RDB:  redis.NewClient(&redis.Options{Addr: addr}),
		Sessions: make(map[models.SessionKey]map[int]*Calls),
		SubscribedChats: make(map[models.SessionKey]bool),
	}
}

func (m *Manager) Close() {
	if m.RDB != nil {
		m.RDB.Close()
	}
}