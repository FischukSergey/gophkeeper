package services

import (
	"context"
	"fmt"
	"log/slog"
)

// PwdKeeper интерфейс для сервиса парольного хранилища
type PwdKeeper interface {
	GetPingDB(ctx context.Context) error
}

// GRPCService структура для сервиса
type GRPCService struct {
	log     *slog.Logger
	storage PwdKeeper
}

// NewGRPCService функция для создания сервиса
func NewGRPCService(log *slog.Logger, storage PwdKeeper) *GRPCService {
	return &GRPCService{log: log, storage: storage}
}

// Ping метод для проверки соединения с сервером
func (g *GRPCService) Ping(ctx context.Context) error {
	err := g.storage.GetPingDB(ctx)	
	if err != nil {
		return fmt.Errorf("failed to ping: %w", err)
	}
	return nil
}

