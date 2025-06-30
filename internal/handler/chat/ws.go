package chat

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	chatModels "github.com/Trecer05/Swiftly/internal/model/chat"
	wsChat "github.com/Trecer05/Swiftly/internal/transport/websocket/chat"

	"github.com/gorilla/websocket"
)

func HandleConnection(w http.ResponseWriter, r *http.Request) {
	var cl chatModels.Client
	// тут получение айди и имени из сессии будет

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	chat := wsChat.NewChat()
	cl.Conn = conn

	chat.Lock.Lock()
	chat.Clients[&cl] = true
	chat.Lock.Unlock()

	chat.WG.Add(1)
	go func() {
		defer chat.WG.Done()
		for {
			_, msgBytes, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				break
			}

			var msg chatModels.Message
			err = json.Unmarshal(msgBytes, &msg)
			if err != nil {
				log.Println(err)
				continue
			}
			msg.Time = time.Now()

			chat.Broadcaster <- msg
		}

		chat.Lock.Lock()
		delete(chat.Clients, &cl)
		chat.Lock.Unlock()
	}()

	chat.WG.Wait()
}
