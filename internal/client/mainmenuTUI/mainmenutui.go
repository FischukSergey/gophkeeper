package mainmenutui

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/FischukSergey/gophkeeper/internal/client/command"
	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
	"github.com/FischukSergey/gophkeeper/internal/client/service"
	"github.com/FischukSergey/gophkeeper/internal/logger"
	"github.com/manifoldco/promptui"
)

var log *slog.Logger

// MainMenuTUI функция для работы с главным меню.
func MainMenuTUI(
	cardService *service.CardService,
	authService *service.AuthService,
	noteService *service.NoteService,
	token *grpcclient.Token,
) {
	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	reader := os.Stdin
	writer := os.Stdout

	// Создаем все команды
	commandRegister := command.NewCommandRegister(authService, token, reader, writer)
	commandLogin := command.NewCommandLogin(authService, token, reader, writer)
	commandFileUpload := command.NewCommandFileUpload(authService, token, reader, writer)
	commandFileDownload := command.NewCommandFileDownload(authService, token, reader, writer)
	commandFileDelete := command.NewCommandFileDelete(authService, token, reader, writer)
	commandFileGetList := command.NewCommandFileGetList(authService, token, reader, writer)
	commandCardAdd := command.NewCommandCardAdd(cardService, token, reader, writer)
	commandCardGetList := command.NewCommandCardGetList(cardService, token, reader, writer)
	commandCardDelete := command.NewCommandCardDelete(cardService, token, reader, writer)
	commandCardAddMetadata := command.NewCommandCardAddMetadata(cardService, token, reader, writer)
	commandNoteAdd := command.NewCommandNoteAdd(noteService, token, reader, writer)
	commandNoteGetList := command.NewCommandNoteGetList(noteService, token, reader, writer)
	commandNoteDelete := command.NewCommandNoteDelete(noteService, token, reader, writer)
	commandExit := command.NewCommandExit(reader, writer)

	// Определяем основное меню и подменю файлов
	mainCommands := []command.ICommand{
		commandLogin,
		commandRegister,
	}
	// команды для работы с файлами
	fileCommands := []command.ICommand{
		commandFileUpload,
		commandFileDownload,
		commandFileGetList,
		commandFileDelete,
	}
	// команды для работы с картами
	cardCommands := []command.ICommand{
		commandCardAdd,
		commandCardGetList,
		commandCardAddMetadata,
		commandCardDelete,
	}
	// команды для работы с заметками
	noteCommands := []command.ICommand{
		commandNoteAdd,
		commandNoteGetList,
		commandNoteDelete,
	}

	// Создаем мапы для команд
	mainCommandsMenu := make(map[string]func())
	fileCommandsMenu := make(map[string]func())
	cardCommandsMenu := make(map[string]func())
	noteCommandsMenu := make(map[string]func())

	// Заполняем мапу основного меню
	for _, cmd := range mainCommands {
		mainCommandsMenu[cmd.Name()] = cmd.Execute
	}
	// добавляем специальную команду для перехода в подменю файлов
	mainCommandsMenu["\tРАБОТА С ФАЙЛАМИ"] = func() {
		handleSubmenu(fileCommands, fileCommandsMenu)
	}
	// добавляем специальную команду для перехода в подменю карт
	mainCommandsMenu["\tРАБОТА С КАРТАМИ"] = func() {
		handleSubmenu(cardCommands, cardCommandsMenu)
	}
	// добавляем специальную команду для перехода в подменю заметок
	mainCommandsMenu["\tРАБОТА С ЗАМЕТКАМИ"] = func() {
		handleSubmenu(noteCommands, noteCommandsMenu)
	}
	// добавляем команду выхода
	mainCommandsMenu[commandExit.Name()] = commandExit.Execute

	// заполняем мапу подменю файлов
	for _, cmd := range fileCommands {
		fileCommandsMenu[cmd.Name()] = cmd.Execute
	}
	fileCommandsMenu["Назад"] = func() {}

	// заполняем мапу подменю карт
	for _, cmd := range cardCommands {
		cardCommandsMenu[cmd.Name()] = cmd.Execute
	}
	cardCommandsMenu["Назад"] = func() {}

	// заполняем мапу подменю заметок
	for _, cmd := range noteCommands {
		noteCommandsMenu[cmd.Name()] = cmd.Execute
	}
	noteCommandsMenu["Назад"] = func() {}

	// формируем список названий команд для главного меню
	mainCommandNames := make([]string, 0, len(mainCommands)+1)
	for _, cmd := range mainCommands {
		mainCommandNames = append(mainCommandNames, cmd.Name())
	}
	mainCommandNames = append(mainCommandNames, "\tРАБОТА С ФАЙЛАМИ")
	mainCommandNames = append(mainCommandNames, "\tРАБОТА С КАРТАМИ")
	mainCommandNames = append(mainCommandNames, "\tРАБОТА С ЗАМЕТКАМИ")
	mainCommandNames = append(mainCommandNames, commandExit.Name())

	for {
		templates := newTemplates()
		prompt := promptui.Select{
			Templates: templates,
			Label:     "Введите команду:",
			Items:     mainCommandNames,
			Size:      10,
		}
		_, result, err := prompt.Run()
		if err != nil {
			log.Error("failed to run prompt", logger.Err(err))
			os.Exit(1)
		}
		mainCommandsMenu[result]()
	}
}

// handleSubmenu функция для обработки подменю.
func handleSubmenu(commands []command.ICommand, commandsMenu map[string]func()) {
	commandNames := make([]string, 0, len(commands)+1)
	for _, cmd := range commands {
		commandNames = append(commandNames, cmd.Name())
	}
	commandNames = append(commandNames, "Назад")

	for {
		templates := newTemplates()
		prompt := promptui.Select{
			Templates: templates,
			Label:     "Меню работы с файлами:",
			Items:     commandNames,
			Size:      10,
		}
		_, result, err := prompt.Run()
		if err != nil {
			fmt.Println("Ошибка при выполнении команды")
			os.Exit(1)
		}
		if result == "Назад" {
			return
		}
		commandsMenu[result]()
	}
}

// newTemplates функция для создания шаблона для меню.
func newTemplates() *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label:    "{{ . | red }}",
		Active:   "\U0001F449 {{ . | green }}",
		Inactive: "  {{ . | cyan }}",
		Selected: "-> {{ . | green }}",
	}
}
