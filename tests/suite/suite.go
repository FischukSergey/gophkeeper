package suite

import (
	"context"
	"fmt"

	"github.com/FischukSergey/gophkeeper/internal/storage/dbstorage"
	"github.com/jackc/pgx/v5/pgxpool"
)


func InitTestStorage() (*dbstorage.Storage) {
	dbconn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		"test", "test", "localhost", "5433", "test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := pgxpool.New(ctx, dbconn)
	if err != nil {
		panic(fmt.Errorf("%w, unable to create connection db:%s", err, "test"))
	}
	err = db.Ping(ctx)
	if err != nil {
		panic(fmt.Errorf("error pinging database: %w", err))
	}
	return &dbstorage.Storage{DB: db}
}