package services

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/FischukSergey/gophkeeper/cmd/server/initial"
	"github.com/FischukSergey/gophkeeper/internal/config"
	"github.com/FischukSergey/gophkeeper/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	// Устанавливаем тестовую конфигурацию JWT
	initial.Cfg = &config.Config{
		JWT: config.JWTConfig{
			SecretKey:  "test_secret_key",
			ExpiresKey: time.Hour * 12,
		},
	}
}

// Mock structures.
type MockPwdKeeper struct {
	mock.Mock
}

func (m *MockPwdKeeper) GetPingDB(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockPwdKeeper) RegisterUser(ctx context.Context, login, password string) (int64, error) {
	args := m.Called(ctx, login, password)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPwdKeeper) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	args := m.Called(ctx, login)
	return args.Get(0).(models.User), args.Error(1)
}

type MockS3Keeper struct {
	mock.Mock
}

func (m *MockS3Keeper) S3UploadFile(
	ctx context.Context,
	fileData []byte,
	filename string,
	bucket string,
) (string, error) {
	args := m.Called(ctx, fileData, filename, bucket)
	return args.String(0), args.Error(1)
}

func (m *MockS3Keeper) S3GetFileList(ctx context.Context, bucketID string, bucket string) ([]models.File, error) {
	args := m.Called(ctx, bucketID, bucket)
	return args.Get(0).([]models.File), args.Error(1)
}

func (m *MockS3Keeper) S3DeleteFile(ctx context.Context, bucketID string, bucket string) error {
	args := m.Called(ctx, bucketID, bucket)
	return args.Error(0)
}

func (m *MockS3Keeper) S3DownloadFile(ctx context.Context, bucketID string, bucket string) ([]byte, error) {
	args := m.Called(ctx, bucketID, bucket)
	return args.Get(0).([]byte), args.Error(1)
}

// Test functions.
func TestPing(t *testing.T) {
	tests := []struct {
		name    string
		dbError error
		wantErr bool
	}{
		{
			name:    "successful ping",
			dbError: nil,
			wantErr: false,
		},
		{
			name:    "failed ping",
			dbError: assert.AnError,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(MockPwdKeeper)
			mockS3 := new(MockS3Keeper)
			logger := slog.Default()
			service := NewGRPCService(logger, mockStorage, mockS3)
			ctx := context.Background()

			mockStorage.On("GetPingDB", ctx).Return(tt.dbError)

			err := service.Ping(ctx)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockStorage.AssertExpectations(t)
		})
	}
}

func TestRegisterUser(t *testing.T) {
	// Сохраняем оригинальную конфигурацию
	originalCfg := initial.Cfg
	// Восстанавливаем в конце теста
	defer func() {
		initial.Cfg = originalCfg
	}()

	tests := []struct {
		name     string
		login    string
		password string
		userID   int64
		dbError  error
		wantErr  bool
	}{
		{
			name:     "successful registration",
			login:    "testuser",
			password: "testpass",
			userID:   1,
			dbError:  nil,
			wantErr:  false,
		},
		{
			name:     "failed registration",
			login:    "testuser",
			password: "testpass",
			userID:   0,
			dbError:  assert.AnError,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(MockPwdKeeper)
			mockS3 := new(MockS3Keeper)
			logger := slog.Default()
			service := NewGRPCService(logger, mockStorage, mockS3)
			ctx := context.Background()

			mockStorage.On("RegisterUser", ctx, tt.login, mock.AnythingOfType("string")).Return(tt.userID, tt.dbError)

			token, err := service.RegisterUser(ctx, tt.login, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, nil)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token.Token)
				assert.Equal(t, tt.userID, token.UserID)
			}
			mockStorage.AssertExpectations(t)
		})
	}
}

