package converters

import (
	"github.com/FischukSergey/gophkeeper/internal/models"
	pb "github.com/FischukSergey/gophkeeper/internal/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToProtoFile преобразует модель File в структуру pb.File
func ToProtoFile(file models.File) *pb.File {
	return &pb.File{
		FileID:    file.FileID,
		UserID:    file.UserID,
		Filename:  file.Filename,
		Size:      file.Size,
		CreatedAt: timestamppb.New(file.CreatedAt),
	}
}
