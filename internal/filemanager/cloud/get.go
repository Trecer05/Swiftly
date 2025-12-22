package cloud

import (
	"net/http"
	"os"
	"path/filepath"

	models "github.com/Trecer05/Swiftly/internal/model/cloud"
)

var (
	cloudFolder = os.Getenv("CLOUD_STORAGE")
)

func GetFileSync(fileModel *models.File) ([]byte, string, error) {
	dir := filepath.Join(cloudFolder, fileModel.StoragePath, fileModel.OriginalFilename)
	file, err := os.ReadFile(dir)
	if err != nil {
		return nil, "", err
	}

	fileType := http.DetectContentType(file)

	return file, fileType, nil
}

func GetFileStream(fileModel *models.File) (*os.File, error) {
	dir := filepath.Join(cloudFolder, fileModel.StoragePath, fileModel.OriginalFilename)
	file, err := os.Open(dir)
	if err != nil {
		return nil, err
	}

	return file, nil
}
