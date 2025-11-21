package chat

import (
	"context"
	logger "github.com/Trecer05/Swiftly/internal/config/logger"
	"net/http"
	"os"
	"strconv"
	"time"

	errors "github.com/Trecer05/Swiftly/internal/errors/auth"
	chatErrors "github.com/Trecer05/Swiftly/internal/errors/chat"
	middleware "github.com/Trecer05/Swiftly/internal/handler"
	models "github.com/Trecer05/Swiftly/internal/model/chat"
	redis "github.com/Trecer05/Swiftly/internal/repository/cache/chat"
	kafka "github.com/Trecer05/Swiftly/internal/repository/kafka/chat"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/chat"
	serviceChat "github.com/Trecer05/Swiftly/internal/service/chat"
	serviceHttp "github.com/Trecer05/Swiftly/internal/transport/http"
	wsChat "github.com/Trecer05/Swiftly/internal/transport/websocket/chat"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
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

var ctx = context.Background()

func InitChatRoutes(router *mux.Router, mgr *manager.Manager, redis *redis.Manager) {
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)

	apiSecure := router.PathPrefix("/api/v1").Subrouter()
	apiSecure.Use(middleware.AuthMiddleware())
	apiSecure.Use(middleware.RateLimitMiddleware(rateLimiter))

	kafkaChangeManager := kafka.NewKafkaManager([]string{os.Getenv("KAFKA_ADDRESS")}, "profile", "user-change-group")
	kafkaChangeManagerTasks := kafka.NewKafkaManager([]string{os.Getenv("KAFKA_ADDRESS")}, "team", "team-user-tasks")

	go kafkaChangeManager.ReadAuthMessages(ctx)
	go kafkaChangeManagerTasks.ReadTaskMessages(ctx)
	
	defer kafkaChangeManager.Close()
	defer kafkaChangeManagerTasks.Close()
	
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		EditProfileHandler(w, r, mgr)
	}).Methods(http.MethodPut)

	apiSecure.HandleFunc("/user/avatar/urls", func(w http.ResponseWriter, r *http.Request) {
		GetProfileAvatarUrlsHandler(w, r)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/user/avatar/{url}", func(w http.ResponseWriter, r *http.Request) {
		GetProfileAvatarHandler(w, r)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/user/{id}/avatar/urls", func(w http.ResponseWriter, r *http.Request) {
		GetUserAvatarUrlsHandler(w, r)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/user/{id}/avatar/{url}", func(w http.ResponseWriter, r *http.Request) {
		GetUserAvatarHandler(w, r)
	}).Methods(http.MethodGet)
	
	apiSecure.HandleFunc("/user/avatar", func(w http.ResponseWriter, r *http.Request) {
		UploadProfileAvatarHandler(w, r)
	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/user/avatar/{url}", func(w http.ResponseWriter, r *http.Request) {
		DeleteProfileAvatarHandler(w, r)
	}).Methods(http.MethodDelete)

	apiSecure.HandleFunc("/user/email", func(w http.ResponseWriter, r *http.Request) {
		EditUserEmailHandler(w, r, mgr, kafkaChangeManager, ctx)
	}).Methods(http.MethodPut)

	apiSecure.HandleFunc("/user/phone", func(w http.ResponseWriter, r *http.Request) {
		EditUserPhoneHandler(w, r, mgr, kafkaChangeManager, ctx)
	}).Methods(http.MethodPut)

	apiSecure.HandleFunc("/user/password", func(w http.ResponseWriter, r *http.Request) {
		EditUserPasswordHandler(w, r, mgr, kafkaChangeManager, ctx)
	}).Methods(http.MethodPut)

	apiSecure.HandleFunc("/chat/{id}", func(w http.ResponseWriter, r *http.Request) {
		ChatHandler(w, r, redis, mgr)
	})

	apiSecure.HandleFunc("/group/{id}", func(w http.ResponseWriter, r *http.Request) {
		GroupHandler(w, r, redis, mgr)
	})

	apiSecure.HandleFunc("/main", func(w http.ResponseWriter, r *http.Request) {
		MainConnectionHandler(w, r, redis, mgr)
	})

	apiSecure.HandleFunc("/group", func(w http.ResponseWriter, r *http.Request) {
		CreateGroupHandler(w, r, mgr)
	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/group/{id}", func(w http.ResponseWriter, r *http.Request) {
		UpdateGroupHandler(w, r, mgr)
	}).Methods(http.MethodPut)

	apiSecure.HandleFunc("/group/{id}/avatar", func(w http.ResponseWriter, r *http.Request) {
		GetGroupAvatarUrlHandler(w, r)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/group/{id}/avatar/{url}", func(w http.ResponseWriter, r *http.Request) {
		GetGroupAvatarHandler(w, r)
	}).Methods(http.MethodGet)
	
	apiSecure.HandleFunc("/group/{id}/avatar", func(w http.ResponseWriter, r *http.Request) {
		UploadGroupAvatarHandler(w, r)
	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/group/{id}/avatar/{url}", func(w http.ResponseWriter, r *http.Request) {
		DeleteGroupAvatarHandler(w, r)
	}).Methods(http.MethodDelete)

	apiSecure.HandleFunc("/group/{id}", func(w http.ResponseWriter, r *http.Request) {
		DeleteGroupHandler(w, r, mgr)
	}).Methods(http.MethodDelete)

	apiSecure.HandleFunc("/chat/{id}", func(w http.ResponseWriter, r *http.Request) {
		DeleteChatHandler(w, r, mgr)
	}).Methods(http.MethodDelete)

	apiSecure.HandleFunc("/group/{id}/users", func(w http.ResponseWriter, r *http.Request) {
		GroupUsersHandler(w, r, mgr)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/chats", func(w http.ResponseWriter, r *http.Request) {
		UserChatsInfoHandler(w, r, mgr)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		UserInfoHandler(w, r, mgr)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		CreateUserHandler(w, r, mgr)
	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		AnotherUserInfoHandler(w, r, mgr)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/chat/{id}", func(w http.ResponseWriter, r *http.Request) {
		CreateChatHandler(w, r, mgr)
	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/chat/{id}/info", func(w http.ResponseWriter, r *http.Request) {
		ChatInfoHandler(w, r, mgr)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/chat/{id}/messages", func(w http.ResponseWriter, r *http.Request) {
		ChatMessagesHandler(w, r, mgr)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/group/{id}/messages", func(w http.ResponseWriter, r *http.Request) {
		GroupMessagesHandler(w, r, mgr)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/group/{id}/info", func(w http.ResponseWriter, r *http.Request) {
		GroupInfoHandler(w, r, mgr)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/group/{id}/add", func(w http.ResponseWriter, r *http.Request) {
		AddUserToGroupHandler(w, r, mgr)
	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/group/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		DeleteUserFromGroupHandler(w, r, mgr)
	}).Methods(http.MethodDelete)

	apiSecure.HandleFunc("/group/{id}/leave", func(w http.ResponseWriter, r *http.Request) {
		ExitGroupHandler(w, r, mgr)
	}).Methods(http.MethodDelete)

	apiSecure.HandleFunc("/group/{id}/img", func(w http.ResponseWriter, r *http.Request) {
		UploadImgHandler(w, r, mgr, models.TypeGroup)
	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/group/{id}/video", func(w http.ResponseWriter, r *http.Request) {
		UploadVideoHandler(w, r, mgr, models.TypeGroup)
	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/group/{id}/file", func(w http.ResponseWriter, r *http.Request) {
		UploadFileHandler(w, r, mgr, models.TypeGroup)
	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/group/{id}/audio", func(w http.ResponseWriter, r *http.Request) {
		UploadAudioHandler(w, r, mgr, models.TypeGroup)
	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/group/{id}/imgvid", func(w http.ResponseWriter, r *http.Request) {
		UploadImgVideoHandler(w, r, mgr, models.TypeGroup)
	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/group/{id}/files", func(w http.ResponseWriter, r *http.Request) {
		GetFilesHandler(w, r, mgr, models.TypeGroup)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/group/{id}/img/{url}", func(w http.ResponseWriter, r *http.Request) {
		GetImgHandler(w, r, mgr, models.TypeGroup)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/group/{id}/video/{url}", func(w http.ResponseWriter, r *http.Request) {
		GetVideoHandler(w, r, mgr, models.TypeGroup)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/group/{id}/audio/{url}", func(w http.ResponseWriter, r *http.Request) {
		GetAudioHandler(w, r, mgr, models.TypeGroup)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/group/{id}/file/{url}", func(w http.ResponseWriter, r *http.Request) {
		GetFileHandler(w, r, mgr, models.TypeGroup)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/chat/{id}/img", func(w http.ResponseWriter, r *http.Request) {
		UploadImgHandler(w, r, mgr, models.TypePrivate)
	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/chat/{id}/video", func(w http.ResponseWriter, r *http.Request) {
		UploadVideoHandler(w, r, mgr, models.TypePrivate)
	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/chat/{id}/file", func(w http.ResponseWriter, r *http.Request) {
		UploadFileHandler(w, r, mgr, models.TypePrivate)
	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/chat/{id}/audio", func(w http.ResponseWriter, r *http.Request) {
		UploadAudioHandler(w, r, mgr, models.TypePrivate)
	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/chat/{id}/imgvid", func(w http.ResponseWriter, r *http.Request) {
		UploadImgVideoHandler(w, r, mgr, models.TypePrivate)
	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/chat/{id}/files", func(w http.ResponseWriter, r *http.Request) {
		GetFilesHandler(w, r, mgr, models.TypePrivate)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/chat/{id}/img/{url}", func(w http.ResponseWriter, r *http.Request) {
		GetImgHandler(w, r, mgr, models.TypePrivate)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/chat/{id}/video/{url}", func(w http.ResponseWriter, r *http.Request) {
		GetVideoHandler(w, r, mgr, models.TypePrivate)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/chat/{id}/audio/{url}", func(w http.ResponseWriter, r *http.Request) {
		GetAudioHandler(w, r, mgr, models.TypePrivate)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/chat/{id}/file/{url}", func(w http.ResponseWriter, r *http.Request) {
		GetFileHandler(w, r, mgr, models.TypePrivate)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/chat/{id}/call", func(w http.ResponseWriter, r *http.Request) {
		HandleChatCallConnection(w, r, redis)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/group/{id}/call", func(w http.ResponseWriter, r *http.Request) {
		HandleGroupCallConnection(w, r, redis)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/group/{id}/message/audio", func(w http.ResponseWriter, r *http.Request) {
		SaveAudioMessageHandler(w, r, models.TypeGroup)
	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/group/{id}/message/audio/{url}", func(w http.ResponseWriter, r *http.Request) {
		GetAudioMessageHandler(w, r, models.TypeGroup)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/group/{id}/message/video", func(w http.ResponseWriter, r *http.Request) {
		SaveVideoMessageHandler(w, r, models.TypeGroup)
	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/group/{id}/message/video/{url}", func(w http.ResponseWriter, r *http.Request) {
		GetVideoMessageHandler(w, r, models.TypeGroup)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/chat/{id}/message/audio", func(w http.ResponseWriter, r *http.Request) {
		SaveAudioMessageHandler(w, r, models.TypePrivate)
	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/chat/{id}/message/audio/{url}", func(w http.ResponseWriter, r *http.Request) {
		GetAudioMessageHandler(w, r, models.TypePrivate)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/chat/{id}/message/video", func(w http.ResponseWriter, r *http.Request) {
		SaveVideoMessageHandler(w, r, models.TypePrivate)
	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/chat/{id}/message/video/{url}", func(w http.ResponseWriter, r *http.Request) {
		GetVideoMessageHandler(w, r, models.TypePrivate)
	}).Methods(http.MethodGet)
	
	// Модуль команды
	apiSecure.HandleFunc("/team/{id}/dashboard", func(w http.ResponseWriter, r *http.Request) {
		GetTeamDashboardHandler(w, r, mgr, kafkaChangeManagerTasks, ctx)
	}).Methods(http.MethodGet)
	
	apiSecure.HandleFunc("/team/{team_id}/user/{username}/add", func(w http.ResponseWriter, r *http.Request) {
		AddUserToTeamByUsernameHandler(w, r, mgr, redis)
	}).Methods(http.MethodPost)
	
	apiSecure.HandleFunc("/team/{team_id}/user/{user_id}/remove", func(w http.ResponseWriter, r *http.Request) {
		RemoveUserFromTeamByIDHandler(w, r, mgr, redis)
	}).Methods(http.MethodDelete)
	
	apiSecure.HandleFunc("/team/{team_id}/user/{user_id}/role", func(w http.ResponseWriter, r *http.Request) {
		UpdateUserRoleHandler(w, r, mgr, redis)
	}).Methods(http.MethodPut)
	
	apiSecure.HandleFunc("/team", func(w http.ResponseWriter, r *http.Request) {
		CreateTeamHandler(w, r, mgr)
	}).Methods(http.MethodPost)
}

