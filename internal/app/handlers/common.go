package handlers

import (
	"context"

	"github.com/FischukSergey/gophkeeper/internal/app/interceptors/auth"
	"github.com/FischukSergey/gophkeeper/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// validateUserID проверяет корректность ID пользователя из контекста.
func validateUserID(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(auth.CtxKeyUserGrpc).(int)
	if !ok {
		return 0, status.Errorf(codes.Unauthenticated, models.UserIDNotFound)
	}
	if userID <= 0 {
		return 0, status.Errorf(codes.InvalidArgument, "invalid user ID")
	}
	return userID, nil
}
