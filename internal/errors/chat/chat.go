package chat

import "errors"

var (
	ErrNoRooms = errors.New("no rooms found")
	ErrNoMessages = errors.New("no messages found")
	ErrUnknownChatType = errors.New("unknown chat type")
	ErrNoGroupFound = errors.New("no group found")

	ErrInvalidGroupData = errors.New("invalid group data")

	ErrUserAlreadyInGroup = errors.New("user already in group")
	ErrUserNotInGroup = errors.New("user not in group")

	ErrInvalidLimit = errors.New("invalid limit")
	ErrInvalidOffset = errors.New("invalid offset")

	ErrNoUser = errors.New("user not found")
	ErrNoUsers = errors.New("no users found")
	ErrUserNotAOwner = errors.New("user not a owner")
	ErrUserAlreadyInTeam = errors.New("user already in team")

	ErrNoData = errors.New("no data found")
)

var (
	ErrorNoCallRoom = errors.New("no room found")
)