package jwt

import (
	"errors"
	"testing"
	"time"

	"github.com/FischukSergey/gophkeeper/cmd/server/initial"
	"github.com/FischukSergey/gophkeeper/internal/config"
	"github.com/FischukSergey/gophkeeper/internal/models"
	"github.com/stretchr/testify/assert"
)

func init() {
	initial.Cfg = &config.Config{
		JWT: config.JWTConfig{
			SecretKey:  "test_secret_key",
			ExpiresKey: time.Hour * 12,
		},
	}
}

func TestGenerateToken(t *testing.T) {
	// Сохраняем оригинальную конфигурацию
	originalCfg := initial.Cfg
	// Восстанавливаем в конце теста
	defer func() {
		initial.Cfg = originalCfg
	}()

	type args struct {
		user models.User
	}
	tests := []struct {
		name    string
		args    args
		want    string
		jwtErr  error
		wantErr bool
	}{
		{
			name:    "valid user",
			args:    args{user: models.User{ID: 1, Login: "test"}},
			want:    "test",
			jwtErr:  nil,
			wantErr: false,
		},
		{
			name:    "invalid user id",
			args:    args{user: models.User{ID: 0, Login: "test"}},
			want:    "",
			jwtErr:  errors.New("can't create JWT, invalid user id or login"),
			wantErr: true,
		},
		{
			name:    "invalid user login",
			args:    args{user: models.User{ID: 1, Login: ""}},
			want:    "",
			jwtErr:  errors.New("can't create JWT, invalid user id or login"),
			wantErr: true,
		},
		{
			name:    "jwt error",
			args:    args{user: models.User{ID: 1, Login: "test"}},
			want:    "",
			jwtErr:  errors.New("jwt error"),
			wantErr: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateToken(tt.args.user)
			if err == nil {
				assert.NotEmpty(t, got)
				return
			}
			assert.Equal(t, tt.jwtErr, err)
		})
	}
}

func TestGetUserID(t *testing.T) {
	// Сохраняем оригинальную конфигурацию
	originalCfg := initial.Cfg
	// Восстанавливаем в конце теста
	defer func() {
		initial.Cfg = originalCfg
	}()

	type args struct {
		token string
	}
	token, _ := GenerateToken(models.User{ID: 18, Login: "test"})
	tests := []struct {
		name    string
		args    args
		want    int
		jwtErr  error
		wantErr bool
	}{
		{
			name:    "valid token",
			args:    args{token: token},
			want:    18,
			jwtErr:  nil,
			wantErr: false,
		},
		{
			name:    "invalid token",
			args:    args{token: "invalid_token"},
			want:    0,
			jwtErr:  errors.New("failed to parse token"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserID(tt.args.token)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.jwtErr.Error())
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
