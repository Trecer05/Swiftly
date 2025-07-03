package chat

import (
	"encoding/json"
	"log"

	"github.com/Trecer05/Swiftly/internal/model/chat"
	"github.com/gorilla/websocket"
)

func NewChat() *chat.ChatRoom {
	return &chat.ChatRoom{
		Broadcast: make(chan chat.Message),
		OnChat:    make(chan *chat.Client),
		OnLeave:   make(chan *chat.Client),
		Users:     make(map[*chat.Client]bool),
		ErrCh:     make(chan error),
	}
}

func Broadcast(ch *chat.ChatRoom, rooms map[int]*chat.ChatRoom) {
	for {
		select {
		case user := <-ch.OnChat:
			ch.Users[user] = true
		case user := <-ch.OnLeave:
			if _, ok := ch.Users[user]; ok {
				ch.Lock.RLock()
				delete(ch.Users, user)
				close(user.Send)
				ch.Lock.RUnlock()
			}
		case msg := <-ch.Broadcast:
			for user := range ch.Users {
				select {
					case user.Send <- msg:
					default:
						log.Printf("User %d disconnected\n", user.ID)
						ch.Lock.RLock()
						close(user.Send)
						delete(ch.Users, user)
						ch.Lock.RUnlock()
				}
			}
		}
	}
}

func WriteMessage(user *chat.Client) {
	defer func() {
		user.Conn.Close()
	}()

	for {
		msg, ok := <-user.Send
		if !ok {
			return
		}

		msgByte, err := json.Marshal(msg)
		if err != nil {
			log.Println(err)
			return
		}

		if err := user.Conn.WriteJSON(msgByte); err != nil {
			log.Println(err)
			return
		}
	}	
}

func ReadMessage(user *chat.Client, room *chat.ChatRoom) {
	defer func() {
		user.Conn.Close()
		room.OnLeave <- user
	}()

	for {
		var message chat.Message

		_, msg, err := user.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println(err)
			break
		}

		room.Broadcast <- message
	}
}