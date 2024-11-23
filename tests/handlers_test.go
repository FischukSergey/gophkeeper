package tests

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/FischukSergey/gophkeeper/cmd/server/initial"
	"github.com/FischukSergey/gophkeeper/internal/client/config"
	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
	"github.com/FischukSergey/gophkeeper/internal/client/service"
	"github.com/brianvoe/gofakeit"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var authService *service.AuthService
var cardService *service.CardService
var noteService *service.NoteService

// TestMain инициализация клиента.
func TestMain(m *testing.M) {
	os.Setenv("CONFIG_PATH", "../config/local.yml")
	initial.InitConfig()
	// создание логгера
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)
	logger.Info("TestMain database connected")
	serverConfig := &config.Config{
		ServerAddress: initial.Cfg.GRPC.Port,
	}
	//создание клиента grpc
	grpcConn, grpcClient, serviceClient, noteServiceClient, err := grpcclient.NewClient(serverConfig, logger)
	if err != nil {
		logger.Error("failed to create grpc client", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := grpcConn.Close(); err != nil {
			logger.Error("failed to close grpc connection", "error", err)
		}
	}()
	// создание сервиса аутентификации
	authService = service.NewAuthService(grpcClient, logger)
	// создание сервиса карт
	cardService = service.NewCardService(serviceClient, logger)
	// создание сервиса заметок
	noteService = service.NewNoteService(noteServiceClient, logger)
	// проверяем, что сервер хранения паролей работает
	err = authService.Check(context.Background())
	if err != nil {
		logger.Error("сервер хранения паролей не работает", "error", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestPing(t *testing.T) {
	err := authService.Check(context.Background())
	if err != nil {
		t.Errorf("failed to ping: %v", err)
	}
}

func TestRegister(t *testing.T) {
	//присваиваем произвольные значения
	login := gofakeit.Username()
	password := gofakeit.Password(true, true, true, false, false, 10)

	//регистрируем пользователя
	token, err := authService.Register(context.Background(), login, password)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	//валидируем токен
	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(initial.Cfg.JWT.SecretKey), nil
	})
	require.NoError(t, err)
	require.NotNil(t, tokenParsed)
	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)
	assert.Equal(t, login, claims["login"])

	//авторизуем пользователя
	user, err := authService.Authorization(context.Background(), login, password)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	//создаем табличные тесты на проверку обработки ошибок
	tests := []struct {
		name     string
		login    string
		password string
		wantErr  bool
	}{
		{name: "empty login", login: "", password: password, wantErr: true},
		{name: "empty password", login: login, password: "", wantErr: true},
		{name: "login exists", login: login, password: password, wantErr: true},
		{name: "short login", login: "a", password: password, wantErr: true},
		{name: "short password", login: login, password: "short", wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := authService.Register(context.Background(), test.login, test.password)
			assert.Error(t, err)
		})
	}
}
func TestAuthorization(t *testing.T) {
	login := gofakeit.Username()
	password := gofakeit.Password(true, true, true, false, false, 10)
	//регистрируем пользователя
	token, err := authService.Register(context.Background(), login, password)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	//авторизуем пользователя
	user, err := authService.Authorization(context.Background(), login, password)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	//создаем табличные тесты на проверку обработки ошибок
	tests := []struct {
		name     string
		login    string
		password string
		wantErr  bool
	}{
		{name: "empty login", login: "", password: password, wantErr: true},
		{name: "empty password", login: login, password: "", wantErr: true},
		{name: "login not exists", login: login + "1", password: password, wantErr: true},
		{name: "wrong password", login: login, password: password + "1", wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := authService.Authorization(context.Background(), test.login, test.password)
			assert.Error(t, err)
		})
	}
}
