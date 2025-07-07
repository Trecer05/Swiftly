package chat

import (
	"encoding/json"
	"log"

	chatErrors "github.com/Trecer05/Swiftly/internal/errors/chat"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/chat"
	models "github.com/Trecer05/Swiftly/internal/model/chat"
	redis "github.com/Trecer05/Swiftly/internal/repository/cache/chat"

	"github.com/gorilla/websocket"
)

func ReadMessage(chatId int, conn *websocket.Conn, rds *redis.Manager, manager *manager.Manager) {
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

		rds.SendToUser(chatId, message)
		if err := manager.SaveChatMessage(message); err != nil {
			log.Println("Failed to save message:", err)
		}
	}
}

func SendChatHistory(conn *websocket.Conn, mgr *manager.Manager, chatId, limit, offset int) error {
	messages, err := mgr.GetChatMessages(chatId, limit, offset)
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
