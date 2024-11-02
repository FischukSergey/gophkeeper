package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/FischukSergey/gophkeeper/cmd/server/initial"
	"github.com/FischukSergey/gophkeeper/internal/lib/jwt"
	"github.com/FischukSergey/gophkeeper/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// PwdKeeper интерфейс для сервиса парольного хранилища.
type PwdKeeper interface {
	GetPingDB(ctx context.Context) error
	RegisterUser(ctx context.Context, login, password string) (int64, error)
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
}

type S3Keeper interface {
	S3UploadFile(ctx context.Context, fileData []byte, filename string, bucket string) (string, error)
}

// GRPCService структура для сервиса.
type GRPCService struct {
	log     *slog.Logger
	storage PwdKeeper
	s3      S3Keeper
}

// NewGRPCService функция для создания сервиса.
func NewGRPCService(log *slog.Logger, storage PwdKeeper, s3 S3Keeper) *GRPCService {
	return &GRPCService{log: log, storage: storage, s3: s3}
}

// Ping метод для проверки соединения с сервером.
func (g *GRPCService) Ping(ctx context.Context) error {
	err := g.storage.GetPingDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping: %w", err)
	}
	return nil
}

// RegisterUser метод для регистрации пользователя.
func (g *GRPCService) RegisterUser(ctx context.Context, login, password string) (models.Token, error) {
	g.log.Info("Service RegisterUser method called")

	// хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.Token{}, fmt.Errorf("failed to hash password: %w", err)
	}

	// регистрируем пользователя
	userID, err := g.storage.RegisterUser(ctx, login, string(hashedPassword))
	if err != nil {
		return models.Token{}, fmt.Errorf("failed to register user: %w", err)
	}

	user := models.User{
		ID:    userID,
		Login: login,
	}
	// генерируем токен
	token, err := jwt.GenerateToken(user)
	if err != nil {
		return models.Token{}, fmt.Errorf("failed to generate token: %w", err)
	}

	tokenInfo := models.Token{
		UserID:    userID,
		Token:     token,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(initial.Cfg.JWT.ExpiresKey),
	}
	return tokenInfo, nil
}

// Authorization метод для авторизации пользователя.
func (g *GRPCService) Authorization(ctx context.Context, login, password string) (models.Token, error) {
	g.log.Info("Service Authorization method called")

	// получаем пользователя из базы данных по логину
	user, err := g.storage.GetUserByLogin(ctx, login)
	if err != nil {
		return models.Token{}, fmt.Errorf("failed to get user by login: %w", err)
	}

	// проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return models.Token{}, fmt.Errorf("invalid password: %w", err)
	}

	// генерируем токен
	token, err := jwt.GenerateToken(user)
	if err != nil {
		return models.Token{}, fmt.Errorf("failed to generate token: %w", err)
	}

	tokenInfo := models.Token{
		UserID:    user.ID,
		Token:     token,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(initial.Cfg.JWT.ExpiresKey),
	}
	return tokenInfo, nil
}

// FileUploadToS3 метод для загрузки файла в S3.
func (g *GRPCService) FileUploadToS3(
	ctx context.Context,
	fileData []byte,
	filename string,
	userID int64,
) (string, error) {
	g.log.Info("Service FileUploadToS3 method called")
	bucket := initial.Cfg.S3.Bucket
	bucketID := fmt.Sprintf("%d/%s", userID, filename)
	g.log.Info("bucketID", slog.String("bucketID", bucketID))
	g.log.Info("bucket", slog.String("bucket", bucket))
	g.log.Info("userID", slog.Int64("userID", userID))

	url, err := g.s3.S3UploadFile(ctx, fileData, bucketID, bucket)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}
	return url, nil
}
