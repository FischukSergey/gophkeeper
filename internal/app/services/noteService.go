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
	NoteAdd(ctx context.Context, note models.Note, metadata []byte) error
	NoteGetList(ctx context.Context, userID int64) ([]models.Note, error)
	NoteDelete(ctx context.Context, userID int64, noteID int64) error
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
	//добавляем заметку в базу данных
	err = s.storage.NoteAdd(ctx, note, metaJSON)
	if err != nil {
		return fmt.Errorf("error during adding note: %w", err)
	}
	return nil
}

// NoteGetListService функция для получения списка заметок.
func (s *NoteService) NoteGetListService(ctx context.Context, userID int64) ([]models.Note, error) {
	s.log.Info("NoteGetListService method called")
	notes, err := s.storage.NoteGetList(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notes from storage: %w", err)
	}

	// Обрабатываем метаданные для каждой заметки
	for i := range notes {
		if notes[i].RawMetadata != "" {
			var rawMetadata map[string]string
			err = json.Unmarshal([]byte(notes[i].RawMetadata), &rawMetadata)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}

			// Преобразуем map в slice of Metadata
			var metadataStruct []models.Metadata
			for k, v := range rawMetadata {
				metadataStruct = append(metadataStruct, models.Metadata{
					Key:   k,
					Value: v,
				})
			}
			notes[i].Metadata = metadataStruct
		}
	}
	return notes, nil
}

// NoteDeleteService функция для удаления заметки.
func (s *NoteService) NoteDeleteService(ctx context.Context, userID int64, noteID int64) error {
	s.log.Info("NoteDeleteService method called")
	err := s.storage.NoteDelete(ctx, userID, noteID)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}
	return nil
}
