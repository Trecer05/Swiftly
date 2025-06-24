package postgres

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(conn *sql.DB, serviceName string) {
	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		panic(fmt.Sprintf("Не удалось создать драйвер миграции: %v", err))
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://db//migrations/%s", serviceName), //- нужно указать норм путь к папке с миграциями
		"postgres",
		driver,
	)
	if err != nil {
		panic(fmt.Sprintf("Не удалось создать мигратора: %v", err))
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(fmt.Sprintf("Не удалось применить миграции: %v", err))
	}
}