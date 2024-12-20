package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationsPath, migrationsTable string

	flag.StringVar(&storagePath, "storage-path", "", "Path to storage")
	flag.StringVar(&migrationsPath, "migrations-path", "", "Path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migration", "name of migrations table")
	flag.Parse()

	if storagePath == "" {
		log.Fatal("storagePath is required")
	}
	if migrationsPath == "" {
		log.Fatal("migrationsPath is required")
	}

	m, err := migrate.New("file://"+migrationsPath, storagePath)
	if err != nil {
		log.Fatal(err)
	}
	defer func(m *migrate.Migrate) {
		err, _ := m.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(m)
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("migration up failed:", err)
			return
		}
		log.Fatal(err)
	}
	fmt.Println("migration up succeeded")
}
