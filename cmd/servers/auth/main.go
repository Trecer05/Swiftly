package main

import (
	"log"
	"os"

	env "github.com/Trecer05/Swiftly/internal/config/environment"
	migrator "github.com/Trecer05/Swiftly/internal/repository/postgres"
	mgr "github.com/Trecer05/Swiftly/internal/repository/postgres/auth"
	server "github.com/Trecer05/Swiftly/internal/transport/http"
	router "github.com/Trecer05/Swiftly/internal/transport/http/auth"
)

func main() {
	if err := env.LoadEnvFile(".env"); err != nil {
		log.Fatalf("Ошибка загрузки env: %v", err)
	}
	
	manager := mgr.NewAuthManager("postgres", os.Getenv("DB_AUTH_CONNECTION_STRING"))

	migrator.Migrate(manager.Conn, "auth")

	r := router.NewAuthRouter(manager)

	s := server.NewServer(os.Getenv("AUTH_SERVER_PORT"), r)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}