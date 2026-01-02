package cloud

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/Trecer05/Swiftly/internal/config/logger"
	errorAuthTypes "github.com/Trecer05/Swiftly/internal/errors/auth"
	errorFileTypes "github.com/Trecer05/Swiftly/internal/errors/file"
	"github.com/Trecer05/Swiftly/internal/model/cloud"
	kafkaManager "github.com/Trecer05/Swiftly/internal/repository/kafka/cloud"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/cloud"
	cloudService "github.com/Trecer05/Swiftly/internal/service/cloud"
	serviceHttp "github.com/Trecer05/Swiftly/internal/transport/http"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetTeamFileByIDHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, kafkaManager *kafkaManager.KafkaManager) {
	fileModel, teamID, err, statusCode := prepareTeamFile(r, mgr)
	if err != nil {
		logger.Logger.Error("Error preparing file", err)
		serviceHttp.NewErrorBody(w, "application/json", err, statusCode)
		return
	}
	fileByte, err := cloudService.GetFileSync(fileModel, teamID, kafkaManager)

	if err != nil {
		if errors.Is(err, errorFileTypes.ErrPermissionDenied) {
			logger.Logger.Error("Permission denied", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusForbidden)
			return
		}
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileModel.DisplayName+"\"")
	w.Header().Set("Content-Type", fileModel.MimeType)

	w.WriteHeader(http.StatusOK)
	w.Write(fileByte)
}

func DownloadTeamFileByIDHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, kafkaManager *kafkaManager.KafkaManager) {
	fileModel, userId, err, statusCode := prepareTeamFile(r, mgr)
	if err != nil {
		logger.Logger.Error("Error preparing file", err)
		serviceHttp.NewErrorBody(w, "application/json", err, statusCode)
		return
	}
	fileByte, err := cloudService.GetFileStream(fileModel, userId, kafkaManager)
	if err != nil {
		if errors.Is(err, errorFileTypes.ErrPermissionDenied) {
			logger.Logger.Error("Permission denied", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusForbidden)
			return
		}
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileModel.DisplayName+"\"")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	defer fileByte.Close()
	if _, err := io.Copy(w, fileByte); err != nil {
		logger.Logger.Error("Error sending file", err)
		http.Error(w, "error sending file", http.StatusInternalServerError)
		return
	}
}

func prepareTeamFile(r *http.Request, mgr *manager.Manager) (*cloud.File, int, error, int) {
	vars := mux.Vars(r)
	teamID, _ := strconv.Atoi(vars["id"])
	fileUUID, err := uuid.Parse(vars["file_id"])
	if err != nil {
		logger.Logger.Error("Error parsing file UUID", err)
		return nil, 0, errorFileTypes.ErrFileNotFound, http.StatusNotFound
	}

	userID, ok := r.Context().Value("id").(int)
	if !ok {
		return nil, 0, errorAuthTypes.ErrUnauthorized, http.StatusUnauthorized
	}

	fileModel, err := mgr.GetTeamFileByID(teamID, fileUUID)
	if err != nil {
		logger.Logger.Error("Error getting team file by ID", err)
		if errors.Is(err, errorFileTypes.ErrFileNotFound) {
			return nil, 0, err, http.StatusNotFound
		}
		return nil, 0, err, http.StatusInternalServerError
	}

	return fileModel, userID, nil, http.StatusOK
}

// /team/{id}/folder/{folder_id}
func GetTeamFolderFilesByIDHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, kafkaManager *kafkaManager.KafkaManager) {
	vars := mux.Vars(r)
	teamId, _ := strconv.Atoi(vars["id"])
	folderId, err := uuid.Parse(vars["folder_id"])
	if err != nil {
		logger.Logger.Error("Error parsing folder UUID", err)
		serviceHttp.NewErrorBody(w, "application/json", errorFileTypes.ErrFolderNotFound, http.StatusNotFound)
		return
	}
	userID, ok := r.Context().Value("id").(int)
	if !ok {
		logger.Logger.Error("Error getting user ID from context", errorAuthTypes.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", errorAuthTypes.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	err = cloudService.CheckUserInTeam(teamId, userID, kafkaManager)
	if err != nil {
		if errors.Is(err, errorAuthTypes.ErrUnauthorized) {
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusForbidden)
			return
		}
		logger.Logger.Error("Error checking user in team", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	folderModel, err := mgr.GetTeamFolderByTeamID(teamId, userID, folderId)
	if err != nil {
		if errors.Is(err, errorFileTypes.ErrFolderNotFound) {
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
			return
		}
		logger.Logger.Error("Error getting team folder by ID", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(folderModel)
}

// "/team/{id}"
func GetTeamFilesAndFoldersHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, kafkaManager *kafkaManager.KafkaManager) {
	vars := mux.Vars(r)
	teamID, _ := strconv.Atoi(vars["id"]) // Валидируем teamId на уровне роутера.

	userID, ok := r.Context().Value("id").(int)
	if !ok {
		logger.Logger.Error("Error getting user ID from context", errorAuthTypes.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", errorAuthTypes.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	err := cloudService.CheckUserInTeam(teamID, userID, kafkaManager)

	if err != nil {
		if errors.Is(err, errorFileTypes.ErrPermissionDenied) {
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusForbidden)
			return
		}
		logger.Logger.Error("Error checking user in team", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	filesAndFolders, err := mgr.GetTeamFilesAndFolders(teamID, userID, cloudService.ValidateQueryDescAsc(r))

	if err != nil {
		logger.Logger.Error("Error getting user files and folders", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": filesAndFolders,
	})
}
