package cloud

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	errorAuthTypes "github.com/Trecer05/Swiftly/internal/errors/auth"
	errorFileTypes "github.com/Trecer05/Swiftly/internal/errors/file"
	"github.com/Trecer05/Swiftly/internal/model/cloud"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/cloud"
	cloudService "github.com/Trecer05/Swiftly/internal/service/cloud"
	serviceHttp "github.com/Trecer05/Swiftly/internal/transport/http"
	"github.com/google/uuid"
)

func GetTeamFileByIDHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	fileModel, userId, err, statusCode := prepareTeamFile(r, mgr)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, statusCode)
		return
	}
	fileByte, err := cloudService.GetFileSync(fileModel, userId)

	if err != nil {
		if errors.Is(err, errorFileTypes.ErrPermissionDenied) {
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

func DownloadTeamFileByIDHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	fileModel, userId, err, statusCode := prepareTeamFile(r, mgr)
	if err != nil {
		serviceHttp.NewErrorBody(w, "application/json", err, statusCode)
		return
	}
	fileByte, err := cloudService.GetFileStream(fileModel, userId)
	if err != nil {
		if errors.Is(err, errorFileTypes.ErrPermissionDenied) {
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
		http.Error(w, "error sending file", http.StatusInternalServerError)
		return
	}
}

func prepareTeamFile(r *http.Request, mgr *manager.Manager) (*cloud.File, int, error, int) {
	teamId, _ := strconv.Atoi(r.URL.Query().Get("id"))
	fileUUID, err := uuid.Parse(r.URL.Query().Get("file_id"))
	if err != nil {
		return nil, 0, errorFileTypes.ErrFileNotFound, http.StatusNotFound
	}

	userID, ok := r.Context().Value("id").(int)
	if !ok {
		return nil, 0, errorAuthTypes.ErrUnauthorized, http.StatusUnauthorized
	}

	fileModel, err := mgr.GetTeamFileByID(teamId, fileUUID)
	if err != nil {
		if errors.Is(err, errorFileTypes.ErrFileNotFound) {
			return nil, 0, err, http.StatusNotFound
		}
		return nil, 0, err, http.StatusInternalServerError
	}

	return fileModel, userID, nil, http.StatusOK
}
