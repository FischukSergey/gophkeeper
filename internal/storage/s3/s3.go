package s3

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)
	
// UploadFile загружает файл в S3 bucket и возвращает URL загруженного файла
func UploadFile(uploader *s3manager.Uploader, fileData []byte, filename string, bucket string) (string, error) {
	// Подготавливаем параметры загрузки
	uploadInput := &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(fileData),
		ACL:    aws.String("public-read"),
	}

	// Загружаем файл
	result, err := uploader.Upload(uploadInput)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	// Возвращаем URL загруженного файла
	return result.Location, nil
}

