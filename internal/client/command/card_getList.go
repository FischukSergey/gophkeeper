package command

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/FischukSergey/gophkeeper/internal/app/interceptors/auth"
	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
	"github.com/FischukSergey/gophkeeper/internal/client/service"
	"github.com/manifoldco/promptui"
)

const cardGetListCommandName = "CardList"

func NewCommandCardGetList(
	cardService *service.CardService,
	token *grpcclient.Token,
	reader io.Reader,
	writer io.Writer,
) *CommandCardGetList {
	return &CommandCardGetList{
		cardService: cardService,
		token:       token,
		reader:      reader,
		writer:      writer,
	}
}

type CommandCardGetList struct {
	cardService *service.CardService
	token       *grpcclient.Token
	reader      io.Reader
	writer      io.Writer
}

// Name возвращает имя команды.
func (c *CommandCardGetList) Name() string {
	return cardGetListCommandName
}

// Execute выполняет команду получения списка карт.
func (c *CommandCardGetList) Execute() {
	// Проверка наличия токена
	if !checkToken(c.token, c.reader) {
		return // Выходим из функции если токен невалидный
	}
	cards, err := c.cardService.GetCardList(context.Background(), c.token.Token)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), auth.ErrNotFound) ||
			strings.Contains(err.Error(), auth.ErrInvalid) ||
			strings.Contains(err.Error(), auth.ErrExpired):
			fprintln(c.writer, errorAuth)
			return
		default:
			fmt.Println("Ошибка при получении списка карт:", err)
			return
		}
	}
	if len(cards) == 0 {
		fmt.Println("Нет карт для отображения")
		waitEnter(c.reader)
		return
	}

	// создаем новый tabwriter
	w := tabwriter.NewWriter(c.writer, 0, 0, 2, ' ', 0)

	// выводим заголовки таблицы
	_, err = fmt.Fprintln(w, "ID\tБанк\tНомер карты\tВладелец\tДата истечения\tCVV")
	if err != nil {
		fmt.Println("Ошибка при выводе списка карт:", err)
	}

	// выводим данные карт
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].CardID < cards[j].CardID
	})
	for _, card := range cards {
		_, err := fmt.Fprintf(
			w, "%d\t%s\t%s\t%s\t%s\t%s\n",
			card.CardID,
			card.CardBank,
			card.CardNumber,
			card.CardHolder,
			card.CardExpirationDate.Format("01/06"),
			card.CardCVV,
		)
		if err != nil {
			fmt.Println("Ошибка при выводе списка карт:", err)
		}
	}
	err = w.Flush()
	if err != nil {
		fmt.Println("Ошибка при выводе списка карт:", err)
	}

	//запрашиваем просмотр метаданных
	prompt := promptui.Prompt{
		Label: "Хотите просмотреть метаданные карт? (y/n)",
	}
	answer, err := prompt.Run()
	if err != nil {
		fmt.Println("Ошибка при вводе ответа:", err)
	}
	if answer == "y" || answer == "Y" {
		for _, card := range cards {
			if card.Metadata != "" {
				fprintf(c.writer, "\nМетаданные карты с ID %d:\n", card.CardID)
				//парсим метаданные
				metadata := make(map[string]interface{})
				err := json.Unmarshal([]byte(card.Metadata), &metadata)
				if err != nil {
					fmt.Println("Ошибка при парсинге метаданных карты:", err)
				}
				//построчно выводим метаданные
				for key, value := range metadata {
					fprintf(c.writer, "key: %s, \tvalue: %v\n", key, value)
				}
			}
		}
		waitEnter(c.reader)
	}
}
