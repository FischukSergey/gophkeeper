package auth

import (
	"context"
	"errors"
	"log/slog"

	"github.com/FischukSergey/gophkeeper/internal/lib/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)
// CtxKey для передачи userID в контекст
type CtxKey string

// WrappedServerStream для передачи измененного контекста
type WrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

const (
	ErrNotFound           = "token not found"
	ErrExpired            = "token expired"
	ErrInvalid            = "invalid token"
	CtxKeyUserGrpc CtxKey = "userID"
)

// UnaryServerInterceptor для обычных unary вызовов
func UnaryAuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	var userID int
	var err error
	switch info.FullMethod {
	case "/server.GophKeeper/Registration": //исключаем регистрацию и авторизацию из проверки токена
		return handler(ctx, req)
	case "/server.GophKeeper/Authorization":
		return handler(ctx, req)
	case "/server.GophKeeper/Ping":
		return handler(ctx, req)
	default:
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			values := md.Get("session_token")
			switch len(values) {
			case 0:
				slog.Info(ErrNotFound)
				return nil, status.Errorf(codes.Unauthenticated, ErrNotFound)
			default:
				token := values[0]
				userID, err = jwt.GetUserID(token)
				if err != nil {
					if errors.Is(err, jwt.ErrTokenExpired) {
						slog.Info(ErrExpired)
						return nil, status.Errorf(codes.Unauthenticated, ErrExpired)
					}
					slog.Info(ErrInvalid)
					return nil, status.Errorf(codes.Unauthenticated, ErrInvalid)
				}
			}
		} else {
				slog.Info(ErrNotFound)
				return nil, status.Errorf(codes.Unauthenticated, ErrNotFound)
			}
		}
		ctxWithUserID := context.WithValue(ctx, CtxKeyUserGrpc, userID)
		return handler(ctxWithUserID, req)
	}

// StreamServerInterceptor для stream вызовов
func StreamAuthInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// пропускаем методы регистрации и авторизации
		if info.FullMethod == "/proto.GophKeeper/Registration" || info.FullMethod == "/proto.GophKeeper/Authorization" {
			return handler(srv, stream)
		}

		// получаем токен из метаданных
		md, ok := metadata.FromIncomingContext(stream.Context())
		if !ok {
			return status.Errorf(codes.Unauthenticated, ErrNotFound)
		}
		var userID int
		var err error
		values := md.Get("session_token")
		switch len(values) {
			case 0:
				slog.Info(ErrNotFound)
				return status.Errorf(codes.Unauthenticated, ErrNotFound)
			default:
				token := values[0]
				userID, err = jwt.GetUserID(token)
				if err != nil {
					if errors.Is(err, jwt.ErrTokenExpired) {
						slog.Info(ErrExpired)
						return status.Errorf(codes.Unauthenticated, ErrExpired)
					}
					slog.Info(ErrInvalid)
					return status.Errorf(codes.Unauthenticated, ErrInvalid)
				}
		}

		// создаем новый контекст с userID
		newCtx := context.WithValue(stream.Context(), CtxKeyUserGrpc, userID)
		wrappedStream := NewWrappedServerStream(newCtx, stream)

		return handler(srv, wrappedStream)
	}
}

// NewWrappedServerStream создает новый WrappedServerStream
func NewWrappedServerStream(ctx context.Context, stream grpc.ServerStream) *WrappedServerStream {
	return &WrappedServerStream{
		ServerStream: stream,
		ctx:          ctx,
	}
}

// добавляем методы для WrappedServerStream
func (w *WrappedServerStream) Context() context.Context {
	return w.ctx
}

func (w *WrappedServerStream) RecvMsg(m interface{}) error {
	return w.ServerStream.RecvMsg(m)
}

func (w *WrappedServerStream) SendMsg(m interface{}) error {
	return w.ServerStream.SendMsg(m)
}
