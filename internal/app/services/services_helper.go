package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/FischukSergey/gophkeeper/internal/models"
)

// ValidateMetadata проверяет валидность метаданных
func ValidateMetadata(metadata []models.Metadata) error {
	metadMap := make(map[string]string)
	for _, m := range metadata {
		if m.Key == "" || m.Value == "" {
			return fmt.Errorf("invalid metadata")
		}
		if len(m.Key) > 255 || len(m.Value) > 255 {
			return fmt.Errorf("metadata key or value is too long")
		}
		//валидируем ключ
		if strings.Contains(m.Key, " ") {
			return fmt.Errorf("metadata key contains spaces")
		}
		//проверяем ключ на уникальность
		if _, ok := metadMap[m.Key]; ok {
			return fmt.Errorf("metadata key already exists, key: %s must be unique", m.Key)
		}
		metadMap[m.Key] = m.Value
	}
	return nil
}

// SerializeMetadata сериализует метаданные в JSON строку
func SerializeMetadata(metadata []models.Metadata) (string, error) {
	if len(metadata) == 0 {
		return "", nil
	}
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return "", fmt.Errorf("failed to serialize metadata: %w", err)
	}
	return string(metadataJSON), nil
}

// DeserializeMetadata десериализует JSON строку в метаданные
func DeserializeMetadata(rawMetadata string) (map[string]string, error) {
	var metadata map[string]string
	err := json.Unmarshal([]byte(rawMetadata), &metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize metadata: %w", err)
	}
	return metadata, nil
}