func TestAuthorization(t *testing.T) {
	// Сохраняем оригинальную конфигурацию
	originalCfg := initial.Cfg
	// Восстанавливаем в конце теста
	defer func() {
		initial.Cfg = originalCfg
	}()

	tests := []struct {
		name       string
		login      string
		password   string
		storedPass string
		userID     int64
		dbError    error
		wantErr    bool
	}{
		{
			name:       "successful authorization",
			login:      "testuser",
			password:   "testpass",
			storedPass: "testpass",
			userID:     1,
			dbError:    nil,
			wantErr:    false,
		},
		{
			name:       "invalid password",
			login:      "testuser",
			password:   "wrongpass",
			storedPass: "correctpass",
			userID:     1,
			dbError:    nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(MockPwdKeeper)
			mockS3 := new(MockS3Keeper)
			logger := slog.Default()
			service := NewGRPCService(logger, mockStorage, mockS3)
			ctx := context.Background()

			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(tt.storedPass), bcrypt.DefaultCost)
			user := models.User{
				ID:             tt.userID,
				Login:          tt.login,
				HashedPassword: string(hashedPassword),
			}

			mockStorage.On("GetUserByLogin", ctx, tt.login).Return(user, tt.dbError)

			token, err := service.Authorization(ctx, tt.login, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token.Token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token.Token)
				assert.Equal(t, tt.userID, token.UserID)
			}
			mockStorage.AssertExpectations(t)
		})
	}
}

func TestFileUploadToS3(t *testing.T) {
	tests := []struct {
		name        string
		fileData    []byte
		filename    string
		userID      int64
		expectedURL string
		s3Error     error
		wantErr     bool
	}{
		{
			name:        "successful upload",
			fileData:    []byte("test data"),
			filename:    "test.txt",
			userID:      1,
			expectedURL: "https://s3.ru-1.storage.selcloud.ru/gophkeeper-bucket/1/test.txt",
			s3Error:     nil,
			wantErr:     false,
		},
		{
			name:        "failed upload",
			fileData:    []byte("test data"),
			filename:    "test.txt",
			userID:      1,
			expectedURL: "",
			s3Error:     assert.AnError,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(MockPwdKeeper)
			mockS3 := new(MockS3Keeper)
			logger := slog.Default()
			service := NewGRPCService(logger, mockStorage, mockS3)
			ctx := context.Background()

			bucketID := fmt.Sprintf("%d/%s", tt.userID, tt.filename)
			mockS3.On(
				"S3UploadFile",
				ctx,
				tt.fileData,
				bucketID,
				mock.AnythingOfType("string"),
			).Return(tt.expectedURL, tt.s3Error)

			url, err := service.FileUploadToS3(ctx, tt.fileData, tt.filename, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, url)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedURL, url)
			}
			mockS3.AssertExpectations(t)
		})
	}
}

func TestFileGetListFromS3(t *testing.T) {
	mockStorage := new(MockPwdKeeper)
	mockS3 := new(MockS3Keeper)
	logger := slog.Default()
	service := NewGRPCService(logger, mockStorage, mockS3)
	ctx := context.Background()

	t.Run("successful get list", func(t *testing.T) {
		userID := int64(1)
		expectedFiles := []models.File{
			{Filename: "test1.txt", Size: 100},
			{Filename: "test2.txt", Size: 200},
		}

		mockS3.On("S3GetFileList", ctx, "1", mock.AnythingOfType("string")).Return(expectedFiles, nil)

		files, err := service.FileGetListFromS3(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedFiles, files)
		mockS3.AssertExpectations(t)
	})
}

func TestFileDeleteFromS3(t *testing.T) {
	mockStorage := new(MockPwdKeeper)
	mockS3 := new(MockS3Keeper)
	logger := slog.Default()
	service := NewGRPCService(logger, mockStorage, mockS3)
	ctx := context.Background()

	t.Run("successful delete", func(t *testing.T) {
		userID := int64(1)
		filename := "test.txt"
		bucketID := "1/test.txt"

		mockS3.On("S3DeleteFile", ctx, bucketID, mock.AnythingOfType("string")).Return(nil)

		err := service.FileDeleteFromS3(ctx, userID, filename)

		assert.NoError(t, err)
		mockS3.AssertExpectations(t)
	})
}
