package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/FischukSergey/gophkeeper/cmd/server/initial"
	"github.com/FischukSergey/gophkeeper/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

var ErrTokenExpired = errors.New("token is expired")

// GenerateToken генерирует токен.
func GenerateToken(user models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("can't create JWT, invalid user id or login")
	}
	if user.ID > 0 && user.Login != "" {
		claims["uid"] = user.ID
		claims["login"] = user.Login
		claims["exp"] = time.Now().Add(initial.Cfg.JWT.ExpiresKey).Unix()
	} else {
		return "", fmt.Errorf("can't create JWT, invalid user id or login")
	}

	tokenString, err := token.SignedString([]byte(initial.Cfg.JWT.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return tokenString, nil
}

// GetUserID получает userID из токена.
func GetUserID(jwtToken string) (int, error) {
	jwtClaims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(jwtToken, jwtClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(initial.Cfg.JWT.SecretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, ErrTokenExpired
		}
		return 0, err
	}
	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}
	uid, ok := jwtClaims["uid"]
	if !ok {
		return 0, fmt.Errorf("user id not found")
	}
	return int(uid.(float64)), nil
}
