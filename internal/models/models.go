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
	FileID    string
	UserID    string
	Filename  string
	CreatedAt time.Time
	DeletedAt time.Time
	Size      int64
}

// Card структура для карты.
type Card struct {
	CardID             string
	UserID             string
	CardNumber         string
	CardHolder         string
	CardExpirationDate time.Time
	CardCVV            string
	CardBank           string
}

// Ошибки для пользователя.
var (
	ErrUserExists   = errors.New("user exists")
	ErrFileExists   = errors.New("file exists")
	ErrFileNotExist = errors.New("file does not exist")
)

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
