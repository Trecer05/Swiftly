package cloud

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Trecer05/Swiftly/internal/config/logger"
	errorAuthTypes "github.com/Trecer05/Swiftly/internal/errors/auth"
	errorCloudTypes "github.com/Trecer05/Swiftly/internal/errors/cloud"
	fileManager "github.com/Trecer05/Swiftly/internal/filemanager/cloud"
	models "github.com/Trecer05/Swiftly/internal/model/cloud"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/cloud"
	cloudService "github.com/Trecer05/Swiftly/internal/service/cloud"
	serviceHttp "github.com/Trecer05/Swiftly/internal/transport/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func UpdateUserFileNameByIDHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	userID, ok := r.Context().Value("id").(int)
	if !ok {
		logger.Logger.Error("Error getting user ID from context", errorAuthTypes.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", errorAuthTypes.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	fileID, err := uuid.Parse(mux.Vars(r)["file_id"])
	if err != nil {
		logger.Logger.Error("Error converting file ID to int", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	var file models.FilenameUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&file); err != nil {
		logger.Logger.Error("Error decoding request body", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	updatedAt, err := cloudService.UpdateFileNameByID(mgr, fileID.String(), file.NewFilename, userID)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.FileUpdateResponse{
		UUID:      fileID,
		UpdatedAt: updatedAt,
		NewFilename: file.NewFilename,
	})
}

func UpdateUserFolderNameByIDHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	userID, ok := r.Context().Value("id").(int)
	if !ok {
		logger.Logger.Error("Error getting user ID from context", errorAuthTypes.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", errorAuthTypes.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	folderID, err := uuid.Parse(mux.Vars(r)["folder_id"])
	if err != nil {
		logger.Logger.Error("Error converting folder ID to int", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	var folder models.FoldernameUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&folder); err != nil {
		logger.Logger.Error("Error decoding request body", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	updatedAt, err := cloudService.UpdateFolderNameByID(mgr, folderID.String(), folder.NewFoldername, userID)
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
		UUID:      folderID,
		UpdatedAt: updatedAt,
		NewFilename: folder.NewFoldername,
	})
}

func ShareUserFileByIDHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	userID, ok := r.Context().Value("id").(int)
	if !ok {
		logger.Logger.Error("Error getting user ID from context", errorAuthTypes.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", errorAuthTypes.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	fileID, err := uuid.Parse(mux.Vars(r)["file_id"])
	if err != nil {
		logger.Logger.Error("Error converting file ID to int", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	if err := mgr.ShareUserFileByID(fileID.String(), userID); err != nil {
		switch {
		case errors.Is(err, errorCloudTypes.ErrFileNotFound):
			logger.Logger.Error("Error sharing file by ID", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
			return
		case err != nil:
			logger.Logger.Error("Error sharing file by ID", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
			return
		}
	}

	link := cloudService.GenerateShareFileLink(fileID.String())
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.ShareLinkResponse{
		Link: link,
	})
}

func UpdateUserFileByIDHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	userID, ok := r.Context().Value("id").(int)
	if !ok {
		logger.Logger.Error("Error getting user ID from context", errorAuthTypes.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", errorAuthTypes.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	fileID, err := uuid.Parse(mux.Vars(r)["file_id"])
	if err != nil {
		logger.Logger.Error("Error converting file ID to int", err)
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

	origFilename, storagePath, err := fileManager.UpdateUserFile(reader, header, userID, fileID, req.ParentID, mgr)
	if err != nil {
		logger.Logger.Error("Error saving user file", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
	
	hash := getHash()

	dbReq := models.File{
		UUID: fileID,
		FolderID: req.ParentID,
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

	if err := mgr.UpdateUserFile(&dbReq); err != nil {
		logger.Logger.Error("Error updating user file", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dbReq)
}

func MoveUserFileByIDHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	userID, ok := r.Context().Value("id").(int)
	if !ok {
		logger.Logger.Error("Error getting user ID from context", errorAuthTypes.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", errorAuthTypes.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	fileID, err := uuid.Parse(mux.Vars(r)["file_id"])
	if err != nil {
		logger.Logger.Error("Error converting file ID to int", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	var req models.MoveUserFileRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Error("Error decoding request body", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	storagePath, err := fileManager.MoveUserFile(&req, userID, fileID.String(), mgr)
	if err != nil {
		switch {
		case errors.Is(err, errorCloudTypes.ErrFileNotFound):
			logger.Logger.Error("Error moving file by ID", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
			return
		case err != nil:
			logger.Logger.Error("Error moving file by ID", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
			return
		}
	}

	if err := mgr.MoveUserFileByID(fileID.String(), req.NewFolderID.String(), userID, storagePath); err != nil {
		logger.Logger.Error("Error moving file by ID", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"new_folder_id": req.NewFolderID.String,
	})
}

func MoveUserFolderByIDHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	userID, ok := r.Context().Value("id").(int)
	if !ok {
		logger.Logger.Error("Error getting user ID from context", errorAuthTypes.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", errorAuthTypes.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	folderID, err := uuid.Parse(mux.Vars(r)["folder_id"])
	if err != nil {
		logger.Logger.Error("Error converting file ID to int", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	var req models.MoveUserFolderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Error("Error decoding request body", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	storagePath, err := fileManager.MoveUserFolder(&req, userID, folderID.String(), mgr)
	if err != nil {
		switch {
		case errors.Is(err, errorCloudTypes.ErrFileNotFound):
			logger.Logger.Error("Error moving file by ID", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
			return
		case err != nil:
			logger.Logger.Error("Error moving file by ID", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
			return
		}
	}

	if err := mgr.MoveUserFolderByID(folderID.String(), req.NewFolderID.String(), userID, storagePath); err != nil {
		logger.Logger.Error("Error moving file by ID", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"new_folder_id": req.NewFolderID.String,
	})
}
