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
	chatFolder = os.Getenv("CHATS_STORAGE")
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

func CreateChatMessagesFolder(id int) error {
	ph := filepath.Join(chatFolder, strconv.Itoa(id), "photos")
	video := filepath.Join(chatFolder, strconv.Itoa(id), "videos")
	if err := os.MkdirAll(ph, os.ModePerm); err != nil {
		if os.IsExist(err) {
			return nil
		} else {
			return err
		}
	}

	if err := os.MkdirAll(video, os.ModePerm); err != nil {
		if os.IsExist(err) {
			return nil
		} else {
			return err
		}
	}

	return nil
}

func CreateGroupMessagesFolder(id int) error {
	ph := filepath.Join(groupFolder, strconv.Itoa(id), "photos")
	video := filepath.Join(groupFolder, strconv.Itoa(id), "videos")

	if err := os.MkdirAll(ph, os.ModePerm); err != nil {
		if os.IsExist(err) {
			return nil
		} else {
			return err
		}
	}

	if err := os.MkdirAll(video, os.ModePerm); err != nil {
		if os.IsExist(err) {
			return nil
		} else {
			return err
		}
	}

	return nil
}
