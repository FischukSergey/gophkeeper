package handlers

import (
	"context"
	"testing"
	"time"

	"github.com/FischukSergey/gophkeeper/internal/app/handlers/mock"
	"github.com/FischukSergey/gophkeeper/internal/app/interceptors/auth"
	"github.com/FischukSergey/gophkeeper/internal/models"
	pb "github.com/FischukSergey/gophkeeper/internal/proto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestProtoCardService_AddCard(t *testing.T) {
	// Создаем контроллер и мок сервиса
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockProtoCardService(ctrl)

	// Тестовые данные
	testCard := &pb.Card{
		CardNumber:         "5272697132101976",
		CardHolder:         "Test Card",
		CardCVV:            "123",
		CardBank:           "Test Bank",
		CardExpirationDate: &timestamppb.Timestamp{Seconds: 1716211200},
	}
	// Структура для хранения результата теста
	type result struct {
		err    error
		status codes.Code
	}
	// Масив тестов
	tests := []struct {
		name      string
		args      *pb.Card
		setupMock func(*mock.MockProtoCardService)
		want      result
		userID    int
	}{
		{
			name: "successful card add",
			args: testCard,
			setupMock: func(mock *mock.MockProtoCardService) {
				mock.EXPECT().
					CardAddService(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			want: result{
				err:    nil,
				status: codes.OK,
			},
			userID: 18,
		},
		{
			name: "unauthorized request",
			setupMock: func(mock *mock.MockProtoCardService) {
				// Для этого кейса мок не нужен, так как ошибка произойдет раньше
			},
			want: result{
				err:    status.Error(codes.Unauthenticated, "unauthorized"),
				status: codes.Unauthenticated,
			},
			userID: 0,
		},
	}

	// Выполняем тесты
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}
			s := &CardServer{CardService: mockService}

			ctx := context.Background()
			if tt.userID != 0 {
				ctx = context.WithValue(ctx, auth.CtxKeyUserGrpc, tt.userID)
			}

			_, err := s.CardAdd(ctx, &pb.CardAddRequest{
				Card: tt.args,
			})

			if tt.want.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.want.status, status.Code(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProtoCardService_CardGetList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockProtoCardService(ctrl)

	testCards := []models.Card{
		{
			CardNumber:         "5272697132101976",
			CardHolder:         "Test Card",
			CardCVV:            "123",
			CardBank:           "Test Bank",
			CardExpirationDate: time.Unix(1716211200, 0),
		},
	}

	type result struct {
		err    error
		status codes.Code
	}

	tests := []struct {
		name      string
		setupMock func(*mock.MockProtoCardService)
		want      result
		userID    int
	}{
		{
			name: "successful get cards list",
			setupMock: func(mock *mock.MockProtoCardService) {
				mock.EXPECT().
					CardGetListService(gomock.Any(), gomock.Any()).
					Return(testCards, nil)
			},
			want: result{
				err:    nil,
				status: codes.OK,
			},
			userID: 18,
		},
		{
			name: "unauthorized request",
			setupMock: func(mock *mock.MockProtoCardService) {
				// Для этого кейса мок не нужен, так как ошибка произойдет раньше
			},
			want: result{
				err:    status.Error(codes.Unauthenticated, "unauthorized"),
				status: codes.Unauthenticated,
			},
			userID: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}
			s := &CardServer{CardService: mockService}

			ctx := context.Background()
			if tt.userID != 0 {
				ctx = context.WithValue(ctx, auth.CtxKeyUserGrpc, tt.userID)
			}

			result, err := s.CardGetList(ctx, &pb.CardGetListRequest{})

			if tt.want.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.want.status, status.Code(err))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestProtoCardService_CardAddMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockProtoCardService(ctrl)

	type result struct {
		err    error
		status codes.Code
	}

	tests := []struct {
		name      string
		request   *pb.CardAddMetadataRequest
		setupMock func(*mock.MockProtoCardService)
		want      result
		userID    int
	}{
		{
			name: "successful metadata add",
			request: &pb.CardAddMetadataRequest{
				CardID: 1,
				Metadata: []*pb.Metadata{
					{Key: "key1", Value: "value1"},
				},
			},
			setupMock: func(mock *mock.MockProtoCardService) {
				mock.EXPECT().
					CardAddMetadataService(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					AnyTimes()
			},
			want: result{
				err:    nil,
				status: codes.OK,
			},
			userID: 18,
		},
		{
			name: "unauthorized request",
			request: &pb.CardAddMetadataRequest{
				CardID: 1,
				Metadata: []*pb.Metadata{
					{Key: "key1", Value: "value1"},
				},
			},
			setupMock: func(mock *mock.MockProtoCardService) {
				// Для этого кейса мок не нужен, так как ошибка произойдет раньше
			},
			want: result{
				err:    status.Error(codes.Unauthenticated, "unauthorized"),
				status: codes.Unauthenticated,
			},
			userID: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}
			s := &CardServer{CardService: mockService}

			ctx := context.Background()
			if tt.userID != 0 {
				ctx = context.WithValue(ctx, auth.CtxKeyUserGrpc, tt.userID)
			}

			_, err := s.CardAddMetadata(ctx, tt.request)

			if tt.want.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.want.status, status.Code(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProtoCardService_DeleteCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockProtoCardService(ctrl)

	type result struct {
		err    error
		status codes.Code
	}

	tests := []struct {
		name      string
		request   *pb.CardDeleteRequest
		setupMock func(*mock.MockProtoCardService)
		want      result
		userID    int
	}{
		{
			name: "successful card delete",
			request: &pb.CardDeleteRequest{
				CardID: 1,
			},
			setupMock: func(mock *mock.MockProtoCardService) {
				mock.EXPECT().
					CardDeleteService(gomock.Any(), gomock.Any()).
					Return(nil).
					AnyTimes()
			},
			want: result{
				err:    nil,
				status: codes.OK,
			},
			userID: 18,
		},
		{
			name: "unauthorized request",
			request: &pb.CardDeleteRequest{
				CardID: 1,
			},
			setupMock: func(mock *mock.MockProtoCardService) {
				// Для этого кейса мок не нужен, так как ошибка произойдет раньше
			},
			want: result{
				err:    status.Error(codes.Unauthenticated, "unauthorized"),
				status: codes.Unauthenticated,
			},
			userID: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}
			s := &CardServer{CardService: mockService}

			ctx := context.Background()
			if tt.userID != 0 {
				ctx = context.WithValue(ctx, auth.CtxKeyUserGrpc, tt.userID)
			}

			_, err := s.CardDelete(ctx, tt.request)

			if tt.want.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.want.status, status.Code(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
