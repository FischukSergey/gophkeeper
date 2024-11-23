package services

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/FischukSergey/gophkeeper/internal/models"
)

type MockNoteKeeper struct {
	mock.Mock
}

func (m *MockNoteKeeper) NoteAdd(ctx context.Context, note models.Note, jsonMetadata []uint8) error {
	args := m.Called(ctx, note, jsonMetadata)
	return args.Error(0)
}

func (m *MockNoteKeeper) NoteGetList(ctx context.Context, userID int64) ([]models.Note, error) {
	args := m.Called(ctx, userID)
	if userID == 1 {
		return []models.Note{
			{NoteText: "test", Metadata: []models.Metadata{{Key: "test", Value: "test"}}},
		}, nil
	}
	return nil, args.Error(1)
}

func (m *MockNoteKeeper) NoteDelete(ctx context.Context, userID int64, noteID int64) error {
	args := m.Called(ctx, userID, noteID)
	return args.Error(0)
}

func TestNoteService_NoteAdd(t *testing.T) {
	logger := slog.Default()
	ctx := context.Background()

	tests := []struct {
		name    string
		note    models.Note
		wantErr error
	}{
		{
			name: "successful add note",
			note: models.Note{
				NoteText: "test",
				Metadata: []models.Metadata{
					{Key: "test", Value: "test"},
				},
			},
			wantErr: nil,
		},
		{
			name: "invalid metadata key",
			note: models.Note{
				NoteText: "test error",
				Metadata: []models.Metadata{
					{Key: "", Value: "test"},
				},
			},
			wantErr: errors.New("invalid metadata"),
		},
		{
			name: "invalid metadata value",
			note: models.Note{
				NoteText: "test error",
				Metadata: []models.Metadata{
					{Key: "test", Value: ""},
				},
			},
			wantErr: errors.New("invalid metadata"),
		},
		{
			name: "duplicate metadata key",
			note: models.Note{
				NoteText: "test error",
				Metadata: []models.Metadata{
					{Key: "test", Value: "test"},
					{Key: "test", Value: "test"},
				},
			},
			wantErr: errors.New("duplicate metadata key"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockNoteKeeper := new(MockNoteKeeper)
			noteService := NewNoteService(logger, mockNoteKeeper)

			if tt.wantErr == nil {
				mockNoteKeeper.On("NoteAdd", mock.Anything, tt.note, mock.AnythingOfType("[]uint8")).
					Return(tt.wantErr)
			}

			err := noteService.NoteAddService(ctx, tt.note)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			mockNoteKeeper.AssertExpectations(t)
		})
	}
}
func TestNoteService_NoteGetList(t *testing.T) {
	logger := slog.Default()
	ctx := context.Background()
	tests := []struct {
		name      string
		userID    int64
		wantNotes []models.Note
		wantErr   error
	}{
		{
			name:   "successful get list",
			userID: 1,
			wantNotes: []models.Note{
				{NoteText: "test", Metadata: []models.Metadata{{Key: "test", Value: "test"}}},
			},
			wantErr: nil,
		},
		{
			name:      "empty list",
			userID:    2,
			wantNotes: nil,
			wantErr:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockNoteKeeper := new(MockNoteKeeper)
			noteService := NewNoteService(logger, mockNoteKeeper)

			if tt.wantErr == nil {
				mockNoteKeeper.On("NoteGetList", ctx, tt.userID).Return(tt.wantNotes, tt.wantErr)
			}

			notes, err := noteService.NoteGetListService(ctx, tt.userID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantNotes, notes)
			}
			mockNoteKeeper.AssertExpectations(t)
		})
	}
}
func TestNoteService_NoteDelete(t *testing.T) {
	logger := slog.Default()
	ctx := context.Background()

	tests := []struct {
		name    string
		userID  int64
		noteID  int64
		wantErr error
	}{
		{
			name:    "successful delete",
			userID:  1,
			noteID:  1,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockNoteKeeper := new(MockNoteKeeper)
			noteService := NewNoteService(logger, mockNoteKeeper)

			if tt.wantErr == nil {
				mockNoteKeeper.On("NoteDelete", ctx, tt.userID, tt.noteID).Return(tt.wantErr)
			}

			err := noteService.NoteDeleteService(ctx, tt.userID, tt.noteID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			mockNoteKeeper.AssertExpectations(t)
		})
	}
}
