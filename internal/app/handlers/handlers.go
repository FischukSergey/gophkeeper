package handlers

import (
	"context"

	pb "github.com/FischukSergey/gophkeeper/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// pwdKeeperServer структура для сервера
type pwdKeeperServer struct {
	pb.UnimplementedGophKeeperServer
	pwdKeeper ProtoKeeperSaver
}

// ProtoKeeperSaver интерфейс для методов сервера	
type ProtoKeeperSaver interface {
	Ping(ctx context.Context) error
}

// RegisterServerAPI регистрация сервера
func RegisterServerAPI(server *grpc.Server, pwdKeeper ProtoKeeperSaver) {
	pb.RegisterGophKeeperServer(server, &pwdKeeperServer{pwdKeeper: pwdKeeper})
}	

// Ping метод для проверки соединения с сервером	
func (s *pwdKeeperServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	err := s.pwdKeeper.Ping(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to ping: %v", err)	
	}
	return &pb.PingResponse{}, nil	
}
