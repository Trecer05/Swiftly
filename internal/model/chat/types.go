package chat

type MessageType string

const (
	Typing      MessageType = "typing"
	StopTyping  MessageType = "stop_typing"
	Default     MessageType = "message"
)