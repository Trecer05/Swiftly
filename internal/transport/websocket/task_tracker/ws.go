package task_tracker

import (
	"encoding/json"

	models "github.com/Trecer05/Swiftly/internal/model/task_tracker"
	redis "github.com/Trecer05/Swiftly/internal/repository/cache/task_tracker"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/task_tracker"
	logger "github.com/Trecer05/Swiftly/internal/config/logger"

	"github.com/gorilla/websocket"
)

func ReadMessage(chatId int, conn *websocket.Conn, rds *redis.WebSocketManager, manager *manager.Manager) {
	defer func() {
		conn.Close()
	}()

	for {
		var message models.Envelope

		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				logger.Logger.Error("Unexpected close error:", err)
			}
			break
		}

		if err := json.Unmarshal(msg, &message); err != nil {
			logger.Logger.Error("Failed to unmarshal message:", err)
			break
		}

		switch message.Type {
		case models.TaskCreateType:
		    _ = rds.SendToUser(message.TeamID, message)
		case models.TaskDeleteType:
		    _ = rds.SendToUser(message.TeamID, message)
		case models.TaskDeveloperUpdateType:
		    _ = rds.SendToUser(message.TeamID, message)
		case models.TaskColumnUpdateType:
		    _ = rds.SendToUser(message.TeamID, message)
		case models.TaskStatusUpdateType:
		    _ = rds.SendToUser(message.TeamID, message)
		case models.TaskDescriptionUpdateType:
		    _ = rds.SendToUser(message.TeamID, message)
		case models.TaskTitleUpdateType:
		    _ = rds.SendToUser(message.TeamID, message)
		default:
			logger.Logger.Error("Unknown message type:", message.Type)
		}
	}
}