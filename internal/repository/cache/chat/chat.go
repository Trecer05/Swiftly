package chat

import (
	"encoding/json"
	"log"
	"strconv"

	chat "github.com/Trecer05/Swiftly/internal/repository/postgres/chat"
	models "github.com/Trecer05/Swiftly/internal/model/chat"
	"github.com/gorilla/websocket"
)

func (manager *Manager) ListenPubSub(chatId int, msgCh chan models.Message, statusCh chan models.Status, notifCh chan models.Notifications, chatType models.ChatType) {
    manager.MU.Lock()
    key := models.SessionKey{Type: chatType, ChatID: chatId}
    if manager.SubscribedChats[key] {
        manager.MU.Unlock()
        return
    }
    manager.SubscribedChats[key] = true
    manager.MU.Unlock()

    pubsub := manager.RDB.Subscribe(ctx, string(key.Type)+":"+strconv.Itoa(key.ChatID))
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

    notifications := manager.RDB.Subscribe(ctx, string(key.Type)+":"+strconv.Itoa(key.ChatID)+":notifications")
    notCh := notifications.Channel()

    go func() {
        for msg := range notCh {
            var notif models.Notifications
            if err := json.Unmarshal([]byte(msg.Payload), &notif); err != nil {
                log.Println("Invalid pubsub message:", err)
                continue
            }
            notifCh <- notif
        }
    }()

    statusSub := manager.RDB.Subscribe(ctx, string(key.Type)+":"+strconv.Itoa(key.ChatID)+":status")
    statusChSub := statusSub.Channel()

    go func() {
        for msg := range statusChSub {
            var status models.Status
            if err := json.Unmarshal([]byte(msg.Payload), &status); err != nil {
                log.Println("Invalid status message:", err)
                continue
            }
            statusCh <- status
        }
    }()
}

func (manager *Manager) AddUser(userId, chatId int, conn *websocket.Conn, chatType models.ChatType, dbManager *chat.Manager) {
    manager.MU.Lock()
    defer manager.MU.Unlock()

    if _, ok := manager.Sessions[models.SessionKey{Type: chatType, ChatID: chatId}]; !ok {
        manager.Sessions[models.SessionKey{Type: chatType, ChatID: chatId}] = make(map[int]*websocket.Conn)
    }

    manager.Sessions[models.SessionKey{Type: chatType, ChatID: chatId}][userId] = conn

    if userRooms, err := dbManager.GetUserRoomsForStatus(userId); err == nil {
        go manager.PublishUserStatus(userId, true, userRooms)
    }
}

func (manager *Manager) RemoveUser(userId, chatId int, chatType models.ChatType, dbManager *chat.Manager) error {
    manager.MU.Lock()
    defer manager.MU.Unlock()

    key := models.SessionKey{Type: chatType, ChatID: chatId}
    conn := manager.Sessions[key][userId]
    conn.Close()
    delete(manager.Sessions[key], userId)
    if len(manager.Sessions[key]) == 0 {
        delete(manager.Sessions, key)
    }

    if userRooms, err := dbManager.GetUserRoomsForStatus(userId); err == nil {
        go manager.PublishUserStatus(userId, false, userRooms)
    }

    return nil
}

func (manager *Manager) SendLocalMessage(userId, chatId int, messages <-chan models.Message, chatType models.ChatType) {
	for message := range messages {
		manager.MU.RLock()
		conn := manager.Sessions[models.SessionKey{Type: chatType, ChatID: chatId}][userId]
		manager.MU.RUnlock()

		if err := conn.WriteJSON(message); err != nil {
			log.Println("write error:", err)
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

func (manager *Manager) SendNotificationToUser(chatId int, message models.Message, chatType models.ChatType, notifType models.NotificationType) error {
	notif := models.Notifications{
		Type: notifType,
		Message: message,
	}

	data, err := json.Marshal(notif)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return err
	}

	channel := string(chatType) + ":" + strconv.Itoa(chatId) + ":notifications"
	return manager.RDB.Publish(ctx, channel, data).Err()
}

func (manager *Manager) SendLocalNotification(userId, chatId int, notifications <-chan models.Notifications, chatType models.ChatType) {
	for notif := range notifications {
		manager.MU.RLock()
		conn := manager.Sessions[models.SessionKey{Type: chatType, ChatID: chatId}][userId]
		manager.MU.RUnlock()

		if err := conn.WriteJSON(notif); err != nil {
			log.Println("Notification write error:", err)
			continue
		}
	}
}

func (manager *Manager) SendLocalStatus(userId, chatId int, statuses <-chan models.Status, chatType models.ChatType) {
	for status := range statuses {
		manager.MU.RLock()
		conn := manager.Sessions[models.SessionKey{Type: chatType, ChatID: chatId}][userId]
		manager.MU.RUnlock()

		if err := conn.WriteJSON(status); err != nil {
			log.Println("Status write error:", err)
		}
	}
}

func (manager *Manager) PublishUserStatus(userId int, online bool, userRooms []models.UserRoom) error {
    status := models.Status{
        Type:    "status",
        User_ID: userId,
        Online:  online,
    }

    data, err := json.Marshal(status)
    if err != nil {
        return err
    }

    for _, room := range userRooms {
        channel := string(room.Type) + ":" + strconv.Itoa(room.ID) + ":status"
        if err := manager.RDB.Publish(ctx, channel, data).Err(); err != nil {
            log.Printf("Failed to publish status to room %d: %v", room.ID, err)
        }
    }

    return nil
}

func (manager *Manager) CallSend(userId, chatId int, chatType models.ChatType) {
	key := models.SessionKey{Type: chatType, ChatID: chatId}
	chatSessions := manager.Sessions[key]

	for _, conn := range chatSessions {
		if err := conn.WriteJSON(models.Message{
			Type:    models.Call,
			Author: models.Client{
				ID: userId,
			},
			ChatID:  chatId,
		}); err != nil {
			log.Println("Error sending call:", err)
		}
	}
}
