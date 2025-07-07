package chat

import (
	"encoding/json"
	"log"

	chatErrors "github.com/Trecer05/Swiftly/internal/errors/chat"
	models "github.com/Trecer05/Swiftly/internal/model/chat"
	redis "github.com/Trecer05/Swiftly/internal/repository/cache/chat"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/chat"

	"github.com/gorilla/websocket"
)

func ReadMessage(chatId int, conn *websocket.Conn, rds *redis.Manager, manager *manager.Manager, chatType models.ChatType) {
	defer func() {
		conn.Close()
	}()

	for {
		var message models.Message

		_, msg, err := conn.ReadMessage()
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
		
		switch message.Type {
		case models.Typing, models.StopTyping:
			_ = rds.SendToUser(chatId, message, chatType)
		case models.Default:
			_ = rds.SendToUser(chatId, message, chatType)

			var dbType models.DBType
			switch chatType {
			case models.TypePrivate:
				dbType = models.DBChat
			case models.TypeGroup:
				dbType = models.DBGroup
			default:
				log.Println("Chat type is not private or group")
			}
			
			if err := manager.SaveMessage(message, dbType); err != nil {
				log.Println("Failed to save message:", err)
			}
		default:
			log.Println("Unknown message type:", message.Type)
		}
	}
}

func SendChatHistory(conn *websocket.Conn, mgr *manager.Manager, chatId, limit, offset int, chatType models.ChatType) error {
	var err error
	var messages []models.Message

	if chatType == models.TypePrivate {
		messages, err = mgr.GetChatMessages(chatId, limit, offset)
	}
	if chatType == models.TypeGroup {
		messages, err = mgr.GetGroupMessages(chatId, limit, offset)
	}
	if err != nil {
		if err == chatErrors.ErrNoMessages {
			return nil
		}
		return err
	}

	for _, msg := range messages {
		if err := conn.WriteJSON(msg); err != nil {
			return err
		}
	}

	return nil
}
