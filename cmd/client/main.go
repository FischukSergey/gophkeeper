package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/FischukSergey/gophkeeper/internal/client/command"
	"github.com/FischukSergey/gophkeeper/internal/client/config"
	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
	mainmenutui "github.com/FischukSergey/gophkeeper/internal/client/mainmenuTUI"
	"github.com/FischukSergey/gophkeeper/internal/client/service"
	"github.com/FischukSergey/gophkeeper/internal/logger"
)

var log *slog.Logger
var (
	noteService command.INoteService
	cardService command.ICardService
	authService command.IAuthService
)

func init() {
	// удаление файла лога, если он существует
	if _, err := os.Stat("client.log"); err == nil {
		if err := os.Remove("client.log"); err != nil {
			slog.Error("failed to remove log file", logger.Err(err))
		}
	}
	// Создаем файл для логов
	logFile, err := os.OpenFile("client.log", os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		slog.Error("failed to open log file", logger.Err(err))
		os.Exit(1)
	}
	// Создаем handler для записи в файл
	log = slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(log)
}

func main() {
	log.Info("start client")
	// получение конфигурации сервера
	serverConfig, err := config.NewConfigServerClient()
	if err != nil {
		log.Error("failed to create server config", logger.Err(err))
		os.Exit(1)
	}
	// создание соединения с сервером
	conn, client, cardClient, noteClient, err := grpcclient.NewClient(serverConfig, log)
	if err != nil {
		log.Error("failed to create client", logger.Err(err))
		os.Exit(1)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Error("failed to close connection", logger.Err(err))
		}
	}()
	// создание сервиса аутентификации
	authService = service.NewAuthService(client, log)
	// создание сервиса карт
	cardService = service.NewCardService(cardClient, log)
	// создание сервиса заметок
	noteService = service.NewNoteService(noteClient, log)

	// проверяем, что сервер хранения паролей работает
	err = authService.Check(context.Background())
	if err != nil {
		log.Error("сервер хранения паролей не работает", logger.Err(err))
		os.Exit(1)
	}
	token := &grpcclient.Token{}

	mainmenutui.MainMenuTUI(cardService, authService, noteService, token)
}
