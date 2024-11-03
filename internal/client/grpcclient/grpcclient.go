package grpcclient

import (
	"fmt"
	"log/slog"

	"github.com/FischukSergey/gophkeeper/internal/client/config"
	pb "github.com/FischukSergey/gophkeeper/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Token структура для токена.
type Token struct {
	Token string
}

// GetToken возвращает токен.
func (t *Token) GetToken() string {
	return t.Token
}

// NewClient создание клиента grpc.
func NewClient(cfg *config.Config, log *slog.Logger) (*grpc.ClientConn, pb.GophKeeperClient, error) {
	log.Info("server address", "address", cfg.ServerAddress)

	conn, err := grpc.NewClient(cfg.ServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create grpc client: %w", err)
	}
	log.Info("connected to server")

	client := pb.NewGophKeeperClient(conn)
	return conn, client, nil
}
