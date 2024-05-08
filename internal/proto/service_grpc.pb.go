// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: internal/proto/service.proto

package proto

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

const (
	SportsbookService_ListPrematch_FullMethodName = "/proto.SportsbookService/ListPrematch"
)

// SportsbookServiceClient is the client API for SportsbookService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SportsbookServiceClient interface {
	// Sends a request to list all permatches available
	ListPrematch(ctx context.Context, in *ListPrematchRequest, opts ...grpc.CallOption) (*ListPrematchResponse, error)
}

type sportsbookServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSportsbookServiceClient(cc grpc.ClientConnInterface) SportsbookServiceClient {
	return &sportsbookServiceClient{cc}
}

func (c *sportsbookServiceClient) ListPrematch(ctx context.Context, in *ListPrematchRequest, opts ...grpc.CallOption) (*ListPrematchResponse, error) {
	out := new(ListPrematchResponse)
	err := c.cc.Invoke(ctx, SportsbookService_ListPrematch_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SportsbookServiceServer is the server API for SportsbookService service.
// All implementations must embed UnimplementedSportsbookServiceServer
// for forward compatibility
type SportsbookServiceServer interface {
	// Sends a request to list all permatches available
	ListPrematch(context.Context, *ListPrematchRequest) (*ListPrematchResponse, error)
	mustEmbedUnimplementedSportsbookServiceServer()
}

// UnimplementedSportsbookServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSportsbookServiceServer struct {
}

func (UnimplementedSportsbookServiceServer) ListPrematch(context.Context, *ListPrematchRequest) (*ListPrematchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPrematch not implemented")
}
func (UnimplementedSportsbookServiceServer) mustEmbedUnimplementedSportsbookServiceServer() {}

// UnsafeSportsbookServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SportsbookServiceServer will
// result in compilation errors.
type UnsafeSportsbookServiceServer interface {
	mustEmbedUnimplementedSportsbookServiceServer()
}

func RegisterSportsbookServiceServer(s grpc.ServiceRegistrar, srv SportsbookServiceServer) {
	s.RegisterService(&SportsbookService_ServiceDesc, srv)
}

func _SportsbookService_ListPrematch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPrematchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SportsbookServiceServer).ListPrematch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SportsbookService_ListPrematch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SportsbookServiceServer).ListPrematch(ctx, req.(*ListPrematchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SportsbookService_ServiceDesc is the grpc.ServiceDesc for SportsbookService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SportsbookService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.SportsbookService",
	HandlerType: (*SportsbookServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListPrematch",
			Handler:    _SportsbookService_ListPrematch_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/proto/service.proto",
}
