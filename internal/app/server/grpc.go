package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/FischukSergey/gophkeeper/cmd/server/initial"
	"github.com/FischukSergey/gophkeeper/internal/app/handlers"
	"github.com/FischukSergey/gophkeeper/internal/app/interceptors/auth"
	"github.com/FischukSergey/gophkeeper/internal/app/services"
	"github.com/FischukSergey/gophkeeper/internal/logger"
	"github.com/FischukSergey/gophkeeper/internal/storage/dbstorage"
	"github.com/FischukSergey/gophkeeper/internal/storage/s3"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GrpcServer структура для хранения gRPC сервера.
type GrpcServer struct {
	grpcServer *grpc.Server
	log        *slog.Logger
	port       string
}

// App для экземпляра	gRPC сервера.
type App struct {
	GrpcServer *GrpcServer
	Storage    *dbstorage.Storage //бд postgres
	S3         *s3.S3             //S3 bucket
}

// NewGrpcServer функция для инициализации gRPC сервера.
func NewGrpcServer(log *slog.Logger, port string) *App {
	// инициализация хранилища бд postgres
	storage, err := initial.InitStorage()
	if err != nil {
		panic("Error initializing storage: " + err.Error())
	}
	log.Info("Database connected")
	err = storage.GetPingDB(context.Background())
	if err != nil {
		panic("Error pinging database: " + err.Error())
	}

	// проверка на имплементацию интерфейса и методов хранилища на этапе компиляции
	var _ services.PwdKeeper = (*dbstorage.Storage)(nil)

	//инициализация S3
	s3Storage, err := initial.InitS3()
	if err != nil {
		panic("Error initializing s3: " + err.Error())
	}
	log.Info("S3 connected")	

	// создание сервиса
	grpcService := services.NewGRPCService(log, storage, s3Storage)

	grpcApp := &GrpcServer{
		log:  log,
		port: port,
	}

	// опции для логирования в middleware
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent),
	}
	// опции для обработки паники в grpc сервере
	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) (err error) {
			log.Error("gRPC server panic", logger.Err(err))
			return status.Errorf(codes.Internal, "panic: %v", p)
		}),
	}

	// создание gRPC сервера
	grpcApp.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
			recovery.UnaryServerInterceptor(recoveryOpts...),
			auth.AuthInterceptor,
		),
	)

	// регистрация сервиса в gRPC сервере
	handlers.RegisterServerAPI(grpcApp.grpcServer, grpcService)

	return &App{
		GrpcServer: grpcApp,
		Storage:    storage,
		S3:         s3Storage,
	}
}

// MustRun запуск grpc сервера.
func (app *App) MustRun() {
	go func() {
		if err := app.GrpcServer.Run(); err != nil {
			app.GrpcServer.log.Error("Error starting gRPC server", logger.Err(err))
			panic(err)
		}
	}()
	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	app.GrpcServer.log.Info("Stopping gRPC server")
	app.GrpcServer.grpcServer.GracefulStop()
	app.GrpcServer.log.Info("gRPC server stopped")
	// закрытие соединения с бд
	app.Storage.Close()
	app.GrpcServer.log.Info("Database closed")
	// закрытие соединения с S3
	if err := app.S3.Close(); err != nil {
		app.GrpcServer.log.Error("Error closing S3", logger.Err(err))
	}
	app.GrpcServer.log.Info("S3 connection closed")
}

// Run запуск grpc сервера.
func (app *GrpcServer) Run() error {
	lis, err := net.Listen("tcp", app.port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", app.port, err)
	}
	app.log.Info("Starting gRPC server on port", slog.String("port", app.port))

	// запускаем обработчик gRPC сообщений
	if err := app.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}

// InterceptorLogger обертка интерцептора для логирования.
// Меняем logging.LevelInfo на slog.LevelInfo.
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.LevelInfo, msg, fields...)
	})
}
