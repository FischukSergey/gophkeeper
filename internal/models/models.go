package models

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	MinLoginLength    = 5
	MaxLoginLength    = 100
	MinPasswordLength = 6
	MaxPasswordLength = 72
)
const (
	UserIDNotFound = "user ID not found in context"
)

// User структура для пользователя.
type User struct {
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
	Login          string
	Password       string
	HashedPassword string
	ID             int64
}

// Token структура для токена.
type Token struct {
	CreatedAt time.Time
	ExpiredAt time.Time
	Token     string
	UserID    int64
}

// File структура для файла.
type File struct {
	CreatedAt time.Time
	DeletedAt time.Time
	FileID    string
	UserID    string
	Filename  string
	Size      int64
}

// Card структура для карты.
type Card struct {
	CardExpirationDate time.Time
	UserID             string
	CardNumber         string
	CardHolder         string
	CardCVV            string
	CardBank           string
	Metadata           string
	CardID             int64
}

// Metadata структура для метаданных.
type Metadata struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Ошибки для пользователя.
var (
	ErrUserExists   = errors.New("user exists")
	ErrFileExists   = errors.New("file exists")
	ErrFileNotExist = errors.New("file does not exist")
)

// Note структура для заметки.
type Note struct {
	NoteText    string     `json:"note_text"`
	RawMetadata string     `json:"-"`
	Metadata    []Metadata `json:"metadata,omitempty"`
	NoteID      int64      `json:"note_id"`
	UserID      int64      `json:"user_id"`
}

// Validate валидация логина и пароля.
func (u *User) Validate() error {
	err := validation.ValidateStruct(
		u,
		validation.Field(
			&u.Login,
			validation.Required,
			validation.Length(MinLoginLength, MaxLoginLength),
		),
		validation.Field(
			&u.Password,
			validation.Required,
			validation.Length(MinPasswordLength, MaxPasswordLength),
			validation.Match(regexp.MustCompile(`^[a-zA-Z0-9]+$`)),
		),
	)
	if err != nil {
		return fmt.Errorf("invalid user data: %w", err)
	}
	return nil
}
