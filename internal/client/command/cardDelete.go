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

const cardDeleteCommandName = "CardDelete"

// CommandCardDelete структура для команды удаления карты.
type CommandCardDelete struct {
	cardService *service.CardService
	token       *grpcclient.Token
	reader      io.Reader
	writer      io.Writer
}

// NewCommandCardDelete создает новый экземпляр команды удаления карты.
func NewCommandCardDelete(
	cardService *service.CardService,
	token *grpcclient.Token,
	reader io.Reader,
	writer io.Writer,
) *CommandCardDelete {
	return &CommandCardDelete{cardService: cardService, reader: reader, writer: writer, token: token}
}

// Name возвращает имя команды.
func (c *CommandCardDelete) Name() string {
	return cardDeleteCommandName
}

// Execute выполняет команду удаления карты.
func (c *CommandCardDelete) Execute() {
	// Проверка наличия токена
	if !checkToken(c.token, c.reader) {
		return // Выходим из функции если токен невалидный
	}
	fprintln(c.writer, "Удаление карты")
	fprintln(c.writer, "Ознакомтесь со списком карт и введите ID карты для удаления")

	//получаем список карт
	cardsList := NewCommandCardGetList(c.cardService, c.token, c.reader, c.writer)
	cardsList.Execute()
	cards, err := c.cardService.GetCardList(context.Background(), c.token.Token)
	if err != nil {
		fprintln(c.writer, "Ошибка при получении списка карт:", err)
		return
	}
	if len(cards) == 0 {
		fprintln(c.writer, "Список карт пуст")
		return
	}
	prompt := promptui.Prompt{
		Label: "Введите ID карты для удаления: ",
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
	//проверяем, что есть такая карта перебором списка карт
	var exist bool
	for _, card := range cards {
		if card.CardNumber == cardID {
			exist = true
			break
		}
	}
	if !exist {
		fprintln(c.writer, "Карта с таким ID не найдена")
		waitEnter(c.reader)
		return
	}

	//удаляем карту
	err = c.cardService.DeleteCard(context.Background(), cardID, c.token.Token)
	if err != nil {
		fprintln(c.writer, "Ошибка при удалении карты:", err)
		return
	}
	fprintln(c.writer, "Карта удалена")
	waitEnter(c.reader)
}
