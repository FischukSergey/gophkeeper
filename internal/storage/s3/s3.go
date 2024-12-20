package s3

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/FischukSergey/gophkeeper/internal/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3Storage структура для S3.
type S3Storage struct {
	S3Session  *session.Session
	BucketName string
}

// S3UploadFile загружает файл в S3 bucket и возвращает URL загруженного файла.
func (s *S3Storage) S3UploadFile(ctx context.Context, fileData io.Reader, filename string) (string, error) {
	// Подготавливаем параметры загрузки
	uploadInput := &s3manager.UploadInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(filename),
		Body:   fileData,
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
func (s *S3Storage) S3GetFileList(ctx context.Context, bucketID string) ([]models.File, error) {
	svc := s3.New(s.S3Session)

	result, err := svc.ListObjectsV2WithContext(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.BucketName),
		Prefix: aws.String(bucketID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get file list: %w", err)
	}

	// Фильтруем файлы, оставляя только те, которые начинаются с точного префикса
	var filteredFiles []models.File
	for _, content := range result.Contents {
		// Проверяем, что ключ начинается с нужного префикса
		if filepath.Dir(*content.Key) == bucketID {
			filteredFiles = append(filteredFiles, models.File{
				FileID:    *content.Key,
				UserID:    bucketID,
				Filename:  filepath.Base(*content.Key),
				CreatedAt: *content.LastModified,
				DeletedAt: time.Time{},
				Size:      *content.Size,
			})
		}
	}

	return filteredFiles, nil
}

// S3DeleteFile удаляет файл из S3.
func (s *S3Storage) S3DeleteFile(ctx context.Context, bucketID string) error {
	svc := s3.New(s.S3Session)
	// Проверяем существование файла
	err := s.S3CheckFileExist(ctx, svc, bucketID)
	if err != nil {
		return err
	}
	// Удаляем файл
	_, err = svc.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(bucketID),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// S3DownloadFile скачивает файл из S3.
func (s *S3Storage) S3DownloadFile(ctx context.Context, bucketID string) ([]byte, error) {
	svc := s3.New(s.S3Session)
	// проверяем существование файла
	err := s.S3CheckFileExist(ctx, svc, bucketID)
	if err != nil {
		return nil, err
	}
	// скачиваем файл
	result, err := svc.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
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
func (s *S3Storage) Close() error {
	s.S3Session = nil
	s.BucketName = ""
	return nil
}

// S3CheckFileExist проверяет существование файла.
func (s *S3Storage) S3CheckFileExist(ctx context.Context, svc *s3.S3, bucketID string) error {
	_, err := svc.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(bucketID),
	})
	if err != nil {
		return models.ErrFileNotExist
	}
	return nil
}
