package handlers

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/FischukSergey/gophkeeper/internal/app/converters"
	"github.com/FischukSergey/gophkeeper/internal/models"
	pb "github.com/FischukSergey/gophkeeper/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ProtoCardService интерфейс для методов сервера.
type ProtoCardService interface {
	CardAdd(ctx context.Context, card models.Card) error
	CardGetList(ctx context.Context, userID int64) ([]models.Card, error)
	CardDelete(ctx context.Context, cardID int64) error
	CardAddMetadata(ctx context.Context, userID int64, cardID int64, metadata []models.Metadata) error
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
	userID, err := validateUserID(ctx)
	if err != nil {
		return nil, err
	}
	log.Info(userFound, slog.Int(user, userID))
	//формируем карту
	card := converters.ToModelCard(req.Card, strconv.Itoa(userID))
	//добавляем карту
	err = h.CardService.CardAdd(ctx, card)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка при добавлении карты: %v", err)
	}
	return &pb.CardAddResponse{}, nil
}

// CardGetList хендлер для получения списка карт пользователя.
func (h *CardServer) CardGetList(ctx context.Context, req *pb.CardGetListRequest) (*pb.CardGetListResponse, error) {
	log.Info("CardGetList", request, req)
	userID, err := validateUserID(ctx)
	if err != nil {
		return nil, err
	}
	log.Info(userFound, slog.Int(user, userID))

	cards, err := h.CardService.CardGetList(ctx, int64(userID))
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
	userID, err := validateUserID(ctx)
	if err != nil {
		return nil, err
	}
	log.Info(userFound, slog.Int(user, userID))
	err = h.CardService.CardDelete(ctx, req.CardID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка при удалении карты: %v", err)
	}
	return &pb.CardDeleteResponse{Success: true}, nil
}

// CardAddMetadata хендлер для добавления метаданных к карте.
func (h *CardServer) CardAddMetadata(ctx context.Context, req *pb.CardAddMetadataRequest) (
	*pb.CardAddMetadataResponse, error) {
	log.Info("CardAddMetadata", request, req)
	userID, err := validateUserID(ctx)
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
	err = h.CardService.CardAddMetadata(ctx, int64(userID), req.CardID, metadata)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка при добавлении метаданных: %v", err)
	}

	return &pb.CardAddMetadataResponse{Success: true}, nil
}
