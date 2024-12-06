package handlers

import (
	"context"
	"log/slog"

	"github.com/FischukSergey/gophkeeper/internal/app/converters"
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
	userID, err := h.validateUserID(ctx)
	if err != nil {
		return nil, err
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
	note := converters.ToModelNote(req.Note, int64(userID), metadata)
	//добавляем заметку
	err = h.NoteService.NoteAddService(ctx, note)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка при добавлении заметки: %v", err)
	}
	return &pb.NoteAddResponse{Success: true}, nil
}

// NoteGetList хендлер для получения списка заметок.
func (h *NoteServer) NoteGetList(ctx context.Context, req *pb.NoteGetListRequest) (
	*pb.NoteGetListResponse, error,
) {
	log.Info("NoteGetList", request, req)
	userID, err := h.validateUserID(ctx)
	if err != nil {
		return nil, err
	}
	log.Info(userFound, slog.Int(user, userID))
	//получаем список заметок
	notes, err := h.NoteService.NoteGetListService(ctx, int64(userID))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка при получении списка заметок: %v", err)
	}
	//формируем ответ
	notesPb := make([]*pb.Note, len(notes))
	for i, n := range notes {
		metadataPb := make([]*pb.Metadata, len(n.Metadata))
		for j, m := range n.Metadata {
			metadataPb[j] = &pb.Metadata{Key: m.Key, Value: m.Value}
		}
		notesPb[i] = converters.ToProtoNote(n, metadataPb)
	}
	return &pb.NoteGetListResponse{Notes: notesPb}, nil
}

// NoteDelete хендлер для удаления заметки.
func (h *NoteServer) NoteDelete(ctx context.Context, req *pb.NoteDeleteRequest) (*pb.NoteDeleteResponse, error) {
	log.Info("NoteDelete", request, req)
	userID, err := h.validateUserID(ctx)
	if err != nil {
		return nil, err
	}
	log.Info(userFound, slog.Int(user, userID))
	//удаляем заметку
	err = h.NoteService.NoteDeleteService(ctx, int64(userID), req.NoteID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка при удалении заметки: %v", err)
	}
	return &pb.NoteDeleteResponse{Success: true}, nil
}

// validateUserID проверяет корректность ID пользователя из контекста
func (h *NoteServer) validateUserID(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(auth.CtxKeyUserGrpc).(int)
	if !ok {
		return 0, status.Errorf(codes.Unauthenticated, models.UserIDNotFound)
	}
	if userID <= 0 {
		return 0, status.Errorf(codes.InvalidArgument, "invalid user ID")
	}
	return userID, nil
}
