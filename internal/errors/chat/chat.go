package chat

import "errors"

var (
	ErrNoRooms = errors.New("no rooms found")
	ErrNoMessages = errors.New("no messages found")
	ErrUnknownChatType = errors.New("unknown chat type")
	ErrNoGroupFound = errors.New("no group found")

	ErrInvalidGroupData = errors.New("invalid group data")
)