package handlers

import (
	"context"
	"log/slog"
	"os"

	"github.com/FischukSergey/gophkeeper/internal/models"
	pb "github.com/FischukSergey/gophkeeper/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var log = slog.New(slog.NewTextHandler(os.Stdout, nil))

// pwdKeeperServer структура для сервера.
type pwdKeeperServer struct {
	pb.UnimplementedGophKeeperServer
	pwdKeeper ProtoKeeperSaver
}

// ProtoKeeperSaver интерфейс для методов сервера.
type ProtoKeeperSaver interface {
	Ping(ctx context.Context) error
	RegisterUser(ctx context.Context, login, password string) (models.Token, error)
	Authorization(ctx context.Context, login, password string) (models.Token, error)
}

// RegisterServerAPI регистрация сервера.
func RegisterServerAPI(server *grpc.Server, pwdKeeper ProtoKeeperSaver) {
	pb.RegisterGophKeeperServer(server, &pwdKeeperServer{pwdKeeper: pwdKeeper})
}

// Ping метод для проверки соединения с сервером.
func (s *pwdKeeperServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	err := s.pwdKeeper.Ping(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to ping: %v", err)
	}
	return &pb.PingResponse{}, nil
}

// Registration метод для регистрации пользователя.
func (s *pwdKeeperServer) Registration(
	ctx context.Context, req *pb.RegistrationRequest) (*pb.RegistrationResponse, error) {
	log.Info("Handler Registration method called")
	login := req.Username
	password := req.Password

	// проводим валидацию данных
	if login == "" || password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "username and password cannot be empty")
	}
	//
	user := models.User{
		Login:    login,
		Password: password,
	}
	err := user.Validate()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "incorrect user data: %v", err)
	}

	// регистрируем пользователя
	token, err := s.pwdKeeper.RegisterUser(ctx, req.Username, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}
	// формируем ответ
	accessToken := &pb.Token{
		UserID:    token.UserID,
		Token:     token.Token,
		CreatedAt: timestamppb.New(token.CreatedAt),
		ExpiredAt: timestamppb.New(token.ExpiredAt),
	}
	return &pb.RegistrationResponse{AccessToken: accessToken}, nil
}

// Authorization метод для авторизации пользователя.
func (s *pwdKeeperServer) Authorization(
	ctx context.Context, req *pb.AuthorizationRequest) (*pb.AuthorizationResponse, error) {
	log.Info("Handler Authorization method called")
	login := req.Username
	password := req.Password

	//проводим валидацию данных
	if login == "" || password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "username or password cannot be empty")
	}
	user := models.User{
		Login:    login,
		Password: password,
	}
	err := user.Validate()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "incorrect user data: %v", err)
	}

	// авторизуем пользователя
	token, err := s.pwdKeeper.Authorization(ctx, login, password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to authorize user: %v", err)
	}

	// формируем ответ
	accessToken := &pb.Token{
		UserID:    token.UserID,
		Token:     token.Token,
		CreatedAt: timestamppb.New(token.CreatedAt),
		ExpiredAt: timestamppb.New(token.ExpiredAt),
	}
	return &pb.AuthorizationResponse{AccessToken: accessToken}, nil
}
