package filemanager

import (
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