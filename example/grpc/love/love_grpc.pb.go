// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: love/love.proto

package love

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

// LoveClient is the client API for Love service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LoveClient interface {
	// 定义Confession方法
	Confession(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type loveClient struct {
	cc grpc.ClientConnInterface
}

// NewLoveClient ...
func NewLoveClient(cc grpc.ClientConnInterface) LoveClient {
	return &loveClient{cc}
}

// Confession ...
func (c *loveClient) Confession(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/love.Love/Confession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LoveServer is the server API for Love service.
// All implementations must embed UnimplementedLoveServer
// for forward compatibility
type LoveServer interface {
	// 定义Confession方法
	Confession(context.Context, *Request) (*Response, error)
	mustEmbedUnimplementedLoveServer()
}

// UnimplementedLoveServer must be embedded to have forward compatible implementations.
type UnimplementedLoveServer struct {
}

// Confession ...
func (UnimplementedLoveServer) Confession(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Confession not implemented")
}
func (UnimplementedLoveServer) mustEmbedUnimplementedLoveServer() {}

// UnsafeLoveServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LoveServer will
// result in compilation errors.
type UnsafeLoveServer interface {
	mustEmbedUnimplementedLoveServer()
}

// RegisterLoveServer ...
func RegisterLoveServer(s grpc.ServiceRegistrar, srv LoveServer) {
	s.RegisterService(&Love_ServiceDesc, srv)
}

func _Love_Confession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoveServer).Confession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/love.Love/Confession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoveServer).Confession(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// Love_ServiceDesc is the grpc.ServiceDesc for Love service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Love_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "love.Love",
	HandlerType: (*LoveServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Confession",
			Handler:    _Love_Confession_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "love/love.proto",
}
