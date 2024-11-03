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


func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
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

	commandRegister := command.NewCommandRegister(authService, token, reader, writer)
	commandLogin := command.NewCommandLogin(authService, token, reader, writer)
	commandFileUpload := command.NewCommandFileUpload(authService, token, reader, writer)
	commandExit := command.NewCommandExit(reader, writer)

	commands := []command.ICommand{
		commandRegister,
		commandLogin,
		commandFileUpload,
		commandExit,
	}

	commandsMenu := make(map[string]func())
	for _, command := range commands {
		commandsMenu[command.Name()] = command.Execute
	}
	commandsNames := make([]string, 0, len(commands))
	for _, command := range commands {
		commandsNames = append(commandsNames, command.Name())
	}
	for {
		promt := promptui.Select{
			Label: "Введите команду",
			Items: commandsNames,
		}
		_, result, err := promt.Run()
		if err != nil {
			log.Error("failed to run prompt", logger.Err(err))
			os.Exit(1)
		}
		commandsMenu[result]()
	}
}
