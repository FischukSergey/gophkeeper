package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/FischukSergey/gophkeeper/internal/app/interceptors/auth"
	"github.com/FischukSergey/gophkeeper/internal/models"
	pb "github.com/FischukSergey/gophkeeper/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProtoCardService интерфейс для методов сервера.
type ProtoCardService interface {
	CardAddService(ctx context.Context, card models.Card) error
	CardGetListService(ctx context.Context, userID int64) ([]models.Card, error)
	CardDeleteService(ctx context.Context, cardID int64) error
}

type CardServer struct {
	pb.UnimplementedCardServiceServer
	CardService ProtoCardService
}

// RegisterCardAPI регистрация сервера.
func RegisterCardAPI(
	server *grpc.Server,
	cardService ProtoCardService,
) {
	pb.RegisterCardServiceServer(server, &CardServer{CardService: cardService})
}

// CardAdd хендлер для добавления карты.
func (h *CardServer) CardAdd(ctx context.Context, req *pb.CardAddRequest) (*pb.CardAddResponse, error) {
	log.Info("CardAdd", "req", req)
	userID, ok := ctx.Value(auth.CtxKeyUserGrpc).(int)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, models.UserIDNotFound)
	}
	log.Info("userID found", slog.Int("userID", userID))

	//формируем карту
	card := models.Card{
		UserID:             strconv.Itoa(userID),
		CardBank:           req.Card.CardBank,
		CardNumber:         req.Card.CardNumber,
		CardHolder:         req.Card.CardHolder,
		CardExpirationDate: req.Card.CardExpirationDate.AsTime(),
		CardCVV:            req.Card.CardCVV,
	}
	err := h.CardService.CardAddService(ctx, card)
	if err != nil {
		return nil, fmt.Errorf("ошибка при добавлении карты: %w", err)
	}
	return &pb.CardAddResponse{}, nil
}

// CardGetList хендлер для получения списка карт пользователя.
func (h *CardServer) CardGetList(ctx context.Context, req *pb.CardGetListRequest) (*pb.CardGetListResponse, error) {
	log.Info("CardGetList", "req", req)
	userID, ok := ctx.Value(auth.CtxKeyUserGrpc).(int)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, models.UserIDNotFound)
	}
	log.Info("userID found", slog.Int("userID", userID))

	cards, err := h.CardService.CardGetListService(ctx, int64(userID))
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка карт: %w", err)
	}

	cardsPb := make([]*pb.Card, len(cards))
	for i, card := range cards {
		cardsPb[i] = &pb.Card{
			CardID:             card.CardID,
			CardBank:           card.CardBank,
			CardNumber:         card.CardNumber,
			CardHolder:         card.CardHolder,
			CardExpirationDate: timestamppb.New(card.CardExpirationDate),
			CardCVV:            card.CardCVV,
		}
	}
	return &pb.CardGetListResponse{Cards: cardsPb}, nil
}

// CardDelete хендлер для удаления карты.
func (h *CardServer) CardDelete(ctx context.Context, req *pb.CardDeleteRequest) (*pb.CardDeleteResponse, error) {
	log.Info("CardDelete", "req", req)
	userID, ok := ctx.Value(auth.CtxKeyUserGrpc).(int)
	if !ok {
		return &pb.CardDeleteResponse{Success: false}, status.Errorf(codes.Unauthenticated, models.UserIDNotFound)
	}
	log.Info("userID found", slog.Int("userID", userID))
	err := h.CardService.CardDeleteService(ctx, req.CardID)
	if err != nil {
		return &pb.CardDeleteResponse{Success: false}, fmt.Errorf("ошибка при удалении карты: %w", err)
	}
	return &pb.CardDeleteResponse{Success: true}, nil
}
