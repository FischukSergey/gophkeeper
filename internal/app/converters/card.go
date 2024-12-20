package converters

import (
	"github.com/FischukSergey/gophkeeper/internal/models"

	pb "github.com/FischukSergey/gophkeeper/internal/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToProtoCard преобразует модель Card в структуру pb.Card.
func ToProtoCard(card models.Card) *pb.Card {
	return &pb.Card{
		CardID:             card.CardID,
		CardBank:           card.CardBank,
		CardNumber:         card.CardNumber,
		CardHolder:         card.CardHolder,
		CardExpirationDate: timestamppb.New(card.CardExpirationDate),
		CardCVV:            card.CardCVV,
		Metadata:           card.Metadata,
	}
}

// ToModelCard преобразует структуру pb.Card в модель Card.
func ToModelCard(protoCard *pb.Card, userID string) models.Card {
	return models.Card{
		UserID:             userID,
		CardID:             protoCard.CardID,
		CardBank:           protoCard.CardBank,
		CardNumber:         protoCard.CardNumber,
		CardHolder:         protoCard.CardHolder,
		CardExpirationDate: protoCard.CardExpirationDate.AsTime(),
		CardCVV:            protoCard.CardCVV,
		Metadata:           protoCard.Metadata,
	}
}
