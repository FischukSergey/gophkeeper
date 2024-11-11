package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/FischukSergey/gophkeeper/internal/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3 структура для S3.
type S3 struct {
	S3Session *session.Session
}

// S3UploadFile загружает файл в S3 bucket и возвращает URL загруженного файла.
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
			return "", fmt.Errorf("failed to create bucket: %w", err)
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
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Возвращаем URL загруженного файла
	return result.Location, nil
}

// S3GetFileList получает список файлов из S3.
func (s *S3) S3GetFileList(ctx context.Context, bucketID string, bucket string) ([]models.File, error) {
	svc := s3.New(s.S3Session)
	// Получаем список объектов в бакете
	result, err := svc.ListObjectsV2WithContext(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(bucketID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get file list: %w", err)
	}
	// формируем ответ
	files := make([]models.File, len(result.Contents))
	prefix := *result.Prefix
	for i, content := range result.Contents {
		files[i] = models.File{
			FileID:    *content.Key,
			UserID:    prefix,
			Filename:  filepath.Base(*content.Key),
			CreatedAt: *content.LastModified,
			DeletedAt: time.Time{},
			Size:      *content.Size,
		}
	}
	return files, nil
}

// S3DeleteFile удаляет файл из S3.
func (s *S3) S3DeleteFile(ctx context.Context, bucketID string, bucket string) error {
	svc := s3.New(s.S3Session)

	// Проверяем существование файла
	_, err := svc.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(bucketID),
	})
	if err != nil {
		return models.ErrFileNotExist
	}

	// Удаляем файл
	result, err := svc.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(bucketID),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	// Проверяем успешность удаления
	if result.DeleteMarker != nil && *result.DeleteMarker {
		return nil
	}

	// Дополнительная проверка, что файл действительно удален
	_, err = svc.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(bucketID),
	})
	if err == nil {
		// Если получаем ошибку, значит файл успешно удален
		return fmt.Errorf("file deletion could not be confirmed")
	}

	return nil
}

// S3DownloadFile скачивает файл из S3.
func (s *S3) S3DownloadFile(ctx context.Context, bucketID string, bucket string) ([]byte, error) {
	svc := s3.New(s.S3Session)
	// проверяем существование файла
	_, err := svc.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(bucketID),
	})
	if err != nil {
		return nil, models.ErrFileNotExist
	}
	// скачиваем файл
	result, err := svc.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(bucketID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	// преобразуем io.ReadCloser в []byte
	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return data, nil
}

// Close закрытие соединения с S3.
func (s *S3) Close() error {
	s.S3Session = nil
	return nil
}
