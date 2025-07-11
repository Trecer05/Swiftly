package chat

import (
	"encoding/json"
	"net/http"

	models "github.com/Trecer05/Swiftly/internal/model/chat"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/chat"
	serviceHttp "github.com/Trecer05/Swiftly/internal/transport/http"
)

func CreateGroupHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	var group models.GroupCreate

	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		serviceHttp.NewHeaderBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	id, err := mgr.CreateGroup(group)
	if err != nil {
		serviceHttp.NewHeaderBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"group_id" : id,
	})
}