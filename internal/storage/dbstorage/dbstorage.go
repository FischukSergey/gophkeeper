package dbstorage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)
type Storage struct {
	DB *pgxpool.Pool
}

// Ping проверка соединения с базой данных
func (s *Storage) GetPingDB(ctx context.Context) error {
	err := s.DB.Ping(ctx)
	if err != nil {
		return fmt.Errorf("error pinging database: %w", err)
	}
	return nil
}

// Close закрытие подключения к базе данных
func (s *Storage) Close() {
	s.DB.Close()	
}		