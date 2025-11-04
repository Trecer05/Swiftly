package chat

import (
	"encoding/json"
	"net/http"
	"time"

	errors "github.com/Trecer05/Swiftly/internal/errors/auth"
	kafkaErrors "github.com/Trecer05/Swiftly/internal/errors/kafka"
	chatErrors "github.com/Trecer05/Swiftly/internal/errors/chat"
	globalErrors "github.com/Trecer05/Swiftly/internal/errors/global"
	fileManager "github.com/Trecer05/Swiftly/internal/filemanager"
	models "github.com/Trecer05/Swiftly/internal/model/chat"
	kafkaModels "github.com/Trecer05/Swiftly/internal/model/kafka"
	"github.com/Trecer05/Swiftly/internal/repository/kafka/chat"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/chat"
	serviceHttp "github.com/Trecer05/Swiftly/internal/transport/http"
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

	var url string
	if _, _, err := r.FormFile("photo"); err == nil {
        if url, err = fileManager.AddUserPhoto(r, id); err != nil {
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
		"url": url,
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

func EditUserPasswordHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, kafkaManager *chat.KafkaManager) {
	var req models.PasswordEdit

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	id, ok := r.Context().Value("id").(int)
	if !ok {
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	if err := kafkaManager.SendMessage(ctx, "password", kafkaModels.PasswordEdit{UserID: id, OldPassword: req.OldPassword, NewPassword: req.NewPassword}); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	resp, err := kafkaManager.WaitForResponse(id, 5*time.Second)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusGatewayTimeout)
		return
	}

	switch resp.Type {
	case "status":
		var status kafkaModels.Status
		json.Unmarshal(resp.Payload, &status)
		json.NewEncoder(w).Encode(status)
	case "error":
		var e kafkaModels.Error
		json.Unmarshal(resp.Payload, &e)
		serviceHttp.NewErrorBody(w, "application/json", e.Err, http.StatusBadRequest)
	default:
		serviceHttp.NewErrorBody(w, "application/json", kafkaErrors.ErrorEnvType, http.StatusInternalServerError)
	}
}

func EditUserEmailHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, kafkaManager *chat.KafkaManager) {
	var req models.EmailEdit

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	id, ok := r.Context().Value("id").(int)
	if !ok {
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	if err := kafkaManager.SendMessage(ctx, "email", kafkaModels.EmailEdit{UserID: id, NewEmail: req.NewEmail}); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	resp, err := kafkaManager.WaitForResponse(id, 5*time.Second)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusGatewayTimeout)
		return
	}

	switch resp.Type {
	case "status":
		var status kafkaModels.Status
		json.Unmarshal(resp.Payload, &status)
		json.NewEncoder(w).Encode(status)
	case "error":
		var e kafkaModels.Error
		json.Unmarshal(resp.Payload, &e)
		serviceHttp.NewErrorBody(w, "application/json", e.Err, http.StatusBadRequest)
	default:
		serviceHttp.NewErrorBody(w, "application/json", kafkaErrors.ErrorEnvType, http.StatusInternalServerError)
	}
}

func EditUserPhoneHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, kafkaManager *chat.KafkaManager) {
	var req models.PhoneEdit

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	id, ok := r.Context().Value("id").(int)
	if !ok {
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	if err := kafkaManager.SendMessage(ctx, "phone", kafkaModels.PhoneEdit{UserID: id, NewPhone: req.NewPhone}); err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	resp, err := kafkaManager.WaitForResponse(id, 5*time.Second)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusGatewayTimeout)
		return
	}

	switch resp.Type {
	case "status":
		var status kafkaModels.Status
		json.Unmarshal(resp.Payload, &status)
		json.NewEncoder(w).Encode(status)
	case "error":
		var e kafkaModels.Error
		json.Unmarshal(resp.Payload, &e)
		serviceHttp.NewErrorBody(w, "application/json", e.Err, http.StatusBadRequest)
	default:
		serviceHttp.NewErrorBody(w, "application/json", kafkaErrors.ErrorEnvType, http.StatusInternalServerError)
	}
}
