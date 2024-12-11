package mainmenutui

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/FischukSergey/gophkeeper/internal/client/command"
	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
	"github.com/FischukSergey/gophkeeper/internal/logger"
	"github.com/manifoldco/promptui"
)

var log *slog.Logger

// MainMenuTUI функция для работы с главным меню.
func MainMenuTUI(
	cardService command.ICardService,
	authService command.IAuthService,
	noteService command.INoteService,
	token *grpcclient.Token,
) {
	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	reader := os.Stdin
	writer := os.Stdout

	commands := initializeCommands(cardService, authService, noteService, token, reader, writer)
	menus := createMenus(commands)
	runMainMenu(commands.mainCommands, menus.mainMenu)
}

type commandGroups struct {
	exitCommand  command.ICommand
	mainCommands []command.ICommand
	fileCommands []command.ICommand
	cardCommands []command.ICommand
	noteCommands []command.ICommand
}

type menuMaps struct {
	mainMenu map[string]func()
	fileMenu map[string]func()
	cardMenu map[string]func()
	noteMenu map[string]func()
}

// initializeCommands функция для инициализации команд.
func initializeCommands(
	cardService command.ICardService,
	authService command.IAuthService,
	noteService command.INoteService,
	token *grpcclient.Token,
	reader *os.File,
	writer *os.File,
) commandGroups {
	return commandGroups{
		mainCommands: []command.ICommand{
			command.NewCommandLogin(authService, token, reader, writer),
			command.NewCommandRegister(authService, token, reader, writer),
		},
		fileCommands: []command.ICommand{
			command.NewCommandFileUpload(authService, token, reader, writer),
			command.NewCommandFileDownload(authService, token, reader, writer),
			command.NewCommandFileGetList(authService, token, reader, writer),
			command.NewCommandFileDelete(authService, token, reader, writer),
		},
		cardCommands: []command.ICommand{
			command.NewCommandCardAdd(cardService, token, reader, writer),
			command.NewCommandCardGetList(cardService, token, reader, writer),
			command.NewCommandCardAddMetadata(cardService, token, reader, writer),
			command.NewCommandCardDelete(cardService, token, reader, writer),
		},
		noteCommands: []command.ICommand{
			command.NewCommandNoteAdd(noteService, token, reader, writer),
			command.NewCommandNoteGetList(noteService, token, reader, writer),
			command.NewCommandNoteDelete(noteService, token, reader, writer),
		},
		exitCommand: command.NewCommandExit(reader, writer),
	}
}

// createMenus функция для создания меню.
func createMenus(commands commandGroups) menuMaps {
	menus := menuMaps{
		mainMenu: make(map[string]func()),
		fileMenu: make(map[string]func()),
		cardMenu: make(map[string]func()),
		noteMenu: make(map[string]func()),
	}

	// инициализируем подменю
	initializeSubmenu(commands.fileCommands, menus.fileMenu)
	initializeSubmenu(commands.cardCommands, menus.cardMenu)
	initializeSubmenu(commands.noteCommands, menus.noteMenu)

	// добавляем команды для работы с главным меню
	for _, cmd := range commands.mainCommands {
		menus.mainMenu[cmd.Name()] = cmd.Execute
	}

	// добавляем команды для работы с подменю
	menus.mainMenu["\tРАБОТА С ФАЙЛАМИ"] = func() {
		handleSubmenu(commands.fileCommands, menus.fileMenu)
	}
	menus.mainMenu["\tРАБОТА С КАРТАМИ"] = func() {
		handleSubmenu(commands.cardCommands, menus.cardMenu)
	}
	menus.mainMenu["\tРАБОТА С ЗАМЕТКАМИ"] = func() {
		handleSubmenu(commands.noteCommands, menus.noteMenu)
	}
	menus.mainMenu["Выход"] = commands.exitCommand.Execute

	return menus
}

// initializeSubmenu функция для инициализации подменю.
func initializeSubmenu(commands []command.ICommand, menu map[string]func()) {
	for _, cmd := range commands {
		menu[cmd.Name()] = cmd.Execute
	}
	menu["Назад"] = func() {}
}

// runMainMenu функция для запуска главного меню.
func runMainMenu(mainCommands []command.ICommand, mainMenu map[string]func()) {
	mainCommandNames := buildMainMenuItems(mainCommands)

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
			fmt.Println("Ошибка при выполнении команды")
			return
		}
		mainMenu[result]()
	}
}

// buildMainMenuItems функция для построения элементов главного меню.
func buildMainMenuItems(mainCommands []command.ICommand) []string {
	items := make([]string, 0, len(mainCommands)+4) // +4 for submenus and exit
	for _, cmd := range mainCommands {
		items = append(items, cmd.Name())
	}
	items = append(items, "\tРАБОТА С ФАЙЛАМИ")
	items = append(items, "\tРАБОТА С КАРТАМИ")
	items = append(items, "\tРАБОТА С ЗАМЕТКАМИ")
	items = append(items, "Выход")
	return items
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
			Label:     "Выберите команду:",
			Items:     commandNames,
			Size:      10,
		}
		_, result, err := prompt.Run()
		if err != nil {
			fmt.Println("Ошибка при выполнении команды")
			return
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
