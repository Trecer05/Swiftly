package main

import (
	"log"
	"os"

	env "github.com/Trecer05/Swiftly/internal/config/environment"
	migrator "github.com/Trecer05/Swiftly/internal/repository/postgres"
	mgr "github.com/Trecer05/Swiftly/internal/repository/postgres/task_tracker"
	// redis "github.com/Trecer05/Swiftly/internal/repository/cache/task_tracker"
	server "github.com/Trecer05/Swiftly/internal/transport/http"
	router "github.com/Trecer05/Swiftly/internal/transport/http/task_tracker"
)

func main() {
	if err := env.LoadEnvFile("./.env"); err != nil {
		log.Fatalf("Ошибка загрузки env: %v", err)
	}
	log.Println("ENV loaded")

	manager := mgr.NewTaskManager("postgres", os.Getenv("DB_TASK_CONNECTION_STRING"))
	log.Println("DB connected")

	// rds := redis.NewTaskManager(os.Getenv("TASK_REDIS_CONNECTION_STRING"))
	// log.Println("Redis connected")

	migrator.Migrate(manager.Conn, "task_tracker")
	log.Println("DB migrated")

	r := router.NewTaskRouter(manager)

	s := server.NewServer(os.Getenv("TASK_SERVER_PORT"), r)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
