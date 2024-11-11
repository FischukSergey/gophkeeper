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
)

// ProtoCardService интерфейс для методов сервера.
type ProtoCardService interface {
	CardAddService(ctx context.Context, card models.Card) error
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

	return &pb.CardGetListResponse{}, nil
}
