package chat

import "github.com/gorilla/websocket"

type Chats struct {
	Conn *websocket.Conn
	Messages chan Message
	Chats map[int]Chat
}