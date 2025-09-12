package chat

import (
	"io"
	"mime/multipart"
	"net/http"
)

type MessageType string

const (
	Typing      MessageType = "typing"
	StopTyping  MessageType = "stop_typing"
	Default     MessageType = "message"
	WithFiles   MessageType = "with_files"
	LastMessage MessageType = "last_message"
	Read        MessageType = "read"
)

type ChatType string

const (
	TypePrivate ChatType = "private"
	TypeGroup   ChatType = "group"
)

type DBType string

const (
	DBChat  DBType = "chat"
	DBGroup DBType = "group"
)

type SessionKey struct {
	Type ChatType
	ID   int
}

type DataType string

const (
	DataTypeImg    DataType = "img"
	DataTypeDoc    DataType = "doc"
	DataTypeVid    DataType = "vid"
	DataTypeAud    DataType = "aud"
	DataTypeImgVid DataType = "imgvid"
)

type FileType string

const (
	FileTypeImage    FileType = "photo"
	FileTypeVideo    FileType = "video"
	FileTypeDocument FileType = "file"
	FileTypeAudio    FileType = "audio"
	FileTypeOther    FileType = "other"
)

var mimeTypeCategories = map[FileType][]string{
	FileTypeImage: {
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/webp",
		"image/svg+xml",
		"image/bmp",
	},
	FileTypeVideo: {
		"video/mp4",
		"video/mpeg",
		"video/quicktime",
		"video/webm",
		"video/x-msvideo",
		"video/x-ms-wmv",
	},
	FileTypeAudio: {
		"audio/mpeg",
		"audio/wav",
		"audio/ogg",
		"audio/webm",
		"audio/aac",
	},
	FileTypeDocument: {
		"application/pdf",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"text/plain",
		"text/csv",
		"application/zip",
		"application/x-rar-compressed",
	},
}

func getFileTypeFromMIME(mimeType string) FileType {
	for category, mimes := range mimeTypeCategories {
		for _, mime := range mimes {
			if mime == mimeType {
				return category
			}
		}
	}
	return FileTypeOther
}

func GetFileType(file multipart.File) (string, FileType, error) {
	buff := make([]byte, 512)

	n, err := file.Read(buff)
	if err != nil && err != io.EOF {
		return "", FileTypeOther, err
	}

	mimeType := http.DetectContentType(buff[:n])

	fileKind := getFileTypeFromMIME(mimeType)

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "", FileTypeOther, err
	}

	return mimeType, fileKind, nil
}
