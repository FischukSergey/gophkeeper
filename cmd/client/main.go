package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/FischukSergey/gophkeeper/internal/client/command"
	"github.com/FischukSergey/gophkeeper/internal/client/config"
	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
	"github.com/FischukSergey/gophkeeper/internal/client/service"
	"github.com/FischukSergey/gophkeeper/internal/logger"
	"github.com/manifoldco/promptui"
)

var log *slog.Logger

func init() {
	// Создаем файл для логов
	logFile, err := os.OpenFile("client.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		slog.Error("failed to open log file", logger.Err(err))
		os.Exit(1)
	}
	// Создаем handler для записи в файл
	log = slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slog.LevelInfo}))
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
	conn, client, err := grpcclient.NewClient(serverConfig, log)
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
	authService := service.NewAuthService(client, log)
	// проверяем, что сервер хранения паролей работает
	err = authService.Check(context.Background())
	if err != nil {
		log.Error("сервер хранения паролей не работает", logger.Err(err))
		os.Exit(1)
	}

	reader := os.Stdin
	writer := os.Stdout
	token := &grpcclient.Token{}

	// Создаем все команды
	commandRegister := command.NewCommandRegister(authService, token, reader, writer)
	commandLogin := command.NewCommandLogin(authService, token, reader, writer)
	commandFileUpload := command.NewCommandFileUpload(authService, token, reader, writer)
	commandFileDownload := command.NewCommandFileDownload(authService, token, reader, writer)
	commandFileDelete := command.NewCommandFileDelete(authService, token, reader, writer)
	commandFileGetList := command.NewCommandFileGetList(authService, token, reader, writer)
	commandExit := command.NewCommandExit(reader, writer)

	// Определяем основное меню и подменю файлов
	mainCommands := []command.ICommand{
		commandLogin,
		commandRegister,
		commandExit,
	}

	fileCommands := []command.ICommand{
		commandFileUpload,
		commandFileDownload,
		commandFileGetList,
		commandFileDelete,
	}

	// Создаем мапы для команд
	mainCommandsMenu := make(map[string]func())
	fileCommandsMenu := make(map[string]func())

	// Заполняем мапу основного меню
	for _, cmd := range mainCommands {
		mainCommandsMenu[cmd.Name()] = cmd.Execute
	}
	// Добавляем специальную команду для перехода в подменю файлов
	mainCommandsMenu["Работа с файлами"] = func() {
		handleFileSubmenu(fileCommands, fileCommandsMenu)
	}

	// Заполняем мапу подменю файлов
	for _, cmd := range fileCommands {
		fileCommandsMenu[cmd.Name()] = cmd.Execute
	}
	// Добавляем команду возврата в главное меню
	fileCommandsMenu["Назад"] = func() {}

	// Формируем список названий команд для главного меню
	mainCommandNames := make([]string, 0, len(mainCommands)+1)
	for _, cmd := range mainCommands {
		mainCommandNames = append(mainCommandNames, cmd.Name())
	}
	mainCommandNames = append(mainCommandNames, "Работа с файлами")

	for {
		templates := promptui.SelectTemplates{
			Label:    "{{ . | red }}",
			Active:   "\U0001F449 {{ . | green }}",
			Inactive: "  {{ . | cyan }}",
			Selected: "-> {{ . | green }}",
		}
		prompt := promptui.Select{
			Templates: &templates,
			Label:     "Введите команду:",
			Items:     mainCommandNames,
		}
		_, result, err := prompt.Run()
		if err != nil {
			log.Error("failed to run prompt", logger.Err(err))
			os.Exit(1)
		}
		mainCommandsMenu[result]()
	}
}

// Функция для обработки подменю файлов
func handleFileSubmenu(fileCommands []command.ICommand, fileCommandsMenu map[string]func()) {
	fileCommandNames := make([]string, 0, len(fileCommands)+1)
	for _, cmd := range fileCommands {
		fileCommandNames = append(fileCommandNames, cmd.Name())
	}
	fileCommandNames = append(fileCommandNames, "Назад")

	for {
		templates := promptui.SelectTemplates{
			Label:    "{{ . | red }}",
			Active:   "\U0001F449 {{ . | green }}",
			Inactive: "  {{ . | cyan }}",
			Selected: "-> {{ . | green }}",
		}
		prompt := promptui.Select{
			Templates: &templates,
			Label:     "Меню работы с файлами:",
			Items:     fileCommandNames,
		}
		_, result, err := prompt.Run()
		if err != nil {
			log.Error("failed to run prompt", logger.Err(err))
			os.Exit(1)
		}
		if result == "Назад" {
			return
		}
		fileCommandsMenu[result]()
	}
}
