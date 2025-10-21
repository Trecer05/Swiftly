package filemanager

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	models "github.com/Trecer05/Swiftly/internal/model/chat"
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
	audio := filepath.Join(chatFolder, strconv.Itoa(id), "audios")
	files := filepath.Join(chatFolder, strconv.Itoa(id), "files")

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

	if err := os.MkdirAll(audio, os.ModePerm); err != nil {
		if os.IsExist(err) {
			return nil
		} else {
			return err
		}
	}

	if err := os.MkdirAll(files, os.ModePerm); err != nil {
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
	audio := filepath.Join(groupFolder, strconv.Itoa(id), "audios")
	files := filepath.Join(groupFolder, strconv.Itoa(id), "files")

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

	if err := os.MkdirAll(audio, os.ModePerm); err != nil {
		if os.IsExist(err) {
			return nil
		} else {
			return err
		}
	}

	if err := os.MkdirAll(files, os.ModePerm); err != nil {
		if os.IsExist(err) {
			return nil
		} else {
			return err
		}
	}

	return nil
}

func SaveMessageFiles(files []*multipart.FileHeader, id int, chatType models.ChatType, dataType models.DataType) ([]string, error) {
	var dir string
	var urls []	string
	var cht string
	var dtp string
	var folder string

	switch chatType {
	case models.TypePrivate:
		dir = filepath.Join(chatFolder, strconv.Itoa(id), "videos")
		folder = chatFolder
		cht = "chat"
	case models.TypeGroup:
		dir = filepath.Join(groupFolder, strconv.Itoa(id), "videos")
		folder = groupFolder
		cht = "group"
	}

	switch dataType {
	case models.DataTypeAud:
		dir = filepath.Join(folder, strconv.Itoa(id), "audios")
		dtp = "audio"
	case models.DataTypeDoc:
		dir = filepath.Join(folder, strconv.Itoa(id), "files")
		dtp = "file"
	case models.DataTypeImg:
		dir = filepath.Join(folder, strconv.Itoa(id), "photos")
		dtp = "photo"
	case models.DataTypeVid:
		dir = filepath.Join(folder, strconv.Itoa(id), "videos")
		dtp = "video"
	case models.DataTypeImgVid:
		dir1 := filepath.Join(folder, strconv.Itoa(id), "photos")
		dir2 := filepath.Join(folder, strconv.Itoa(id), "videos")
		dtp1 := "photo"
		dtp2 := "video"

		if err := os.MkdirAll(dir1, os.ModePerm); err != nil {
			return nil, err
		}
		if err := os.MkdirAll(dir2, os.ModePerm); err != nil {
			return nil, err
		}

		for _, file := range files {
			fl, err := file.Open()
			if err != nil {
				return nil, err
			}
			defer fl.Close()

			_, fileType, err := models.GetFileType(fl)
			if err != nil {
				return nil, err
			}

			fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)

			var filePath string
			switch fileType {
			case models.FileTypeVideo:
				filePath = filepath.Join(dir2, fileName)
			case models.FileTypeImage:
				filePath = filepath.Join(dir1, fileName)
			}
			dst, err := os.Create(filePath)
			if err != nil {
				return nil, err
			}
			defer dst.Close()

			if _, err := io.Copy(dst, fl); err != nil {
				return nil, err
			}

			switch fileType {
			case models.FileTypeVideo:
				urls = append(urls, fmt.Sprintf("/%s/%s/%s/%s", cht, strconv.Itoa(id), dtp2, fileName))
			case models.FileTypeImage:
				urls = append(urls, fmt.Sprintf("/%s/%s/%s/%s", cht, strconv.Itoa(id), dtp1, fileName))
			}
		}

		return urls, nil
	}

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	for _, file := range files {
		fl, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer fl.Close()

		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)

		filePath := filepath.Join(dir, fileName)
		dst, err := os.Create(filePath)
		if err != nil {
			return nil, err
		}
		defer dst.Close()

		if _, err := io.Copy(dst, fl); err != nil {
			return nil, err
		}

		urls = append(urls, fmt.Sprintf("/%s/%s/%s/%s", cht, strconv.Itoa(id), dtp, fileName))
	}

	return urls, nil
}
