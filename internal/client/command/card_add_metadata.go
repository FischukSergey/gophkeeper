package command

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
	"github.com/FischukSergey/gophkeeper/internal/models"
	"github.com/manifoldco/promptui"
)

const (
	cardAddMetadataCommandName = "AddMetadata"
)

// NewCommandCardAddMetadata создает новый экземпляр команды добавления метаданных к карте.
func NewCommandCardAddMetadata(
	cardService ICardService,
	token *grpcclient.Token,
	reader io.Reader,
	writer io.Writer,
) *CommandCardAddMetadata {
	return &CommandCardAddMetadata{cardService, token, reader, writer}
}

// CommandCardAddMetadata структура команды добавления метаданных к карте.
type CommandCardAddMetadata struct {
	cardService ICardService
	token       *grpcclient.Token
	reader      io.Reader
	writer      io.Writer
}

// Name возвращает имя команды.
func (c *CommandCardAddMetadata) Name() string {
	return cardAddMetadataCommandName
}

// Execute выполняет команду добавления метаданных к карте.
func (c *CommandCardAddMetadata) Execute() {
	//проверка наличия токена
	if !checkToken(c.token, c.reader) {
		return
	}
	//вывод заголовка команды
	fprintln(c.writer, "Добавление метаданных к карте")
	//получение списка карт
	cards, err := getCardList(c.cardService, c.token, c.reader, c.writer)
	if err != nil {
		fprintln(c.writer, errGetCardsMessage, err)
		waitEnter(c.reader)
		return
	}
	//ввод ID карты
	cardIDInt, err := promptCardID()
	if err != nil {
		fprintln(c.writer, "Ошибка при вводе ID карты:", err)
		return
	}
	//проверяем, что есть такая карта перебором списка карт
	var exist bool
	for _, card := range cards {
		if card.CardID == int64(cardIDInt) {
			exist = true
			break
		}
	}
	if !exist {
		fprintln(c.writer, "Карта с таким ID не найдена")
		waitEnter(c.reader)
		return
	}
	//создаем map для хранения метаданных
	metadata, err := c.promptMetadata()
	if err != nil {
		fprintln(c.writer, "Ошибка при вводе метаданных:", err)
		return
	}
	//добавление метаданных к карте
	err = c.cardService.AddCardMetadata(context.Background(), int64(cardIDInt), metadata, c.token.Token)
	if err != nil {
		fprintln(c.writer, "Ошибка при добавлении метаданных:", err)
	}
	fprintln(c.writer, "Метаданные успешно добавлены")
	waitEnter(c.reader)
}

// getCardList получение списка карт.
func getCardList(
	cardService ICardService,
	token *grpcclient.Token,
	reader io.Reader,
	writer io.Writer,
) ([]models.Card, error) {
	cardsList := NewCommandCardGetList(cardService, token, reader, writer)
	cardsList.Execute()
	cards, err := cardService.GetCardList(context.Background(), token.Token)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errGetCardsMessage, err)
	}
	if len(cards) == 0 {
		return nil, errors.New("список карт пуст")
	}
	return cards, nil
}

// promptCardID ввод ID карты.
func promptCardID() (int, error) {
	prompt := promptui.Prompt{
		Label: "Введите ID карты: ",
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
		return 0, fmt.Errorf("ошибка при вводе ID карты: %w", err)
	}
	cardIDInt, err := strconv.Atoi(cardID)
	if err != nil {
		return 0, fmt.Errorf("ошибка при преобразовании ID карты: %w", err)
	}
	return cardIDInt, nil
}

// promptMetadata ввод метаданных.
func (c *CommandCardAddMetadata) promptMetadata() (map[string]string, error) {
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
			Label:   "Добавить новую пару key-value? (y/n)",
			Default: "y",
			Validate: func(input string) error {
				if input != "y" && input != "n" {
					return errors.New("пожалуйста, введите 'y' или 'n'")
				}
				return nil
			},
		}

		shouldContinue, err := continuePrompt.Run()
		if err != nil {
			return nil, fmt.Errorf("ошибка при вводе продолжения: %w", err)
		}
		if shouldContinue == "n" {
			break //выход из цикла, если пользователь не хочет добавлять больше метаданных
		}
		// Ввод ключа
		key, err := c.inputKey()
		if err != nil {
			return nil, fmt.Errorf("ошибка при вводе ключа: %w", err)
		}
		// Ввод значения
		value, err := c.inputValue()
		if err != nil {
			return nil, fmt.Errorf("ошибка при вводе значения: %w", err)
		}
		//проверяем, нет ли такого ключа в map
		if _, ok := metadata[key]; ok {
			fprintln(c.writer, "Такой ключ уже существует")
			continue
		}
		metadata[key] = value //добавление пары ключ-значение в map
	}
	return metadata, nil
}

// inputKey ввод ключа.
func (c *CommandCardAddMetadata) inputKey() (string, error) {
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
		return "", fmt.Errorf("ошибка при вводе ключа: %w", err)
	}
	return key, nil
}

// inputValue ввод значения.
func (c *CommandCardAddMetadata) inputValue() (string, error) {
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
		return "", fmt.Errorf("ошибка при вводе значения: %w", err)
	}
	return value, nil
}
