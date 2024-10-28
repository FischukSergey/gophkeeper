package modelsclient

import (
	"errors"
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	minLoginLength    = 6
	maxLoginLength    = 100
	minPasswordLength = 7
	maxPasswordLength = 72
)

// User структура для пользователя.
type User struct {
	Login    string
	Password string
}

// ErrInvalidLogin ошибка, если логин не валидный.
var ErrInvalidLogin = errors.New("invalid login")

// ValidateLogin валидация логина.
func ValidateLogin(login string) error {
	if login == "" {
		return fmt.Errorf("login cannot be empty")
	}
	u := User{Login: login}
	err := validation.ValidateStruct(
		&u,
		validation.Field(&u.Login, validation.Required, validation.Length(minLoginLength, maxLoginLength)),
	)
	if err != nil {
		return fmt.Errorf("invalid login: %w", err)
	}
	return nil
}

// ValidatePassword валидация пароля.
func ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	u := User{Password: password}
	err := validation.ValidateStruct(
		&u,
		validation.Field(&u.Password,
			validation.Required,
			validation.Length(minPasswordLength, maxPasswordLength),
			validation.Match(regexp.MustCompile(`^[a-zA-Z0-9]+$`)),
		),
	)
	if err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}
	return nil
}
