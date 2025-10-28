package filemanager

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	fileErrors "github.com/Trecer05/Swiftly/internal/errors/file"
	models "github.com/Trecer05/Swiftly/internal/model/chat"
)

func GetFile(url string, id int, ct models.ChatType, ft models.DataType) ([]byte, string, error) {
	url = filepath.Base(url)
	var dir string

	switch ct {
	case models.TypePrivate:
		dir = filepath.Join(chatFolder, strconv.Itoa(id))
	case models.TypeGroup:
		dir = filepath.Join(groupFolder, strconv.Itoa(id))
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
		dir = filepath.Join(chatFolder, strconv.Itoa(id))
	case models.TypeGroup:
		dir = filepath.Join(groupFolder, strconv.Itoa(id))
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

func getFileHelper(url string, dir string) ([]byte, string, error) {
	file, err := os.ReadFile(dir)
	if err != nil {
		return nil, "", err
	}

	fileType := http.DetectContentType(file)
	
	return file, fileType, nil
}

func getPhoto(url string, dir string) ([]byte, string, error) {
	url = filepath.Base(url)
	dir = filepath.Join(dir, "photos", url)

	return getFileHelper(url, dir)
}

func getVideo(url string, dir string) ([]byte, string, error) {
	url = filepath.Base(url)
	dir = filepath.Join(dir, "videos", url)

	return getFileHelper(url, dir)
}

func getAudio(url string, dir string) ([]byte, string, error) {
	url = filepath.Base(url)
	dir = filepath.Join(dir, "audios", url)

	return getFileHelper(url, dir)
}

func getDoc(url string, dir string) ([]byte, string, error) {
	url = filepath.Base(url)
	dir = filepath.Join(dir, "files", url)
	
	return getFileHelper(url, dir)
}

func GetUserPhotosUrls(userId int) ([]string, error) {
	userDir := filepath.Join(userFolder, strconv.Itoa(userId), "avatars")
	var names []string

	dir, err := os.ReadDir(userDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	for _, entry := range dir {
		names = append(names, entry.Name())
	}

	return names, nil
}

func GetUserPhotoByUrl(userId int, url string) ([]byte, string, error) {
	url = filepath.Base(url)
	dir := filepath.Join(userFolder, strconv.Itoa(userId), "avatars", url)

	return getFileHelper(url, dir)
}

func GetLatestGroupAvatarUrl(groupId int) (string, error) {
	dir := filepath.Join(groupFolder, strconv.Itoa(groupId), "avatars")

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	var latestFile string
	var latestTime int64

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()
		ext := strings.ToLower(filepath.Ext(name))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
			continue
		}

		parts := strings.SplitN(name, "_", 2)
		if len(parts) < 2 {
			continue
		}

		timestamp, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			continue
		}

		if timestamp > latestTime {
			latestTime = timestamp
			latestFile = name
		}
	}

	if latestFile == "" {
		return "", nil
	}

	return latestFile, nil
}

func GetGroupPhotoByUrl(groupId int, url string) ([]byte, string, error) {
	url = filepath.Base(url)
	photoUrl := filepath.Join(groupFolder, strconv.Itoa(groupId), "avatars", url)

	return getFileHelper(url, photoUrl)
}

func GetAudioMessageByUrl(id int, url string, chatType models.ChatType) ([]byte, string, error) {
	if chatType == models.TypePrivate {
		return getFileHelper(url, filepath.Join(chatFolder, strconv.Itoa(id), "messages", "audio", url))
	} else {
		return getFileHelper(url, filepath.Join(groupFolder, strconv.Itoa(id), "messages", "audio", url))
	}
}

func GetVideoMessageByUrl(id int, url string, chatType models.ChatType) ([]byte, string, error) {
	if chatType == models.TypePrivate {
		return getFileHelper(url, filepath.Join(chatFolder, strconv.Itoa(id), "messages", "video", url))
	} else {
		return getFileHelper(url, filepath.Join(groupFolder, strconv.Itoa(id), "messages", "video", url))
	}
}
