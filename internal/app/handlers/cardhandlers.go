package handlers

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/FischukSergey/gophkeeper/internal/app/converters"
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
	CardGetListService(ctx context.Context, userID int64) ([]models.Card, error)
	CardDeleteService(ctx context.Context, cardID int64) error
	CardAddMetadataService(ctx context.Context, userID int64, cardID int64, metadata []models.Metadata) error
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
	log.Info("CardAdd", request, req)
	userID, ok := ctx.Value(auth.CtxKeyUserGrpc).(int)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, models.UserIDNotFound)
	}
	log.Info(userFound, slog.Int(user, userID))

	//формируем карту
	card := converters.ToModelCard(req.Card, strconv.Itoa(userID))
	//добавляем карту
	err := h.CardService.CardAddService(ctx, card)	
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка при добавлении карты: %v", err)
	}
	return &pb.CardAddResponse{}, nil
}

// CardGetList хендлер для получения списка карт пользователя.
func (h *CardServer) CardGetList(ctx context.Context, req *pb.CardGetListRequest) (*pb.CardGetListResponse, error) {
	log.Info("CardGetList", request, req)
	userID, err := h.validateUserID(ctx)
	if err != nil {
		return nil, err
	}
	log.Info(userFound, slog.Int(user, userID))

	cards, err := h.CardService.CardGetListService(ctx, int64(userID))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка при получении списка карт: %v", err)
	}

	cardsPb := make([]*pb.Card, len(cards))
	for i, card := range cards {
		cardsPb[i] = converters.ToProtoCard(card)
	}
	return &pb.CardGetListResponse{Cards: cardsPb}, nil
}

// CardDelete хендлер для удаления карты.
func (h *CardServer) CardDelete(ctx context.Context, req *pb.CardDeleteRequest) (*pb.CardDeleteResponse, error) {
	log.Info("CardDelete", request, req)
	userID, err := h.validateUserID(ctx)
	if err != nil {
		return nil, err
	}
	log.Info(userFound, slog.Int(user, userID))
	err = h.CardService.CardDeleteService(ctx, req.CardID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка при удалении карты: %v", err)
	}
	return &pb.CardDeleteResponse{Success: true}, nil
}

// CardAddMetadata хендлер для добавления метаданных к карте.
func (h *CardServer) CardAddMetadata(ctx context.Context, req *pb.CardAddMetadataRequest) (
	*pb.CardAddMetadataResponse, error) {
	log.Info("CardAddMetadata", request, req)
	userID, err := h.validateUserID(ctx)
	if err != nil {
		return nil, err
	}
	log.Info(userFound, slog.Int(user, userID))
	//формируем массив метаданных
	metadata := make([]models.Metadata, len(req.Metadata))
	for i, m := range req.Metadata {
		metadata[i] = models.Metadata{Key: m.Key, Value: m.Value}
	}
	log.Info("metadata", "metadata", metadata)
	err = h.CardService.CardAddMetadataService(ctx, int64(userID), req.CardID, metadata)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка при добавлении метаданных: %v", err)
	}

	return &pb.CardAddMetadataResponse{Success: true}, nil
}

// validateUserID проверяет корректность ID пользователя из контекста
func (h *CardServer) validateUserID(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(auth.CtxKeyUserGrpc).(int)
	if !ok {
		return 0, status.Errorf(codes.Unauthenticated, models.UserIDNotFound)
	}
	if userID <= 0 {
		return 0, status.Errorf(codes.InvalidArgument, "invalid user ID")
	}
	return userID, nil
}
