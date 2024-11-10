package handlers

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/FischukSergey/gophkeeper/internal/app/interceptors/auth"
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
	FileUploadToS3(ctx context.Context, fileData []byte, filename string, userID int64) (string, error)
	FileGetListFromS3(ctx context.Context, userID int64) ([]models.File, error)
	FileDeleteFromS3(ctx context.Context, userID int64, filename string) error
	FileDownloadFromS3(ctx context.Context, userID int64, filename string) ([]byte, error)
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
		return nil, status.Errorf(codes.InvalidArgument, "username and password cannot be empty")
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

// GetList метод для получения списка записей пользователя.
func (s *pwdKeeperServer) NoteGetList(
	ctx context.Context, req *pb.NoteGetListRequest) (*pb.NoteGetListResponse, error) {
	// log.Info("Handler NoteGetList method called")
	// userID := req.UserID
	// notes, err := s.pwdKeeper.GetList(ctx, userID)
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, "failed to get list: %v", err)
	// }
	return &pb.NoteGetListResponse{}, nil
}

// FileUpload метод для загрузки файла в S3.
func (s *pwdKeeperServer) FileUpload(
	ctx context.Context, req *pb.FileUploadRequest) (*pb.FileUploadResponse, error) {
	log.Info("Handler FileUpload method called")
	userID, ok := ctx.Value(auth.CtxKeyUserGrpc).(int)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user ID not found in context")
	}
	log.Info("userID found", slog.Int("userID", userID))
	// загружаем файл в S3
	url, err := s.pwdKeeper.FileUploadToS3(ctx, req.Data, req.Filename, int64(userID))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to upload file: %v", err)
	}
	return &pb.FileUploadResponse{Message: url}, nil
}

// FileGetList метод для получения списка файлов пользователя.
func (s *pwdKeeperServer) FileGetList(
	ctx context.Context, req *pb.FileGetListRequest) (*pb.FileGetListResponse, error) {
	log.Info("Handler FileGetList method called")
	userID, ok := ctx.Value(auth.CtxKeyUserGrpc).(int)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user ID not found in context")
	}
	log.Info("userID found", slog.Int("userID", userID))
	files, err := s.pwdKeeper.FileGetListFromS3(ctx, int64(userID))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get file list: %v", err)
	}
	// формируем ответ
	filesPb := make([]*pb.File, len(files))
	for i, file := range files {
		filesPb[i] = &pb.File{
			FileID:    file.FileID,
			UserID:    file.UserID,
			Filename:  file.Filename,
			Size:      file.Size,
			CreatedAt: timestamppb.New(file.CreatedAt),
		}
	}
	return &pb.FileGetListResponse{Files: filesPb}, nil
}

// FileDelete метод для удаления файла из S3.
func (s *pwdKeeperServer) FileDelete(
	ctx context.Context, req *pb.FileDeleteRequest) (*pb.FileDeleteResponse, error) {
	log.Info("Handler FileDelete method called")
	userID, ok := ctx.Value(auth.CtxKeyUserGrpc).(int)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user ID not found in context")
	}
	log.Info("userID found", slog.Int("userID", userID))
	err := s.pwdKeeper.FileDeleteFromS3(ctx, int64(userID), req.Filename)
	if err != nil {
		if errors.Is(err, models.ErrFileNotExist) {
			return nil, status.Errorf(codes.NotFound, "file does not exist: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to delete file: %v", err)
	}
	return &pb.FileDeleteResponse{}, nil
}

// FileDownload метод для скачивания файла из S3.
func (s *pwdKeeperServer) FileDownload(
	ctx context.Context, req *pb.FileDownloadRequest) (*pb.FileDownloadResponse, error) {
	log.Info("Handler FileDownload method called")
	userID, ok := ctx.Value(auth.CtxKeyUserGrpc).(int)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user ID not found in context")
	}
	log.Info("userID found", slog.Int("userID", userID))
	// скачиваем файл из S3
	data, err := s.pwdKeeper.FileDownloadFromS3(ctx, int64(userID), req.Filename)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to download file: %v", err)
	}
	return &pb.FileDownloadResponse{Data: data}, nil
}
