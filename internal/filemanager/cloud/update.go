package cloud

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	cloudErrors "github.com/Trecer05/Swiftly/internal/errors/cloud"
	models "github.com/Trecer05/Swiftly/internal/model/cloud"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/cloud"
	"github.com/google/uuid"
)

func UpdateFileName(filePath, newFilename string) (string, string, error) {
	dir := filepath.Dir(filePath)
	originalExt := filepath.Ext(filePath)

	newName := newFilename

	timestamp := fmt.Sprintf("%d_", time.Now().UnixNano())
	newName = timestamp + newName

	if originalExt != "" && filepath.Ext(newName) == "" {
		newName += originalExt
	}

	newPath := filepath.Join(dir, newName)

	if _, err := os.Stat(newPath); err == nil {
		newName = fmt.Sprintf("%s_%d", newName[:len(newName)-len(filepath.Ext(newName))], time.Now().UnixNano())
		if originalExt != "" {
			newName += originalExt
		}
		newPath = filepath.Join(dir, newName)
	}

	if err := os.Rename(filePath, newPath); err != nil {
		return "", "", err
	}

	return newName, newPath, nil
}

func UpdateFolderName(folderPath, newFolderName string) (string, string, error) {
	fileInfo, err := os.Stat(folderPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", "", err
		}
		return "", "", err
	}

	if !fileInfo.IsDir() {
		return "", "", err
	}

	parentDir := filepath.Dir(folderPath)

	cleanedName := cleanFolderName(newFolderName)
	if cleanedName == "" {
		return "", "", err
	}

	newPath := filepath.Join(parentDir, cleanedName)

	if _, err := os.Stat(newPath); err == nil {
		timestamp := fmt.Sprintf("_%d", time.Now().UnixNano())
		cleanedName += timestamp
		newPath = filepath.Join(parentDir, cleanedName)
	}

	if err := os.Rename(folderPath, newPath); err != nil {
		return "", "", err
	}

	return cleanedName, newPath, nil
}

