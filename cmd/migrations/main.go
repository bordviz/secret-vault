package main

import (
	"errors"
	"flag"
	"fmt"
	"vault/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var action string

	flag.StringVar(&action, "action", "", "migrations action")
	flag.Parse()

	if action != "up" && action != "down" {
		fmt.Println("action flag is required (example: --action=up)")
		return
	}

	cfg := config.MustLoad()

	fmt.Println(action)

	dsn := createDSN(
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	fmt.Println("path: ", cfg.MigrationsPath)

	m, err := migrate.New(
		fmt.Sprintf("file://%s", cfg.MigrationsPath),
		dsn,
	)

	if err != nil {
		panic(err)
	}

	if action == "up" {
		migrationsUp(m)
	}

	if action == "down" {
		migrationsDown(m)
	}

	fmt.Println("migrations applied")
}

func createDSN(user string, pass string, host string, port int, name string, sslmode string) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s&x-migrations-table=migrations",
		user,
		pass,
		host,
		port,
		name,
		sslmode,
	)
}

func migrationsUp(m *migrate.Migrate) {
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}
}

func migrationsDown(m *migrate.Migrate) {
	if err := m.Down(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}
}
