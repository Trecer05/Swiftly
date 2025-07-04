package chat

import (
	"encoding/json"
	"log"

	models "github.com/Trecer05/Swiftly/internal/model/chat"
	redis "github.com/Trecer05/Swiftly/internal/repository/cache/chat"

	"github.com/gorilla/websocket"
)

func ReadMessage(chatId int, conn *websocket.Conn, mgr *redis.Manager) {
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

		mgr.SendToUser(chatId, message)
	}
}