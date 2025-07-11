package chat

import (
	"encoding/json"
	"net/http"

	errors "github.com/Trecer05/Swiftly/internal/errors/auth"
	chatErrors "github.com/Trecer05/Swiftly/internal/errors/chat"
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

	if group.Name == "" || group.OwnerID == 0 {
		serviceHttp.NewHeaderBody(w, "application/json", chatErrors.ErrInvalidGroupData, http.StatusBadRequest)
		return
	}

	ownerID, ok := r.Context().Value("id").(int)
	if !ok {
		serviceHttp.NewHeaderBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}
	group.OwnerID = ownerID

	id, err := mgr.CreateGroup(group)
	if err != nil {
		serviceHttp.NewHeaderBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"group" : map[string]interface{}{
			"id" : id,
			"name" : group.Name,
		},
	})
}