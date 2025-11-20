package chat

import (
	"encoding/json"
	"net/http"
	"time"
	"context"

	"github.com/Trecer05/Swiftly/internal/config/logger"
	errors "github.com/Trecer05/Swiftly/internal/errors/auth"
	chatErrors "github.com/Trecer05/Swiftly/internal/errors/chat"
	globalErrors "github.com/Trecer05/Swiftly/internal/errors/global"
	kafkaErrors "github.com/Trecer05/Swiftly/internal/errors/kafka"
	fileManager "github.com/Trecer05/Swiftly/internal/filemanager"
	models "github.com/Trecer05/Swiftly/internal/model/chat"
	kafkaModels "github.com/Trecer05/Swiftly/internal/model/kafka"
	kafka "github.com/Trecer05/Swiftly/internal/repository/kafka/chat"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/chat"
	serviceHttp "github.com/Trecer05/Swiftly/internal/transport/http"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	var user models.RegisterUser

	r.ParseMultipartForm(10 << 20)
	jsonData := r.FormValue("json")
	if jsonData == ""{
		logger.Logger.Error("No JSON data provided", globalErrors.ErrNoJsonData)
		serviceHttp.NewErrorBody(w, "application/json", globalErrors.ErrNoJsonData, http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal([]byte(jsonData), &user); err != nil {
		logger.Logger.Error("Error unmarshaling JSON data", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	id, ok := r.Context().Value("id").(int)
	if !ok {
		logger.Logger.Error("Error getting user ID from context", errors.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var url string
	if _, _, err := r.FormFile("photo"); err == nil {
        if url, err = fileManager.AddUserPhoto(r, id); err != nil {
            logger.Logger.Error("Error adding user photo", err)
            serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
            return
        }
    } else if err != http.ErrMissingFile {
        logger.Logger.Error("Error getting user photo", err)
        serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
        return
    }

	err := mgr.CreateUser(user, id)
	if err != nil {
		logger.Logger.Error("Error creating user", err)
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
		logger.Logger.Error("Error getting user ID from context", errors.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var user models.ProfileEdit
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.Logger.Error("Error decoding profile edit request", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	var err error
	if user.Description != nil && user.Name == nil {
		err = mgr.EditProfileDescription(*user.Description, id)
	} else if user.Description == nil && user.Name != nil {
		err = mgr.EditProfileName(*user.Name, id)
	} else {
		logger.Logger.Error("Error editing profile", err)
		serviceHttp.NewErrorBody(w, "application/json", chatErrors.ErrNoData, http.StatusBadRequest)
		return
	}

	switch err {
		case chatErrors.ErrNoUser:
			logger.Logger.Error("Error editing profile", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
			return
		default:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "ok",
			})
	}
}

func EditUserPasswordHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, kafkaManager *kafka.KafkaManager, ctx context.Context) {
	var req models.PasswordEdit

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Error("Error decoding password edit request", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	id, ok := r.Context().Value("id").(int)
	if !ok {
		logger.Logger.Error("Error getting user ID from context", errors.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	if err := kafkaManager.SendMessage(ctx, "password", kafkaModels.PasswordEdit{UserID: id, OldPassword: req.OldPassword, NewPassword: req.NewPassword}); err != nil {
		logger.Logger.Error("Error sending password edit message", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	resp, err := kafkaManager.WaitForResponse(id, 5*time.Second)
	if err != nil {
		logger.Logger.Error("Error waiting for password edit response", err)
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
		logger.Logger.Error("Error decoding password edit response", err)
		json.Unmarshal(resp.Payload, &e)
		serviceHttp.NewErrorBody(w, "application/json", e.Err, http.StatusBadRequest)
	default:
		serviceHttp.NewErrorBody(w, "application/json", kafkaErrors.ErrorEnvType, http.StatusInternalServerError)
	}
}

func EditUserEmailHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, kafkaManager *kafka.KafkaManager, ctx context.Context) {
	var req models.EmailEdit

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Error("Error decoding email edit request", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	id, ok := r.Context().Value("id").(int)
	if !ok {
		logger.Logger.Error("Invalid user ID", errors.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	if err := kafkaManager.SendMessage(ctx, "email", kafkaModels.EmailEdit{UserID: id, NewEmail: req.NewEmail}); err != nil {
		logger.Logger.Error("Error sending email edit message", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	resp, err := kafkaManager.WaitForResponse(id, 5*time.Second)
	if err != nil {
		logger.Logger.Error("Error waiting for response", err)
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
		logger.Logger.Error("Error decoding error response", err)
		json.Unmarshal(resp.Payload, &e)
		serviceHttp.NewErrorBody(w, "application/json", e.Err, http.StatusBadRequest)
	default:
		serviceHttp.NewErrorBody(w, "application/json", kafkaErrors.ErrorEnvType, http.StatusInternalServerError)
	}
}

func EditUserPhoneHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, kafkaManager *kafka.KafkaManager, ctx context.Context) {
	var req models.PhoneEdit

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Error("Error decoding request body", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	id, ok := r.Context().Value("id").(int)
	if !ok {
		logger.Logger.Error("Invalid user ID")
		serviceHttp.NewErrorBody(w, "application/json", errors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	if err := kafkaManager.SendMessage(ctx, "phone", kafkaModels.PhoneEdit{UserID: id, NewPhone: req.NewPhone}); err != nil {
		logger.Logger.Error("Error sending phone edit message", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	resp, err := kafkaManager.WaitForResponse(id, 5*time.Second)
	if err != nil {
		logger.Logger.Error("Error waiting for response", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusGatewayTimeout)
		return
	}

	switch resp.Type {
	case "status":
		var status kafkaModels.Status
		logger.Logger.Info("Received status response")
		json.Unmarshal(resp.Payload, &status)
		json.NewEncoder(w).Encode(status)
	case "error":
		var e kafkaModels.Error
		logger.Logger.Error("Error unmarshalling error response", err)
		json.Unmarshal(resp.Payload, &e)
		serviceHttp.NewErrorBody(w, "application/json", e.Err, http.StatusBadRequest)
	default:
		serviceHttp.NewErrorBody(w, "application/json", kafkaErrors.ErrorEnvType, http.StatusInternalServerError)
	}
}
