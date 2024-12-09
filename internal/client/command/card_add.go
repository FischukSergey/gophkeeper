package command

import (
	"bufio"
	"context"
	"io"
	"strings"
	"time"

	"github.com/FischukSergey/gophkeeper/internal/app/interceptors/auth"
	"github.com/FischukSergey/gophkeeper/internal/client/grpcclient"
	"github.com/FischukSergey/gophkeeper/internal/models"
)

const nameCommandCardAdd = "CardAdd"
const errorAuth = "Ошибка авторизации. Пожалуйста, войдите в систему заново"

type CommandCardAdd struct {
	cardAddService ICardService
	token          *grpcclient.Token
	reader         io.Reader
	writer         io.Writer
}

// NewCommandCardAdd создание новой команды добавления карты.
func NewCommandCardAdd(
	cardAddService ICardService,
	token *grpcclient.Token,
	reader io.Reader,
	writer io.Writer,
) *CommandCardAdd {
	return &CommandCardAdd{
		cardAddService: cardAddService,
		token:          token,
		reader:         reader,
		writer:         writer,
	}
}

func (c *CommandCardAdd) Name() string {
	return nameCommandCardAdd
}

func (c *CommandCardAdd) Execute() {
	// Проверка наличия токена
	if !checkToken(c.token, c.reader) {
		return
	}

	fprintln(c.writer, "Добавление карты")

	var card models.Card
	var input string

	for {
		fprintln(c.writer, "\nВыберите действие:")
		fprintln(c.writer, "1. Ввести название банка \033[32m"+card.CardBank+"\033[0m")
		fprintln(c.writer, "2. Ввести номер карты (16 цифр) \033[32m"+card.CardNumber+"\033[0m")
		fprintln(c.writer, "3. Ввести срок действия (MM/YY) \033[32m"+
			time.Unix(card.CardExpirationDate.Unix(), 0).Format("01/06")+"\033[0m")
		fprintln(c.writer, "4. Ввести CVV (3 цифры) \033[32m"+card.CardCVV+"\033[0m")
		fprintln(c.writer, "5. Ввести владельца карты \033[32m"+card.CardHolder+"\033[0m")
		fprintln(c.writer, "6. Сохранить и выйти")
		fprintln(c.writer, "0. Отмена")

		fscanln(c.reader, &input)

		switch input {
		case "1":
			fprint(c.writer, "Введите название банка: ")
			fscanln(c.reader, &card.CardBank)
		case "2":
			fprint(c.writer, "Введите номер карты (16 цифр): ")
			fscanln(c.reader, &card.CardNumber)
		case "3":
			fprint(c.writer, "Введите срок действия (MM/YY): ")
			var expirationDate string
			fscanln(c.reader, &expirationDate)
			// Parse expiration date
			t, err := time.Parse("01/06", expirationDate)
			if err != nil {
				fprintln(c.writer, "Неверный формат даты. Используйте MM/YY")
				continue
			}
			card.CardExpirationDate = t
		case "4":
			fprint(c.writer, "Введите CVV (3 цифры): ")
			fscanln(c.reader, &card.CardCVV)
		case "5":
			fprint(c.writer, "Введите владельца карты: ")
			reader := bufio.NewReader(c.reader)
			cardHolder, _ := reader.ReadString('\n')
			card.CardHolder = strings.TrimSpace(cardHolder) // убираем символ новой строки
		case "6":
			if card.CardBank == "" ||
				len(card.CardNumber) != 16 ||
				card.CardExpirationDate.IsZero() ||
				len(card.CardCVV) != 3 ||
				card.CardHolder == "" {
				fprintln(c.writer, "Необходимо правильно заполнить все поля!")
				//ожидание нажатия любой клавиши
				fscanln(c.reader)
				continue
			}
			err := c.cardAddService.CardAdd(context.Background(), card, c.token.Token)
			if err != nil {
				switch {
				case strings.Contains(err.Error(), auth.ErrNotFound) ||
					strings.Contains(err.Error(), auth.ErrInvalid) ||
					strings.Contains(err.Error(), auth.ErrExpired):
					fprintln(c.writer, errorAuth)
					return
				default:
					fprintf(c.writer, "Ошибка при добавлении карты: %v\n", err)
				}
				return
			}
			fprintln(c.writer, "Карта успешно добавлена")
			return
		case "0":
			fprintln(c.writer, "Операция отменена")
			return
		default:
			fprintln(c.writer, "Неверный выбор. Попробуйте снова")
		}
	}
}
