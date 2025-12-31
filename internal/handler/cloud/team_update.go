package cloud

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Trecer05/Swiftly/internal/config/logger"
	errorAuthTypes "github.com/Trecer05/Swiftly/internal/errors/auth"
	errorCloudTypes "github.com/Trecer05/Swiftly/internal/errors/cloud"
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

// "/team/{id}/file/{file_id}"
func UpdateTeamFileByIDHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, kafkaManager *kafkaManager.KafkaManager, rds *redis.WebSocketManager) {
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
		logger.Logger.Error("Error parsing folder UUID", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	var req models.UpdateFileRequest

	r.Body = http.MaxBytesReader(w, r.Body, models.MaxUploadSize) // 50 MB

	if err := r.ParseMultipartForm(models.MaxUploadSize); err != nil {
		logger.Logger.Error("Error parsing multipart form", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusRequestEntityTooLarge)
		return
	}

	file, header, err := cloudService.GetFileAndMetadataFromRequest(r, &req)
	if err != nil {
		logger.Logger.Error("Error getting file and metadata from request", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	defer file.Close()

	reader, getHash := fileManager.HashingReader(file)

	mimeType, reader, err := fileManager.DetectMimeType(reader)
	if err != nil {
		logger.Logger.Error("Error detecting mime type", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	origFilename, storagePath, err := fileManager.UpdateTeamFile(reader, header, teamID, userID, fileUUID, req.ParentID, mgr)
	if err != nil {
		logger.Logger.Error("Error saving user file", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	hash := getHash()

	dbReq := models.File{
		UUID:             fileUUID,
		FolderID:         req.ParentID,
		OriginalFilename: origFilename,
		DisplayName:      req.DisplayName,
		StoragePath:      storagePath,
		CreatedBy:        userID,
		OwnerID:          userID,
		OwnerType:        models.OwnerTypeUser,
		Hash:             hash,
		MimeType:         mimeType,
		Visibility:       req.Visibility,
		Size:             header.Size,
	}

	if err := mgr.UpdateTeamFile(&dbReq); err != nil {
		logger.Logger.Error("Error updating user file", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	dataBytes, err := json.Marshal(dbReq)

	if err != nil {
		logger.Logger.Error("Error marshaling file for websocket message", err)
	} else {
		msg := models.Envelope{
			TeamID: teamID,
			Data:   json.RawMessage(dataBytes),
			Type:   models.FileUpdateType,
		}
		rds.SendToTeam(teamID, msg)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dbReq)
}

// "/team/{id}/file/{file_id}/name"
func UpdateTeamFileNameByIDHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, kafkaManager *kafkaManager.KafkaManager, rds *redis.WebSocketManager) {
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
		logger.Logger.Error("Error parsing folder UUID", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	var file models.FilenameUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&file); err != nil {
		logger.Logger.Error("Error decoding request body", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	updatedAt, err := cloudService.UpdateTeamFileNameByID(mgr, fileUUID.String(), file.NewFilename, teamID, userID)
	switch {
	case errors.Is(err, errorCloudTypes.ErrFileNotFound):
		logger.Logger.Error("Error update user filename by ID", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
		return
	case err != nil:
		logger.Logger.Error("Error update user filename by ID", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	dataBytes, err := json.Marshal(file.NewFilename)

	if err != nil {
		logger.Logger.Error("Error marshaling file for websocket message", err)
	} else {
		msg := models.Envelope{
			TeamID: teamID,
			Data:   json.RawMessage(dataBytes),
			Type:   models.FileUpdateType,
		}
		rds.SendToTeam(teamID, msg)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.FileUpdateResponse{
		UUID:        fileUUID,
		UpdatedAt:   updatedAt,
		NewFilename: file.NewFilename,
	})
}

// "/team/{id:[0-9]+}/file/{file_id}/share"
func ShareTeamFileByIDHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, kafkaManager *kafkaManager.KafkaManager) {
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
		logger.Logger.Error("Error parsing folder UUID", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	if err := mgr.ShareTeamFileByID(fileUUID.String(), teamID, userID); err != nil {
		switch {
		case errors.Is(err, errorCloudTypes.ErrFileNotFound):
			logger.Logger.Error("Error sharing file by ID", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
			return
		}
		logger.Logger.Error("Error sharing file by ID", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	link := cloudService.GenerateShareFileLink(fileUUID.String())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.ShareLinkResponse{
		Link: link,
	})

}

// "/team/{id:[0-9]+}/folder/{folder_id}/name"
func UpdateTeamFolderNameByIDHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, kafkaManager *kafkaManager.KafkaManager, rds *redis.WebSocketManager) {
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

	var folder models.FoldernameUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&folder); err != nil {
		logger.Logger.Error("Error decoding request body", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	updatedAt, err := cloudService.UpdateTeamFolderNameByID(mgr, folderUUID.String(), folder.NewFoldername, teamID, userID)
	switch {
	case errors.Is(err, errorCloudTypes.ErrFolderNotFound):
		logger.Logger.Error("Error updating file name", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
		return
	case err != nil:
		logger.Logger.Error("Error updating file name", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.FileUpdateResponse{
		UUID:        folderUUID,
		UpdatedAt:   updatedAt,
		NewFilename: folder.NewFoldername,
	})
}

// "/team/{id:[0-9]+}/folder/{folder_id}/move"
func MoveTeamFolderByIDHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, kafkaManager *kafkaManager.KafkaManager, rds *redis.WebSocketManager) {
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

	var req models.MoveFolderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Error("Error decoding request body", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	storagePath, err := fileManager.MoveTeamFolder(&req, teamID, folderUUID.String(), mgr)
	if err != nil {
		switch {
		case errors.Is(err, errorCloudTypes.ErrFileNotFound):
			logger.Logger.Error("Error moving file by ID", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
			return
		}
		logger.Logger.Error("Error moving file by ID", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	if err := mgr.MoveTeamFolderByID(folderUUID.String(), req.NewFolderID.String(), userID, storagePath); err != nil {
		logger.Logger.Error("Error moving folder by ID", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"new_folder_id": req.NewFolderID.String,
	})
}

// "/team/{id:[0-9]+}/file/{file_id}/move"
func MoveTeamFileByIDHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, kafkaManager *kafkaManager.KafkaManager, rds *redis.WebSocketManager) {
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

	fileUUID, err := uuid.Parse(vars["folder_id"])
	if err != nil {
		logger.Logger.Error("Error parsing folder UUID", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	var req models.MoveFileRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Error("Error decoding request body", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	storagePath, err := fileManager.MoveTeamFile(&req, userID, fileUUID.String(), mgr)
	if err != nil {
		switch {
		case errors.Is(err, errorCloudTypes.ErrFileNotFound):
			logger.Logger.Error("Error moving file by ID", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
			return
		}
		logger.Logger.Error("Error moving file by ID", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	if err := mgr.MoveTeamFileByID(fileUUID.String(), req.NewFolderID.String(), teamID, storagePath); err != nil {
		logger.Logger.Error("Error moving file by ID", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"new_folder_id": req.NewFolderID.String,
	})
}
