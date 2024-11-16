package services

import (
	"context"
	"encoding/json"
	"fmt"

	"log/slog"

	"github.com/FischukSergey/gophkeeper/internal/models"
)

// NoteKeeper интерфейс для сервиса заметки.
type NoteKeeper interface {
	NoteAdd(ctx context.Context, note models.Note, metadata string) error
}

// NoteService структура для сервиса заметки.
type NoteService struct {
	log     *slog.Logger
	storage NoteKeeper
}

// NewNoteService функция для создания сервиса заметки.
func NewNoteService(log *slog.Logger, storage NoteKeeper) *NoteService {
	return &NoteService{log: log, storage: storage}
}

// NoteAddService функция для добавления заметки.
func (s *NoteService) NoteAddService(ctx context.Context, note models.Note) error {
	s.log.Info("NoteAddService method called")
	//валидируем данные
	if note.NoteText == "" {
		return fmt.Errorf("invalid note data")
	}
	metaMap := make(map[string]string)
	var metadata string
	//проверяем наличие метаданных
	if len(note.Metadata) > 0 {
		//валидируем метаданные
		for _, m := range note.Metadata {
			if m.Key == "" || m.Value == "" {
				return fmt.Errorf("invalid metadata")
			}
			if _, ok := metaMap[m.Key]; ok {
				return fmt.Errorf("duplicate metadata key: %s", m.Key)
			}
			metaMap[m.Key] = m.Value
		}
	}
	//сериализуем метаданные
	metaJSON, err := json.Marshal(metaMap)
	if err != nil {
		return fmt.Errorf("error during metadata serialization: %w", err)
	}
	metadata = string(metaJSON)
	//добавляем заметку в базу данных
	err = s.storage.NoteAdd(ctx, note, metadata)
	if err != nil {
		return fmt.Errorf("error during adding note: %w", err)
	}
	return nil
}
