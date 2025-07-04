package chat

import (
	mgr "github.com/Trecer05/Swiftly/internal/repository/postgres/chat"
	redis "github.com/Trecer05/Swiftly/internal/repository/cache/chat"
	handlers "github.com/Trecer05/Swiftly/internal/handler/chat"

	"github.com/gorilla/mux"
)

func NewChatRouter(manager *mgr.Manager, rds *redis.Manager) *mux.Router {
	r := mux.NewRouter()

	handlers.InitChatRoutes(r, manager, rds)

	return r
}
