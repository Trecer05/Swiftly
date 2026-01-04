package cloud

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Trecer05/Swiftly/internal/config/logger"
	errorAuthTypes "github.com/Trecer05/Swiftly/internal/errors/auth"
	cloudErrors "github.com/Trecer05/Swiftly/internal/errors/cloud"
	fileManager "github.com/Trecer05/Swiftly/internal/filemanager/cloud"
	models "github.com/Trecer05/Swiftly/internal/model/cloud"
	redis "github.com/Trecer05/Swiftly/internal/repository/cache/cloud"
	kafkaManager "github.com/Trecer05/Swiftly/internal/repository/kafka/cloud"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/cloud"
	cloudService "github.com/Trecer05/Swiftly/internal/service/cloud"
	serviceHttp "github.com/Trecer05/Swiftly/internal/transport/http"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// "/team/{id:[0-9]+}/file/{file_id}"
func DeleteTeamFileByIDHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, kafkaManager *kafkaManager.KafkaManager, rds *redis.WebSocketManager) {

	// Получаем переменные из запроса
	vars := mux.Vars(r)
	teamID, _ := strconv.Atoi(vars["id"]) // Валидируем teamId на уровне роутера.

	userID, err := cloudService.AuthorizeUserInTeam(r, teamID, kafkaManager)
	if err != nil {
		if errors.Is(err, errorAuthTypes.ErrUnauthorized) {
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusForbidden)
			return
		}
		logger.Logger.Error("Error checking user in team", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	fileUUID, err := uuid.Parse(vars["file_id"])
	if err != nil {
		logger.Logger.Error("Error parsing file UUID", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	filePath, err := mgr.DeleteTeamFileByID(teamID, userID, fileUUID)
	if err != nil {
		switch {
		case errors.Is(err, cloudErrors.ErrFileNotFound):
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
			return
		case errors.Is(err, cloudErrors.ErrNoPermissions):
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusForbidden)
			return
		}
		logger.Logger.Error("Error deleting team file by ID", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	err = fileManager.DeleteFile(filePath)
	if err != nil {
		logger.Logger.Error("Error deleting file from storage", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	dataBytes, err := json.Marshal(fileUUID)

	if err != nil {
		logger.Logger.Error("Error marshaling file for websocket message", err)
	} else {
		msg := models.Envelope{
			TeamID: teamID,
			Data:   json.RawMessage(dataBytes),
			Type:   models.FileDeleteType,
		}
		rds.SendToTeam(teamID, msg)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}

// "/team/{id:[0-9]+}/folder/{folder_id}"
func DeleteTeamFolderByIDHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, kafkaManager *kafkaManager.KafkaManager, rds *redis.WebSocketManager) {
	// Получаем переменные из запроса
	vars := mux.Vars(r)
	teamID, _ := strconv.Atoi(vars["id"]) // Валидируем teamId на уровне роутера.

	userID, err := cloudService.AuthorizeUserInTeam(r, teamID, kafkaManager)
	if err != nil {
		if errors.Is(err, errorAuthTypes.ErrUnauthorized) {
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusForbidden)
			return
		}
		logger.Logger.Error("Error checking user in team", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	folderUUID, err := uuid.Parse(vars["folder_id"])
	if err != nil {
		logger.Logger.Error("Error parsing folder UUID", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	folderPath, err := mgr.DeleteTeamFolderByID(teamID, userID, folderUUID)
	if err != nil {
		switch {
		case errors.Is(err, cloudErrors.ErrFileNotFound):
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
			return
		case errors.Is(err, cloudErrors.ErrNoPermissions):
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusForbidden)
			return
		}
		logger.Logger.Error("Error deleting team folder by ID", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	err = fileManager.DeleteFolder(folderPath)
	if err != nil {
		logger.Logger.Error("Error deleting file from storage", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	dataBytes, err := json.Marshal(folderUUID)

	if err != nil {
		logger.Logger.Error("Error marshaling file for websocket message", err)
	} else {
		msg := models.Envelope{
			TeamID: teamID,
			Data:   json.RawMessage(dataBytes),
			Type:   models.FileDeleteType,
		}
		rds.SendToTeam(teamID, msg)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}
