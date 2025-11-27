package task_tracker

import (
	"net/http"
	"strconv"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/Trecer05/Swiftly/internal/config/logger"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/task_tracker"
	serviceHttp "github.com/Trecer05/Swiftly/internal/transport/http"
)

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		logger.Logger.Error("Error converting task ID to integer", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

    task, err := mgr.GetTaskByID(id)
	if err != nil {
		logger.Logger.Error("Error getting task", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"task": task,
	})
}
