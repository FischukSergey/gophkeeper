package handlers

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/FischukSergey/gophkeeper/cmd/server/initial"
	mock_handlers "github.com/FischukSergey/gophkeeper/internal/app/handlers/mock"
	"github.com/FischukSergey/gophkeeper/internal/app/interceptors/auth"
	"github.com/FischukSergey/gophkeeper/internal/config"
	"github.com/FischukSergey/gophkeeper/internal/models"
	pb "github.com/FischukSergey/gophkeeper/internal/proto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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
func Test_pwdKeeperServer_Registration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_handlers.NewMockProtoKeeperSaver(ctrl)

	type args struct {
		req *pb.RegistrationRequest
	}
	type result struct {
		got    *pb.RegistrationResponse
		err    error
		status codes.Code
	}
	tests := []struct {
		name string
		args args
		want result
		mock func()
	}{
		{
			name: "successful registration",
			args: args{
				req: &pb.RegistrationRequest{
					Username: "testuser",
					Password: "testpass123",
				},
			},
			want: result{
				got: &pb.RegistrationResponse{
					AccessToken: &pb.Token{
						CreatedAt: timestamppb.New(time.Unix(0, 0)),
						ExpiredAt: timestamppb.New(time.Unix(0, 0)),
					},
				},
				err:    nil,
				status: codes.OK,
			},
			mock: func() {
				mockService.EXPECT().
					RegisterUser(gomock.Any(), "testuser", "testpass123").
					Return(models.Token{CreatedAt: time.Unix(0, 0), ExpiredAt: time.Unix(0, 0)}, nil)
			},
		},
		{
			name: "login is empty",
			args: args{
				req: &pb.RegistrationRequest{
					Username: "",
					Password: "test",
				},
			},
			want: result{
				got:    nil,
				err:    status.Errorf(codes.InvalidArgument, "username and password cannot be empty"),
				status: codes.InvalidArgument,
			},
			mock: nil,
		},
		{
			name: "password is empty",
			args: args{
				req: &pb.RegistrationRequest{
					Username: "test",
					Password: "",
				},
			},
			want: result{
				got:    nil,
				err:    status.Errorf(codes.InvalidArgument, "username and password cannot be empty"),
				status: codes.InvalidArgument,
			},
			mock: nil,
		},
		{
			name: "incorrect login",
			args: args{
				req: &pb.RegistrationRequest{
					Username: "1234",
					Password: "test123",
				},
			},
			want: result{
				got: nil,
				err: status.Errorf(
					codes.InvalidArgument,
					"incorrect user data: %v",
					"invalid user data: Login: the length must be between 5 and 100.",
				),
				status: codes.InvalidArgument,
			},
			mock: nil,
		},
		{
			name: "incorrect password",
			args: args{
				req: &pb.RegistrationRequest{
					Username: "testuser",
					Password: "1234",
				},
			},
			want: result{
				got: nil,
				err: status.Errorf(
					codes.InvalidArgument,
					"incorrect user data: %v",
					"invalid user data: Password: the length must be between 6 and 72.",
				),
				status: codes.InvalidArgument,
			},
			mock: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}
			s := &PwdKeeperServer{
				PwdKeeper: mockService,
			}
			got, err := s.Registration(context.Background(), tt.args.req)
			if !assert.ErrorIs(t, err, tt.want.err) {
				t.Errorf("pwdKeeperServer.Registration() error = %v, wantErr %v", err, tt.want.err)
				return
			}
			if !assert.Equal(t, got, tt.want.got) {
				t.Errorf("pwdKeeperServer.Registration() = %v, want %v", got, tt.want.got)
			}
			if status.Code(err) != tt.want.status {
				t.Errorf("pwdKeeperServer.Registration() status = %v, wantStatus %v", status.Code(err), tt.want.status)
			}
		})
	}
}

