package main

import (
	"context"
	logger "github.com/Trecer05/Swiftly/internal/config/logger"
	"os"

	env "github.com/Trecer05/Swiftly/internal/config/environment"
	kafka "github.com/Trecer05/Swiftly/internal/repository/kafka/auth"
	migrator "github.com/Trecer05/Swiftly/internal/repository/postgres"
	mgr "github.com/Trecer05/Swiftly/internal/repository/postgres/auth"
	server "github.com/Trecer05/Swiftly/internal/transport/http"
	router "github.com/Trecer05/Swiftly/internal/transport/http/auth"
)

var ctx = context.Background()

func main() {
	if err := env.LoadEnvFile("./.env"); err != nil {
		logger.Logger.Fatalf("Ошибка загрузки env: %v", err)
	}
	logger.Logger.Println("ENV loaded")

	manager := mgr.NewAuthManager("postgres", os.Getenv("DB_AUTH_CONNECTION_STRING"))
	logger.Logger.Println("DB connected")

	migrator.Migrate(manager.Conn, "auth")
	logger.Logger.Println("DB migrated")

	kfk := kafka.NewKafkaManager([]string{os.Getenv("KAFKA_ADDRESS")}, "profile", "user-change-group")

	go kfk.ReadUserEditMessages(ctx, manager)
	
	defer kfk.Close()

	r := router.NewAuthRouter(manager)

	s := server.NewServer(os.Getenv("AUTH_SERVER_PORT"), r)
	if err := s.ListenAndServe(); err != nil {
		logger.Logger.Fatal(err)
	}
}
