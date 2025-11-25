package main

import (
	logger "github.com/Trecer05/Swiftly/internal/config/logger"
	"os"
	"log"
	"context"

	env "github.com/Trecer05/Swiftly/internal/config/environment"
	migrator "github.com/Trecer05/Swiftly/internal/repository/postgres"
	mgr "github.com/Trecer05/Swiftly/internal/repository/postgres/task_tracker"
	kfk "github.com/Trecer05/Swiftly/internal/repository/kafka/task_tracker"
	// redis "github.com/Trecer05/Swiftly/internal/repository/cache/task_tracker"
	server "github.com/Trecer05/Swiftly/internal/transport/http"
	router "github.com/Trecer05/Swiftly/internal/transport/http/task_tracker"
)

var ctx context.Context

func main() {
	if err := env.LoadEnvFile("./.env"); err != nil {
		log.Fatalf("Ошибка загрузки env: %v", err)
	}
	log.Println("ENV loaded")
	
	logger.LogInit()
	logger.Logger.Info("Логгер запущен без ошибок")

	manager := mgr.NewTaskManager("postgres", os.Getenv("DB_TASK_CONNECTION_STRING"))
	logger.Logger.Println("DB connected")

	// rds := redis.NewTaskManager(os.Getenv("TASK_REDIS_CONNECTION_STRING"))
	// logger.Logger.Println("Redis connected")
	
	kfMgr := kfk.NewKafkaManager([]string{os.Getenv("KAFKA_ADDRESS")}, "team", "team-user-tasks")
	
	go kfMgr.ReadChatMessages(ctx, manager)

	migrator.Migrate(manager.Conn, "task_tracker")
	logger.Logger.Println("DB migrated")

	r := router.NewTaskRouter(manager)

	s := server.NewServer(os.Getenv("TASK_SERVER_PORT"), r)
	if err := s.ListenAndServe(); err != nil {
		logger.Logger.Fatal(err)
	}
}
