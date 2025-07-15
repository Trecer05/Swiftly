package chat

import (
	"encoding/json"
	"net/http"

	errors "github.com/Trecer05/Swiftly/internal/errors/auth"
	chatErrors "github.com/Trecer05/Swiftly/internal/errors/chat"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/chat"
	serviceHttp "github.com/Trecer05/Swiftly/internal/transport/http"
	serviceChat "github.com/Trecer05/Swiftly/internal/service/chat"
)

func UserChatsInfoHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	limit, offset, err := serviceChat.ValidateLimitOffset(r)
	if err != nil {
		serviceHttp.NewHeaderBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	userId, ok := r.Context().Value("id").(int)
	if !ok {
		serviceHttp.NewHeaderBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	rooms, err := mgr.GetUserRooms(userId, limit, offset)
	switch {
	case err == chatErrors.ErrNoRooms:
		serviceHttp.NewHeaderBody(w, "application/json", err, http.StatusNotFound)
		return
	case err != nil:
		serviceHttp.NewHeaderBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms)
}

func UserInfoHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	userId, ok := r.Context().Value("id").(int)
	if !ok {
		serviceHttp.NewHeaderBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	userInfo, err := mgr.GetUserInfo(userId)
	switch {
	case err == chatErrors.ErrNoUser:
		serviceHttp.NewHeaderBody(w, "application/json", err, http.StatusNotFound)
		return
	case err != nil:
		serviceHttp.NewHeaderBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}
