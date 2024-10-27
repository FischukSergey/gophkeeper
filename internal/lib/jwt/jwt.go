package jwt

import (
	"errors"
	"time"

	"github.com/FischukSergey/gophkeeper/cmd/server/initial"
	"github.com/FischukSergey/gophkeeper/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken генерирует токен
func GenerateToken(user models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	if user.ID > 0 && user.Login != "" {
		claims["uid"] = user.ID
		claims["login"] = user.Login
		claims["exp"] = time.Now().Add(time.Hour * initial.Cfg.JWT.ExpiresKey).Unix()
	} else {
		return "", errors.New("can't create JWT, invalid user id or login")
	}

	tokenString, err := token.SignedString([]byte(initial.Cfg.JWT.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, err
}
	