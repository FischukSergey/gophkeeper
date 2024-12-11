package tests

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/FischukSergey/gophkeeper/internal/models"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// тест на добавление карты.
func TestCardAdd(t *testing.T) {
	login := gofakeit.Username()
	password := gofakeit.Password(true, true, true, false, false, 10)
	//регистрируем пользователя
	token, err := authService.Register(context.Background(), login, password)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	//создаем табличные тесты на проверку обработки ошибок
	tests := []struct {
		name    string
		card    models.Card
		token   string
		wantErr bool
	}{
		{name: "empty token", card: models.Card{}, token: "", wantErr: true},
		{name: "empty card", card: models.Card{}, token: token, wantErr: true},
		{name: "empty card number", card: models.Card{CardNumber: ""}, token: token, wantErr: true},
		{name: "empty card holder", card: models.Card{CardHolder: ""}, token: token, wantErr: true},
		{name: "empty card expiration date", card: models.Card{CardExpirationDate: time.Time{}}, token: token, wantErr: true},
		{name: "empty card cvv", card: models.Card{CardCVV: ""}, token: token, wantErr: true},
		{name: "invalid card cvv", card: models.Card{
			CardCVV:            "1234",
			CardNumber:         "5272697132101976",
			CardHolder:         "John Doe",
			CardExpirationDate: time.Now().AddDate(1, 0, 0),
			CardBank:           "Tinkoff",
		}, token: token, wantErr: true},
		{name: "invalid card expiration date", card: models.Card{
			CardExpirationDate: time.Now().AddDate(-1, 0, 0),
			CardNumber:         "5272697132101976",
			CardHolder:         "John Doe",
			CardCVV:            "123",
			CardBank:           "Tinkoff",
		}, token: token, wantErr: true},
		{name: "invalid number", card: models.Card{
			CardNumber:         "5272697132101970",
			CardHolder:         "John Doe",
			CardExpirationDate: time.Now().AddDate(1, 0, 0),
			CardCVV:            "123",
			CardBank:           "Tinkoff",
		}, token: token, wantErr: true},
		{name: "success", card: models.Card{
			CardNumber:         "5272697132101976",
			CardHolder:         "John Doe",
			CardExpirationDate: time.Now().AddDate(1, 0, 0),
			CardCVV:            "123",
			CardBank:           "Tinkoff",
		}, token: token, wantErr: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := cardService.CardAdd(context.Background(), test.card, test.token)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// тест на получение списка карт.
func TestCardGetList(t *testing.T) {
	login := gofakeit.Username()
	password := gofakeit.Password(true, true, true, false, false, 10)
	//регистрируем пользователя
	token, err := authService.Register(context.Background(), login, password)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	//создаем карту
	card1 := models.Card{
		CardNumber:         "5272697132101976",
		CardHolder:         "John Doe",
		CardExpirationDate: time.Now().AddDate(1, 0, 0),
		CardCVV:            "123",
		CardBank:           "Tinkoff",
	}
	err = cardService.CardAdd(context.Background(), card1, token)
	require.NoError(t, err)
	//получаем список карт
	var notes []models.Card
	notes, err = cardService.GetCardList(context.Background(), token)
	require.NoError(t, err)
	require.NotEmpty(t, notes)
	//проверяем, что карты совпадают
	assert.Equal(t, card1.CardNumber, notes[0].CardNumber)
	assert.Equal(t, strings.ToUpper(card1.CardHolder), notes[0].CardHolder)
	assert.Equal(t, card1.CardCVV, notes[0].CardCVV)
	assert.Equal(t, card1.CardBank, notes[0].CardBank)
}

// тест на удаление карты.
func TestCardDelete(t *testing.T) {
	login := gofakeit.Username()
	password := gofakeit.Password(true, true, true, false, false, 10)
	//регистрируем пользователя
	token, err := authService.Register(context.Background(), login, password)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	//создаем карту
	card1 := models.Card{
		CardNumber:         "5272697132101976",
		CardHolder:         "John Doe",
		CardExpirationDate: time.Now().AddDate(1, 0, 0),
		CardCVV:            "123",
		CardBank:           "Tinkoff",
	}
	err = cardService.CardAdd(context.Background(), card1, token)
	require.NoError(t, err)
	//получаем список карт
	var notes []models.Card
	notes, err = cardService.GetCardList(context.Background(), token)
	require.NoError(t, err)
	require.NotEmpty(t, notes)
	//удаляем карту
	err = cardService.DeleteCard(context.Background(), notes[0].CardID, token)
	require.NoError(t, err)
	//проверяем, что карты удалены
	notes, err = cardService.GetCardList(context.Background(), token)
	require.NoError(t, err)
	require.Empty(t, notes)
}

// тест на добавление метаданных к карте.
func TestCardAddMetadata(t *testing.T) {
	login := gofakeit.Username()
	password := gofakeit.Password(true, true, true, false, false, 10)
	//регистрируем пользователя
	token, err := authService.Register(context.Background(), login, password)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	//создаем карту
	card1 := models.Card{
		CardNumber:         "5272697132101976",
		CardHolder:         "John Doe",
		CardExpirationDate: time.Now().AddDate(1, 0, 0),
		CardCVV:            "123",
		CardBank:           "Tinkoff",
	}
	metadata := make(map[string]string)
	metadata["key1"] = "value1"
	metadata["key2"] = "value2"
	//создаем карту
	err = cardService.CardAdd(context.Background(), card1, token)
	require.NoError(t, err)
	//получаем список карт
	var notes []models.Card
	notes, err = cardService.GetCardList(context.Background(), token)
	require.NoError(t, err)
	require.NotEmpty(t, notes)

	//добавляем метаданные
	err = cardService.AddCardMetadata(context.Background(), notes[0].CardID, metadata, token)
	require.NoError(t, err)

	// Получаем обновленный список карт для проверки метаданных
	notes, err = cardService.GetCardList(context.Background(), token)
	require.NoError(t, err)
	require.NotEmpty(t, notes)

	// Проверяем, что метаданные были успешно добавлены
	require.NotEmpty(t, notes[0].Metadata, "Metadata should not be empty")
	// Дальнейшие проверки зависят от формата хранения метаданных
}
