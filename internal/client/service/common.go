package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/FischukSergey/gophkeeper/internal/lib/luhn"
	"github.com/FischukSergey/gophkeeper/internal/models"
	pb "github.com/FischukSergey/protos/gen/gophkeeper/gophkeeper"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProtoToModel преобразует proto в модель.
func ProtoToModel(proto *pb.File) models.File {
	return models.File{
		FileID:    proto.GetFileID(),
		UserID:    proto.GetUserID(),
		Filename:  proto.GetFilename(),
		CreatedAt: proto.GetCreatedAt().AsTime(),
		DeletedAt: proto.GetDeletedAt().AsTime(),
		Size:      proto.GetSize(),
	}
}

// ProtoToMetadata преобразует proto в модель.
func ProtoToMetadata(proto *pb.Metadata) models.Metadata {
	return models.Metadata{
		Key:   proto.GetKey(),
		Value: proto.GetValue(),
	}
}

// ProtoToNote преобразует proto в модель.
func ProtoToNote(proto *pb.Note, metadata []models.Metadata) models.Note {
	return models.Note{
		NoteID:   proto.GetNoteID(),
		NoteText: proto.GetNoteText(),
		Metadata: metadata,
	}
}

// CardToProto преобразует модель в proto.
func CardToProto(card models.Card) *pb.Card {
	return &pb.Card{
		CardBank:           card.CardBank,
		CardNumber:         card.CardNumber,
		CardHolder:         card.CardHolder,
		CardExpirationDate: timestamppb.New(card.CardExpirationDate),
		CardCVV:            card.CardCVV,
	}
}

// ProtoToCard преобразует proto в модель.
func ProtoToCard(proto *pb.Card) models.Card {
	return models.Card{
		CardBank:           proto.GetCardBank(),
		CardNumber:         proto.GetCardNumber(),
		CardHolder:         proto.GetCardHolder(),
		CardExpirationDate: proto.GetCardExpirationDate().AsTime(),
		CardCVV:            proto.GetCardCVV(),
		CardID:             proto.GetCardID(),
		Metadata:           proto.GetMetadata(),
	}
}

// ValidateCard валидация карты.
func ValidateCard(card models.Card) error {
	if card.CardBank == "" ||
		card.CardNumber == "" ||
		card.CardExpirationDate.IsZero() ||
		card.CardCVV == "" ||
		card.CardHolder == "" {
		return fmt.Errorf("необходимо заполнить все поля")
	}
	card.CardNumber = strings.ReplaceAll(card.CardNumber, "-", "")
	//валидируем номер карты
	if !luhn.Valid(card.CardNumber) || len(card.CardNumber) != 16 {
		return fmt.Errorf("неверный номер карты")
	}
	card.CardHolder = strings.ToUpper(card.CardHolder)
	//валидируем CVV
	if len(card.CardCVV) != 3 {
		return fmt.Errorf("неверный CVV")
	}
	//валидируем дату
	if card.CardExpirationDate.IsZero() || card.CardExpirationDate.Before(time.Now()) {
		return fmt.Errorf("неверная дата")
	}
	return nil
}

// MapToMetadata преобразует map в массив структур Metadata.
func MapToMetadata(metaData map[string]string) []*pb.Metadata {
	jsonData := make([]*pb.Metadata, 0, len(metaData))
	for key, value := range metaData {
		jsonData = append(jsonData, &pb.Metadata{Key: key, Value: value})
	}
	return jsonData
}
