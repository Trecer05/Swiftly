package cloud

import (
	"encoding/json"
	"errors"
	"io"
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

func GetUserFilesAndFoldersHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	userID, ok := r.Context().Value("id").(int)
	if !ok {
		logger.Logger.Error("Error getting user ID from context", errorAuthTypes.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", errorAuthTypes.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	filesAndFolders, err := mgr.GetUserFilesAndFolders(userID, cloudService.ValidateQueryDescAsc(r))
	if err != nil {
		logger.Logger.Error("Error getting user files and folders", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":   filesAndFolders,
	})
}

func GetUserFileByIDHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
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

	fileModel, err := mgr.GetShortUserFileByID(userID, fileID)
	switch {
	case errors.Is(err, errorCloudTypes.ErrFileNotFound):
		logger.Logger.Error("Error getting user file path by ID", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
		return
	case err != nil:
		logger.Logger.Error("Error getting user file path by ID", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	byteData, err := fileManager.GetUserFileSync(fileModel)
	if err != nil {
		logger.Logger.Error("Error getting user file", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileModel.DisplayName+"\"")
	w.Header().Set("Content-Type", fileModel.MimeType)

	w.WriteHeader(http.StatusOK)
	w.Write(byteData)
}

func GetUserFolderFilesByIDHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
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

	filesAndFolders, err := mgr.GetUserFilesAndFoldersByFolderID(userID, folderID.String(), cloudService.ValidateQueryDescAsc(r))
	if err != nil {
		logger.Logger.Error("Error getting user files and folders", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":   filesAndFolders,
	})
}

func DownloadUserFileByIDHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
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

	fileModel, err := mgr.GetShortUserFileByID(userID, fileID)
	switch {
	case errors.Is(err, errorCloudTypes.ErrFileNotFound):
		logger.Logger.Error("Error getting user file path by ID", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
		return
	case err != nil:
		logger.Logger.Error("Error getting user file path by ID", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	fileByte, err := fileManager.GetUserFileStream(fileModel)
	if err != nil {
		logger.Logger.Error("Error getting user file", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileModel.DisplayName+"\"")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	defer fileByte.Close()
	if _, err := io.Copy(w, fileByte); err != nil {
		http.Error(w, "error sending file", http.StatusInternalServerError)
		return
	}
}

func GetSharedFileHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	_, ok := r.Context().Value("id").(int)
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

	var file models.File

	file.StoragePath, file.DisplayName, err = mgr.GetSharedFile(fileID.String())
	switch {
	case errors.Is(err, errorCloudTypes.ErrFileNotFound):
		logger.Logger.Error("Error getting shared file", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
		return
	case errors.Is(err, errorCloudTypes.ErrFileNotShared):
		logger.Logger.Error("Error getting shared file", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusForbidden)
		return
	case err != nil:
		logger.Logger.Error("Error getting shared file", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	var byteData []byte

	byteData, file.MimeType, err = fileManager.GetFileSync(&file)
	if err != nil {
		logger.Logger.Error("Error getting shared file", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+file.DisplayName+"\"")
	w.Header().Set("Content-Type", file.MimeType)

	w.WriteHeader(http.StatusOK)
	w.Write(byteData)
}

func GetSharedFolderHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	_, ok := r.Context().Value("id").(int)
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

	filesAndFolders, err := mgr.GetSharedFilesAndFoldersByFolderID(folderID.String(), cloudService.ValidateQueryDescAsc(r))
	if err != nil {
		logger.Logger.Error("Error getting user files and folders", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":   filesAndFolders,
	})
}
