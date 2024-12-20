package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/FischukSergey/gophkeeper/internal/models"
	pb "github.com/FischukSergey/gophkeeper/internal/proto"
	"google.golang.org/grpc/metadata"
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
	if err := ValidateCard(card); err != nil {
		return err
	}
	// добавление токена авторизации в контекст
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(sessionToken, token))
	//добавляем карту на сервер
	_, err := s.client.CardAdd(ctx, &pb.CardAddRequest{
		Card: CardToProto(card),
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
		cards = append(cards, ProtoToCard(card))
	}
	s.log.Info("Card list received", "cards", cards)
	return cards, nil
}

// DeleteCard метод для удаления карты.
func (s *CardService) DeleteCard(ctx context.Context, cardID int64, token string) error {
	// добавление токена авторизации в контекст
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(sessionToken, token))
	// удаление карты
	_, err := s.client.CardDelete(ctx, &pb.CardDeleteRequest{CardID: cardID})
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
	jsonData := MapToMetadata(metaData)
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
