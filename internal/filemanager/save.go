package filemanager

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	userFolder = os.Getenv("USERS_STORAGE")
	groupFolder = os.Getenv("GROUPS_STORAGE")
)

func AddUserPhoto(r *http.Request, id int) error {
	file, handler, err := r.FormFile("photo")
	if err != nil {
		return err
	}
	defer file.Close()

	uploadDir := filepath.Join(userFolder, strconv.Itoa(id), "avatars")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return err
	}

	fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), handler.Filename)

	dst, err := os.Create(filepath.Join(uploadDir, fileName))
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return err
	}
	return nil
}

func AddGroupPhoto(r *http.Request, id int) error {
	file, handler, err := r.FormFile("photo")
	if err != nil {
		return err
	}
	defer file.Close()

	uploadDir := filepath.Join(groupFolder, strconv.Itoa(id), "avatars")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return err
	}

	fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), handler.Filename)

	dst, err := os.Create(filepath.Join(uploadDir, fileName))
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return err
	}
	return nil
}
