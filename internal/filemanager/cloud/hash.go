package cloud

import (
	"crypto/sha256"
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

func FileSHA256FromFiles(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

func FileSHA256FromMultipartFile(file *multipart.File) (string, error) {
	if file != nil {
		hasher := sha256.New()
		if _, err := io.Copy(hasher, *file); err != nil {
			return "", err
		}
		return fmt.Sprintf("%x", hasher.Sum(nil)), nil
	}

	return "", nil
}

func HashingReader(r io.Reader) (io.Reader, func() string) {
	hasher := sha256.New()

	reader := io.TeeReader(r, hasher)

	getHash := func() string {
		return fmt.Sprintf("%x", hasher.Sum(nil))
	}

	return reader, getHash
}
