package cloud

import (
	"encoding/json"
	"net/http"

	"github.com/Trecer05/Swiftly/internal/config/logger"
	cloudService "github.com/Trecer05/Swiftly/internal/service/cloud"
	errorAuthTypes "github.com/Trecer05/Swiftly/internal/errors/auth"
	fileManager "github.com/Trecer05/Swiftly/internal/filemanager/cloud"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/cloud"
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

	filesAndFolders, err := mgr.GetUserFilesAndFolders(userID, cloudService.ValidateDescAsc(r))
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
	if err != nil {
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

	filesAndFolders, err := mgr.GetUserFilesAndFoldersByFolderID(userID, folderID.String(), cloudService.ValidateDescAsc(r))
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
