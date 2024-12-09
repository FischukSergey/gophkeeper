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

// NoteAdd функция для добавления заметки.
func (s *NoteService) NoteAdd(ctx context.Context, note models.Note) error {
	s.log.Info("NoteAdd method called")
	//валидируем данные
	if note.NoteText == "" {
		return fmt.Errorf("invalid note data")
	}
	//проверяем наличие метаданных
	if len(note.Metadata) > 0 {
		//валидируем метаданные
		err := ValidateMetadata(note.Metadata)
		if err != nil {
			return err
		}
	}
	//сериализуем метаданные
	metaJSON, err := SerializeMetadata(note.Metadata)
	if err != nil {
		return err
	}
	//добавляем заметку в базу данных
	err = s.storage.NoteAdd(ctx, note, []byte(metaJSON))
	if err != nil {
		return err
	}
	return nil
}

// NoteGetListService функция для получения списка заметок.
func (s *NoteService) NoteGetList(ctx context.Context, userID int64) ([]models.Note, error) {
	s.log.Info("NoteGetList method called")
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
func (s *NoteService) NoteDelete(ctx context.Context, userID int64, noteID int64) error {
	s.log.Info("NoteDelete method called")
	err := s.storage.NoteDelete(ctx, userID, noteID)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}
	return nil
}