package services

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/FischukSergey/gophkeeper/cmd/server/initial"
	"github.com/FischukSergey/gophkeeper/internal/lib/jwt"
	"github.com/FischukSergey/gophkeeper/internal/lib/luhn"
	"github.com/FischukSergey/gophkeeper/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// DBKeeper интерфейс для сервиса парольного хранилища.
type DBKeeper interface {
	GetPingDB(ctx context.Context) error
	RegisterUser(ctx context.Context, login, password string) (int64, error)
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
}

// S3Keeper интерфейс для сервиса S3.
type S3Keeper interface {
	S3UploadFile(ctx context.Context, fileData io.Reader, filename string, bucket string) (string, error)
	S3GetFileList(ctx context.Context, bucketID string, bucket string) ([]models.File, error)
	S3DeleteFile(ctx context.Context, bucketID string, bucket string) error
	S3DownloadFile(ctx context.Context, bucketID string, bucket string) ([]byte, error)
}

// CardKeeper интерфейс для сервиса карт.
type CardKeeper interface {
	CardAdd(ctx context.Context, card models.Card) error
	CardGetList(ctx context.Context, userID int64) ([]models.Card, error)
	CardDelete(ctx context.Context, cardID int64) error
	CardAddMetadata(ctx context.Context, cardID int64, metadata string) error
}

// GRPCService структура для сервиса.
type GRPCService struct {
	log     *slog.Logger
	storage DBKeeper
	s3      S3Keeper
}

// CardService структура для сервиса карт.
type CardService struct {
	log     *slog.Logger
	storage CardKeeper
}

// NewGRPCService функция для создания сервиса.
func NewGRPCService(log *slog.Logger, storage DBKeeper, s3 S3Keeper) *GRPCService {
	return &GRPCService{log: log, storage: storage, s3: s3}
}

// NewCardService функция для создания сервиса карт.
func NewCardService(log *slog.Logger, storage CardKeeper) *CardService {
	return &CardService{log: log, storage: storage}
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
	token, err := g.GenerateToken(ctx, user)
	if err != nil {
		return models.Token{}, err
	}
	return token, nil
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
	token, err := g.GenerateToken(ctx, user)
	if err != nil {
		return models.Token{}, err
	}
	return token, nil
}

// FileUploadToS3 метод для загрузки файла в S3.
func (g *GRPCService) FileUploadToS3(
	ctx context.Context,
	fileData io.Reader,
	filename string,
	userID int64,
) (string, error) {
	g.log.Info("Service FileUploadToS3 method called")
	bucket := initial.Cfg.GetS3Bucket()
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

// FileGetListFromS3 метод для получения списка файлов пользователя из S3.
func (g *GRPCService) FileGetListFromS3(ctx context.Context, userID int64) ([]models.File, error) {
	g.log.Info("Service FileGetListFromS3 method called")
	bucket := initial.Cfg.GetS3Bucket()
	bucketID := fmt.Sprintf("%d", userID)
	files, err := g.s3.S3GetFileList(ctx, bucketID, bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to get file list: %w", err)
	}
	return files, nil
}

// FileDeleteFromS3 метод для удаления файла из S3.
func (g *GRPCService) FileDeleteFromS3(ctx context.Context, userID int64, filename string) error {
	g.log.Info("Service FileDeleteFromS3 method called")
	bucket := initial.Cfg.GetS3Bucket()
	bucketID := fmt.Sprintf("%d/%s", userID, filename)
	err := g.s3.S3DeleteFile(ctx, bucketID, bucket)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// FileDownloadFromS3 метод для скачивания файла из S3.
func (g *GRPCService) FileDownloadFromS3(ctx context.Context, userID int64, filename string) ([]byte, error) {
	g.log.Info("Service FileDownloadFromS3 method called")
	bucket := initial.Cfg.GetS3Bucket()
	bucketID := fmt.Sprintf("%d/%s", userID, filename)
	data, err := g.s3.S3DownloadFile(ctx, bucketID, bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	return data, nil
}

// GenerateToken метод для генерации токена.
func (g *GRPCService) GenerateToken(ctx context.Context, user models.User) (models.Token, error) {
	g.log.Info("Service GenerateToken method called")
	token, err := jwt.GenerateToken(user)
	if err != nil {
		return models.Token{}, fmt.Errorf("failed to generate token: %w", err)
	}
	return models.Token{
		UserID:    user.ID,
		Token:     token,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(initial.Cfg.JWT.ExpiresKey),
	}, nil
}

// CardAdd метод для добавления карты.
func (g *CardService) CardAdd(ctx context.Context, card models.Card) error {
	g.log.Info("Service CardAdd method called")
	//валидируем данные
	if card.CardNumber == "" || card.CardHolder == "" || card.CardCVV == "" {
		return fmt.Errorf("invalid card data")
	}
	card.CardNumber = strings.ReplaceAll(card.CardNumber, "-", "")
	//валидируем номер карты
	if !luhn.Valid(card.CardNumber) || len(card.CardNumber) != 16 {
		return fmt.Errorf("invalid card number")
	}
	card.CardHolder = strings.ToUpper(card.CardHolder)
	//добавляем карту в базу данных
	err := g.storage.CardAdd(ctx, card)
	if err != nil {
		return fmt.Errorf("failed to add card: %w", err)
	}
	return nil
}

// CardGetList метод для получения списка карт.
func (g *CardService) CardGetList(ctx context.Context, userID int64) ([]models.Card, error) {
	g.log.Info("Service CardGetList method called")
	cards, err := g.storage.CardGetList(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get card list: %w", err)
	}
	g.log.Info("cards", slog.Any("cards", cards))
	return cards, nil
}

// CardDelete метод для удаления карты.
func (g *CardService) CardDelete(ctx context.Context, cardID int64) error {
	g.log.Info("Service CardDelete method called")
	err := g.storage.CardDelete(ctx, cardID)
	if err != nil {
		return fmt.Errorf("failed to delete card: %w", err)
	}
	return nil
}

// CardAddMetadata метод для добавления метаданных к карте.
func (g *CardService) CardAddMetadata(
	ctx context.Context,
	userID int64,
	cardID int64,
	metadata []models.Metadata,
) error {
	g.log.Info("Service CardAddMetadata method called")
	//валидируем метаданные
	err := ValidateMetadata(metadata)
	if err != nil {
		return err
	}
	//сериализуем данные
	metadataJSON, err := SerializeMetadata(metadata)
	if err != nil {
		return err
	}
	g.log.Info("metadataJSON", slog.String("metadataJSON", metadataJSON))
	err = g.storage.CardAddMetadata(ctx, cardID, metadataJSON)
	if err != nil {
		return fmt.Errorf("failed to add metadata: %w", err)
	}
	return nil
}
