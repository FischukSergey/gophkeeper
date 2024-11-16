package service

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/FischukSergey/gophkeeper/internal/proto"
	"google.golang.org/grpc/metadata"
)

// NoteService сервис для работы с заметками.
type NoteService struct {
	client pb.NoteServiceClient
	log    *slog.Logger
}

// NewNoteService создание сервиса для работы с заметками.
func NewNoteService(client pb.NoteServiceClient, log *slog.Logger) *NoteService {
	return &NoteService{client: client, log: log}
}

// NoteAdd метод для добавления заметки.
func (s *NoteService) NoteAdd(ctx context.Context, note string, metaData map[string]string, token string) error {
	s.log.Info("Service NoteAdd method called")
	// добавление токена авторизации в контекст
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("session_token", token))
	//преобразование map в массив структур Metadata
	jsonData := make([]*pb.Metadata, 0, len(metaData))
	for key, value := range metaData {
		jsonData = append(jsonData, &pb.Metadata{Key: key, Value: value})
	}
	// добавление заметки на сервер
	resp, err := s.client.NoteAdd(ctx, &pb.NoteAddRequest{
		Note: &pb.Note{
			NoteText: note,
			Metadata: jsonData,
		},
	})
	if resp.GetSuccess() || err == nil {
		return nil
	}
	return fmt.Errorf("failed to add note: %w", err)
}
