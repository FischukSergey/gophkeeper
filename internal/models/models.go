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
	ID             int64
	Login          string
	Password       string
	HashedPassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
}

// Token структура для токена.
type Token struct {
	UserID    int64
	Token     string
	CreatedAt time.Time
	ExpiredAt time.Time
}

// ErrUserExists ошибка, если пользователь уже существует.
var ErrUserExists = errors.New("user exists")

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