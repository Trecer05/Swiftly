package task_tracker

import (
	redis "github.com/Trecer05/Swiftly/internal/repository/cache/cloud"
	mgr "github.com/Trecer05/Swiftly/internal/repository/postgres/cloud"
	handlers "github.com/Trecer05/Swiftly/internal/handler/cloud"

	"github.com/gorilla/mux"
)

func NewCloudRouter(manager *mgr.Manager, rds *redis.WebSocketManager) *mux.Router {
	r := mux.NewRouter()

	handlers.InitCloudRoutes(r, manager, rds)

	return r
}
