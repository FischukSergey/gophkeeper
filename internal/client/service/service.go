package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/FischukSergey/gophkeeper/internal/app/interceptors/auth"
	"github.com/FischukSergey/gophkeeper/internal/logger"
	"github.com/FischukSergey/gophkeeper/internal/models"
	pb "github.com/FischukSergey/protos/gen/gophkeeper/gophkeeper"
)

const (
	chunkSize    = 1024
	sessionToken = "session_token"
	percent      = 100
	percentStep  = 5
	kb           = 1024
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
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(sessionToken, token))
	// проверка наличия файла в S3
	files, err := s.GetFileList(ctx, token)
	if err == nil || len(files) > 0 {
		for _, file := range files {
			if file.Filename == filename {
				// файл уже существует спрашиваем пользователя хочет ли он его перезаписать
				fmt.Println("Файл уже существует. Хотите перезаписать? (y/n)")
				var input string
				_, err = fmt.Scanln(&input)
				if err != nil {
					fmt.Printf("Ошибка ввода: %s\n", err)
				}
				if input != "y" && input != "Y" {
					return "", fmt.Errorf("file already exists")
				}
			}
		}
	}
	//открываем стрим для загрузки файла
	stream, err := s.client.FileUpload(ctx)
	if err != nil {
		s.log.Error("ошибка открытия стрима", logger.Err(err))
		return "", fmt.Errorf("failed to open stream: %w", err)
	}
	//сначала отправляем имя файла
	err = stream.Send(&pb.FileUploadRequest{
		File: &pb.FileUploadRequest_Info{
			Info: &pb.FileInfo{
				Filename: filename,
				Size:     int64(len(fileData)),
			},
		},
	})
	if err != nil {
		s.log.Error("ошибка отправки файла", logger.Err(err))
		return "", fmt.Errorf("failed to send file info: %w", err)
	}
	totalSize := len(fileData)
	var lastPercent int
	_, _ = fmt.Printf("Uploading %s (%d bytes):\n", filename, totalSize)

	// отправляем файл частями
	for i := 0; i < totalSize; i += chunkSize {
		end := i + chunkSize
		if end > totalSize {
			end = totalSize
		}
		err = stream.Send(&pb.FileUploadRequest{
			File: &pb.FileUploadRequest_Chunk{
				Chunk: fileData[i:end],
			},
		})
		if err != nil {
			fmt.Print("\n")
			s.log.Error("ошибка загрузки файла", logger.Err(err))
			return "", fmt.Errorf("failed to send chunk: %w", err)
		}
		// Показываем процент каждые 5%
		currentPercent := (i * percent) / totalSize
		if currentPercent >= lastPercent+percentStep {
			fmt.Printf("\r%d%%", currentPercent)
			lastPercent = currentPercent
		}
	}
	// Закрываем стрим и получаем ответ
	response, err := stream.CloseAndRecv()
	if err != nil && !errors.Is(err, io.EOF) {
		fmt.Print("\n")
		s.log.Error("ошибка закрытия стрима", logger.Err(err))
		return "", fmt.Errorf("failed to close stream: %w", err)
	}

	fmt.Printf("\r100%%\nUpload complete: %s\n", filename)

	// Проверяем response только если он не nil
	if response != nil {
		s.log.Debug("файл загружен", "filename", filename, "size", totalSize)
		return response.GetMessage(), nil
	} else {
		s.log.Debug("файл загружен", "filename", filename, "size", totalSize)
	}
	return "", nil
}

// GetFileList получение списка файлов.
func (s *AuthService) GetFileList(ctx context.Context, token string) ([]models.File, error) {
	// добавление токена авторизации в контекст
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(sessionToken, token))

	// получение списка файлов
	response, err := s.client.FileGetList(ctx, &pb.FileGetListRequest{})
	if err != nil {
		switch {
		case errors.Is(err, status.Error(codes.Unauthenticated, auth.ErrNotFound)):
			return nil, fmt.Errorf("токен не найден: %w", err)
		default:
			return nil, fmt.Errorf("failed to get file list: %w", err)
		}
	}
	files := make([]models.File, 0, len(response.GetFiles()))
	for _, file := range response.GetFiles() {
		files = append(files, ProtoToModel(file))
	}
	return files, nil
}

// S3FileDelete удаление файла.
func (s *AuthService) S3FileDelete(ctx context.Context, token string, filename string) error {
	// добавление токена авторизации в контекст
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(sessionToken, token))
	// удаление файла
	_, err := s.client.FileDelete(ctx, &pb.FileDeleteRequest{
		Filename: filename,
	})
	if err != nil {
		s.log.Error("ошибка удаления файла", logger.Err(err))
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// S3FileDownload загрузка файла с сервера.
func (s *AuthService) S3FileDownload(ctx context.Context, token string, filename string) ([]byte, error) {
	// добавление токена авторизации в контекст
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(sessionToken, token))
	//открываем стрим для загрузки файла
	stream, err := s.client.FileDownload(ctx, &pb.FileDownloadRequest{
		Filename: filename,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open stream: %w", err)
	}
	fmt.Printf("Downloading %s:\n", filename)
	var fileData []byte
	var lastPercent int
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Print("\n")
			s.log.Error("ошибка получения чанка", logger.Err(err))
			return nil, fmt.Errorf("failed to receive chunk: %w", err)
		}
		fileData = append(fileData, response.GetChunk()...)
		// Показываем процент каждые 5%.
		currentPercent := (len(fileData) / kb) % percent // примерный процент от каждого мегабайта
		if currentPercent >= lastPercent+percentStep {
			fmt.Printf("\r%d%%", currentPercent)
			lastPercent = currentPercent
		}
	}
	fmt.Printf("\r100%%\nDownload complete: %s (%d bytes)\n", filename, len(fileData))
	s.log.Debug("файл загружен", "filename", filename)
	return fileData, nil
}
