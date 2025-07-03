package main

import (
	"log"
	"os"

	env "github.com/Trecer05/Swiftly/internal/config/environment"
	migrator "github.com/Trecer05/Swiftly/internal/repository/postgres"
	mgr "github.com/Trecer05/Swiftly/internal/repository/postgres/chat"
	server "github.com/Trecer05/Swiftly/internal/transport/http"
	router "github.com/Trecer05/Swiftly/internal/transport/http/chat"
)

func main() {
	if err := env.LoadEnvFile("./.env"); err != nil {
		log.Fatalf("Ошибка загрузки env: %v", err)
	}
	log.Println("ENV loaded")

	manager := mgr.NewChatManager("postgres", os.Getenv("DB_AUTH_CONNECTION_STRING"))
	log.Println("DB connected")

	migrator.Migrate(manager.Conn, "chat")
	log.Println("DB migrated")

	r := router.NewChatRouter(manager)

	s := server.NewServer(os.Getenv("CHAT_SERVER_PORT"), r)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
