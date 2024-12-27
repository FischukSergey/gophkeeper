package converters

import (
	"github.com/FischukSergey/gophkeeper/internal/models"
	pb "github.com/FischukSergey/protos/gen/gophkeeper/gophkeeper"
)

// ToModelNote преобразует структуру pb.Note в модель Note.
func ToModelNote(protoNote *pb.Note, userID int64, metadata []models.Metadata) models.Note {
	return models.Note{
		UserID:   userID,
		NoteText: protoNote.NoteText,
		Metadata: metadata,
	}
}

// ToProtoNote преобразует модель Note в структуру pb.Note.
func ToProtoNote(note models.Note, metadataPb []*pb.Metadata) *pb.Note {
	return &pb.Note{
		NoteID:   note.NoteID,
		NoteText: note.NoteText,
		Metadata: metadataPb,
	}
}
