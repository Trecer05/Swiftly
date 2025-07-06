package chat

import (
	"log"
	"net/http"
	"strconv"

	errors "github.com/Trecer05/Swiftly/internal/errors/auth"
	middleware "github.com/Trecer05/Swiftly/internal/handler"
	redis "github.com/Trecer05/Swiftly/internal/repository/cache/chat"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/chat"
	serviceHttp "github.com/Trecer05/Swiftly/internal/transport/http"
	models "github.com/Trecer05/Swiftly/internal/model/chat"
	wsChat "github.com/Trecer05/Swiftly/internal/transport/websocket/chat"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// TODO: заменить на:
// CheckOrigin: func(r *http.Request) bool {
// 	origin := r.Header.Get("Origin")
// 	return origin == "https://домен"
// }
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func InitChatRoutes(r *mux.Router, mgr *manager.Manager, redis *redis.Manager) {
	r.Use(middleware.AuthMiddleware())

	r.HandleFunc("/chat/{id}", func(w http.ResponseWriter, r *http.Request) {
		ChatHandler(w, r, redis)
	})
}

func ChatHandler(w http.ResponseWriter, r *http.Request, rds *redis.Manager) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		serviceHttp.NewHeaderBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	chatId, err := strconv.Atoi(vars["id"])
	if err != nil {
		serviceHttp.NewHeaderBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	userId, ok := r.Context().Value("id").(int)
	if !ok {
		serviceHttp.NewHeaderBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	rds.AddUser(userId, chatId, conn)
	log.Println("User connected to chat", chatId, userId)
	msgCh := make(chan models.Message)

	go rds.ListenPubSub(chatId, msgCh)

	defer func() {
		rds.RemoveUser(userId, chatId)
	}()

	go rds.SendLocalMessage(userId, chatId, msgCh)
	wsChat.ReadMessage(chatId, conn, rds)

	close(msgCh)
}
