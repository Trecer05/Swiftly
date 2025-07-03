package chat

import (
	mgr "github.com/Trecer05/Swiftly/internal/repository/postgres/chat"
	handlers "github.com/Trecer05/Swiftly/internal/handler/chat"

	"github.com/gorilla/mux"
)

func NewChatRouter(manager *mgr.Manager) *mux.Router {
	r := mux.NewRouter()

	handlers.InitChatRoutes(r, manager)

	return r
}
