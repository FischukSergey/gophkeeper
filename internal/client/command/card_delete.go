package command

import (
	"context"
	"io"

	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
)

const cardDeleteCommandName = "CardDelete"

// CommandCardDelete структура для команды удаления карты.
type CommandCardDelete struct {
	cardService ICardService
	token       *grpcclient.Token
	reader      io.Reader
	writer      io.Writer
}

// NewCommandCardDelete создает новый экземпляр команды удаления карты.
func NewCommandCardDelete(
	cardService ICardService,
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
	cards, err := getCardList(c.cardService, c.token, c.reader, c.writer)
	if err != nil {
		fprintln(c.writer, "Ошибка при получении списка карт:", err)
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

	//удаляем карту
	err = c.cardService.DeleteCard(context.Background(), int64(cardIDInt), c.token.Token)
	if err != nil {
		fprintln(c.writer, "Ошибка при удалении карты:", err)
		return
	}
	fprintln(c.writer, "Карта удалена")
	waitEnter(c.reader)
}
