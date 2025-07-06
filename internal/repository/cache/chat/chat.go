package chat

import (
	"encoding/json"
	"log"
	"strconv"

	models "github.com/Trecer05/Swiftly/internal/model/chat"
	"github.com/gorilla/websocket"
)

func (manager *Manager) ListenPubSub(chatId int, msgCh chan models.Message) {
	manager.MU.Lock()
	if manager.SubscribedChats[chatId] {
		manager.MU.Unlock()
		return
	}
	manager.SubscribedChats[chatId] = true
	manager.MU.Unlock()

	pubsub := manager.RDB.Subscribe(ctx, "chat:"+strconv.Itoa(chatId))
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

func (manager *Manager) AddUser(userId, chatId int, conn *websocket.Conn) {
	manager.MU.Lock()
	defer manager.MU.Unlock()

	if _, ok := manager.Sessions[chatId]; !ok {
		manager.Sessions[chatId] = make(map[int]*websocket.Conn)
	}

	manager.Sessions[chatId][userId] = conn
}

func (manager *Manager) RemoveUser(userId, chatId int) error {
	manager.MU.Lock()
	defer manager.MU.Unlock()

	if conn, ok := manager.Sessions[chatId][userId]; ok {
		if err := conn.Close(); err != nil {
			return err
		} else {
			delete(manager.Sessions, userId)
			return nil
		}
	}

	return nil
}

func (manager *Manager) SendLocalMessage(userId, chatId int, messages <-chan models.Message) {
	for message := range messages {
		manager.MU.RLock()
		conn, ok := manager.Sessions[chatId][userId]
		manager.MU.RUnlock()

		if ok {
			if err := conn.WriteJSON(message); err != nil {
				log.Println("write error:", err)
			}
		}
	}
}

func (manager *Manager) SendToUser(chatId int, message models.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return err
	}

	if err := manager.RDB.Publish(ctx, "chat:"+strconv.Itoa(chatId), data).Err(); err != nil {
		log.Println("Error publishing message:", err)
		return err
	}
	return nil
}
