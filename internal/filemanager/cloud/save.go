package cloud

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func SaveUserFile(reader io.Reader, handler *multipart.FileHeader, userID int) (string, string, error) {
	dir := filepath.Join(userFolder, strconv.Itoa(userID))
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
