package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/FischukSergey/gophkeeper/internal/models"
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
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(sessionToken, token))
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

// NoteGetList метод для получения списка заметок.
func (s *NoteService) NoteGetList(ctx context.Context, token string) ([]models.Note, error) {
	s.log.Info("Service NoteGetList method called")
	// добавление токена авторизации в контекст
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(sessionToken, token))
	// получение списка заметок с сервера
	resp, err := s.client.NoteGetList(ctx, &pb.NoteGetListRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get note list: %w", err)
	}
	s.log.Debug("resp", "resp", resp.GetNotes())
	// преобразование ответа сервера в список заметок
	notes := make([]models.Note, 0, len(resp.GetNotes()))
	for _, note := range resp.GetNotes() {
		metadata := make([]models.Metadata, 0, len(note.GetMetadata()))
		for _, meta := range note.GetMetadata() {
			metadata = append(metadata, ProtoToMetadata(meta))
		}
		notes = append(notes, ProtoToNote(note, metadata))
	}
	return notes, nil
}

// NoteDelete метод для удаления заметки.
func (s *NoteService) NoteDelete(ctx context.Context, noteID int64, token string) error {
	s.log.Info("Service NoteDelete method called")
	// добавление токена авторизации в контекст
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(sessionToken, token))
	// удаление заметки на сервере
	resp, err := s.client.NoteDelete(ctx, &pb.NoteDeleteRequest{NoteID: noteID})
	if resp.GetSuccess() || err == nil {
		return nil
	}
	return fmt.Errorf("failed to delete note: %w", err)
}
