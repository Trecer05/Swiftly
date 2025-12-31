package cloud

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"

	models "github.com/Trecer05/Swiftly/internal/model/cloud"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/cloud"
	"github.com/google/uuid"
)

func SaveUserFile(reader io.Reader, handler *multipart.FileHeader, userID int, parentID *uuid.UUID, mgr *manager.Manager) (string, string, error) {
	if parentID == nil {
		dir := filepath.Join(userFolder, strconv.Itoa(userID))
		return saveFilesHelper(reader, handler, dir)
	}

	storagePath, err := mgr.GetUserFolderpathByID(userID, parentID.String())
	if err != nil {
		return "", "", err
	}

	return saveFilesHelper(reader, handler, storagePath)
}

func SaveTeamFile(reader io.Reader, handler *multipart.FileHeader, teamID int) (string, string, error) {
	dir := filepath.Join(teamFolder, strconv.Itoa(teamID))
	return saveFilesHelper(reader, handler, dir)
}

func saveFilesHelper(reader io.Reader, handler *multipart.FileHeader, dir string) (string, string, error) {
	fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), handler.Filename)

	filePath := filepath.Join(dir, fileName)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, reader); err != nil {
		return "", "", err
	}
	return fileName, filePath, nil
}

func CreateUserFolder(req *models.CreateFolderRequest, userID int, mgr *manager.Manager) (string, error) {
	if req.ParentID == nil {
		storagePath := filepath.Join(userFolder, strconv.Itoa(userID))

		return saveFolderHelper(storagePath, req.DisplayName)
	}

	storagePath, err := mgr.GetUserFolderpathByID(userID, req.ParentID.String())
	if err != nil {
		return "", err
	}

	return saveFolderHelper(storagePath, req.DisplayName)
}

func CreateTeamFolder(req *models.CreateFolderRequest, teamID int, mgr *manager.Manager) (string, error) {
	if req.ParentID == nil {
		storagePath := filepath.Join(teamFolder, strconv.Itoa(teamID))

		return saveFolderHelper(storagePath, req.DisplayName)
	}

	storagePath, err := mgr.GetTeamFolderpathByID(teamID, req.ParentID.String())
	if err != nil {
		return "", err
	}

	return saveFolderHelper(storagePath, req.DisplayName)
}

func saveFolderHelper(storagePath, name string) (string, error) {
	if err := os.MkdirAll(filepath.Join(storagePath, name), os.ModePerm); err != nil {
		if os.IsExist(err) {
			storagePath = filepath.Join(storagePath, name+fmt.Sprintf("_%d", time.Now().UnixNano()))
			os.MkdirAll(storagePath, os.ModePerm)
			return storagePath, nil
		} else {
			return "", err
		}
	}

	return filepath.Join(storagePath, name), nil
}
