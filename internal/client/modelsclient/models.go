package modelsclient

import (
	"errors"
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
)

type User struct {
	Login string
	Password string
}

// ErrInvalidLogin ошибка, если логин не валидный
var ErrInvalidLogin = errors.New("invalid login")

// ValidateLogin валидация логина
func ValidateLogin(login string) error {
	if login == "" {
		return fmt.Errorf("login cannot be empty")
	}
	u := User{Login: login}
	err := validation.ValidateStruct(	
		&u,
		validation.Field(&u.Login, validation.Required, validation.Length(6, 100)),
	)
	if err != nil {
		return err
	}
	return nil
}

// ValidatePassword валидация пароля
func ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	u := User{Password: password}
	err := validation.ValidateStruct(	
		&u,
		validation.Field(&u.Password, validation.Required, validation.Length(6, 72), validation.Match(regexp.MustCompile(`^[a-zA-Z0-9]+$`))),
	)
	if err != nil {
		return err
	}	
	return nil
}
