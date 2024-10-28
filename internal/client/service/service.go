package service

import (
	"context"
	"log/slog"

	pb "github.com/FischukSergey/gophkeeper/internal/proto"
)

type AuthService struct {
	client pb.GophKeeperClient
	log    *slog.Logger
}

func NewAuthService(client pb.GophKeeperClient, log *slog.Logger) *AuthService {
	return &AuthService{client: client, log: log}
}


// Register implements command.IRegisterService.
// Регистрация нового клиента
func (s *AuthService) Register(ctx context.Context, login string, password string) (string, error) {
		token, err := s.client.Registration(ctx, &pb.RegistrationRequest{	
		Username: login,
		Password: password,
	})
	if err != nil {
		return "", err
	}
	s.log.Debug("регистрация нового клиента", "login", login, "password", password)
	s.log.Debug("токен", "token", token)
	return token.GetAccessToken().Token, nil
}

// Check проверка работоспособности сервера
func (s *AuthService) Check(ctx context.Context) error {
	_, err := s.client.Ping(ctx, &pb.PingRequest{})
	return err
}

// Authorization имплементирует интерфейс command.IAuthService
// Authorization авторизация клиента
func (s *AuthService) Authorization(ctx context.Context, login, password string) (string, error) {
	token, err := s.client.Authorization(ctx, &pb.AuthorizationRequest{
		Username: login,
		Password: password,
	})
	if err != nil {
		return "", err
	}
	s.log.Debug("авторизация клиента", "login", login, "password", password)
	s.log.Debug("токен", "token", token)
	return token.GetAccessToken().Token, nil
}
