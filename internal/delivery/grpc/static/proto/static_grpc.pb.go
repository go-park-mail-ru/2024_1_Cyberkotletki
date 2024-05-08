// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: static.proto

package static

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// StaticServiceClient is the client API for StaticService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StaticServiceClient interface {
	GetStatic(ctx context.Context, in *Static, opts ...grpc.CallOption) (*Static, error)
	UploadAvatar(ctx context.Context, opts ...grpc.CallOption) (StaticService_UploadAvatarClient, error)
	GetStaticFile(ctx context.Context, in *Static, opts ...grpc.CallOption) (StaticService_GetStaticFileClient, error)
	Ping(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*Nothing, error)
}

type staticServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStaticServiceClient(cc grpc.ClientConnInterface) StaticServiceClient {
	return &staticServiceClient{cc}
}

func (c *staticServiceClient) GetStatic(ctx context.Context, in *Static, opts ...grpc.CallOption) (*Static, error) {
	out := new(Static)
	err := c.cc.Invoke(ctx, "/static.StaticService/GetStatic", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *staticServiceClient) UploadAvatar(ctx context.Context, opts ...grpc.CallOption) (StaticService_UploadAvatarClient, error) {
	stream, err := c.cc.NewStream(ctx, &StaticService_ServiceDesc.Streams[0], "/static.StaticService/UploadAvatar", opts...)
	if err != nil {
		return nil, err
	}
	x := &staticServiceUploadAvatarClient{stream}
	return x, nil
}

type StaticService_UploadAvatarClient interface {
	Send(*StaticUpload) error
	CloseAndRecv() (*Static, error)
	grpc.ClientStream
}

type staticServiceUploadAvatarClient struct {
	grpc.ClientStream
}

func (x *staticServiceUploadAvatarClient) Send(m *StaticUpload) error {
	return x.ClientStream.SendMsg(m)
}

func (x *staticServiceUploadAvatarClient) CloseAndRecv() (*Static, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Static)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *staticServiceClient) GetStaticFile(ctx context.Context, in *Static, opts ...grpc.CallOption) (StaticService_GetStaticFileClient, error) {
	stream, err := c.cc.NewStream(ctx, &StaticService_ServiceDesc.Streams[1], "/static.StaticService/GetStaticFile", opts...)
	if err != nil {
		return nil, err
	}
	x := &staticServiceGetStaticFileClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type StaticService_GetStaticFileClient interface {
	Recv() (*StaticUpload, error)
	grpc.ClientStream
}

type staticServiceGetStaticFileClient struct {
	grpc.ClientStream
}

func (x *staticServiceGetStaticFileClient) Recv() (*StaticUpload, error) {
	m := new(StaticUpload)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *staticServiceClient) Ping(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*Nothing, error) {
	out := new(Nothing)
	err := c.cc.Invoke(ctx, "/static.StaticService/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StaticServiceServer is the server API for StaticService service.
// All implementations must embed UnimplementedStaticServiceServer
// for forward compatibility
type StaticServiceServer interface {
	GetStatic(context.Context, *Static) (*Static, error)
	UploadAvatar(StaticService_UploadAvatarServer) error
	GetStaticFile(*Static, StaticService_GetStaticFileServer) error
	Ping(context.Context, *Nothing) (*Nothing, error)
	mustEmbedUnimplementedStaticServiceServer()
}

// UnimplementedStaticServiceServer must be embedded to have forward compatible implementations.
type UnimplementedStaticServiceServer struct {
}

func (UnimplementedStaticServiceServer) GetStatic(context.Context, *Static) (*Static, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStatic not implemented")
}
func (UnimplementedStaticServiceServer) UploadAvatar(StaticService_UploadAvatarServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadAvatar not implemented")
}
func (UnimplementedStaticServiceServer) GetStaticFile(*Static, StaticService_GetStaticFileServer) error {
	return status.Errorf(codes.Unimplemented, "method GetStaticFile not implemented")
}
func (UnimplementedStaticServiceServer) Ping(context.Context, *Nothing) (*Nothing, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedStaticServiceServer) mustEmbedUnimplementedStaticServiceServer() {}

// UnsafeStaticServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StaticServiceServer will
// result in compilation errors.
type UnsafeStaticServiceServer interface {
	mustEmbedUnimplementedStaticServiceServer()
}

func RegisterStaticServiceServer(s grpc.ServiceRegistrar, srv StaticServiceServer) {
	s.RegisterService(&StaticService_ServiceDesc, srv)
}

func _StaticService_GetStatic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Static)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StaticServiceServer).GetStatic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/static.StaticService/GetStatic",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StaticServiceServer).GetStatic(ctx, req.(*Static))
	}
	return interceptor(ctx, in, info, handler)
}

func _StaticService_UploadAvatar_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StaticServiceServer).UploadAvatar(&staticServiceUploadAvatarServer{stream})
}

type StaticService_UploadAvatarServer interface {
	SendAndClose(*Static) error
	Recv() (*StaticUpload, error)
	grpc.ServerStream
}

type staticServiceUploadAvatarServer struct {
	grpc.ServerStream
}

func (x *staticServiceUploadAvatarServer) SendAndClose(m *Static) error {
	return x.ServerStream.SendMsg(m)
}

func (x *staticServiceUploadAvatarServer) Recv() (*StaticUpload, error) {
	m := new(StaticUpload)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _StaticService_GetStaticFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Static)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StaticServiceServer).GetStaticFile(m, &staticServiceGetStaticFileServer{stream})
}

type StaticService_GetStaticFileServer interface {
	Send(*StaticUpload) error
	grpc.ServerStream
}

type staticServiceGetStaticFileServer struct {
	grpc.ServerStream
}

func (x *staticServiceGetStaticFileServer) Send(m *StaticUpload) error {
	return x.ServerStream.SendMsg(m)
}

func _StaticService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Nothing)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StaticServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/static.StaticService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StaticServiceServer).Ping(ctx, req.(*Nothing))
	}
	return interceptor(ctx, in, info, handler)
}

// StaticService_ServiceDesc is the grpc.ServiceDesc for StaticService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StaticService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "static.StaticService",
	HandlerType: (*StaticServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetStatic",
			Handler:    _StaticService_GetStatic_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _StaticService_Ping_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadAvatar",
			Handler:       _StaticService_UploadAvatar_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "GetStaticFile",
			Handler:       _StaticService_GetStaticFile_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "static.proto",
}