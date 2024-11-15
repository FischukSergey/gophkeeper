package command

import (
	"context"
	"errors"
	"io"
	"strconv"

	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
	"github.com/FischukSergey/gophkeeper/internal/client/service"
	"github.com/manifoldco/promptui"
)

const (
	cardAddMetadataCommandName = "AddMetadata"
)

// NewCommandCardAddMetadata создает новый экземпляр команды добавления метаданных к карте
func NewCommandCardAddMetadata(
	cardService *service.CardService,
	token *grpcclient.Token,
	reader io.Reader,
	writer io.Writer,
) *CommandCardAddMetadata {
	return &CommandCardAddMetadata{cardService, token, reader, writer}
}

// CommandCardAddMetadata структура команды добавления метаданных к карте
type CommandCardAddMetadata struct {
	cardService *service.CardService
	token       *grpcclient.Token
	reader      io.Reader
	writer      io.Writer
}

// Name возвращает имя команды
func (c *CommandCardAddMetadata) Name() string {
	return cardAddMetadataCommandName
}

// Execute выполняет команду добавления метаданных к карте
func (c *CommandCardAddMetadata) Execute() {
	//проверка наличия токена
	if !checkToken(c.token, c.reader) {
		return
	}
	//вывод заголовка команды
	fprintln(c.writer, "Добавление метаданных к карте")
	//получение списка карт
	cardsList := NewCommandCardGetList(c.cardService, c.token, c.reader, c.writer)
	cardsList.Execute()
	//ввод ID карты
	prompt := promptui.Prompt{
		Label: "Введите ID карты для добавления метаданных: ",
		Validate: func(input string) error {
			if input == "" {
				return errors.New("ID карты не может быть пустым")
			}
			if _, err := strconv.Atoi(input); err != nil {
				return errors.New("ID карты должен быть числом")
			}
			return nil
		},
	}
	cardID, err := prompt.Run()
	if err != nil {
		fprintln(c.writer, "Ошибка при вводе ID карты:", err)
		return
	}	
	// Создаем map для хранения метаданных
	metadata := make(map[string]string)
	
	for {
		// Показываем текущие метаданные
		fprintln(c.writer, "\nТекущие введенные метаданные:")
		if len(metadata) == 0 {
			fprintln(c.writer, "Пока нет введенных метаданных")
		} else {
			for k, v := range metadata {
				fprintf(c.writer, "%s: %s\n", k, v)
			}
		}

		// Спрашиваем, хочет ли пользователь добавить еще пару
		continuePrompt := promptui.Prompt{
			Label:     "Добавить новую пару key-value? (y/n)",
			Default:   "y",
			Validate: func(input string) error {
				if input != "y" && input != "n" {
					return errors.New("пожалуйста, введите 'y' или 'n'")
				}
				return nil
			},
		}

		shouldContinue, err := continuePrompt.Run()
		if err != nil || shouldContinue == "n" {
			break //выход из цикла, если пользователь не хочет добавлять больше метаданных
		}

		// Ввод ключа
		keyPrompt := promptui.Prompt{
			Label: "Введите ключ",
			Validate: func(input string) error {
				if input == "" {
					return errors.New("ключ не может быть пустым")
				}
				return nil
			},
		}

		key, err := keyPrompt.Run()
		if err != nil {
			fprintln(c.writer, "Ошибка при вводе ключа:", err)
			return
		}

		// Ввод значения
		valuePrompt := promptui.Prompt{
			Label: "Введите значение",
			Validate: func(input string) error {
				if input == "" {
					return errors.New("значение не может быть пустым")
				}
				return nil
			},
		}

		value, err := valuePrompt.Run()
		if err != nil {
			fprintln(c.writer, "Ошибка при вводе значения:", err)
			return
		}
		//проверяем, нет ли такого ключа в map
		if _, ok := metadata[key]; ok {
			fprintln(c.writer, "Такой ключ уже существует")
			continue
		}
		metadata[key] = value //добавление пары ключ-значение в map
	}
	//добавление метаданных к карте
	cardIDInt, err := strconv.Atoi(cardID)
	if err != nil {
		fprintln(c.writer, "Ошибка при преобразовании ID карты:", err)
		return
	}
	err = c.cardService.AddCardMetadata(context.Background(), int64(cardIDInt), metadata, c.token.Token)
	if err != nil {
		fprintln(c.writer, "Ошибка при добавлении метаданных:", err)
	}
	fprintln(c.writer, "Метаданные успешно добавлены")
	waitEnter(c.reader)
}
