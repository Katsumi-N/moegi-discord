// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package grpc

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

// ConohaServiceClient is the client API for ConohaService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConohaServiceClient interface {
	Minecraft(ctx context.Context, in *MinecraftRequest, opts ...grpc.CallOption) (ConohaService_MinecraftClient, error)
}

type conohaServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewConohaServiceClient(cc grpc.ClientConnInterface) ConohaServiceClient {
	return &conohaServiceClient{cc}
}

func (c *conohaServiceClient) Minecraft(ctx context.Context, in *MinecraftRequest, opts ...grpc.CallOption) (ConohaService_MinecraftClient, error) {
	stream, err := c.cc.NewStream(ctx, &ConohaService_ServiceDesc.Streams[0], "/ConohaService/Minecraft", opts...)
	if err != nil {
		return nil, err
	}
	x := &conohaServiceMinecraftClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ConohaService_MinecraftClient interface {
	Recv() (*MinecraftResponse, error)
	grpc.ClientStream
}

type conohaServiceMinecraftClient struct {
	grpc.ClientStream
}

func (x *conohaServiceMinecraftClient) Recv() (*MinecraftResponse, error) {
	m := new(MinecraftResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ConohaServiceServer is the server API for ConohaService service.
// All implementations must embed UnimplementedConohaServiceServer
// for forward compatibility
type ConohaServiceServer interface {
	Minecraft(*MinecraftRequest, ConohaService_MinecraftServer) error
	mustEmbedUnimplementedConohaServiceServer()
}

// UnimplementedConohaServiceServer must be embedded to have forward compatible implementations.
type UnimplementedConohaServiceServer struct {
}

func (UnimplementedConohaServiceServer) Minecraft(*MinecraftRequest, ConohaService_MinecraftServer) error {
	return status.Errorf(codes.Unimplemented, "method Minecraft not implemented")
}
func (UnimplementedConohaServiceServer) mustEmbedUnimplementedConohaServiceServer() {}

// UnsafeConohaServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConohaServiceServer will
// result in compilation errors.
type UnsafeConohaServiceServer interface {
	mustEmbedUnimplementedConohaServiceServer()
}

func RegisterConohaServiceServer(s grpc.ServiceRegistrar, srv ConohaServiceServer) {
	s.RegisterService(&ConohaService_ServiceDesc, srv)
}

func _ConohaService_Minecraft_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(MinecraftRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ConohaServiceServer).Minecraft(m, &conohaServiceMinecraftServer{stream})
}

type ConohaService_MinecraftServer interface {
	Send(*MinecraftResponse) error
	grpc.ServerStream
}

type conohaServiceMinecraftServer struct {
	grpc.ServerStream
}

func (x *conohaServiceMinecraftServer) Send(m *MinecraftResponse) error {
	return x.ServerStream.SendMsg(m)
}

// ConohaService_ServiceDesc is the grpc.ServiceDesc for ConohaService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ConohaService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ConohaService",
	HandlerType: (*ConohaServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Minecraft",
			Handler:       _ConohaService_Minecraft_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/conoha.proto",
}