package cloud

import (
	"net/http"
	"time"

	middleware "github.com/Trecer05/Swiftly/internal/handler"
	redis "github.com/Trecer05/Swiftly/internal/repository/cache/cloud"
	kafkaManager "github.com/Trecer05/Swiftly/internal/repository/kafka/cloud"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/cloud"

	"github.com/gorilla/mux"
)

func InitCloudRoutes(r *mux.Router, manager *manager.Manager, rds *redis.WebSocketManager, kafkaManager *kafkaManager.KafkaManager) {
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)

	apiSecure := r.PathPrefix("/api/v1/cloud").Subrouter()
	apiSecure.Use(middleware.AuthMiddleware())
	apiSecure.Use(middleware.RateLimitMiddleware(rateLimiter))

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/team/{id}/dashboard", func(w http.ResponseWriter, r *http.Request) {
		DashboardWSHandler(w, r, manager, rds)
	})

	apiSecure.HandleFunc("/team/{id}/dashboard", func(w http.ResponseWriter, r *http.Request) {
		DashboardWSHandler(w, r, manager, rds)
	})

	apiSecure.HandleFunc("/team/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		GetTeamFilesAndFoldersHandler(w, r, manager, kafkaManager)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		GetUserFilesAndFoldersHandler(w, r, manager)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/team/{id:[0-9]+}/file", func(w http.ResponseWriter, r *http.Request) {
		CreateTeamFileHandler(w, r, manager, kafkaManager, rds)
	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/team/{id:[0-9]+}/file/{file_id}", func(w http.ResponseWriter, r *http.Request) {
		GetTeamFileByIDHandler(w, r, manager, kafkaManager)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/user/file", func(w http.ResponseWriter, r *http.Request) {
		CreateUserFileHandler(w, r, manager)
	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/user/file/{file_id}", func(w http.ResponseWriter, r *http.Request) {
		GetUserFileByIDHandler(w, r, manager)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/team/{id:[0-9]+}/file/{file_id}", func(w http.ResponseWriter, r *http.Request) {
		UpdateTeamFileByIDHandler(w, r, manager, kafkaManager, rds)
	}).Methods(http.MethodPut)

	apiSecure.HandleFunc("/user/file/{file_id}", func(w http.ResponseWriter, r *http.Request) {
		UpdateUserFileByIDHandler(w, r, manager)
	}).Methods(http.MethodPut)

	apiSecure.HandleFunc("/team/{id:[0-9]+}/file/{file_id}", func(w http.ResponseWriter, r *http.Request) {
		DeleteTeamFileByIDHandler(w, r, manager, kafkaManager, rds)
	}).Methods(http.MethodDelete)

	apiSecure.HandleFunc("/user/file/{file_id}", func(w http.ResponseWriter, r *http.Request) {
		DeleteUserFileByIDHandler(w, r, manager)
	}).Methods(http.MethodDelete)

	apiSecure.HandleFunc("/team/{id:[0-9]+}/file/{file_id}/download", func(w http.ResponseWriter, r *http.Request) {
		DownloadTeamFileByIDHandler(w, r, manager, kafkaManager)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/user/file/{file_id}/download", func(w http.ResponseWriter, r *http.Request) {
		DownloadUserFileByIDHandler(w, r, manager)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/team/{id:[0-9]+}/file/{file_id}/name", func(w http.ResponseWriter, r *http.Request) {
		UpdateTeamFileNameByIDHandler(w, r, manager, kafkaManager, rds)
	}).Methods(http.MethodPut)

	apiSecure.HandleFunc("/user/file/{file_id}/name", func(w http.ResponseWriter, r *http.Request) {
		UpdateUserFileNameByIDHandler(w, r, manager)
	}).Methods(http.MethodPut)

	apiSecure.HandleFunc("/team/{id:[0-9]+}/file/{file_id}/share", func(w http.ResponseWriter, r *http.Request) {
		ShareTeamFileByIDHandler(w, r, manager, kafkaManager)
	}).Methods(http.MethodPatch)

	apiSecure.HandleFunc("/user/file/{file_id}/share", func(w http.ResponseWriter, r *http.Request) {
		ShareUserFileByIDHandler(w, r, manager)
	}).Methods(http.MethodPatch)

	apiSecure.HandleFunc("/team/{id:[0-9]+}/folder", func(w http.ResponseWriter, r *http.Request) {
		CreateTeamFolderHandler(w, r, manager, kafkaManager, rds)
	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/user/folder", func(w http.ResponseWriter, r *http.Request) {
		CreateUserFolderHandler(w, r, manager)
	}).Methods(http.MethodPost)

	apiSecure.HandleFunc("/team/{id:[0-9]+}/folder/{folder_id}", func(w http.ResponseWriter, r *http.Request) {
		GetTeamFolderFilesByIDHandler(w, r, manager, kafkaManager)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/user/folder/{folder_id}", func(w http.ResponseWriter, r *http.Request) {
		GetUserFolderFilesByIDHandler(w, r, manager)
	}).Methods(http.MethodGet)

	apiSecure.HandleFunc("/team/{id:[0-9]+}/folder/{folder_id}", func(w http.ResponseWriter, r *http.Request) {
		DeleteTeamFolderByIDHandler(w, r, manager, kafkaManager, rds)
	}).Methods(http.MethodDelete)

	apiSecure.HandleFunc("/user/folder/{folder_id}", func(w http.ResponseWriter, r *http.Request) {
		DeleteUserFolderByIDHandler(w, r, manager)
	}).Methods(http.MethodDelete)

	// apiSecure.HandleFunc("/user/folder/{folder_id}/share", func(w http.ResponseWriter, r *http.Request) {
	// 	ShareUserFolderByIDHandler(w, r, manager)
	// }).Methods(http.MethodPatch)

	// apiSecure.HandleFunc("/team/{id:[0-9]+}/folder/{folder_id}/share", func(w http.ResponseWriter, r *http.Request) {
	// 	ShareTeamFolderByIDHandler(w, r, manager)
	// }).Methods(http.MethodPatch)

	apiSecure.HandleFunc("/user/folder/{folder_id}/name", func(w http.ResponseWriter, r *http.Request) {
		UpdateUserFolderNameByIDHandler(w, r, manager)
	}).Methods(http.MethodPut)

	apiSecure.HandleFunc("/team/{id:[0-9]+}/folder/{folder_id}/name", func(w http.ResponseWriter, r *http.Request) {
		UpdateTeamFolderNameByIDHandler(w, r, manager, kafkaManager, rds)
	}).Methods(http.MethodPut)

	apiSecure.HandleFunc("/user/folder/{folder_id}/move", func(w http.ResponseWriter, r *http.Request) {
		MoveUserFolderByIDHandler(w, r, manager)
	}).Methods(http.MethodPut)

	apiSecure.HandleFunc("/team/{id:[0-9]+}/folder/{folder_id}/move", func(w http.ResponseWriter, r *http.Request) {
		MoveTeamFolderByIDHandler(w, r, manager, kafkaManager, rds)
	}).Methods(http.MethodPut)

	apiSecure.HandleFunc("/user/file/{file_id}/move", func(w http.ResponseWriter, r *http.Request) {
		MoveUserFileByIDHandler(w, r, manager)
	}).Methods(http.MethodPut)

	apiSecure.HandleFunc("/team/{id:[0-9]+}/file/{file_id}/move", func(w http.ResponseWriter, r *http.Request) {
		MoveTeamFileByIDHandler(w, r, manager, kafkaManager, rds)
	}).Methods(http.MethodPut)

	apiSecure.HandleFunc("/sharedfiles/file/{file_id}", func(w http.ResponseWriter, r *http.Request) {
		GetSharedFileHandler(w, r, manager)
	}).Methods(http.MethodGet)

	// apiSecure.HandleFunc("/sharedfiles/folder/{folder_id}", func(w http.ResponseWriter, r *http.Request) {
	// 	GetSharedFolderHandler(w, r, manager)
	// }).Methods(http.MethodGet)
}
