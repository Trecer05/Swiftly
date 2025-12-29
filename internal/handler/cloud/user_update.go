package cloud

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Trecer05/Swiftly/internal/config/logger"
	errorAuthTypes "github.com/Trecer05/Swiftly/internal/errors/auth"
	errorCloudTypes "github.com/Trecer05/Swiftly/internal/errors/cloud"
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
