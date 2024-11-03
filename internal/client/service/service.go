package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/FischukSergey/gophkeeper/internal/app/interceptors/auth"
	"github.com/FischukSergey/gophkeeper/internal/models"
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

// S3FileUpload загрузка файла на сервер.
func (s *AuthService) S3FileUpload(
	ctx context.Context,
	token string,
	fileData []byte,
	filename string,
) (string, error) {
	// добавление токена авторизации в контекст
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("session_token", token))
	// проверка наличия файла в S3
	files, err := s.GetFileList(ctx, token)
	if err == nil || len(files) > 0 {
		for _, file := range files {
			if file.Filename == filename {
				// файл уже существует спрашиваем пользователя хочет ли он его перезаписать
				fmt.Println("Файл уже существует. Хотите перезаписать? (y/n)")
				var input string
				fmt.Scanln(&input)
				if input != "y" && input != "Y" {
					return "", fmt.Errorf("file already exists")
				}
			}
		}
	}
	// загрузка файла на сервер
	response, err := s.client.FileUpload(ctx, &pb.FileUploadRequest{
		Filename: filename,
		Data:     fileData,
	})
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			return "", fmt.Errorf("запрос отменен: %w", err)
		case errors.Is(err, status.Error(codes.Unauthenticated, "user ID not found in context")):
			return "", fmt.Errorf("не авторизован: %w", err)
		case errors.Is(err, status.Error(codes.Unauthenticated, auth.ErrNotFound)):
			return "", fmt.Errorf("токен не найден: %w", err)
		case errors.Is(err, status.Error(codes.Unauthenticated, auth.ErrInvalid)):
			return "", fmt.Errorf("токен не валиден: %w", err)
		case errors.Is(err, status.Error(codes.Unauthenticated, auth.ErrExpired)):
			return "", fmt.Errorf("токен просрочен: %w", err)
		default:
			s.log.Error("ошибка загрузки файла", "error", err)
			return "", fmt.Errorf("failed to upload file: %w", err)
		}
	}
	s.log.Debug("файл загружен", "filename", response.GetMessage())
	return response.GetMessage(), nil
}

// GetFileList получение списка файлов.
func (s *AuthService) GetFileList(ctx context.Context, token string) ([]models.File, error) {
	// добавление токена авторизации в контекст
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("session_token", token))

	// получение списка файлов
	response, err := s.client.FileGetList(ctx, &pb.FileGetListRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file list: %w", err)
	}
	files := make([]models.File, 0, len(response.GetFiles()))
	for _, file := range response.GetFiles() {
		files = append(files, models.File{
			FileID:    file.GetFileID(),
			UserID:    file.GetUserID(),
			Filename:  file.GetFilename(),
			CreatedAt: file.GetCreatedAt().AsTime(),
			DeletedAt: file.GetDeletedAt().AsTime(),
			Size:      file.GetSize(),
		})
	}
	return files, nil
}
