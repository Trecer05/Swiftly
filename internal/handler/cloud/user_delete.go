package cloud

import (
	"encoding/json"
	"net/http"

	"github.com/Trecer05/Swiftly/internal/config/logger"
	errorAuthTypes "github.com/Trecer05/Swiftly/internal/errors/auth"
	errorCloudTypes "github.com/Trecer05/Swiftly/internal/errors/cloud"
	fileManager "github.com/Trecer05/Swiftly/internal/filemanager/cloud"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/cloud"
	serviceHttp "github.com/Trecer05/Swiftly/internal/transport/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func DeleteUserFileByIDHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
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

	filepath, err := mgr.DeleteUserFileByID(userID, fileID.String())
	if err != nil {
		if err == errorCloudTypes.ErrFileNotFound {
			logger.Logger.Error("Error deleting user file", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
		} else {
			logger.Logger.Error("Error deleting user file", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		}
		return
	}

	if err := fileManager.DeleteFile(filepath); err != nil {
		logger.Logger.Error("Error deleting file from disk", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}

func DeleteUserFolderByIDHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	userID, ok := r.Context().Value("id").(int)
	if !ok {
		logger.Logger.Error("Error getting user ID from context", errorAuthTypes.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", errorAuthTypes.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	folderID, err := uuid.Parse(mux.Vars(r)["folderID"])
	if err != nil {
		logger.Logger.Error("Error getting folder ID from request", errorCloudTypes.ErrInvalidRequest)
		serviceHttp.NewErrorBody(w, "application/json", errorCloudTypes.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	folderpath, err := mgr.DeleteUserFolderByID(userID, folderID.String())
	if err != nil {
		if err == errorCloudTypes.ErrFolderNotFound {
			logger.Logger.Error("Error deleting user file", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
		} else {
			logger.Logger.Error("Error deleting user file", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		}
		return
	}

	if err := fileManager.DeleteFolder(folderpath); err != nil {
		logger.Logger.Error("Error deleting file from disk", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}
