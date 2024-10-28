package service

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/FischukSergey/gophkeeper/internal/proto"
)

// AuthService сервис авторизации.
type AuthService struct {
	client pb.GophKeeperClient
	log    *slog.Logger
}

// NewAuthService создание сервиса авторизации.
func NewAuthService(client pb.GophKeeperClient, log *slog.Logger) *AuthService {
	return &AuthService{client: client, log: log}
}

// Register регистрация нового клиента.
func (s *AuthService) Register(ctx context.Context, login string, password string) (string, error) {
	token, err := s.client.Registration(ctx, &pb.RegistrationRequest{
		Username: login,
		Password: password,
	})
	if err != nil {
		return "", fmt.Errorf("failed to register: %w", err)
	}
	s.log.Debug("регистрация нового клиента", "login", login, "password", password)
	s.log.Debug("токен", "token", token)
	return token.GetAccessToken().Token, nil
}

// Check проверка работоспособности сервера.
func (s *AuthService) Check(ctx context.Context) error {
	_, err := s.client.Ping(ctx, &pb.PingRequest{})
	if err != nil {
		return fmt.Errorf("failed to check server: %w", err)
	}
	return nil
}

// Authorization авторизация клиента.
func (s *AuthService) Authorization(ctx context.Context, login, password string) (string, error) {
	token, err := s.client.Authorization(ctx, &pb.AuthorizationRequest{
		Username: login,
		Password: password,
	})
	if err != nil {
		return "", fmt.Errorf("failed to authorization: %w", err)
	}
	s.log.Debug("авторизация клиента", "login", login, "password", password)
	s.log.Debug("токен", "token", token)
	return token.GetAccessToken().Token, nil
}
