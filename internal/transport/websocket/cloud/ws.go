package cloud

import (
	"encoding/json"

	models "github.com/Trecer05/Swiftly/internal/model/cloud"
	redis "github.com/Trecer05/Swiftly/internal/repository/cache/cloud"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/cloud"
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
		case models.FileCreateType:
		    _ = rds.SendToTeam(message.TeamID, message)
		case models.FileDeleteType:
		    _ = rds.SendToTeam(message.TeamID, message)
		case models.FileUpdateType:
		    _ = rds.SendToTeam(message.TeamID, message)
		case models.FileNameUpdateType:
		    _ = rds.SendToTeam(message.TeamID, message)
		case models.FolderCreateType:
		    _ = rds.SendToTeam(message.TeamID, message)
		case models.FolderDeleteType:
		    _ = rds.SendToTeam(message.TeamID, message)
		case models.FolderNameUpdateType:
		    _ = rds.SendToTeam(message.TeamID, message)
		case models.FolderMoveType:
		    _ = rds.SendToTeam(message.TeamID, message)
		case models.FileMoveType:
		    _ = rds.SendToTeam(message.TeamID, message)
		default:
			logger.Logger.Error("Unknown message type:", message.Type)
		}
	}
}