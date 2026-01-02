package cloud

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Trecer05/Swiftly/internal/config/logger"
	errorAuthTypes "github.com/Trecer05/Swiftly/internal/errors/auth"
	fileManager "github.com/Trecer05/Swiftly/internal/filemanager/cloud"
	models "github.com/Trecer05/Swiftly/internal/model/cloud"
	redis "github.com/Trecer05/Swiftly/internal/repository/cache/cloud"
	kafkaManager "github.com/Trecer05/Swiftly/internal/repository/kafka/cloud"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/cloud"
	cloudService "github.com/Trecer05/Swiftly/internal/service/cloud"
	serviceHttp "github.com/Trecer05/Swiftly/internal/transport/http"
	"github.com/gorilla/mux"
)

// "/team/{id}/file"
func CreateTeamFileHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, kafkaManager *kafkaManager.KafkaManager, rds *redis.WebSocketManager) {
	vars := mux.Vars(r)
	// Получаем переменные из запроса
	userID, ok := r.Context().Value("id").(int)
	if !ok {
		logger.Logger.Error("Error getting user ID from context", errorAuthTypes.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", errorAuthTypes.ErrUnauthorized, http.StatusUnauthorized)
		return
	}
	teamID, _ := strconv.Atoi(vars["id"]) // Валидируем teamId на уровне роутера.

	// Проверяем, что пользователь состоит в команде
	err := cloudService.CheckUserInTeam(teamID, userID, kafkaManager)
	if err != nil {
		if errors.Is(err, errorAuthTypes.ErrUnauthorized) {
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusForbidden)
			return
		}
		logger.Logger.Error("Error checking user in team", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	// Обрабатываем файл из запроса
	var req models.CreateFileRequest

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

	// Сохраняем файл на диск
	originalFilename, storagePath, err := fileManager.SaveTeamFile(reader, header, teamID, req.FolderID, mgr)
	if err != nil {
		logger.Logger.Error("Error saving team file", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	hash := getHash()

	// Сохраняем информацию о файле в БД
	dbReq := models.File{
		FolderID:         req.FolderID,
		OriginalFilename: originalFilename,
		DisplayName:      req.DisplayName,
		StoragePath:      storagePath,
		CreatedBy:        userID,
		OwnerID:          userID,
		OwnerType:        models.OwnerTypeTeam,
		Hash:             hash,
		MimeType:         mimeType,
		Visibility:       req.Visibility,
		Size:             header.Size,
	}

	err = mgr.CreateTeamFile(&dbReq)

	if err != nil {
		logger.Logger.Error("Error saving info about user file", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)

		if err := fileManager.DeleteFile(storagePath); err != nil {
			logger.Logger.Error("Error deleting user file", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
			return
		}
		return
	}

	dataBytes, err := json.Marshal(dbReq)

	if err != nil {
		logger.Logger.Error("Error marshaling file for websocket message", err)
	} else {
		msg := models.Envelope{
			TeamID: teamID,
			Data:   json.RawMessage(dataBytes),
			Type:   models.FileCreateType,
		}
		rds.SendToTeam(teamID, msg)
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dbReq)
}

// /team/{id:[0-9]+}/folder
func CreateTeamFolderHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, kafkaManager *kafkaManager.KafkaManager, rds *redis.WebSocketManager) {
	// Получаем переменные из запроса
	vars := mux.Vars(r)
	teamID, _ := strconv.Atoi(vars["id"]) // Валидируем teamId на уровне роутера.

	userID, ok := r.Context().Value("id").(int)
	if !ok {
		logger.Logger.Error("Error getting user ID from context", errorAuthTypes.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", errorAuthTypes.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	// Проверяем, что пользователь состоит в команде
	err := cloudService.CheckUserInTeam(teamID, userID, kafkaManager)
	if err != nil {
		if errors.Is(err, errorAuthTypes.ErrUnauthorized) {
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusForbidden)
			return
		}
		logger.Logger.Error("Error checking user in team", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	var req models.CreateFolderRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Logger.Error("Error decoding request body", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	storagePath, err := fileManager.CreateTeamFolder(&req, teamID, mgr)
	if err != nil {
		logger.Logger.Error("Error creating team folder", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	dbReq := models.Folder{
		Name:           req.DisplayName,
		OwnerID:        teamID,
		CreatedBy:      userID,
		StoragePath:    storagePath,
		ParentFolderID: req.ParentID,
		OwnerType:      req.OwnerType,
		Visibility:     req.Visibility,
	}

	if err := mgr.CreateTeamFolder(&dbReq, storagePath); err != nil {
		if err := fileManager.DeleteFolder(storagePath); err != nil {
			logger.Logger.Error("Error deleting team folder", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		}

		logger.Logger.Error("Error creating team folder", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dbReq)
}
