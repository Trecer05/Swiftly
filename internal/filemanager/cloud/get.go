package cloud

import (
	"net/http"
	"os"
	"path/filepath"

	models "github.com/Trecer05/Swiftly/internal/model/cloud"
)

var (
	teamFolder = os.Getenv("TEAM_STORAGE")
	userFolder  = os.Getenv("USER_STORAGE")
)

func GetFileSync(fileModel *models.File) ([]byte, string, error) {
	dir := filepath.Join(fileModel.StoragePath)
	file, err := os.ReadFile(dir)
	if err != nil {
		return nil, "", err
	}

	fileType := http.DetectContentType(file)

	return file, fileType, nil
}

func GetUserFileSync(fileModel *models.FileShort) ([]byte, error) {
	dir := filepath.Join(fileModel.StoragePath)
	file, err := os.ReadFile(dir)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func GetUserFileStream(fileModel *models.FileShort) (*os.File, error) {
	dir := filepath.Join(fileModel.StoragePath)
	file, err := os.Open(dir)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func GetFileStream(fileModel *models.File) (*os.File, error) {
	dir := filepath.Join(teamFolder, fileModel.StoragePath, fileModel.OriginalFilename)
	file, err := os.Open(dir)
	if err != nil {
		return nil, err
	}

	return file, nil
}
