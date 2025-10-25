package chat

import (
	"encoding/json"
	"net/http"
	"strconv"

	errors "github.com/Trecer05/Swiftly/internal/errors/auth"
	chatErrors "github.com/Trecer05/Swiftly/internal/errors/chat"
	globalErrors "github.com/Trecer05/Swiftly/internal/errors/global"
	fileManager "github.com/Trecer05/Swiftly/internal/filemanager"
	models "github.com/Trecer05/Swiftly/internal/model/chat"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/chat"
	serviceChat "github.com/Trecer05/Swiftly/internal/service/chat"
	serviceHttp "github.com/Trecer05/Swiftly/internal/transport/http"

	"github.com/gorilla/mux"
)

func CreateGroupHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	var group models.GroupCreate

	r.ParseMultipartForm(10 << 20)
	jsonData := r.FormValue("json")
	if jsonData == ""{
		serviceHttp.NewErrorBody(w, "application/json", globalErrors.ErrNoJsonData, http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal([]byte(jsonData), &group); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	if group.Name == "" || group.OwnerID == 0 {
		serviceHttp.NewErrorBody(w, "application/json", chatErrors.ErrInvalidGroupData, http.StatusBadRequest)
		return
	}

	ownerID, ok := r.Context().Value("id").(int)
	if !ok {
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}
	group.OwnerID = ownerID

	id, err := mgr.CreateGroup(group)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	err = fileManager.CreateGroupMessagesFolder(id)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	if _, _, err := r.FormFile("photo"); err == nil {
        if err := fileManager.AddGroupPhoto(r, id); err != nil {
            serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
            return
        }
    } else if err != http.ErrMissingFile {
        serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
        return
    }

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"group": map[string]interface{}{
			"id":   id,
			"name": group.Name,
		},
	})
}

func DeleteGroupHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	vars := mux.Vars(r)
	groupId, err := strconv.Atoi(vars["id"])
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	status, err := serviceChat.ValidateGroupOwner(groupId, r, mgr)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, status)
		return
	}

	err = mgr.DeleteGroup(groupId)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	err = fileManager.DeleteGroupPhoto(groupId)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	err = fileManager.DeleteGroupFiles(groupId)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}

func AddUserToGroupHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	vars := mux.Vars(r)
	groupId, err := strconv.Atoi(vars["id"])
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	status, err := serviceChat.ValidateGroupOwner(groupId, r, mgr)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, status)
		return
	}

	var users models.Users
	if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	err = mgr.AddUsersToGroup(users, groupId)
	if err != nil {
		if err == chatErrors.ErrUserAlreadyInGroup {
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusConflict)
		} else {
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}

func DeleteUserFromGroupHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	vars := mux.Vars(r)
	groupId, err := strconv.Atoi(vars["id"])
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	status, err := serviceChat.ValidateGroupOwner(groupId, r, mgr)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, status)
		return
	}

	var users models.Users
	if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	err = mgr.DeleteUsersFromGroup(users, groupId)
	if err != nil {
		if err == chatErrors.ErrUserNotInGroup {
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusConflict)
		} else {
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}

func CreateChatHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	us1Id, ok := r.Context().Value("id").(int)
	if !ok {
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	us2Id, err := strconv.Atoi(vars["id"])
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	chatId, err := mgr.CreateChat(us1Id, us2Id)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	err = fileManager.CreateChatMessagesFolder(chatId)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"chat": map[string]interface{}{
			"id": chatId,
		},
	})
}

func ExitGroupHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	userId, ok := r.Context().Value("id").(int)
	if !ok {
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	groupId, err := strconv.Atoi(vars["id"])
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	err = mgr.ExitGroup(userId, groupId)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}

func DeleteChatHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	id, err := serviceChat.GetIdFromVars(r)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	if err := mgr.DeleteChat(id); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	err = fileManager.DeleteChatFiles(id)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}

func UpdateGroupHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	var req models.GroupEdit

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	id, err := serviceChat.GetIdFromVars(r)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	err = mgr.UpdateGroup(id, req)
	switch err {
	case chatErrors.ErrNoGroupFound:
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
		return
	case nil:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "ok",
		})
		return
	default:
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
}
