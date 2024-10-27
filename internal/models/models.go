package models

import (
	"errors"
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

// User структура для пользователя
type User struct {
	ID int64
	Login string
	Password string
	HashedPassword string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

// Token структура для токена
type Token struct {
	UserID int64
	Token string
	CreatedAt time.Time
	ExpiredAt time.Time
}

// ErrUserExists ошибка, если пользователь уже существует
var ErrUserExists = errors.New("user exists")

// Validate валидация логина и пароля
func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Login, validation.Required, validation.Length(6, 100)),
		validation.Field(&u.Password, validation.Required, validation.Length(6, 72), validation.Match(regexp.MustCompile(`^[a-zA-Z0-9]+$`))),
	)
}
