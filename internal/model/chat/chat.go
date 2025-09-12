package chat

import (
	"time"
)

type Message struct {
	ID     int    `json:"id"`
	ChatID int    `json:"chat_id"`
	Type   MessageType `json:"type"`
	Read   bool `json:"read"`
	Text   string    `json:"text"`
	Status Status `json:"status"`
	Author Client    `json:"author"`
	Time   time.Time `json:"time"`
	FileURL   *string    `json:"file_url,omitempty"`
	FileUrls  []string `json:"file_urls,omitempty"`
    FileName  string    `json:"file_name"`
    FileMIME  string    `json:"file_mime"`
	FileType  FileType `json:"file_type"`
    FileSize  int64     `json:"file_size"`
}

type Status struct {
	Type string `json:"type"`
	User_ID int `json:"user_id"`
	Online bool `json:"online"`
}

type ChatRoom struct {
	ID int `json:"id"`
	Name string `json:"name"`
	LastMessage *Message `json:"last_message"`
	Type ChatType `json:"type"`
}

type Client struct {
	ID    int	`json:"id"`
	Name  string `json:"name"`
}

type Group struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
}
