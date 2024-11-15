package command

import (
	"context"
	"io"

	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
	"github.com/FischukSergey/gophkeeper/internal/client/service"
)

const cardDeleteCommandName = "cardDelete"

// CommandCardDelete структура для команды удаления карты.
type CommandCardDelete struct {
	cardService *service.CardService
	token       *grpcclient.Token
	reader      io.Reader
	writer      io.Writer
}

// NewCommandCardDelete создает новый экземпляр команды удаления карты.
func NewCommandCardDelete(cardService *service.CardService, token *grpcclient.Token, reader io.Reader, writer io.Writer) *CommandCardDelete {
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
	fprintf(c.writer, "Введите ID карты для удаления: ")

	var cardID string
	fscanln(c.reader, &cardID)

	//удаляем карту
	err := c.cardService.DeleteCard(context.Background(), cardID, c.token.Token)
	if err != nil {
		fprintln(c.writer, "Ошибка при удалении карты:", err)
		return
	}
	fprintln(c.writer, "Карта удалена")
	waitEnter(c.reader)
}
