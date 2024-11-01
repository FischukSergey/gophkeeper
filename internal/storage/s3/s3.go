package s3

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3 структура для S3.
type S3 struct {
	S3Session *session.Session
}

// UploadFile загружает файл в S3 bucket и возвращает URL загруженного файла
func (s *S3) S3UploadFile(ctx context.Context, fileData []byte, filename string, bucket string) (string, error) {
	svc := s3.New(s.S3Session)
	// Проверяем существование бакета
	_, err := svc.HeadBucketWithContext(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		// Создаем новый приватный бакет
		_, err = svc.CreateBucketWithContext(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucket),
	})
		if err != nil {
			return "", fmt.Errorf("failed to create bucket: %v", err)
		}
	}
	slog.Info("bucket created", slog.String("bucket", bucket))
	
	// Подготавливаем параметры загрузки
	uploadInput := &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(fileData),
		ACL:    aws.String("public-read"),
	}

	// Загружаем файл
	result, err := s3manager.NewUploader(s.S3Session).Upload(uploadInput)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	// Возвращаем URL загруженного файла
	return result.Location, nil
}

// Close закрытие соединения с S3.
func (s *S3) Close() error {
	s.S3Session = nil
	return nil
}
