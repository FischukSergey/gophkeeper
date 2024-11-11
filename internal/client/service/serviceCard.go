package service

import (
	"context"
	"fmt"
	"log/slog"
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
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("session_token", token))
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
