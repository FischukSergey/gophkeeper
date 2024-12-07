package service

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/FischukSergey/gophkeeper/internal/lib/luhn"
	"github.com/FischukSergey/gophkeeper/internal/models"
	pb "github.com/FischukSergey/gophkeeper/internal/proto"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CardService сервис карт.
type CardService struct {
	client pb.CardServiceClient
	log    *slog.Logger
}

// NewCardService создание сервиса карт.
func NewCardService(client pb.CardServiceClient, log *slog.Logger) *CardService {
	return &CardService{client: client, log: log}
}

// CardAdd метод для добавления карты.
func (s *CardService) CardAdd(ctx context.Context, card models.Card, token string) error {
	s.log.Info("Service CardAdd method called")
	//валидация карты
	if card.CardBank == "" ||
		card.CardNumber == "" ||
		card.CardExpirationDate.IsZero() ||
		card.CardCVV == "" ||
		card.CardHolder == "" {
		return fmt.Errorf("необходимо заполнить все поля")
	}
	card.CardNumber = strings.ReplaceAll(card.CardNumber, "-", "")
	//валидируем номер карты
	if !luhn.Valid(card.CardNumber) || len(card.CardNumber) != 16 {
		return fmt.Errorf("неверный номер карты")
	}
	card.CardHolder = strings.ToUpper(card.CardHolder)
	//валидируем CVV
	if len(card.CardCVV) != 3 {
		return fmt.Errorf("неверный CVV")
	}
	//валидируем дату
	if card.CardExpirationDate.IsZero() || card.CardExpirationDate.Before(time.Now()) {
		return fmt.Errorf("неверная дата")
	}
	// добавление токена авторизации в контекст
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(sessionToken, token))
	//добавляем карту на сервер
	_, err := s.client.CardAdd(ctx, &pb.CardAddRequest{
		Card: &pb.Card{
			CardBank:           card.CardBank,
			CardNumber:         card.CardNumber,
			CardHolder:         card.CardHolder,
			CardExpirationDate: timestamppb.New(card.CardExpirationDate),
			CardCVV:            card.CardCVV,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to add card: %w", err)
	}
	return nil
}

// GetCardList метод для получения списка карт.
func (s *CardService) GetCardList(ctx context.Context, token string) ([]models.Card, error) {
	// добавление токена авторизации в контекст
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(sessionToken, token))
	// получение списка карт
	response, err := s.client.CardGetList(ctx, &pb.CardGetListRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get card list: %w", err)
	}
	cards := make([]models.Card, 0, len(response.GetCards()))
	for _, card := range response.GetCards() {
		cards = append(cards, models.Card{
			CardBank:           card.GetCardBank(),
			CardNumber:         card.GetCardNumber(),
			CardHolder:         card.GetCardHolder(),
			CardExpirationDate: card.GetCardExpirationDate().AsTime(),
			CardCVV:            card.GetCardCVV(),
			CardID:             card.GetCardID(),
			Metadata:           card.GetMetadata(),
		})
	}
	s.log.Info("Card list received", "cards", cards)
	return cards, nil
}

// DeleteCard метод для удаления карты.
func (s *CardService) DeleteCard(ctx context.Context, cardID string, token string) error {
	// добавление токена авторизации в контекст
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(sessionToken, token))
	// удаление карты
	cardIDInt, err := strconv.ParseInt(cardID, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse card ID: %w", err)
	}
	_, err = s.client.CardDelete(ctx, &pb.CardDeleteRequest{CardID: cardIDInt})
	if err != nil {
		return fmt.Errorf("failed to delete card: %w", err)
	}
	return nil
}

// AddCardMetadata метод для добавления метаданных к карте.
func (s *CardService) AddCardMetadata(
	ctx context.Context, cardID int64, metaData map[string]string, token string) error {
	// добавление токена авторизации в контекст
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(sessionToken, token))
	//преобразование map в массив структур Metadata
	jsonData := make([]*pb.Metadata, 0, len(metaData))
	for key, value := range metaData {
		jsonData = append(jsonData, &pb.Metadata{Key: key, Value: value})
	}
	s.log.Info("Metadata to add", "metadata", jsonData)
	// добавление метаданных к карте
	_, err := s.client.CardAddMetadata(ctx, &pb.CardAddMetadataRequest{
		CardID:   cardID,
		Metadata: jsonData,
	})
	if err != nil {
		return fmt.Errorf("failed to add card metadata: %w", err)
	}
	return nil
}
