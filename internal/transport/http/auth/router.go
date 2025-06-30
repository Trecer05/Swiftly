package auth

import (
	handlers "github.com/Trecer05/Swiftly/internal/handler/auth"
	mgr "github.com/Trecer05/Swiftly/internal/repository/postgres/auth"

	"github.com/gorilla/mux"
)

func NewAuthRouter(manager *mgr.Manager) *mux.Router {
	r := mux.NewRouter()

	handlers.InitAuthRoutes(r, manager)

	return r
}
