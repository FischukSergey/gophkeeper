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

type CtxKey int

const (
	CtxKeyUserGrpc CtxKey = iota + 1
)


// AuthInterceptor интерцептор для проверки токена.
func AuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	var userID int
	var err error
	switch info.FullMethod {
	case "/server.GophKeeper/Register": //исключаем регистрацию и авторизацию из проверки токена	
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
				slog.Info("token not found")
				return nil, status.Error(codes.Unauthenticated, "token not found")
			default:
				token := values[0]
				slog.Info("token found")
				userID, err = jwt.GetUserID(token)
				if err != nil {
					if errors.Is(err, jwt.ErrTokenExpired) {
						slog.Info("token expired")
						return nil, status.Error(codes.Unauthenticated, "token expired")
					}
					return nil, status.Error(codes.Unauthenticated, "invalid token")
				}
			}
		} else {
			slog.Info("token not found")
			return nil, status.Error(codes.Unauthenticated, "token not found")
		}
	}
	slog.Info("userID found", slog.Int("userID", userID))
	ctx = context.WithValue(ctx, CtxKeyUserGrpc, userID)
	return handler(ctx, req)
}
