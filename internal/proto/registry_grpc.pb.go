// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: internal/proto/registry.proto

package gophkeeper

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	GophKeeper_Ping_FullMethodName          = "/server.GophKeeper/Ping"
	GophKeeper_Registration_FullMethodName  = "/server.GophKeeper/Registration"
	GophKeeper_Authorization_FullMethodName = "/server.GophKeeper/Authorization"
	GophKeeper_FileUpload_FullMethodName    = "/server.GophKeeper/FileUpload"
	GophKeeper_FileDownload_FullMethodName  = "/server.GophKeeper/FileDownload"
	GophKeeper_FileDelete_FullMethodName    = "/server.GophKeeper/FileDelete"
	GophKeeper_FileGetList_FullMethodName   = "/server.GophKeeper/FileGetList"
	GophKeeper_NoteAdd_FullMethodName       = "/server.GophKeeper/NoteAdd"
	GophKeeper_NoteGetList_FullMethodName   = "/server.GophKeeper/NoteGetList"
	GophKeeper_NoteUpdate_FullMethodName    = "/server.GophKeeper/NoteUpdate"
	GophKeeper_NoteDelete_FullMethodName    = "/server.GophKeeper/NoteDelete"
)

// GophKeeperClient is the client API for GophKeeper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GophKeeperClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	Registration(ctx context.Context, in *RegistrationRequest, opts ...grpc.CallOption) (*RegistrationResponse, error)
	Authorization(ctx context.Context, in *AuthorizationRequest, opts ...grpc.CallOption) (*AuthorizationResponse, error)
	FileUpload(ctx context.Context, in *FileUploadRequest, opts ...grpc.CallOption) (*FileUploadResponse, error)
	FileDownload(ctx context.Context, in *FileDownloadRequest, opts ...grpc.CallOption) (*FileDownloadResponse, error)
	FileDelete(ctx context.Context, in *FileDeleteRequest, opts ...grpc.CallOption) (*FileDeleteResponse, error)
	FileGetList(ctx context.Context, in *FileGetListRequest, opts ...grpc.CallOption) (*FileGetListResponse, error)
	NoteAdd(ctx context.Context, in *NoteAddRequest, opts ...grpc.CallOption) (*NoteAddResponse, error)
	NoteGetList(ctx context.Context, in *NoteGetListRequest, opts ...grpc.CallOption) (*NoteGetListResponse, error)
	NoteUpdate(ctx context.Context, in *NoteUpdateRequest, opts ...grpc.CallOption) (*NoteUpdateResponse, error)
	NoteDelete(ctx context.Context, in *NoteDeleteRequest, opts ...grpc.CallOption) (*NoteDeleteResponse, error)
}

type gophKeeperClient struct {
	cc grpc.ClientConnInterface
}

func NewGophKeeperClient(cc grpc.ClientConnInterface) GophKeeperClient {
	return &gophKeeperClient{cc}
}

func (c *gophKeeperClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, GophKeeper_Ping_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) Registration(ctx context.Context, in *RegistrationRequest, opts ...grpc.CallOption) (*RegistrationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegistrationResponse)
	err := c.cc.Invoke(ctx, GophKeeper_Registration_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) Authorization(ctx context.Context, in *AuthorizationRequest, opts ...grpc.CallOption) (*AuthorizationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AuthorizationResponse)
	err := c.cc.Invoke(ctx, GophKeeper_Authorization_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) FileUpload(ctx context.Context, in *FileUploadRequest, opts ...grpc.CallOption) (*FileUploadResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FileUploadResponse)
	err := c.cc.Invoke(ctx, GophKeeper_FileUpload_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) FileDownload(ctx context.Context, in *FileDownloadRequest, opts ...grpc.CallOption) (*FileDownloadResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FileDownloadResponse)
	err := c.cc.Invoke(ctx, GophKeeper_FileDownload_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) FileDelete(ctx context.Context, in *FileDeleteRequest, opts ...grpc.CallOption) (*FileDeleteResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FileDeleteResponse)
	err := c.cc.Invoke(ctx, GophKeeper_FileDelete_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) FileGetList(ctx context.Context, in *FileGetListRequest, opts ...grpc.CallOption) (*FileGetListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FileGetListResponse)
	err := c.cc.Invoke(ctx, GophKeeper_FileGetList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) NoteAdd(ctx context.Context, in *NoteAddRequest, opts ...grpc.CallOption) (*NoteAddResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NoteAddResponse)
	err := c.cc.Invoke(ctx, GophKeeper_NoteAdd_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) NoteGetList(ctx context.Context, in *NoteGetListRequest, opts ...grpc.CallOption) (*NoteGetListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NoteGetListResponse)
	err := c.cc.Invoke(ctx, GophKeeper_NoteGetList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) NoteUpdate(ctx context.Context, in *NoteUpdateRequest, opts ...grpc.CallOption) (*NoteUpdateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NoteUpdateResponse)
	err := c.cc.Invoke(ctx, GophKeeper_NoteUpdate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) NoteDelete(ctx context.Context, in *NoteDeleteRequest, opts ...grpc.CallOption) (*NoteDeleteResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NoteDeleteResponse)
	err := c.cc.Invoke(ctx, GophKeeper_NoteDelete_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GophKeeperServer is the server API for GophKeeper service.
// All implementations must embed UnimplementedGophKeeperServer
// for forward compatibility.
type GophKeeperServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	Registration(context.Context, *RegistrationRequest) (*RegistrationResponse, error)
	Authorization(context.Context, *AuthorizationRequest) (*AuthorizationResponse, error)
	FileUpload(context.Context, *FileUploadRequest) (*FileUploadResponse, error)
	FileDownload(context.Context, *FileDownloadRequest) (*FileDownloadResponse, error)
	FileDelete(context.Context, *FileDeleteRequest) (*FileDeleteResponse, error)
	FileGetList(context.Context, *FileGetListRequest) (*FileGetListResponse, error)
	NoteAdd(context.Context, *NoteAddRequest) (*NoteAddResponse, error)
	NoteGetList(context.Context, *NoteGetListRequest) (*NoteGetListResponse, error)
	NoteUpdate(context.Context, *NoteUpdateRequest) (*NoteUpdateResponse, error)
	NoteDelete(context.Context, *NoteDeleteRequest) (*NoteDeleteResponse, error)
	mustEmbedUnimplementedGophKeeperServer()
}

// UnimplementedGophKeeperServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedGophKeeperServer struct{}

func (UnimplementedGophKeeperServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedGophKeeperServer) Registration(context.Context, *RegistrationRequest) (*RegistrationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Registration not implemented")
}
func (UnimplementedGophKeeperServer) Authorization(context.Context, *AuthorizationRequest) (*AuthorizationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authorization not implemented")
}
func (UnimplementedGophKeeperServer) FileUpload(context.Context, *FileUploadRequest) (*FileUploadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FileUpload not implemented")
}
func (UnimplementedGophKeeperServer) FileDownload(context.Context, *FileDownloadRequest) (*FileDownloadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FileDownload not implemented")
}
func (UnimplementedGophKeeperServer) FileDelete(context.Context, *FileDeleteRequest) (*FileDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FileDelete not implemented")
}
func (UnimplementedGophKeeperServer) FileGetList(context.Context, *FileGetListRequest) (*FileGetListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FileGetList not implemented")
}
func (UnimplementedGophKeeperServer) NoteAdd(context.Context, *NoteAddRequest) (*NoteAddResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NoteAdd not implemented")
}
func (UnimplementedGophKeeperServer) NoteGetList(context.Context, *NoteGetListRequest) (*NoteGetListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NoteGetList not implemented")
}
func (UnimplementedGophKeeperServer) NoteUpdate(context.Context, *NoteUpdateRequest) (*NoteUpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NoteUpdate not implemented")
}
func (UnimplementedGophKeeperServer) NoteDelete(context.Context, *NoteDeleteRequest) (*NoteDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NoteDelete not implemented")
}
func (UnimplementedGophKeeperServer) mustEmbedUnimplementedGophKeeperServer() {}
func (UnimplementedGophKeeperServer) testEmbeddedByValue()                    {}

// UnsafeGophKeeperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GophKeeperServer will
// result in compilation errors.
type UnsafeGophKeeperServer interface {
	mustEmbedUnimplementedGophKeeperServer()
}

func RegisterGophKeeperServer(s grpc.ServiceRegistrar, srv GophKeeperServer) {
	// If the following call pancis, it indicates UnimplementedGophKeeperServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&GophKeeper_ServiceDesc, srv)
}

func _GophKeeper_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_Registration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegistrationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).Registration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_Registration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).Registration(ctx, req.(*RegistrationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_Authorization_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthorizationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).Authorization(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_Authorization_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).Authorization(ctx, req.(*AuthorizationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_FileUpload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileUploadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).FileUpload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_FileUpload_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).FileUpload(ctx, req.(*FileUploadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_FileDownload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileDownloadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).FileDownload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_FileDownload_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).FileDownload(ctx, req.(*FileDownloadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_FileDelete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).FileDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_FileDelete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).FileDelete(ctx, req.(*FileDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_FileGetList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileGetListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).FileGetList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_FileGetList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).FileGetList(ctx, req.(*FileGetListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_NoteAdd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NoteAddRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).NoteAdd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_NoteAdd_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).NoteAdd(ctx, req.(*NoteAddRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_NoteGetList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NoteGetListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).NoteGetList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_NoteGetList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).NoteGetList(ctx, req.(*NoteGetListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_NoteUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NoteUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).NoteUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_NoteUpdate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).NoteUpdate(ctx, req.(*NoteUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_NoteDelete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NoteDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).NoteDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_NoteDelete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).NoteDelete(ctx, req.(*NoteDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GophKeeper_ServiceDesc is the grpc.ServiceDesc for GophKeeper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GophKeeper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "server.GophKeeper",
	HandlerType: (*GophKeeperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _GophKeeper_Ping_Handler,
		},
		{
			MethodName: "Registration",
			Handler:    _GophKeeper_Registration_Handler,
		},
		{
			MethodName: "Authorization",
			Handler:    _GophKeeper_Authorization_Handler,
		},
		{
			MethodName: "FileUpload",
			Handler:    _GophKeeper_FileUpload_Handler,
		},
		{
			MethodName: "FileDownload",
			Handler:    _GophKeeper_FileDownload_Handler,
		},
		{
			MethodName: "FileDelete",
			Handler:    _GophKeeper_FileDelete_Handler,
		},
		{
			MethodName: "FileGetList",
			Handler:    _GophKeeper_FileGetList_Handler,
		},
		{
			MethodName: "NoteAdd",
			Handler:    _GophKeeper_NoteAdd_Handler,
		},
		{
			MethodName: "NoteGetList",
			Handler:    _GophKeeper_NoteGetList_Handler,
		},
		{
			MethodName: "NoteUpdate",
			Handler:    _GophKeeper_NoteUpdate_Handler,
		},
		{
			MethodName: "NoteDelete",
			Handler:    _GophKeeper_NoteDelete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/proto/registry.proto",
}
