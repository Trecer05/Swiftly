package filemanager

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	fileErrors "github.com/Trecer05/Swiftly/internal/errors/file"
	models "github.com/Trecer05/Swiftly/internal/model/chat"
)

func GetFile(url string, id int, ct models.ChatType, ft models.DataType) ([]byte, string, error) {
	var dir string

	switch ct {
	case models.TypePrivate:
		dir = filepath.Join(groupFolder, strconv.Itoa(id))
	case models.TypeGroup:
		dir = filepath.Join(userFolder, strconv.Itoa(id))
	}

	switch ft {
	case models.DataTypeImg:
		return getPhoto(url, dir)
	case models.DataTypeVid:
		return getVideo(url, dir)
	case models.DataTypeAud:
		return getAudio(url, dir)
	case models.DataTypeDoc:
		return getDoc(url, dir)
	}

	return nil, "", fileErrors.ErrFileNotFound
}

func GetMedias(id int, ct models.ChatType) ([]string, error) {
	var dir string
	var files []string

	switch ct {
	case models.TypePrivate:
		dir = filepath.Join(groupFolder, strconv.Itoa(id))
	case models.TypeGroup:
		dir = filepath.Join(userFolder, strconv.Itoa(id))
	}

	photoDir := filepath.Join(dir, "photos")
	videoDir := filepath.Join(dir, "videos")

	photos, err := os.ReadDir(photoDir)
	if err != nil {
		return nil, err
	}
	videos, err := os.ReadDir(videoDir)
	if err != nil {
		return nil, err
	}

	for _, photo := range photos {
		files = append(files, photo.Name())
	}

	for _, video := range videos {
		files = append(files, video.Name())
	}

	return files, nil
}

func GetDocs(id int, ct models.ChatType) ([]byte, string, error) {

	return nil, "", fileErrors.ErrFilesNotFound
}

func getPhoto(url string, dir string) ([]byte, string, error) {
	dir = filepath.Join(dir, "photos", url)

	file, err := os.ReadFile(dir)
	if err != nil {
		return nil, "", err
	}

	fileType := http.DetectContentType(file)
	
	return file, fileType, nil
}

func getVideo(url string, dir string) ([]byte, string, error) {
	dir = filepath.Join(dir, "videos", url)

	file, err := os.ReadFile(dir)
	if err != nil {
		return nil, "", err
	}

	fileType := http.DetectContentType(file)
	
	return file, fileType, nil
}

func getAudio(url string, dir string) ([]byte, string, error) {
	dir = filepath.Join(dir, "photos", url)

	file, err := os.ReadFile(dir)
	if err != nil {
		return nil, "", err
	}

	fileType := http.DetectContentType(file)
	
	return file, fileType, nil
}

func getDoc(url string, dir string) ([]byte, string, error) {
	dir = filepath.Join(dir, "audios", url)

	file, err := os.ReadFile(dir)
	if err != nil {
		return nil, "", err
	}

	fileType := http.DetectContentType(file)
	
	return file, fileType, nil
}
