package cloud

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/Trecer05/Swiftly/internal/config/logger"
	wsTracker "github.com/Trecer05/Swiftly/internal/transport/websocket/cloud"
	authErrors "github.com/Trecer05/Swiftly/internal/errors/auth"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/cloud"
	redis "github.com/Trecer05/Swiftly/internal/repository/cache/cloud"
	serviceHttp "github.com/Trecer05/Swiftly/internal/transport/http"
	models "github.com/Trecer05/Swiftly/internal/model/cloud"
)

// TODO: заменить на:
//
//	CheckOrigin: func(r *http.Request) bool {
//		origin := r.Header.Get("Origin")
//		return origin == "https://домен"
//	}
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        // origin := r.Header.Get("Origin")
        // allowedOrigins := []string{
        //     "http://localhost:3000",
        //     "https://yourdomain.com",
        // }
        
        // for _, allowed := range allowedOrigins {
        //     if origin == allowed {
        //         return true
        //     }
        // }
        // return false

		return true
    },
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func DashboardWSHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, rds *redis.WebSocketManager) {
	teamID, err := strconv.Atoi(mux.Vars(r)["team_id"])
	if err != nil {
		logger.Logger.Error("Error getting team ID", authErrors.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", authErrors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}
	
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Logger.Error("Error upgrading connection", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	
	userId, ok := r.Context().Value("id").(int)
	if !ok {
		logger.Logger.Error("Error getting user ID", authErrors.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", authErrors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}
	
	rds.AddUser(userId, teamID, conn)
	logger.Logger.Println("User connected to dashboard", userId, teamID)
	msgCh := make(chan models.Envelope)
	
	go rds.ListenPubSub(teamID, msgCh)
	
	defer func() {
		rds.RemoveUser(userId, teamID)
		logger.Logger.Println("User disconnected from dashboard", userId, teamID)
	}()
	
	go rds.SendLocalMessage(userId, teamID, msgCh)
	wsTracker.ReadMessage(teamID, conn, rds, mgr)

	conn.Close()
	close(msgCh)
}