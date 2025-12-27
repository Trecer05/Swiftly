package cloud

import (
	"encoding/json"
	logger "github.com/Trecer05/Swiftly/internal/config/logger"
	"strconv"

	models "github.com/Trecer05/Swiftly/internal/model/cloud"
	"github.com/gorilla/websocket"
)

func (manager *WebSocketManager) ListenPubSub(teamID int, msgCh chan models.Envelope) {
    manager.MU.Lock()
    if manager.SubscribedSessions[teamID] {
        manager.MU.Unlock()
        return
    }
    manager.SubscribedSessions[teamID] = true
    manager.MU.Unlock()

    pubsub := manager.RDB.Subscribe(ctx, "cloud:team:"+strconv.Itoa(teamID))
    ch := pubsub.Channel()

    go func() {
        for msg := range ch {
            var m models.Envelope
            if err := json.Unmarshal([]byte(msg.Payload), &m); err != nil {
                logger.Logger.Println("Invalid pubsub message:", err)
                continue
            }
            m.TeamID = teamID
            msgCh <- m
        }
    }()
}

func (manager *WebSocketManager) AddUser(userID, teamID int, conn *websocket.Conn) {
    manager.MU.Lock()
    defer manager.MU.Unlock()

    if _, ok := manager.Sessions[teamID]; !ok {
        manager.Sessions[teamID] = make(map[int]*websocket.Conn)
    }

    manager.Sessions[teamID][userID] = conn
}

func (manager *WebSocketManager) RemoveUser(userID, teamID int) error {
    manager.MU.Lock()
    defer manager.MU.Unlock()

    conn := manager.Sessions[teamID][userID]
    conn.Close()
    delete(manager.Sessions[teamID], userID)
    if len(manager.Sessions[teamID]) == 0 {
        delete(manager.Sessions, teamID)
    }
    
    return nil
}

func (manager *WebSocketManager) SendLocalMessage(userID, teamID int, messages <-chan models.Envelope) {
	for message := range messages {
		manager.MU.RLock()
		conn := manager.Sessions[teamID][userID]
		manager.MU.RUnlock()

		if err := conn.WriteJSON(message); err != nil {
			logger.Logger.Println("write error:", err)
		}
	}
}

func (manager *WebSocketManager) SendToUser(teamID int, message models.Envelope) error {
	data, err := json.Marshal(message)
	if err != nil {
		logger.Logger.Println("Error marshalling message:", err)
		return err
	}

	channel := "team" + ":" + strconv.Itoa(teamID)
	return manager.RDB.Publish(ctx, channel, data).Err()
}
