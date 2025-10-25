package chat

import (
	"encoding/json"
	"net/http"

	errors "github.com/Trecer05/Swiftly/internal/errors/auth"
	chatErrors "github.com/Trecer05/Swiftly/internal/errors/chat"
	globalErrors "github.com/Trecer05/Swiftly/internal/errors/global"
	models "github.com/Trecer05/Swiftly/internal/model/chat"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/chat"
	serviceHttp "github.com/Trecer05/Swiftly/internal/transport/http"
	fileManager "github.com/Trecer05/Swiftly/internal/filemanager"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	var user models.RegisterUser

	r.ParseMultipartForm(10 << 20)
	jsonData := r.FormValue("json")
	if jsonData == ""{
		serviceHttp.NewErrorBody(w, "application/json", globalErrors.ErrNoJsonData, http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal([]byte(jsonData), &user); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	id, ok := r.Context().Value("id").(int)
	if !ok {
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	if _, _, err := r.FormFile("photo"); err == nil {
        if err := fileManager.AddUserPhoto(r, id); err != nil {
            serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
            return
        }
    } else if err != http.ErrMissingFile {
        serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
        return
    }

	err := mgr.CreateUser(user, id)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}

func EditProfileHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	id, ok := r.Context().Value("id").(int)
	if !ok {
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var user models.ProfileEdit
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	var err error
	if user.Description != nil && user.Name == nil {
		err = mgr.EditProfileDescription(*user.Description, id)
	} else if user.Description == nil && user.Name != nil {
		err = mgr.EditProfileName(*user.Name, id)
	} else {
		serviceHttp.NewErrorBody(w, "application/json", chatErrors.ErrNoData, http.StatusBadRequest)
		return
	}

	switch err {
		case chatErrors.ErrNoUser:
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
			return
		default:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "ok",
			})
	}
}
