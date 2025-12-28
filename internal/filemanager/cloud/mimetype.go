package cloud

import (
	"bytes"
	"io"
	"net/http"
)

func DetectMimeType(r io.Reader) (string, io.Reader, error) {
	buf := make([]byte, 512)

	n, err := io.ReadFull(r, buf)
	if err != nil && err != io.ErrUnexpectedEOF {
		return "", nil, err
	}

	mimeType := http.DetectContentType(buf[:n])

	reader := io.MultiReader(bytes.NewReader(buf[:n]), r)

	return mimeType, reader, nil
}