package chat

import (
	"log"
	"net/http"
	"strconv"

	errors "github.com/Trecer05/Swiftly/internal/errors/auth"
	chatErrors "github.com/Trecer05/Swiftly/internal/errors/chat"
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
		ChatHandler(w, r, redis, mgr)
	})

	r.HandleFunc("/group/{id}", func(w http.ResponseWriter, r *http.Request) {
		GroupHandler(w, r, redis, mgr)
	})

	r.HandleFunc("/main", func(w http.ResponseWriter, r *http.Request) {
		MainConnectionHandler(w, r, redis, mgr)
	})

	r.HandleFunc("/group", func(w http.ResponseWriter, r *http.Request) {
		CreateGroupHandler(w, r, mgr)
	}).Methods(http.MethodPost)

	r.HandleFunc("/group/{id}", func(w http.ResponseWriter, r *http.Request) {
		DeleteGroupHandler(w, r, mgr)
	}).Methods(http.MethodDelete)
}

func ChatHandler(w http.ResponseWriter, r *http.Request, rds *redis.Manager, mgr *manager.Manager) {
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

	var limit, offset int
	strLimit, strOffset := r.URL.Query()["limit"], r.URL.Query()["offset"]
	if strLimit == nil {
		limit = 100
	} else {
		if limit, err = strconv.Atoi(strLimit[0]); err != nil {
			serviceHttp.NewHeaderBody(w, "application/json", err, http.StatusBadRequest)
			return
		}
	}
	if strOffset == nil {
		offset = 0
	} else {
		if offset, err = strconv.Atoi(strOffset[0]); err != nil {
			serviceHttp.NewHeaderBody(w, "application/json", err, http.StatusBadRequest)
			return
		}
	}

	userId, ok := r.Context().Value("id").(int)
	if !ok {
		serviceHttp.NewHeaderBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	chatType := models.TypePrivate

	rds.AddUser(userId, chatId, conn, chatType)
	log.Println("User connected to chat", chatId, userId)
	msgCh := make(chan models.Message)

	if err := wsChat.SendChatHistory(conn, mgr, chatId, limit, offset, chatType); err != nil {
		if err == chatErrors.ErrNoMessages {
			conn.WriteJSON(map[string]interface{}{
				"type": "history",
				"messages": []models.Message{},
				"error": err.Error(),
			})
		} else {
			serviceHttp.NewHeaderBody(w, "application/json", err, http.StatusInternalServerError)
			return
		}
	}

	go rds.ListenPubSub(chatId, msgCh, chatType)

	defer func() {
		rds.RemoveUser(userId, chatId, chatType)
	}()

	go rds.SendLocalMessage(userId, chatId, msgCh, chatType)
	wsChat.ReadMessage(chatId, conn, rds, mgr, chatType)

	close(msgCh)
}

func GroupHandler(w http.ResponseWriter, r *http.Request, rds *redis.Manager, mgr *manager.Manager) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		serviceHttp.NewHeaderBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	groupId, err := strconv.Atoi(vars["id"])
	if err != nil {
		serviceHttp.NewHeaderBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	var limit, offset int
	strLimit, strOffset := r.URL.Query()["limit"], r.URL.Query()["offset"]
	if strLimit == nil {
		limit = 100
	} else {
		if limit, err = strconv.Atoi(strLimit[0]); err != nil {
			serviceHttp.NewHeaderBody(w, "application/json", err, http.StatusBadRequest)
			return
		}
	}
	if strOffset == nil {
		offset = 0
	} else {
		if offset, err = strconv.Atoi(strOffset[0]); err != nil {
			serviceHttp.NewHeaderBody(w, "application/json", err, http.StatusBadRequest)
			return
		}
	}

	userId, ok := r.Context().Value("id").(int)
	if !ok {
		serviceHttp.NewHeaderBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	chatType := models.TypeGroup

	rds.AddUser(userId, groupId, conn, chatType)
	log.Println("User connected to chat", groupId, userId)
	msgCh := make(chan models.Message)

	if err := wsChat.SendChatHistory(conn, mgr, groupId, limit, offset, chatType); err != nil {
		if err == chatErrors.ErrNoMessages {
			conn.WriteJSON(map[string]interface{}{
				"type": "history",
				"messages": []models.Message{},
				"error": err.Error(),
			})
		} else {
			serviceHttp.NewHeaderBody(w, "application/json", err, http.StatusInternalServerError)
			return
		}
	}

	go rds.ListenPubSub(groupId, msgCh, chatType)

	defer func() {
		rds.RemoveUser(userId, groupId, chatType)
	}()

	go rds.SendLocalMessage(userId, groupId, msgCh, chatType)
	wsChat.ReadMessage(groupId, conn, rds, mgr, chatType)

	close(msgCh)
}

func MainConnectionHandler(w http.ResponseWriter, r *http.Request, rds *redis.Manager, mgr *manager.Manager) {
	var err error

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		serviceHttp.NewHeaderBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	userId, ok := r.Context().Value("id").(int)
	if !ok {
		serviceHttp.NewHeaderBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	msgCh := make(chan models.Message)

	chats, err := mgr.GetUserRooms(userId, 0, 0)
	switch {
		case err == chatErrors.ErrNoRooms:
			serviceHttp.NewHeaderBody(w, "application/json", err, http.StatusNotFound)
			return
		case err != nil:
			serviceHttp.NewHeaderBody(w, "application/json", err, http.StatusInternalServerError)
			return
	}

	for _, room := range chats.Rooms {
		go rds.ListenPubSub(room.ID, msgCh, room.Type)

		if room.LastMessage != nil {
			room.LastMessage.Type = models.LastMessage
			room.LastMessage.ChatID = room.ID
			msgCh <- *room.LastMessage
		}
	}

	go wsChat.SendAllUserMessages(conn, msgCh)

	conn.Close()
	close(msgCh)
}
