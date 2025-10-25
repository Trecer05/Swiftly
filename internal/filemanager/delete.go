package filemanager

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func DeleteUserProfilePhotos(id int) error {
	path := filepath.Join(userFolder, strconv.Itoa(id), "avatars")

	if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil
    }

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			filePath := filepath.Join(path, entry.Name())
            if err := os.Remove(filePath); err != nil {
                return err
            }
		}
	}

	return nil
}

func DeleteGroupPhoto(id int) error {
	path := filepath.Join(groupFolder, strconv.Itoa(id), "avatars")

	if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil
    }

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			filePath := filepath.Join(path, entry.Name())
            if err := os.Remove(filePath); err != nil {
                return err
            }
		}
	}

	return nil
}

func deleteFilesFromFolder(baseFolder string, id int) error {
    path := filepath.Join(baseFolder, strconv.Itoa(id))
    
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil
    }
    
    err := os.RemoveAll(path)
    if err != nil {
        return fmt.Errorf("failed to delete files from %s: %w", baseFolder, err)
    }
    
    return nil
}

func DeleteGroupFiles(id int) error {
    return deleteFilesFromFolder(groupFolder, id)
}

func DeleteChatFiles(id int) error {
    return deleteFilesFromFolder(chatFolder, id)
}

func DeleteUserAvatar(url string, userId int) error {
	url = filepath.Base(url)
	userPhotoDir := filepath.Join(userFolder, strconv.Itoa(userId), "avatars", url)

	if err := os.Remove(userPhotoDir); err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}

	return nil
}

func DeleteGroupAvatar(url string, groupId int) error {
	url = filepath.Base(url)
	groupPhotoDir := filepath.Join(groupFolder, strconv.Itoa(groupId), "avatars", url)

	if err := os.Remove(groupPhotoDir); err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}

	return nil
}
