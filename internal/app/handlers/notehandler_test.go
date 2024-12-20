package handlers

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/FischukSergey/gophkeeper/internal/app/handlers/mock"
	"github.com/FischukSergey/gophkeeper/internal/app/interceptors/auth"
	"github.com/FischukSergey/gophkeeper/internal/models"
	pb "github.com/FischukSergey/gophkeeper/internal/proto"
)

func TestProtoNoteService_NoteAdd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockProtoNoteService(ctrl)

	type result struct {
		err    error
		status codes.Code
	}

	tests := []struct {
		name      string
		request   *pb.NoteAddRequest
		setupMock func(*mock.MockProtoNoteService)
		want      result
		userID    int
	}{
		{
			name: "successful note creation",
			request: &pb.NoteAddRequest{
				Note: &pb.Note{
					NoteText: "Test Note",
					Metadata: []*pb.Metadata{
						{Key: "key1", Value: "value1"},
					},
				},
			},
			setupMock: func(mock *mock.MockProtoNoteService) {
				mock.EXPECT().
					NoteAdd(gomock.Any(), gomock.Any()).
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
			request: &pb.NoteAddRequest{
				Note: &pb.Note{
					NoteText: "Test Note",
					Metadata: []*pb.Metadata{
						{Key: "key1", Value: "value1"},
					},
				},
			},
			setupMock: func(mock *mock.MockProtoNoteService) {
				// Мок не нужен, так как ошибка произойдет раньше
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
			s := &NoteServer{NoteService: mockService}

			ctx := context.Background()
			if tt.userID != 0 {
				ctx = context.WithValue(ctx, auth.CtxKeyUserGrpc, tt.userID)
			}

			_, err := s.NoteAdd(ctx, tt.request)

			if tt.want.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.want.status, status.Code(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProtoNoteService_NoteDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockProtoNoteService(ctrl)

	type result struct {
		err    error
		status codes.Code
	}

	tests := []struct {
		name      string
		request   *pb.NoteDeleteRequest
		setupMock func(*mock.MockProtoNoteService)
		want      result
		userID    int
	}{
		{
			name: "successful note deletion",
			request: &pb.NoteDeleteRequest{
				NoteID: 1,
			},
			setupMock: func(mock *mock.MockProtoNoteService) {
				mock.EXPECT().
					NoteDelete(gomock.Any(), gomock.Any(), gomock.Any()).
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
			request: &pb.NoteDeleteRequest{
				NoteID: 1,
			},
			setupMock: func(mock *mock.MockProtoNoteService) {
				// Мок не нужен, так как ошибка произойдет раньше
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
			s := &NoteServer{NoteService: mockService}

			ctx := context.Background()
			if tt.userID != 0 {
				ctx = context.WithValue(ctx, auth.CtxKeyUserGrpc, tt.userID)
			}

			_, err := s.NoteDelete(ctx, tt.request)

			if tt.want.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.want.status, status.Code(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProtoNoteService_NoteGetList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockProtoNoteService(ctrl)

	type result struct {
		err    error
		status codes.Code
	}
	notes := []models.Note{
		{NoteID: 1, NoteText: "Test Note", Metadata: []models.Metadata{{Key: "key1", Value: "value1"}}},
	}

	tests := []struct {
		name      string
		request   *pb.NoteGetListRequest
		setupMock func(*mock.MockProtoNoteService)
		want      result
		userID    int
	}{
		{
			name:    "successful note update",
			request: &pb.NoteGetListRequest{},
			setupMock: func(mock *mock.MockProtoNoteService) {
				mock.EXPECT().
					NoteGetList(gomock.Any(), gomock.Any()).
					Return(notes, nil)
			},
			want: result{
				err:    nil,
				status: codes.OK,
			},
			userID: 18,
		},
		{
			name:    "unauthorized request",
			request: &pb.NoteGetListRequest{},
			setupMock: func(mock *mock.MockProtoNoteService) {
				// Мок не нужен, так как ошибка произойдет раньше
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
			s := &NoteServer{NoteService: mockService}

			ctx := context.Background()
			if tt.userID != 0 {
				ctx = context.WithValue(ctx, auth.CtxKeyUserGrpc, tt.userID)
			}

			_, err := s.NoteGetList(ctx, tt.request)

			if tt.want.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.want.status, status.Code(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
