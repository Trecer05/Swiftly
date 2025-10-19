package chat

import (
	"net/http"
	"strconv"

	errors "github.com/Trecer05/Swiftly/internal/errors/auth"
	chatErrors "github.com/Trecer05/Swiftly/internal/errors/chat"
	models "github.com/Trecer05/Swiftly/internal/model/chat"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/chat"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func ValidateGroupOwner(groupId int, r *http.Request, mgr *manager.Manager) (int, error) {
	userId, ok := r.Context().Value("id").(int)
	if !ok {
		return http.StatusUnauthorized, errors.ErrUnauthorized
	}

	ok, err := mgr.ValidateOwnerId(groupId, userId)
	switch {
	case err == chatErrors.ErrNoGroupFound:
		return http.StatusNotFound, err
	case err != nil:
		return http.StatusInternalServerError, err
	}

	if !ok {
		return http.StatusForbidden, errors.ErrGroupForbidden
	}
	return http.StatusOK, nil
}

func ValidateLimitOffset(r *http.Request) (int, int, error) {
	var limit, offset int
	strLimit, strOffset := r.URL.Query()["limit"], r.URL.Query()["offset"]
	if strLimit != nil {
		if l, err := strconv.Atoi(strLimit[0]); err == nil {
			limit = l
		} else {
			return 0, 0, chatErrors.ErrInvalidLimit
		}
	}
	if strOffset != nil {
		if o, err := strconv.Atoi(strOffset[0]); err == nil {
			offset = o
		} else {
			return 0, 0, chatErrors.ErrInvalidOffset
		}
	}

	return limit, offset, nil
}

func GetIdFromVars(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	return strconv.Atoi(id)
}

func CreateCurPS(conn *websocket.Conn, sessionID string) *models.PeerState {
	return &models.PeerState{
		WS:        conn,
		SessionID: sessionID,
		Tracks:    make(map[string]*models.Track),
	}
}

func NewRoom() *models.Room {
	return &models.Room{
		Peers:     make(map[string]*models.PeerState),
		Published: make(map[string]*models.PublishedTrack),
	}
}
