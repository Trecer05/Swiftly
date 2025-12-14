package main

import (
	logger "github.com/Trecer05/Swiftly/internal/config/logger"
	"os"
	"log"
	"context"

	kfk "github.com/Trecer05/Swiftly/internal/repository/kafka/cloud"
	env "github.com/Trecer05/Swiftly/internal/config/environment"
	migrator "github.com/Trecer05/Swiftly/internal/repository/postgres"
	mgr "github.com/Trecer05/Swiftly/internal/repository/postgres/cloud"
	server "github.com/Trecer05/Swiftly/internal/transport/http"
	router "github.com/Trecer05/Swiftly/internal/transport/http/cloud"
	"github.com/Trecer05/Swiftly/internal/filemanager/cloud"
)

var ctx = context.Background()

func main() {
	if _, err := os.Stat("./.env"); err == nil {
	    env.LoadEnvFile("./.env")
	    log.Println("ENV loaded from .env")
	} else {
	    log.Println("Running inside Docker – using environment variables")
	}
	log.Println("ENV loaded")
	
	logger.LogInit()
	logger.Logger.Info("Логгер запущен без ошибок")

	manager := mgr.NewCloudManager("postgres", os.Getenv("DB_CLOUD_CONNECTION_STRING"))
	logger.Logger.Println("DB connected")

	cloud.CreateStartDirs()
	logger.Logger.Println("Start dirs created")

	// rds := redis.NewCloudManager(os.Getenv("CLOUD_REDIS_CONNECTION_STRING"))
	// logger.Logger.Println("Redis connected")
	
	kfMgr := kfk.NewKafkaManager([]string{os.Getenv("KAFKA_ADDRESS")}, "cloud", "cloud-team")
	
	go kfMgr.ReadChatMessages(ctx)

	migrator.Migrate(manager.Conn, "cloud")
	logger.Logger.Println("DB migrated")

	r := router.NewCloudRouter(manager)

	s := server.NewServer(os.Getenv("CLOUD_SERVER_PORT"), r)
	if err := s.ListenAndServe(); err != nil {
		logger.Logger.Fatal(err)
	}
}