func Test_PwdKeeperServer_Authorization(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_handlers.NewMockProtoKeeperSaver(ctrl)

	type args struct {
		req *pb.AuthorizationRequest
	}
	type result struct {
		got    *pb.AuthorizationResponse
		err    error
		status codes.Code
	}
	tests := []struct {
		name string
		args args
		want result
		mock func()
	}{
		{
			name: "successful authorization",
			args: args{
				req: &pb.AuthorizationRequest{
					Username: "testuser",
					Password: "testpass123",
				},
			},
			want: result{
				got: &pb.AuthorizationResponse{
					AccessToken: &pb.Token{
						CreatedAt: timestamppb.New(time.Unix(0, 0)),
						ExpiredAt: timestamppb.New(time.Unix(0, 0)),
					},
				},
				err:    nil,
				status: codes.OK,
			},
			mock: func() {
				mockService.EXPECT().
					Authorization(gomock.Any(), "testuser", "testpass123").
					Return(models.Token{CreatedAt: time.Unix(0, 0), ExpiredAt: time.Unix(0, 0)}, nil)
			},
		},
		{
			name: "login is empty",
			args: args{
				req: &pb.AuthorizationRequest{
					Username: "",
					Password: "test",
				},
			},
			want: result{
				got:    nil,
				err:    status.Errorf(codes.InvalidArgument, "username and password cannot be empty"),
				status: codes.InvalidArgument,
			},
			mock: nil,
		},
		{
			name: "password is empty",
			args: args{
				req: &pb.AuthorizationRequest{
					Username: "test",
					Password: "",
				},
			},
			want: result{
				got:    nil,
				err:    status.Errorf(codes.InvalidArgument, "username and password cannot be empty"),
				status: codes.InvalidArgument,
			},
			mock: nil,
		},
		{
			name: "incorrect login",
			args: args{
				req: &pb.AuthorizationRequest{
					Username: "test",
					Password: "test123",
				},
			},
			want: result{
				got: nil,
				err: status.Errorf(
					codes.InvalidArgument,
					"incorrect user data: %v",
					"invalid user data: Login: the length must be between 5 and 100.",
				),
				status: codes.InvalidArgument,
			},
			mock: nil,
		},
		{
			name: "incorrect password",
			args: args{
				req: &pb.AuthorizationRequest{
					Username: "test123",
					Password: "test",
				},
			},
			want: result{
				got: nil,
				err: status.Errorf(
					codes.InvalidArgument,
					"incorrect user data: %v",
					"invalid user data: Password: the length must be between 6 and 72.",
				),
				status: codes.InvalidArgument,
			},
			mock: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}
			s := &PwdKeeperServer{
				PwdKeeper: mockService,
			}
			got, err := s.Authorization(context.Background(), tt.args.req)
			if !assert.ErrorIs(t, err, tt.want.err) {
				t.Errorf("pwdKeeperServer.Authorization() error = %v, wantErr %v", err, tt.want.err)
				return
			}
			if !assert.Equal(t, got, tt.want.got) {
				t.Errorf("pwdKeeperServer.Authorization() = %v, want %v", got, tt.want.got)
			}
			if status.Code(err) != tt.want.status {
				t.Errorf("pwdKeeperServer.Authorization() error = %v, wantErr %v", err, tt.want.status)
			}
		})
	}
}