func ChatHandler(w http.ResponseWriter, r *http.Request, rds *redis.Manager, mgr *manager.Manager) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	chatId, err := strconv.Atoi(vars["id"])
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	limit, offset, err := serviceChat.ValidateLimitOffset(r)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	userId, ok := r.Context().Value("id").(int)
	if !ok {
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	chatType := models.TypePrivate

	rds.AddUser(userId, chatId, conn, chatType, mgr)
	logger.Logger.Println("User connected to chat", chatId, userId)
	msgCh := make(chan models.Message)
	statusCh := make(chan models.Status)
	notifCh := make(chan models.Notifications)

	if err := wsChat.SendChatHistory(conn, mgr, chatId, limit, offset, chatType); err != nil {
		if err == chatErrors.ErrNoMessages {
			conn.WriteJSON(map[string]interface{}{
				"type":     "history",
				"messages": []models.Message{},
				"error":    err.Error(),
			})
		} else {
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
			return
		}
	}

	go rds.ListenPubSub(chatId, msgCh, statusCh, notifCh, chatType)

	defer func() {
		rds.RemoveUser(userId, chatId, chatType, mgr)
	}()

	go rds.SendLocalMessage(userId, chatId, msgCh, chatType)
	go rds.SendLocalStatus(userId, chatId, statusCh, chatType)
	go rds.SendLocalNotification(userId, chatId, notifCh, chatType)
	wsChat.ReadMessage(chatId, conn, rds, mgr, chatType)

	conn.Close()
	close(notifCh)
	close(msgCh)
	close(statusCh)
}

