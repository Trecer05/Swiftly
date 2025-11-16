package task_tracker

import (
	mgr "github.com/Trecer05/Swiftly/internal/repository/postgres/task_tracker"
	// redis "github.com/Trecer05/Swiftly/internal/repository/cache/task_tracker"
	handlers "github.com/Trecer05/Swiftly/internal/handler/task_tracker"

	"github.com/gorilla/mux"
)

func NewTaskRouter(manager *mgr.Manager) *mux.Router {
	r := mux.NewRouter()

	handlers.InitTaskRoutes(r, manager)

	return r
}