func Test_PwdKeeperServer_Ping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_handlers.NewMockProtoKeeperSaver(ctrl)

	type args struct {
		req *pb.PingRequest
	}
	type result struct {
		got    *pb.PingResponse
		err    error
		status codes.Code
	}
	tests := []struct {
		name string
		args args
		want result
		mock func()
	}{
		{
			name: "successful ping",
			args: args{
				req: &pb.PingRequest{},
			},
			want: result{
				got:    &pb.PingResponse{},
				err:    nil,
				status: codes.OK,
			},
			mock: func() {
				mockService.EXPECT().
					Ping(gomock.Any()).
					Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}
			s := &PwdKeeperServer{
				PwdKeeper: mockService,
			}
			got, err := s.Ping(context.Background(), tt.args.req)
			if !assert.ErrorIs(t, err, tt.want.err) {
				t.Errorf("pwdKeeperServer.Ping() error = %v, wantErr %v", err, tt.want.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want.got) {
				t.Errorf("pwdKeeperServer.Ping() = %v, want %v", got, tt.want.got)
			}
			if status.Code(err) != tt.want.status {
				t.Errorf("pwdKeeperServer.Ping() error = %v, wantStatus %v", status.Code(err), tt.want.status)
			}
		})
	}
}

	func Test_PwdKeeperServer_FileUploadToS3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_handlers.NewMockProtoKeeperSaver(ctrl)

	userID := int(18)

	type args struct {
		req *pb.FileUploadRequest
		id  int
	}
	type result struct {
		got    *pb.FileUploadResponse
		err    error
		status codes.Code
	}
	tests := []struct {
		name string
		args args
		want result
		mock func()
	}{
		{
			name: "successful file upload",
			args: args{
				req: &pb.FileUploadRequest{
					Filename: "test.txt",
					Data:     []byte("test data"),
				},
				id: userID,
			},
			want: result{
				got:    &pb.FileUploadResponse{Message: "File uploaded successfully"},
				err:    nil,
				status: codes.OK,
			},
			mock: func() {
				mockService.EXPECT().
					FileUploadToS3(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("File uploaded successfully", nil)
			},
		},
		{
			name: "user id is not set",
			args: args{
				req: &pb.FileUploadRequest{
					Filename: "test.txt",
					Data:     []byte("test data"),
				},
				id: 0,
			},
			want: result{
				got:    nil,
				err:    status.Errorf(codes.Unauthenticated, "user ID not found in context"),
				status: codes.Unauthenticated,
			},
			mock: nil,
		},
		{
			name: "file upload error",
			args: args{
				req: &pb.FileUploadRequest{
					Filename: "",
					Data:     []byte("test data"),
				},
				id: userID,
			},
			want: result{
				got:    nil,
				err:    status.Errorf(codes.Internal, "failed to upload file: %v", "file upload error"),
				status: codes.Internal,
			},
			mock: func() {
				mockService.EXPECT().
					FileUploadToS3(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", errors.New("file upload error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}
			s := &PwdKeeperServer{
				PwdKeeper: mockService,
			}
			var got *pb.FileUploadResponse
			var err error
			if tt.args.id != 0 {
				ctx := context.WithValue(context.Background(), auth.CtxKeyUserGrpc, tt.args.id)
				got, err = s.FileUpload(ctx, tt.args.req)
			} else {
				got, err = s.FileUpload(context.Background(), tt.args.req)
			}
			if !assert.ErrorIs(t, err, tt.want.err) {
				t.Errorf("pwdKeeperServer.FileUpload() error = %v, wantErr %v", err, tt.want.err)
				return
			}
			if !assert.Equal(t, got, tt.want.got) {
				t.Errorf("pwdKeeperServer.FileUpload() = %v, want %v", got, tt.want.got)
			}
			if status.Code(err) != tt.want.status {
				t.Errorf("pwdKeeperServer.FileUpload() error = %v, wantStatus %v", status.Code(err), tt.want.status)
			}
		})
	}
}

func Test_PwdKeeperServer_FileDeleteFromS3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_handlers.NewMockProtoKeeperSaver(ctrl)
	userID := int(18)
	type args struct {
		req *pb.FileDeleteRequest
		id  int
	}
	type result struct {
		got    *pb.FileDeleteResponse
		err    error
		status codes.Code
	}
	tests := []struct {
		name string
		args args
		want result
		mock func()
	}{
		{
			name: "successful file delete",
			args: args{
				req: &pb.FileDeleteRequest{
					Filename: "test.txt",
				},
				id: userID,
			},
			want: result{
				got:    &pb.FileDeleteResponse{},
				err:    nil,
				status: codes.OK,
			},
			mock: func() {
				mockService.EXPECT().
					FileDeleteFromS3(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
		},
		{
			name: "file not found",
			args: args{
				req: &pb.FileDeleteRequest{
					Filename: "test.txt",
				},
				id: userID,
			},
			want: result{
				got:    nil,
				err:    status.Errorf(codes.NotFound, "file does not exist: %v", models.ErrFileNotExist),
				status: codes.NotFound,
			},
			mock: func() {
				mockService.EXPECT().
					FileDeleteFromS3(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(models.ErrFileNotExist)
			},
		},
		{
			name: "user id is not set",
			args: args{
				req: &pb.FileDeleteRequest{
					Filename: "test.txt",
				},
				id: 0,
			},
			want: result{
				got:    nil,
				err:    status.Errorf(codes.Unauthenticated, "user ID not found in context"),
				status: codes.Unauthenticated,
			},
			mock: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}
			s := &PwdKeeperServer{
				PwdKeeper: mockService,
			}
			var got *pb.FileDeleteResponse
			var err error
			if tt.args.id != 0 {
				ctx := context.WithValue(context.Background(), auth.CtxKeyUserGrpc, tt.args.id)
				got, err = s.FileDelete(ctx, tt.args.req)
			} else {
				got, err = s.FileDelete(context.Background(), tt.args.req)
			}
			if !assert.ErrorIs(t, err, tt.want.err) {
				t.Errorf("pwdKeeperServer.FileDelete() error = %v, wantErr %v", err, tt.want.err)
				return
			}
			if !assert.Equal(t, got, tt.want.got) {
				t.Errorf("pwdKeeperServer.FileDelete() = %v, want %v", got, tt.want.got)
			}
			if status.Code(err) != tt.want.status {
				t.Errorf("pwdKeeperServer.FileDelete() error = %v, wantStatus %v", status.Code(err), tt.want.status)
			}
		})
	}
}