func GroupHandler(w http.ResponseWriter, r *http.Request, rds *redis.Manager, mgr *manager.Manager) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	groupId, err := strconv.Atoi(vars["id"])
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	limit, offset, err := serviceChat.ValidateLimitOffset(r)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	userId, ok := r.Context().Value("id").(int)
	if !ok {
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	chatType := models.TypeGroup

	rds.AddUser(userId, groupId, conn, chatType, mgr)
	logger.Logger.Println("User connected to chat", groupId, userId)
	msgCh := make(chan models.Message)
	statusCh := make(chan models.Status)
	notifCh := make(chan models.Notifications)

	if err := wsChat.SendChatHistory(conn, mgr, groupId, limit, offset, chatType); err != nil {
		if err == chatErrors.ErrNoMessages {
			conn.WriteJSON(map[string]interface{}{
				"type":     "history",
				"messages": []models.Message{},
				"error":    err.Error(),
			})
		} else {
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
			return
		}
	}

	go rds.ListenPubSub(groupId, msgCh, statusCh, notifCh, chatType)

	defer func() {
		rds.RemoveUser(userId, groupId, chatType, mgr)
	}()

	go rds.SendLocalMessage(userId, groupId, msgCh, chatType)
	go rds.SendLocalStatus(userId, groupId, statusCh, chatType)
	go rds.SendLocalNotification(userId, groupId, notifCh, chatType)
	wsChat.ReadMessage(groupId, conn, rds, mgr, chatType)

	conn.Close()
	close(notifCh)
	close(msgCh)
	close(statusCh)
}

