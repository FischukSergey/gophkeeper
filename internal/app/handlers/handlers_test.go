package handlers

import (
	"context"
	"reflect"
	"testing"

	pb "github.com/FischukSergey/gophkeeper/internal/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fields struct {
	UnimplementedGophKeeperServer pb.UnimplementedGophKeeperServer
	pwdKeeper                     ProtoKeeperSaver
}

func Test_pwdKeeperServer_Registration(t *testing.T) {
	type args struct {
		req *pb.RegistrationRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.RegistrationResponse
		wantErr bool
		status  codes.Code
	}{
		{
			name: "login is empty",
			args: args{
				req: &pb.RegistrationRequest{
					Username: "",
					Password: "test",
				},
			},
			want:    nil,
			wantErr: true,
			status:  codes.InvalidArgument,
		},
		{
			name: "password is empty",
			args: args{
				req: &pb.RegistrationRequest{
					Username: "test",
					Password: "",
				},
			},
			want:    nil,
			wantErr: true,
			status:  codes.InvalidArgument,
		},
		{
			name: "incorrect login",
			args: args{
				req: &pb.RegistrationRequest{
					Username: "1234",
					Password: "test",
				},
			},
			want:    nil,
			wantErr: true,
			status:  codes.InvalidArgument,
		},
		{
			name: "incorrect password",
			args: args{
				req: &pb.RegistrationRequest{
					Username: "test",
					Password: "1234",
				},
			},
			want:    nil,
			wantErr: true,
			status:  codes.InvalidArgument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &pwdKeeperServer{
				UnimplementedGophKeeperServer: tt.fields.UnimplementedGophKeeperServer,
				pwdKeeper:                     tt.fields.pwdKeeper,
			}
			got, err := s.Registration(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("pwdKeeperServer.Registration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pwdKeeperServer.Registration() = %v, want %v", got, tt.want)
			}
			if status.Code(err) != tt.status {
				t.Errorf("pwdKeeperServer.Registration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_pwdKeeperServer_Authorization(t *testing.T) {
	type args struct {
		req *pb.AuthorizationRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.AuthorizationResponse
		wantErr bool
		status  codes.Code
	}{
		{
			name: "login is empty",
			args: args{
				req: &pb.AuthorizationRequest{
					Username: "",
					Password: "test",
				},
			},
			want:    nil,
			wantErr: true,
			status:  codes.InvalidArgument,
		},
		{
			name: "password is empty",
			args: args{
				req: &pb.AuthorizationRequest{
					Username: "test",
					Password: "",
				},
			},
			want:    nil,
			wantErr: true,
			status:  codes.InvalidArgument,
		},
		{
			name: "incorrect login",
			args: args{
				req: &pb.AuthorizationRequest{
					Username: "test",
					Password: "test123",
				},
			},
			want:    nil,
			wantErr: true,
			status:  codes.InvalidArgument,
		},
		{
			name: "incorrect password",
			args: args{
				req: &pb.AuthorizationRequest{
					Username: "test123",
					Password: "test",
				},
			},
			want:    nil,
			wantErr: true,
			status:  codes.InvalidArgument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &pwdKeeperServer{
				UnimplementedGophKeeperServer: tt.fields.UnimplementedGophKeeperServer,
				pwdKeeper:                     tt.fields.pwdKeeper,
			}
			got, err := s.Authorization(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("pwdKeeperServer.Authorization() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pwdKeeperServer.Authorization() = %v, want %v", got, tt.want)
			}
			if status.Code(err) != tt.status {
				t.Errorf("pwdKeeperServer.Authorization() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_pwdKeeperServer_Ping(t *testing.T) {
	type args struct {
		req *pb.PingRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.PingResponse
		wantErr bool
		status  codes.Code
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &pwdKeeperServer{
				UnimplementedGophKeeperServer: tt.fields.UnimplementedGophKeeperServer,
				pwdKeeper:                     tt.fields.pwdKeeper,
			}
			got, err := s.Ping(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("pwdKeeperServer.Ping() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pwdKeeperServer.Ping() = %v, want %v", got, tt.want)
			}
			if status.Code(err) != tt.status {
				t.Errorf("pwdKeeperServer.Ping() error = %v, wantStatus %v", status.Code(err), tt.status)
			}
		})
	}
}