func Test_PwdKeeperServer_FileGetListS3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_handlers.NewMockProtoKeeperSaver(ctrl)
	userID := int(18)
	type args struct {
		req *pb.FileGetListRequest
		id  int
	}
	type result struct {
		got    *pb.FileGetListResponse
		err    error
		status codes.Code
	}
	tests := []struct {
		name string
		args args
		want result
		mock func()
	}{
		{
			name: "successful file get list",
			args: args{
				req: &pb.FileGetListRequest{},
				id:  userID,
			},
			want: result{
				got: &pb.FileGetListResponse{
					Files: []*pb.File{
						{
							Filename:  "test.txt",
							FileID:    "1",
							UserID:    "18",
							Size:      100,
							CreatedAt: timestamppb.New(time.Unix(0, 0)),
						},
						{
							Filename:  "test2.txt",
							FileID:    "2",
							UserID:    "18",
							Size:      200,
							CreatedAt: timestamppb.New(time.Unix(0, 0)),
						},
					},
				},
				err:    nil,
				status: codes.OK,
			},
			mock: func() {
				mockService.EXPECT().
					FileGetListFromS3(gomock.Any(), gomock.Any()).
					Return([]models.File{
						{
							Filename:  "test.txt",
							FileID:    "1",
							UserID:    "18",
							Size:      100,
							CreatedAt: time.Unix(0, 0),
						},
						{
							Filename:  "test2.txt",
							FileID:    "2",
							UserID:    "18",
							Size:      200,
							CreatedAt: time.Unix(0, 0),
						},
					}, nil)
			},
		},
		{
			name: "user id is not set",
			args: args{
				req: &pb.FileGetListRequest{},
				id:  0,
			},
			want: result{
				got:    nil,
				err:    status.Errorf(codes.Unauthenticated, "user ID not found in context"),
				status: codes.Unauthenticated,
			},
			mock: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}
			s := &PwdKeeperServer{
				PwdKeeper: mockService,
			}
			var got *pb.FileGetListResponse
			var err error
			if tt.args.id != 0 {
				ctx := context.WithValue(context.Background(), auth.CtxKeyUserGrpc, tt.args.id)
				got, err = s.FileGetList(ctx, tt.args.req)
			} else {
				got, err = s.FileGetList(context.Background(), tt.args.req)
			}
			if !assert.ErrorIs(t, err, tt.want.err) {
				t.Errorf("pwdKeeperServer.FileGetList() error = %v, wantErr %v", err, tt.want.err)
				return
			}
			if !assert.Equal(t, got, tt.want.got) {
				t.Errorf("pwdKeeperServer.FileGetList() = %v, want %v", got, tt.want.got)
			}
			if status.Code(err) != tt.want.status {
				t.Errorf("pwdKeeperServer.FileGetList() error = %v, wantStatus %v", status.Code(err), tt.want.status)
			}
		})
	}
}
