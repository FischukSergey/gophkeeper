package handlers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/FischukSergey/gophkeeper/internal/app/interceptors/auth"
	"github.com/FischukSergey/gophkeeper/internal/models"
	pb "github.com/FischukSergey/gophkeeper/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ProtoNoteService интерфейс для методов сервера.
type ProtoNoteService interface {
	NoteAddService(ctx context.Context, note models.Note) error
	NoteGetListService(ctx context.Context, userID int64) ([]models.Note, error)
	NoteDeleteService(ctx context.Context, userID int64, noteID int64) error
}

// NoteServer сервер для методов заметки.
type NoteServer struct {
	pb.UnimplementedNoteServiceServer
	NoteService ProtoNoteService
}

// RegisterNoteAPI регистрация сервера.
func RegisterNoteAPI(
	server *grpc.Server,
	noteService ProtoNoteService,
) {
	pb.RegisterNoteServiceServer(server, &NoteServer{NoteService: noteService})
}

// NoteAdd хендлер для добавления заметки.
func (h *NoteServer) NoteAdd(ctx context.Context, req *pb.NoteAddRequest) (*pb.NoteAddResponse, error) {
	log.Info("NoteAdd", request, req)
	userID, ok := ctx.Value(auth.CtxKeyUserGrpc).(int)
	if !ok {
		return &pb.NoteAddResponse{Success: false}, status.Errorf(codes.Unauthenticated, models.UserIDNotFound)
	}
	log.Info(userFound, slog.Int(user, userID))
	//формируем массив метаданных
	var metadata []models.Metadata
	if req.Note.Metadata != nil {
		metadata = make([]models.Metadata, len(req.Note.Metadata))
		for i, m := range req.Note.Metadata {
			metadata[i] = models.Metadata{Key: m.Key, Value: m.Value}
		}
	}
	//формируем заметку
	note := models.Note{
		UserID:   int64(userID),
		NoteText: req.Note.NoteText,
		Metadata: metadata,
	}
	//добавляем заметку
	err := h.NoteService.NoteAddService(ctx, note)
	if err != nil {
		return &pb.NoteAddResponse{Success: false}, fmt.Errorf("ошибка при добавлении заметки: %w", err)
	}
	return &pb.NoteAddResponse{Success: true}, nil
}

// NoteGetList хендлер для получения списка заметок.
func (h *NoteServer) NoteGetList(ctx context.Context, req *pb.NoteGetListRequest) (
	*pb.NoteGetListResponse, error,
) {
	log.Info("NoteGetList", request, req)
	userID, ok := ctx.Value(auth.CtxKeyUserGrpc).(int)
	if !ok {
		return &pb.NoteGetListResponse{Notes: nil}, status.Errorf(codes.Unauthenticated, models.UserIDNotFound)
	}
	log.Info(userFound, slog.Int(user, userID))
	//получаем список заметок
	notes, err := h.NoteService.NoteGetListService(ctx, int64(userID))
	if err != nil {
		return &pb.NoteGetListResponse{Notes: nil}, fmt.Errorf("ошибка при получении списка заметок: %w", err)
	}
	//формируем ответ
	notesPb := make([]*pb.Note, len(notes))
	for i, n := range notes {
		metadataPb := make([]*pb.Metadata, len(n.Metadata))
		for j, m := range n.Metadata {
			metadataPb[j] = &pb.Metadata{Key: m.Key, Value: m.Value}
		}
		notesPb[i] = &pb.Note{
			NoteID:   n.NoteID,
			NoteText: n.NoteText,
			Metadata: metadataPb,
		}
	}
	return &pb.NoteGetListResponse{Notes: notesPb}, nil
}

// NoteDelete хендлер для удаления заметки.
func (h *NoteServer) NoteDelete(ctx context.Context, req *pb.NoteDeleteRequest) (*pb.NoteDeleteResponse, error) {
	log.Info("NoteDelete", request, req)
	userID, ok := ctx.Value(auth.CtxKeyUserGrpc).(int)
	if !ok {
		return &pb.NoteDeleteResponse{Success: false}, status.Errorf(codes.Unauthenticated, models.UserIDNotFound)
	}
	log.Info(userFound, slog.Int(user, userID))
	//удаляем заметку
	err := h.NoteService.NoteDeleteService(ctx, int64(userID), req.NoteID)
	if err != nil {
		return &pb.NoteDeleteResponse{Success: false}, fmt.Errorf("ошибка при удалении заметки: %w", err)
	}
	return &pb.NoteDeleteResponse{Success: true}, nil
}
