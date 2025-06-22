package chat

import (
	"log"

	"github.com/Trecer05/Swiftly/internal/model/chat"
)

func NewChat() *chat.Chat {
	return &chat.Chat{
		Broadcaster: make(chan chat.Message),
		Clients:     make(map[*chat.Client]bool),
		ErrCh:       make(chan error),
	}
}

func Broadcast(ch *chat.Chat) {
	for {
		msg := <-ch.Broadcaster

		ch.Lock.Lock()
		for client := range ch.Clients {
			err := client.Conn.WriteJSON(msg)
			if err != nil {
				log.Printf("error sending message: %v", err)
				ch.ErrCh <- err
				continue
			}
		}
		ch.Lock.Unlock()
	}
}