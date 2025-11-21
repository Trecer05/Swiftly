package task_tracker

import (
	mgr "github.com/Trecer05/Swiftly/internal/repository/postgres/task_tracker"
	// rds "github.com/Trecer05/Swiftly/internal/repository/cache/task_tracker"
	
	"github.com/gorilla/mux"
	"net/http"
)

func InitTaskRoutes(r *mux.Router, manager *mgr.Manager) {
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)
}
