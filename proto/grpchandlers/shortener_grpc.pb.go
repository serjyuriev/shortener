// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.2
// source: proto/shortener.proto

package grpchandlers

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ShortenerClient is the client API for Shortener service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShortenerClient interface {
	DeleteURLs(ctx context.Context, in *DeleteURLsRequest, opts ...grpc.CallOption) (*DeleteURLsResponse, error)
	GetURL(ctx context.Context, in *GetURLRequest, opts ...grpc.CallOption) (*GetURLResponse, error)
	GetUserURLs(ctx context.Context, in *GetUserURLsRequest, opts ...grpc.CallOption) (*GetUserURLsResponse, error)
	Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PingResponse, error)
	PostBatch(ctx context.Context, in *PostBatchRequest, opts ...grpc.CallOption) (*PostBatchResponse, error)
	PostURL(ctx context.Context, in *PostURLRequest, opts ...grpc.CallOption) (*PostURLResponse, error)
}

type shortenerClient struct {
	cc grpc.ClientConnInterface
}

func NewShortenerClient(cc grpc.ClientConnInterface) ShortenerClient {
	return &shortenerClient{cc}
}

func (c *shortenerClient) DeleteURLs(ctx context.Context, in *DeleteURLsRequest, opts ...grpc.CallOption) (*DeleteURLsResponse, error) {
	out := new(DeleteURLsResponse)
	err := c.cc.Invoke(ctx, "/proto.grpchandlers.Shortener/DeleteURLs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) GetURL(ctx context.Context, in *GetURLRequest, opts ...grpc.CallOption) (*GetURLResponse, error) {
	out := new(GetURLResponse)
	err := c.cc.Invoke(ctx, "/proto.grpchandlers.Shortener/GetURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) GetUserURLs(ctx context.Context, in *GetUserURLsRequest, opts ...grpc.CallOption) (*GetUserURLsResponse, error) {
	out := new(GetUserURLsResponse)
	err := c.cc.Invoke(ctx, "/proto.grpchandlers.Shortener/GetUserURLs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/proto.grpchandlers.Shortener/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) PostBatch(ctx context.Context, in *PostBatchRequest, opts ...grpc.CallOption) (*PostBatchResponse, error) {
	out := new(PostBatchResponse)
	err := c.cc.Invoke(ctx, "/proto.grpchandlers.Shortener/PostBatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) PostURL(ctx context.Context, in *PostURLRequest, opts ...grpc.CallOption) (*PostURLResponse, error) {
	out := new(PostURLResponse)
	err := c.cc.Invoke(ctx, "/proto.grpchandlers.Shortener/PostURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShortenerServer is the server API for Shortener service.
// All implementations must embed UnimplementedShortenerServer
// for forward compatibility
type ShortenerServer interface {
	DeleteURLs(context.Context, *DeleteURLsRequest) (*DeleteURLsResponse, error)
	GetURL(context.Context, *GetURLRequest) (*GetURLResponse, error)
	GetUserURLs(context.Context, *GetUserURLsRequest) (*GetUserURLsResponse, error)
	Ping(context.Context, *emptypb.Empty) (*PingResponse, error)
	PostBatch(context.Context, *PostBatchRequest) (*PostBatchResponse, error)
	PostURL(context.Context, *PostURLRequest) (*PostURLResponse, error)
	mustEmbedUnimplementedShortenerServer()
}

// UnimplementedShortenerServer must be embedded to have forward compatible implementations.
type UnimplementedShortenerServer struct {
}

func (UnimplementedShortenerServer) DeleteURLs(context.Context, *DeleteURLsRequest) (*DeleteURLsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteURLs not implemented")
}
func (UnimplementedShortenerServer) GetURL(context.Context, *GetURLRequest) (*GetURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetURL not implemented")
}
func (UnimplementedShortenerServer) GetUserURLs(context.Context, *GetUserURLsRequest) (*GetUserURLsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserURLs not implemented")
}
func (UnimplementedShortenerServer) Ping(context.Context, *emptypb.Empty) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedShortenerServer) PostBatch(context.Context, *PostBatchRequest) (*PostBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostBatch not implemented")
}
func (UnimplementedShortenerServer) PostURL(context.Context, *PostURLRequest) (*PostURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostURL not implemented")
}
func (UnimplementedShortenerServer) mustEmbedUnimplementedShortenerServer() {}

// UnsafeShortenerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShortenerServer will
// result in compilation errors.
type UnsafeShortenerServer interface {
	mustEmbedUnimplementedShortenerServer()
}

func RegisterShortenerServer(s grpc.ServiceRegistrar, srv ShortenerServer) {
	s.RegisterService(&Shortener_ServiceDesc, srv)
}

func _Shortener_DeleteURLs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteURLsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).DeleteURLs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.grpchandlers.Shortener/DeleteURLs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).DeleteURLs(ctx, req.(*DeleteURLsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_GetURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).GetURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.grpchandlers.Shortener/GetURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).GetURL(ctx, req.(*GetURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_GetUserURLs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserURLsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).GetUserURLs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.grpchandlers.Shortener/GetUserURLs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).GetUserURLs(ctx, req.(*GetUserURLsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.grpchandlers.Shortener/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).Ping(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_PostBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).PostBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.grpchandlers.Shortener/PostBatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).PostBatch(ctx, req.(*PostBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_PostURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).PostURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.grpchandlers.Shortener/PostURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).PostURL(ctx, req.(*PostURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Shortener_ServiceDesc is the grpc.ServiceDesc for Shortener service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Shortener_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.grpchandlers.Shortener",
	HandlerType: (*ShortenerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DeleteURLs",
			Handler:    _Shortener_DeleteURLs_Handler,
		},
		{
			MethodName: "GetURL",
			Handler:    _Shortener_GetURL_Handler,
		},
		{
			MethodName: "GetUserURLs",
			Handler:    _Shortener_GetUserURLs_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _Shortener_Ping_Handler,
		},
		{
			MethodName: "PostBatch",
			Handler:    _Shortener_PostBatch_Handler,
		},
		{
			MethodName: "PostURL",
			Handler:    _Shortener_PostURL_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/shortener.proto",
}
