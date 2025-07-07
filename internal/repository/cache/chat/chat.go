package chat

import (
	"encoding/json"
	"log"
	"strconv"

	models "github.com/Trecer05/Swiftly/internal/model/chat"
	"github.com/gorilla/websocket"
)

func (manager *Manager) ListenPubSub(chatId int, msgCh chan models.Message, chatType models.ChatType) {
	manager.MU.Lock()
	key := models.SessionKey{Type: chatType, ID: chatId}
	if manager.SubscribedChats[key] {
		manager.MU.Unlock()
		return
	}
	manager.SubscribedChats[key] = true
	manager.MU.Unlock()

	pubsub := manager.RDB.Subscribe(ctx, string(key.Type) + ":" + strconv.Itoa(key.ID))
	ch := pubsub.Channel()

	go func() {
		for msg := range ch {
			var m models.Message
			if err := json.Unmarshal([]byte(msg.Payload), &m); err != nil {
				log.Println("Invalid pubsub message:", err)
				continue
			}
			m.ChatID = chatId

			msgCh <- m
		}
	}()
}

func (manager *Manager) AddUser(userId, chatId int, conn *websocket.Conn, chatType models.ChatType) {
	manager.MU.Lock()
	defer manager.MU.Unlock()

	if _, ok := manager.Sessions[models.SessionKey{Type: chatType, ID: chatId}]; !ok {
		manager.Sessions[models.SessionKey{Type: chatType, ID: chatId}] = make(map[int]*websocket.Conn)
	}

	manager.Sessions[models.SessionKey{Type: chatType, ID: chatId}][userId] = conn
}

func (manager *Manager) RemoveUser(userId, chatId int, chatType models.ChatType) error {
	manager.MU.Lock()
	defer manager.MU.Unlock()

	key := models.SessionKey{Type: chatType, ID: chatId}
	if conn, ok := manager.Sessions[key][userId]; ok {
		conn.Close()
		delete(manager.Sessions[key], userId)
		if len(manager.Sessions[key]) == 0 {
			delete(manager.Sessions, key)
		}
	}

	return nil
}

func (manager *Manager) SendLocalMessage(userId, chatId int, messages <-chan models.Message, chatType models.ChatType) {
	for message := range messages {
		manager.MU.RLock()
		conn, ok := manager.Sessions[models.SessionKey{Type: chatType, ID: chatId}][userId]
		manager.MU.RUnlock()

		if ok {
			if err := conn.WriteJSON(message); err != nil {
				log.Println("write error:", err)
			}
		}
	}
}

func (manager *Manager) SendToUser(chatId int, message models.Message, chatType models.ChatType) error {
	data, err := json.Marshal(message)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return err
	}

	channel := string(chatType) + ":" + strconv.Itoa(chatId)
    return manager.RDB.Publish(ctx, channel, data).Err()
}
