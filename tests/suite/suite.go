package suite

import (
	"context"
	"fmt"

	"github.com/FischukSergey/gophkeeper/internal/storage/dbstorage"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbUser     = "test"
	dbPassword = "test"
	dbHost     = "localhost"
	dbPort     = "5433"
	dbName     = "test"
)

// InitTestStorage инициализация тестового хранилища.
func InitTestStorage() (*dbstorage.Storage, error) {
	dbconn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := pgxpool.New(ctx, dbconn)
	if err != nil {
		return nil, fmt.Errorf("%w, unable to create connection db:%s", err, "test")
	}
	err = db.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}
	return &dbstorage.Storage{DB: db}, nil
}
