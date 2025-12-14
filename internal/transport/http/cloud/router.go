package task_tracker

import (
	mgr "github.com/Trecer05/Swiftly/internal/repository/postgres/cloud"
	handlers "github.com/Trecer05/Swiftly/internal/handler/cloud"

	"github.com/gorilla/mux"
)

func NewCloudRouter(manager *mgr.Manager) *mux.Router {
	r := mux.NewRouter()

	handlers.InitCloudRoutes(r, manager)

	return r
}