func MainConnectionHandler(w http.ResponseWriter, r *http.Request, rds *redis.Manager, mgr *manager.Manager) {
	var err error

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	userId, ok := r.Context().Value("id").(int)
	if !ok {
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	msgCh := make(chan models.Message)
	statusCh := make(chan models.Status)
	notifCh := make(chan models.Notifications)

	chats, err := mgr.GetUserRooms(userId, 0, 0)
	switch {
	case err == chatErrors.ErrNoRooms:
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
		return
	case err != nil:
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	for _, room := range chats.Rooms {
		go rds.ListenPubSub(room.ID, msgCh, statusCh, notifCh, room.Type)

		if room.LastMessage != nil {
			room.LastMessage.Type = models.LastMessage
			room.LastMessage.ChatID = room.ID
			msgCh <- *room.LastMessage
		}
	}

	rds.WG.Add(3)
	go func() {
		defer rds.WG.Done()
		wsChat.SendAllUserMessages(conn, msgCh)
	}()
	go func() {
		defer rds.WG.Done()
		wsChat.SendAllUserStatuses(conn, statusCh)
	}()
	go func() {
		defer rds.WG.Done()
		wsChat.SendAllUserNotifications(conn, notifCh)
	}()
	rds.WG.Wait()

	conn.Close()
	close(notifCh)
	close(msgCh)
	close(statusCh)
}
