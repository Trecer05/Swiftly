package main

import (
	logger "github.com/Trecer05/Swiftly/internal/config/logger"
	"os"

	env "github.com/Trecer05/Swiftly/internal/config/environment"
	migrator "github.com/Trecer05/Swiftly/internal/repository/postgres"
	mgr "github.com/Trecer05/Swiftly/internal/repository/postgres/chat"
	redis "github.com/Trecer05/Swiftly/internal/repository/cache/chat"
	server "github.com/Trecer05/Swiftly/internal/transport/http"
	router "github.com/Trecer05/Swiftly/internal/transport/http/chat"
)

func main() {
	if err := env.LoadEnvFile("./.env"); err != nil {
		logger.Logger.Fatalf("Ошибка загрузки env: %v", err)
	}
	logger.Logger.Println("ENV loaded")

	manager := mgr.NewChatManager("postgres", os.Getenv("DB_AUTH_CONNECTION_STRING"))
	logger.Logger.Println("DB connected")

	rds := redis.NewChatManager(os.Getenv("CHAT_REDIS_CONNECTION_STRING"))
	logger.Logger.Println("Redis connected")

	migrator.Migrate(manager.Conn, "chat")
	logger.Logger.Println("DB migrated")

	r := router.NewChatRouter(manager, rds)

	s := server.NewServer(os.Getenv("CHAT_SERVER_PORT"), r)
	if err := s.ListenAndServe(); err != nil {
		logger.Logger.Fatal(err)
	}
}
