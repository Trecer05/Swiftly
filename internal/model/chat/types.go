package chat

type MessageType string

const (
	Typing      MessageType = "typing"
	StopTyping  MessageType = "stop_typing"
	Default     MessageType = "message"
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
