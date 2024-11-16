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
	log.Info("NoteAdd", "req", req)
	userID, ok := ctx.Value(auth.CtxKeyUserGrpc).(int)
	if !ok {
		return &pb.NoteAddResponse{Success: false}, status.Errorf(codes.Unauthenticated, models.UserIDNotFound)
	}
	log.Info("userID found", slog.Int("userID", userID))
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
