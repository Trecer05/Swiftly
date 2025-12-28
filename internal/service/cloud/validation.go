package cloud

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/Trecer05/Swiftly/internal/config/logger"
	errorCloudTypes "github.com/Trecer05/Swiftly/internal/errors/cloud"
	models "github.com/Trecer05/Swiftly/internal/model/cloud"
)

func ValidateQueryDescAsc(r *http.Request) string {
	sort := strings.ToUpper(r.URL.Query().Get("sort"))
	if sort != "ASC" && sort != "DESC" {
	sort = "DESC"
	}

	return sort
}

func GetFileAndMetadataFromRequest(r *http.Request, req *models.CreateFileRequest) (multipart.File, *multipart.FileHeader, error) {
	file, header, err := r.FormFile("file")
	if err != nil {
		logger.Logger.Error("Error getting file from form", err)
		return nil, nil, err
	}

	metaStr := r.FormValue("metadata")
	if metaStr == "" {
		logger.Logger.Error("Error getting metadata from form", errorCloudTypes.ErrEmptyMetadata)
		return nil, nil, errorCloudTypes.ErrEmptyMetadata
	}

	if err := json.Unmarshal([]byte(metaStr), req); err != nil {
		logger.Logger.Error("Error unmarshalling metadata", err)
		return nil, nil, err
	}

	if header.Size > models.MaxUploadSize {
		return nil, nil, errorCloudTypes.ErrFileTooLarge
	}

	return file, header, nil
}
