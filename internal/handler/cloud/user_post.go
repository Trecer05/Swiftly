package cloud

import (
	"encoding/json"
	"net/http"

	"github.com/Trecer05/Swiftly/internal/config/logger"
	errorAuthTypes "github.com/Trecer05/Swiftly/internal/errors/auth"
	fileManager "github.com/Trecer05/Swiftly/internal/filemanager/cloud"
	models "github.com/Trecer05/Swiftly/internal/model/cloud"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/cloud"
	cloudService "github.com/Trecer05/Swiftly/internal/service/cloud"
	serviceHttp "github.com/Trecer05/Swiftly/internal/transport/http"
)

func CreateUserFileHandler(w http.ResponseWriter, r *http.Request, manager *manager.Manager) {
	userID, ok := r.Context().Value("id").(int)
	if !ok {
		logger.Logger.Error("Error getting user ID from context", errorAuthTypes.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", errorAuthTypes.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

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
	
	origFilename, storagePath, err := fileManager.SaveUserFile(reader, header, userID, req.FolderID, manager)
	if err != nil {
		logger.Logger.Error("Error saving user file", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
	
	hash := getHash()

	dbReq := models.File{
		FolderID: req.FolderID,
		OriginalFilename: origFilename,
		DisplayName: req.DisplayName,
		StoragePath: storagePath,
		CreatedBy: userID,
		OwnerID: userID,
		OwnerType: models.OwnerTypeUser,
		Hash: hash,
		MimeType: mimeType,
		Visibility: req.Visibility,
		Size: header.Size,
	}

	err = manager.CreateUserFile(&dbReq)
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

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dbReq)
}

func CreateUserFolderHandler(w http.ResponseWriter, r *http.Request, manager *manager.Manager) {
	userID, ok := r.Context().Value("id").(int)
	if !ok {
		logger.Logger.Error("Error getting user ID from context", errorAuthTypes.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", errorAuthTypes.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	var req models.CreateFolderRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Logger.Error("Error decoding request body", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	storagePath, err := fileManager.CreateUserFolder(&req, userID, manager)
	if err != nil {
		logger.Logger.Error("Error creating team folder", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	dbReq := models.Folder{
		Name: req.DisplayName,
		OwnerID: userID,
		CreatedBy: userID,
		StoragePath: storagePath,
		ParentFolderID: req.ParentID,
		OwnerType: req.OwnerType,
		Visibility: req.Visibility,
	}

	if err := manager.CreateUserFolder(&dbReq); err != nil {
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
