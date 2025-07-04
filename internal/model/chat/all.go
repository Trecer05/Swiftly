package chat

import "github.com/gorilla/websocket"

type ChatRooms struct {
	Conn *websocket.Conn
	Messages chan Message
	Rooms []*ChatRoom
}