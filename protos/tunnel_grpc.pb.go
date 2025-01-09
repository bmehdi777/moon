// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.19.6
// source: protos/tunnel.proto

package protos

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
	Tunnel_Stream_FullMethodName = "/tunnel.Tunnel/Stream"
)

// TunnelClient is the client API for Tunnel service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TunnelClient interface {
	Stream(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[StreamClient, StreamServer], error)
}

type tunnelClient struct {
	cc grpc.ClientConnInterface
}

func NewTunnelClient(cc grpc.ClientConnInterface) TunnelClient {
	return &tunnelClient{cc}
}

func (c *tunnelClient) Stream(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[StreamClient, StreamServer], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Tunnel_ServiceDesc.Streams[0], Tunnel_Stream_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[StreamClient, StreamServer]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Tunnel_StreamClient = grpc.BidiStreamingClient[StreamClient, StreamServer]

// TunnelServer is the server API for Tunnel service.
// All implementations must embed UnimplementedTunnelServer
// for forward compatibility.
type TunnelServer interface {
	Stream(grpc.BidiStreamingServer[StreamClient, StreamServer]) error
	mustEmbedUnimplementedTunnelServer()
}

// UnimplementedTunnelServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedTunnelServer struct{}

func (UnimplementedTunnelServer) Stream(grpc.BidiStreamingServer[StreamClient, StreamServer]) error {
	return status.Errorf(codes.Unimplemented, "method Stream not implemented")
}
func (UnimplementedTunnelServer) mustEmbedUnimplementedTunnelServer() {}
func (UnimplementedTunnelServer) testEmbeddedByValue()                {}

// UnsafeTunnelServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TunnelServer will
// result in compilation errors.
type UnsafeTunnelServer interface {
	mustEmbedUnimplementedTunnelServer()
}

func RegisterTunnelServer(s grpc.ServiceRegistrar, srv TunnelServer) {
	// If the following call pancis, it indicates UnimplementedTunnelServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Tunnel_ServiceDesc, srv)
}

func _Tunnel_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TunnelServer).Stream(&grpc.GenericServerStream[StreamClient, StreamServer]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Tunnel_StreamServer = grpc.BidiStreamingServer[StreamClient, StreamServer]

// Tunnel_ServiceDesc is the grpc.ServiceDesc for Tunnel service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Tunnel_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tunnel.Tunnel",
	HandlerType: (*TunnelServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _Tunnel_Stream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "protos/tunnel.proto",
}
