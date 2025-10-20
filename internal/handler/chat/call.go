package chat

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	errors "github.com/Trecer05/Swiftly/internal/errors/auth"
	models "github.com/Trecer05/Swiftly/internal/model/chat"
	cache "github.com/Trecer05/Swiftly/internal/repository/cache/chat"
	service "github.com/Trecer05/Swiftly/internal/service/chat"
	serviceHttp "github.com/Trecer05/Swiftly/internal/transport/http"
	calls "github.com/Trecer05/Swiftly/internal/transport/webrtc"
)

var (
	roomMutex sync.RWMutex
)

func HandleChatCallConnection(w http.ResponseWriter, r *http.Request, rds *cache.Manager) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	defer conn.Close()

	chatId, err := service.GetIdFromVars(r)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	userId, ok := r.Context().Value("id").(int)
	if !ok {
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	sessionID := fmt.Sprintf("%s-%d", conn.RemoteAddr().String(), time.Now().UnixNano())
	log.Println("New chat client:", conn.RemoteAddr().String(), "session:", sessionID)

	var room *models.Room

	key := models.CallsKey{
		Type: models.TypePrivate,
		RoomID: chatId,
	}

	currentPeerState := service.CreateCurPS(conn, sessionID)

	calls.ReadWS(chatId, userId, rds, rds.Calls, room, key, currentPeerState, &roomMutex)

	if chatId != 0 {
		roomMutex.Lock()
		if r := rds.Calls[key]; r != nil {
			r.RemovePeer(sessionID)
			if len(r.Peers) == 0 {
				delete(rds.Calls, key)
				log.Printf("Chat Room %d deleted (empty)", chatId)
				currentPeerState.PeerConnection.Close()
			}
		}
		roomMutex.Unlock()
	}
	log.Printf("WS chat handler finished for %s", sessionID)
}

func HandleGroupCallConnection(w http.ResponseWriter, r *http.Request, rds *cache.Manager) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	defer conn.Close()

	chatId, err := service.GetIdFromVars(r)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	userId, ok := r.Context().Value("id").(int)
	if !ok {
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	sessionID := fmt.Sprintf("%s-%d", conn.RemoteAddr().String(), time.Now().UnixNano())
	log.Println("New group client:", conn.RemoteAddr().String(), "session:", sessionID)

	var room *models.Room

	key := models.CallsKey{
		Type: models.TypeGroup,
		RoomID: chatId,
	}

	currentPeerState := service.CreateCurPS(conn, sessionID)

	calls.ReadWS(chatId, userId, rds, rds.Calls, room, key, currentPeerState, &roomMutex)

	if chatId != 0 {
		roomMutex.Lock()
		if r := rds.Calls[key]; r != nil {
			r.RemovePeer(sessionID)
			if len(r.Peers) == 0 {
				delete(rds.Calls, key)
				log.Printf("Group Room %d deleted (empty)", chatId)
				currentPeerState.PeerConnection.Close()
			}
		}
		roomMutex.Unlock()
	}
	log.Printf("WS group handler finished for %s", sessionID)
}
