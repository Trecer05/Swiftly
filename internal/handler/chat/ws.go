package chat

import (
	"net/http"
	"strconv"

	errors "github.com/Trecer05/Swiftly/internal/errors/auth"
	middleware "github.com/Trecer05/Swiftly/internal/handler"
	models "github.com/Trecer05/Swiftly/internal/model/chat"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/chat"
	serviceHttp "github.com/Trecer05/Swiftly/internal/transport/http"
	wsChat "github.com/Trecer05/Swiftly/internal/transport/websocket/chat"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func InitChatRoutes(r *mux.Router, mgr *manager.Manager) {
	r.Use(middleware.AuthMiddleware())

	r.HandleFunc("/chat/{id}", func(w http.ResponseWriter, r *http.Request) {

	})
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	intId, err := strconv.Atoi(vars["id"])
	if err != nil {
		serviceHttp.NewHeaderBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	// TODO: сделать функцию для получения комнаты
	room := GetUserRoom(intId)
	chatConnection(room, w, r)
}

func chatConnection(room models.ChatRoom, w http.ResponseWriter, r *http.Request) {
	var id int
	var ok bool

	if id, ok = r.Context().Value("id").(int); !ok {
		serviceHttp.NewHeaderBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		serviceHttp.NewHeaderBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	user := &models.Client{
		ID: id,
		Conn: conn,
		Send: make(chan models.Message, 256),
	}

	room.OnChat <- user

	go wsChat.ReadMessage(user, &room)
	go wsChat.WriteMessage(user)
}