func cleanFolderName(name string) string {
	name = strings.TrimSpace(name)

	invalidChars := []string{"<", ">", ":", "\"", "|", "?", "*", "\\", "/", "\x00"}

	for _, char := range invalidChars {
		name = strings.ReplaceAll(name, char, "_")
	}

	name = strings.Trim(name, ".")

	reservedNames := []string{"CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}

	upperName := strings.ToUpper(name)
	for _, reserved := range reservedNames {
		if upperName == reserved || strings.HasPrefix(upperName, reserved+".") {
			return "_" + name
		}
	}

	return name
}

func UpdateUserFile(reader io.Reader, handler *multipart.FileHeader, userID int, fileID uuid.UUID, parentID *uuid.UUID, mgr *manager.Manager) (string, string, error) {
	var storagePath string
	var olderParentDir string
	var err error

	if parentID == nil {
		origFilename, err := mgr.GetOriginalUserFilenameByID(userID, fileID.String())
		if err != nil {
			return "", "", err
		}

		olderParentDir = filepath.Join(userFolder, strconv.Itoa(userID))
		storagePath = filepath.Join(userFolder, strconv.Itoa(userID), origFilename)
	} else {
		olderParentDir, err = mgr.GetUserFolderpathByID(userID, parentID.String())
		if err != nil {
			return "", "", err
		}

		origFilename, err := mgr.GetOriginalUserFilenameByID(userID, fileID.String())
		if err != nil {
			return "", "", err
		}
		storagePath = filepath.Join(olderParentDir, origFilename)
	}

	newFilename := fmt.Sprintf(
		"%d_%s",
		time.Now().UnixNano(),
		handler.Filename,
	)

	newStoragePath := filepath.Join(olderParentDir, newFilename)

	tmpFile, err := os.CreateTemp(olderParentDir, ".tmp-*")
	if err != nil {
		return "", "", err
	}
	tmpPath := tmpFile.Name()

	defer func() {
		tmpFile.Close()
		os.Remove(tmpPath)
	}()

	if _, err := io.Copy(tmpFile, reader); err != nil {
		return "", "", err
	}

	if err := tmpFile.Sync(); err != nil {
		return "", "", err
	}

	if err := tmpFile.Close(); err != nil {
		return "", "", err
	}

	if err := os.Rename(tmpPath, newStoragePath); err != nil {
		return "", "", err
	}

	if dirFd, err := os.Open(olderParentDir); err == nil {
		_ = dirFd.Sync()
		_ = dirFd.Close()
	}

	if err := os.Remove(storagePath); err != nil && !os.IsNotExist(err) {
		return "", "", err
	}

	return newFilename, newStoragePath, nil
}

func UpdateTeamFile(reader io.Reader, handler *multipart.FileHeader, teamID int, userID int, fileID uuid.UUID, parentID *uuid.UUID, mgr *manager.Manager) (string, string, error) {
	var storagePath string
	var olderParentDir string
	var err error

	if parentID == nil {
		origFilename, err := mgr.GetOriginalTeamFilenameByID(teamID, fileID.String())
		if err != nil {
			return "", "", err
		}

		olderParentDir = filepath.Join(teamFolder, strconv.Itoa(teamID))
		storagePath = filepath.Join(teamFolder, strconv.Itoa(teamID), origFilename)
	} else {
		olderParentDir, err = mgr.GetTeamFolderpathByID(teamID, parentID.String())
		if err != nil {
			return "", "", err
		}

		origFilename, err := mgr.GetOriginalTeamFilenameByID(teamID, fileID.String())
		if err != nil {
			return "", "", err
		}
		storagePath = filepath.Join(olderParentDir, origFilename)
	}

	newFilename := fmt.Sprintf(
		"%d_%s",
		time.Now().UnixNano(),
		handler.Filename,
	)

	newStoragePath := filepath.Join(olderParentDir, newFilename)

	tmpFile, err := os.CreateTemp(olderParentDir, ".tmp-*")
	if err != nil {
		return "", "", err
	}
	tmpPath := tmpFile.Name()

	defer func() {
		tmpFile.Close()
		os.Remove(tmpPath)
	}()

	if _, err := io.Copy(tmpFile, reader); err != nil {
		return "", "", err
	}

	if err := tmpFile.Sync(); err != nil {
		return "", "", err
	}

	if err := tmpFile.Close(); err != nil {
		return "", "", err
	}

	if err := os.Rename(tmpPath, newStoragePath); err != nil {
		return "", "", err
	}

	if dirFd, err := os.Open(olderParentDir); err == nil {
		_ = dirFd.Sync()
		_ = dirFd.Close()
	}

	if err := os.Remove(storagePath); err != nil && !os.IsNotExist(err) {
		return "", "", err
	}

	return newFilename, newStoragePath, nil
}

func MoveUserFile(req *models.MoveFileRequest, userID int, fileID string, mgr *manager.Manager) (string, error) {
	var olderStoragePath, newStoragePath, origFilename string
	var err error

	if req.NewFolderID == nil {
		origFilename, olderStoragePath, err = mgr.GetOriginalUserFilenameAndStoragePathByID(userID, fileID)
		if err != nil {
			return "", err
		}

		newStoragePath = filepath.Join(userFolder, strconv.Itoa(userID), origFilename)
	} else {
		newStoragePath, err = mgr.GetUserFolderpathByID(userID, req.NewFolderID.String())
		if err != nil {
			return "", err
		}

		origFilename, olderStoragePath, err = mgr.GetOriginalUserFilenameAndStoragePathByID(userID, fileID)
		if err != nil {
			return "", err
		}

		newStoragePath = filepath.Join(newStoragePath, origFilename)
	}

	if _, err := os.Stat(newStoragePath); err == nil {
		return "", cloudErrors.ErrFileAlreadyExists
	}

	if olderStoragePath == newStoragePath {
		return newStoragePath, nil
	}

	if err := os.MkdirAll(filepath.Dir(newStoragePath), 0755); err != nil {
		return "", err
	}

	if err := os.Rename(olderStoragePath, newStoragePath); err != nil {
		return "", err
	}

	if dirFd, err := os.Open(filepath.Dir(newStoragePath)); err == nil {
		_ = dirFd.Sync()
		_ = dirFd.Close()
	}

	return newStoragePath, nil
}

func MoveUserFolder(req *models.MoveFolderRequest, userID int, fileID string, mgr *manager.Manager) (string, error) {
	var olderStoragePath, newStoragePath string
	var err error

	if req.NewFolderID == nil {
		olderStoragePath, err = mgr.GetUserFolderpathByID(userID, fileID)
		if err != nil {
			return "", err
		}

		newStoragePath = filepath.Join(userFolder, strconv.Itoa(userID), req.FolderName)
	} else {
		newStoragePath, err = mgr.GetUserFolderpathByID(userID, req.NewFolderID.String())
		if err != nil {
			return "", err
		}

		olderStoragePath, err = mgr.GetUserFolderpathByID(userID, fileID)
		if err != nil {
			return "", err
		}

		newStoragePath = filepath.Join(newStoragePath, req.FolderName)
	}

	if err := os.MkdirAll(filepath.Dir(newStoragePath), 0755); err != nil {
		return "", err
	}

	if err := os.Rename(olderStoragePath, newStoragePath); err != nil {
		return "", err
	}

	if dirFd, err := os.Open(filepath.Dir(newStoragePath)); err == nil {
		_ = dirFd.Sync()
		_ = dirFd.Close()
	}

	return "", nil
}

func MoveTeamFolder(req *models.MoveFolderRequest, teamID int, folderID string, mgr *manager.Manager) (string, error) {
	var olderStoragePath, newStoragePath string
	var err error

	if req.NewFolderID == nil {
		olderStoragePath, err = mgr.GetTeamFolderpathByID(teamID, folderID)
		if err != nil {
			return "", err
		}

		newStoragePath = filepath.Join(teamFolder, strconv.Itoa(teamID), req.FolderName)
	} else {
		newStoragePath, err = mgr.GetTeamFolderpathByID(teamID, req.NewFolderID.String())
		if err != nil {
			return "", err
		}

		olderStoragePath, err = mgr.GetTeamFolderpathByID(teamID, folderID)
		if err != nil {
			return "", err
		}

		newStoragePath = filepath.Join(newStoragePath, req.FolderName)
	}

	if err := os.MkdirAll(filepath.Dir(newStoragePath), 0755); err != nil {
		return "", err
	}

	if err := os.Rename(olderStoragePath, newStoragePath); err != nil {
		return "", err
	}

	if dirFd, err := os.Open(filepath.Dir(newStoragePath)); err == nil {
		_ = dirFd.Sync()
		_ = dirFd.Close()
	}

	return newStoragePath, nil
}

func MoveTeamFile(req *models.MoveFileRequest, teamID int, fileID string, mgr *manager.Manager) (string, error) {
	var olderStoragePath, newStoragePath, origFilename string
	var err error

	if req.NewFolderID == nil {
		origFilename, olderStoragePath, err = mgr.GetOriginalTeamFilenameAndStoragePathByID(teamID, fileID)
		if err != nil {
			return "", err
		}

		newStoragePath = filepath.Join(userFolder, strconv.Itoa(teamID), origFilename)
	} else {
		newStoragePath, err = mgr.GetTeamFolderpathByID(teamID, req.NewFolderID.String())
		if err != nil {
			return "", err
		}

		origFilename, olderStoragePath, err = mgr.GetOriginalTeamFilenameAndStoragePathByID(teamID, fileID)
		if err != nil {
			return "", err
		}

		newStoragePath = filepath.Join(newStoragePath, origFilename)
	}

	if _, err := os.Stat(newStoragePath); err == nil {
		return "", cloudErrors.ErrFileAlreadyExists
	}

	if olderStoragePath == newStoragePath {
		return newStoragePath, nil
	}

	if err := os.MkdirAll(filepath.Dir(newStoragePath), 0755); err != nil {
		return "", err
	}

	if err := os.Rename(olderStoragePath, newStoragePath); err != nil {
		return "", err
	}

	if dirFd, err := os.Open(filepath.Dir(newStoragePath)); err == nil {
		_ = dirFd.Sync()
		_ = dirFd.Close()
	}

	return newStoragePath, nil
}
