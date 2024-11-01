package services

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/FischukSergey/gophkeeper/cmd/server/initial"
	"github.com/FischukSergey/gophkeeper/internal/models"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type mockStorage struct{}

func (m *mockStorage) RegisterUser(ctx context.Context, login, hashedPassword string) (int64, error) {
	if login == "qwerty" && hashedPassword == "$2a$10$5a1BYih/bXvIJnkquBqKAeZ/8mBpecQ4HivNQ2AisiU5GeTP0MLem" {
		return 18, nil
	}
	return 0, fmt.Errorf("invalid login or password")
}

func (m *mockStorage) GetPingDB(ctx context.Context) error {
	return nil
}

func (m *mockStorage) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	if login == "qwerty" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			return models.User{}, err
		}
		return models.User{
			ID:             18,
			Login:          "qwerty",
			HashedPassword: string(hashedPassword),
			CreatedAt:      time.Now(),
		}, nil
	}
	return models.User{}, fmt.Errorf("user not found")
}


func TestAuthorization(t *testing.T) {
	os.Setenv("CONFIG_PATH", "../../../config/local.yml")
	initial.InitConfig()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	mockStorage := &mockStorage{}
	service := NewGRPCService(logger, mockStorage, nil)

	type arg struct {
		name          string
		login         string
		password      string
		expectedUser  models.User
		expectedError error
	}
	args := []arg{
		{
			name:     "success",
			login:    "qwerty",
			password: "password",
			expectedUser: models.User{
				ID: 18,
			},
			expectedError: nil,
		},
		{
			name:          "invalid password",
			login:         "qwerty",
			password:      "password123",
			expectedUser:  models.User{},
			expectedError: fmt.Errorf("invalid password"),
		},
	}
	for _, arg := range args {
		t.Run(arg.name, func(t *testing.T) {
			token, err := service.Authorization(context.Background(), arg.login, arg.password)
			switch arg.name {
			case "success":
				assert.Equal(t, arg.expectedUser.ID, token.UserID)
			case "invalid password":
				assert.Contains(t, err.Error(), arg.expectedError.Error())
			}
		})
	}
}

func TestRegisterUser(t *testing.T) {
	os.Setenv("CONFIG_PATH", "../../../config/local.yml")
	initial.InitConfig()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	mockStorage := &mockStorage{}
	service := NewGRPCService(logger, mockStorage, nil)
	type arg struct {
		name          string
		login         string
		password      string
		expectedError error
		expectedID    int64
	}
	args := []arg{
		{
			name:          "success",
			login:         "qwerty",
			password:      "password",
			expectedError: nil,
			expectedID:    18,
		},
		{
			name:          "invalid password",
			login:         "qwerty",
			password:      "password123",
			expectedError: fmt.Errorf("invalid password"),
		},
	}

	for _, arg := range args {
		t.Run(arg.name, func(t *testing.T) {
			token, err := service.RegisterUser(context.Background(), arg.login, arg.password)
			switch arg.name {
			case "success":
				assert.Equal(t, arg.expectedID, token.UserID)
				assert.Equal(t, arg.expectedError, err)
			case "invalid password":
				assert.Contains(t, err.Error(), arg.expectedError.Error())
			}
		})
	}
}

//TODO: добавить тесты для S3
