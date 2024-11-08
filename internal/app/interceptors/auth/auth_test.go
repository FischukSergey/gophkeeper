package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/FischukSergey/gophkeeper/cmd/server/initial"
	"github.com/FischukSergey/gophkeeper/internal/config"
	"github.com/FischukSergey/gophkeeper/internal/lib/jwt"
	"github.com/FischukSergey/gophkeeper/internal/models"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestMain(m *testing.M) {
	// Устанавливаем тестовую конфигурацию JWT
	initial.Cfg = &config.Config{
		JWT: config.JWTConfig{
			SecretKey:  "test_secret_key",
			ExpiresKey: time.Hour * 12,
		},
	}
	m.Run()
}

func TestAuthInterceptor(t *testing.T) {
	//генерируем токен
	userID := int64(18)
	token, err := jwt.GenerateToken(models.User{ID: userID, Login: "test"})
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	//создаем аргументы для теста
	type args struct {
		token   string
		req     any
		info    *grpc.UnaryServerInfo
		handler grpc.UnaryHandler
	}

	mockHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		// Получаем ID из контекста
		if id := ctx.Value(CtxKeyUserGrpc); id != nil {
			return id, nil
		}
		return nil, errors.New(ErrNotFound)
	}

	tests := []struct {
		name    string
		args    args
		want    any
		err     error
		wantErr bool
	}{
		{
			name: "successful auth",
			args: args{
				token:   token,
				req:     "test request",
				info:    &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"},
				handler: mockHandler,
			},
			want:    int(userID),
			err:     nil,
			wantErr: false,
		},
		{
			name: "missing auth token",
			args: args{
				token:   "",
				req:     "test request",
				info:    &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"},
				handler: mockHandler,
			},
			want:    nil,
			err:     errors.New(ErrNotFound),
			wantErr: true,
		},
		{
			name: "invalid auth token",
			args: args{
				token:   "invalid_token",
				req:     "test request",
				info:    &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"},
				handler: mockHandler,
			},
			want:    nil,
			err:     errors.New(ErrInvalid),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//пишем в метаданные токен
			ctx := context.Background()
			ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("session_token", tt.args.token))
			got, err := AuthInterceptor(ctx, tt.args.req, tt.args.info, tt.args.handler)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthInterceptor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, got, tt.want) {
				t.Errorf("AuthInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}
