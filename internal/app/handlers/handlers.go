package handlers

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"

	"github.com/FischukSergey/gophkeeper/internal/app/converters"
	"github.com/FischukSergey/gophkeeper/internal/logger"
	"github.com/FischukSergey/gophkeeper/internal/models"
	pb "github.com/FischukSergey/protos/gen/gophkeeper/gophkeeper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	userFound = "userID found"
	request   = "request"
	user      = "userID"
)

var log = slog.New(slog.NewTextHandler(os.Stdout, nil))

// PwdKeeperServer структура для сервера.
type PwdKeeperServer struct {
	pb.UnimplementedGophKeeperServer
	PwdKeeper ProtoKeeperSaver
}

// ProtoKeeperSaver интерфейс для методов сервера.
type ProtoKeeperSaver interface {
	Ping(ctx context.Context) error
	RegisterUser(ctx context.Context, user models.User) (models.Token, error)
	Authorization(ctx context.Context, user models.User) (models.Token, error)
	FileUploadToS3(ctx context.Context, fileData io.Reader, filename string, userID int64) (string, error)
	FileGetListFromS3(ctx context.Context, userID int64) ([]models.File, error)
	FileDeleteFromS3(ctx context.Context, userID int64, filename string) error
	FileDownloadFromS3(ctx context.Context, userID int64, filename string) ([]byte, error)
}

// RegisterServerAPI регистрация сервера.
func RegisterServerAPI(
	server *grpc.Server,
	pwdKeeper ProtoKeeperSaver,
) {
	pb.RegisterGophKeeperServer(server, &PwdKeeperServer{PwdKeeper: pwdKeeper})
}

// Ping метод для проверки соединения с сервером.
func (s *PwdKeeperServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	err := s.PwdKeeper.Ping(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to ping: %v", err)
	}
	return &pb.PingResponse{}, nil
}

// Registration метод для регистрации пользователя.
func (s *PwdKeeperServer) Registration(
	ctx context.Context, req *pb.RegistrationRequest) (*pb.RegistrationResponse, error) {
	log.Info("Handler Registration method called")
	login := req.GetUsername()
	password := req.GetPassword()
	applicationName := req.GetAppName()
	role := req.GetRole()

	// проводим валидацию данных
	if login == "" || password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "username and password cannot be empty")
	}
	//
	user := models.User{
		Login:           login,
		Password:        password,
		ApplicationName: applicationName,
		Role:            role,
	}
	err := user.Validate()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "incorrect user data: %v", err)
	}

	// регистрируем пользователя
	token, err := s.PwdKeeper.RegisterUser(ctx, user)
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
func (s *PwdKeeperServer) Authorization(
	ctx context.Context, req *pb.AuthorizationRequest) (*pb.AuthorizationResponse, error) {
	log.Info("Handler Authorization method called")
	login := req.GetUsername()
	password := req.GetPassword()
	applicationName := req.GetAppName()

	//проводим валидацию данных
	if login == "" || password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "username and password cannot be empty")
	}
	user := models.User{
		Login:           login,
		Password:        password,
		ApplicationName: applicationName,
	}
	err := user.Validate()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "incorrect user data: %v", err)
	}

	// авторизуем пользователя
	token, err := s.PwdKeeper.Authorization(ctx, user)
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

// FileUpload метод для загрузки файла в S3.
func (s *PwdKeeperServer) FileUpload(
	stream pb.GophKeeper_FileUploadServer) error {
	log.Info("Handler FileUpload method called")
	ctx := stream.Context()
	userID, err := validateUserID(ctx)
	if err != nil {
		return err
	}
	log.Info(userFound, slog.Int(user, userID))

	// Получаем первый чанк с информацией о файле
	req, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Internal, "failed to receive file info: %v", err)
	}
	fileInfo := req.GetInfo()
	log.Info("fileInfo", slog.String("filename", fileInfo.Filename), slog.Int64("size", fileInfo.Size))
	if fileInfo.Filename == "" {
		return status.Errorf(codes.InvalidArgument, "filename cannot be empty")
	}
	if fileInfo.Size == 0 {
		return status.Errorf(codes.InvalidArgument, "file size cannot be 0")
	}

	// Создаем pipe для потоковой передачи данных
	pr, pw := io.Pipe()
	errCh := make(chan error, 1)

	// Запускаем горутину для записи данных в pipe
	go func() {
		defer func() {
			if err := pw.Close(); err != nil {
				log.Error("failed to close pipe writer", logger.Err(err))
			}
		}()
		for {
			req, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				errCh <- status.Errorf(codes.Internal, "failed to receive file chunk: %v", err)
				return
			}
			chunk := req.GetChunk()
			if chunk == nil {
				errCh <- status.Errorf(codes.InvalidArgument, "chunk cannot be nil")
				return
			}
			if _, err := pw.Write(chunk); err != nil {
				errCh <- err
				return
			}
		}
	}()

	// Загружаем файл в S3 потоково
	_, err = s.PwdKeeper.FileUploadToS3(ctx, io.Reader(pr), fileInfo.Filename, int64(userID))
	if err != nil {
		return status.Errorf(codes.Internal, "failed to upload file: %v", err)
	}

	// Проверяем ошибки из горутины
	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

// FileGetList метод для получения списка файлов пользователя.
func (s *PwdKeeperServer) FileGetList(
	ctx context.Context, req *pb.FileGetListRequest) (*pb.FileGetListResponse, error) {
	log.Info("Handler FileGetList method called")
	userID, err := validateUserID(ctx)
	if err != nil {
		return nil, err
	}
	log.Info(userFound, slog.Int(user, userID))
	files, err := s.PwdKeeper.FileGetListFromS3(ctx, int64(userID))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get file list: %v", err)
	}
	// формируем ответ
	filesPb := make([]*pb.File, len(files))
	for i, file := range files {
		filesPb[i] = converters.ToProtoFile(file)
	}
	return &pb.FileGetListResponse{Files: filesPb}, nil
}

// FileDelete метод для удаления файла из S3.
func (s *PwdKeeperServer) FileDelete(
	ctx context.Context, req *pb.FileDeleteRequest) (*pb.FileDeleteResponse, error) {
	log.Info("Handler FileDelete method called")
	userID, err := validateUserID(ctx)
	if err != nil {
		return nil, err
	}
	log.Info(userFound, slog.Int(user, userID))
	err = s.PwdKeeper.FileDeleteFromS3(ctx, int64(userID), req.Filename)
	if err != nil {
		if errors.Is(err, models.ErrFileNotExist) {
			return nil, status.Errorf(codes.NotFound, "file does not exist: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to delete file: %v", err)
	}
	return &pb.FileDeleteResponse{}, nil
}

// FileDownload метод для скачивания файла из S3.
func (s *PwdKeeperServer) FileDownload(
	req *pb.FileDownloadRequest, stream pb.GophKeeper_FileDownloadServer) error {
	log.Info("Handler FileDownload method called")
	ctx := stream.Context()
	userID, err := validateUserID(ctx)
	if err != nil {
		return err
	}
	log.Info(userFound, slog.Int(user, userID))
	// скачиваем файл из S3

	data, err := s.PwdKeeper.FileDownloadFromS3(ctx, int64(userID), req.Filename)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to download file: %v", err)
	}
	const chunkSize = 1024 * 1024 // 1MB
	// отправляем файл по частям
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunk := data[i:end]
		if err := stream.Send(&pb.FileDownloadResponse{
			Chunk: chunk,
		}); err != nil {
			return status.Errorf(codes.Internal, "failed to send file chunk: %v", err)
		}
	}
	return nil
}
