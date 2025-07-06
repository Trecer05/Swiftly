package chat

import "errors"

var (
	ErrNoChats = errors.New("no chats found")
	ErrNoMessages = errors.New("no messages found")
)