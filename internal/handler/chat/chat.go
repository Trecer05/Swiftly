package chat

import (
	"encoding/json"
	"net/http"

	errors "github.com/Trecer05/Swiftly/internal/errors/auth"
	chatErrors "github.com/Trecer05/Swiftly/internal/errors/chat"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/chat"
	serviceChat "github.com/Trecer05/Swiftly/internal/service/chat"
	serviceHttp "github.com/Trecer05/Swiftly/internal/transport/http"
)

func UserChatsInfoHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	limit, offset, err := serviceChat.ValidateLimitOffset(r)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	userId, ok := r.Context().Value("id").(int)
	if !ok {
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	rooms, err := mgr.GetUserRooms(userId, limit, offset)
	switch {
	case err == chatErrors.ErrNoRooms:
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
		return
	case err != nil:
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms)
}

func UserInfoHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	userId, ok := r.Context().Value("id").(int)
	if !ok {
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	userInfo, err := mgr.GetUserInfo(userId)
	switch {
	case err == chatErrors.ErrNoUser:
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
		return
	case err != nil:
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}

func AnotherUserInfoHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	id, err := serviceChat.GetIdFromVars(r)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	userInfo, err := mgr.GetUserInfo(id)
	switch {
	case err == chatErrors.ErrNoUser:
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
		return
	case err != nil:
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}

func GroupUsersHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	id, err := serviceChat.GetIdFromVars(r)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	users, err := mgr.GetGroupUsers(id)
	switch {
	case err == chatErrors.ErrNoUsers:
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
		return
	case err != nil:
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GroupInfoHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	id, err := serviceChat.GetIdFromVars(r)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	info, err := mgr.GetGroupInfo(id)
	switch {
	case err == chatErrors.ErrNoGroupFound:
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
		return
	case err != nil:
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func ChatInfoHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	id, err := serviceChat.GetIdFromVars(r)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

    info, err := mgr.GetUserInfo(id)
	switch {
	case err == chatErrors.ErrNoGroupFound:
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
		return
	case err != nil:
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func ReadChatMessagesHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {

}

func ReadGroupMessagesHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	
}
