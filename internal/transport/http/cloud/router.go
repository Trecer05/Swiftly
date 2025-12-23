package task_tracker

import (
	handlers "github.com/Trecer05/Swiftly/internal/handler/cloud"
	redis "github.com/Trecer05/Swiftly/internal/repository/cache/cloud"
	"github.com/Trecer05/Swiftly/internal/repository/kafka/cloud"
	mgr "github.com/Trecer05/Swiftly/internal/repository/postgres/cloud"

	"github.com/gorilla/mux"
)

func NewCloudRouter(manager *mgr.Manager, rds *redis.WebSocketManager, kafka *cloud.KafkaManager) *mux.Router {
	r := mux.NewRouter()

	handlers.InitCloudRoutes(r, manager, rds, kafka)

	return r
}
