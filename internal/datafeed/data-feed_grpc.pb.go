// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: internal/datafeed/data-feed.proto

package datafeed

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

const (
	FeedService_GetSports_FullMethodName            = "/datafeed.FeedService/GetSports"
	FeedService_GetCountries_FullMethodName         = "/datafeed.FeedService/GetCountries"
	FeedService_GetTournaments_FullMethodName       = "/datafeed.FeedService/GetTournaments"
	FeedService_GetMarketDefinitions_FullMethodName = "/datafeed.FeedService/GetMarketDefinitions"
	FeedService_GetMatches_FullMethodName           = "/datafeed.FeedService/GetMatches"
	FeedService_GetMatchSnapshots_FullMethodName    = "/datafeed.FeedService/GetMatchSnapshots"
	FeedService_SubscribeToFeed_FullMethodName      = "/datafeed.FeedService/SubscribeToFeed"
	FeedService_SyncData_FullMethodName             = "/datafeed.FeedService/SyncData"
	FeedService_BetControl_FullMethodName           = "/datafeed.FeedService/BetControl"
)

// FeedServiceClient is the client API for FeedService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FeedServiceClient interface {
	GetSports(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SportResponse, error)
	GetCountries(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*CountryResponse, error)
	GetTournaments(ctx context.Context, in *TournamentRequest, opts ...grpc.CallOption) (*TournamentResponse, error)
	GetMarketDefinitions(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*MarketDefinitionsResponse, error)
	GetMatches(ctx context.Context, in *MatchesRequest, opts ...grpc.CallOption) (*MatchesResponse, error)
	GetMatchSnapshots(ctx context.Context, in *MatchSnapshotsRequest, opts ...grpc.CallOption) (*MatchSnapshotsResponse, error)
	SubscribeToFeed(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (FeedService_SubscribeToFeedClient, error)
	SyncData(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SyncDataResponse, error)
	BetControl(ctx context.Context, in *BetControlRequest, opts ...grpc.CallOption) (*BetControlResponse, error)
}

type feedServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFeedServiceClient(cc grpc.ClientConnInterface) FeedServiceClient {
	return &feedServiceClient{cc}
}

func (c *feedServiceClient) GetSports(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SportResponse, error) {
	out := new(SportResponse)
	err := c.cc.Invoke(ctx, FeedService_GetSports_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *feedServiceClient) GetCountries(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*CountryResponse, error) {
	out := new(CountryResponse)
	err := c.cc.Invoke(ctx, FeedService_GetCountries_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *feedServiceClient) GetTournaments(ctx context.Context, in *TournamentRequest, opts ...grpc.CallOption) (*TournamentResponse, error) {
	out := new(TournamentResponse)
	err := c.cc.Invoke(ctx, FeedService_GetTournaments_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *feedServiceClient) GetMarketDefinitions(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*MarketDefinitionsResponse, error) {
	out := new(MarketDefinitionsResponse)
	err := c.cc.Invoke(ctx, FeedService_GetMarketDefinitions_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *feedServiceClient) GetMatches(ctx context.Context, in *MatchesRequest, opts ...grpc.CallOption) (*MatchesResponse, error) {
	out := new(MatchesResponse)
	err := c.cc.Invoke(ctx, FeedService_GetMatches_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *feedServiceClient) GetMatchSnapshots(ctx context.Context, in *MatchSnapshotsRequest, opts ...grpc.CallOption) (*MatchSnapshotsResponse, error) {
	out := new(MatchSnapshotsResponse)
	err := c.cc.Invoke(ctx, FeedService_GetMatchSnapshots_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *feedServiceClient) SubscribeToFeed(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (FeedService_SubscribeToFeedClient, error) {
	stream, err := c.cc.NewStream(ctx, &FeedService_ServiceDesc.Streams[0], FeedService_SubscribeToFeed_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &feedServiceSubscribeToFeedClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type FeedService_SubscribeToFeedClient interface {
	Recv() (*FeedUpdateData, error)
	grpc.ClientStream
}

type feedServiceSubscribeToFeedClient struct {
	grpc.ClientStream
}

func (x *feedServiceSubscribeToFeedClient) Recv() (*FeedUpdateData, error) {
	m := new(FeedUpdateData)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *feedServiceClient) SyncData(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SyncDataResponse, error) {
	out := new(SyncDataResponse)
	err := c.cc.Invoke(ctx, FeedService_SyncData_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *feedServiceClient) BetControl(ctx context.Context, in *BetControlRequest, opts ...grpc.CallOption) (*BetControlResponse, error) {
	out := new(BetControlResponse)
	err := c.cc.Invoke(ctx, FeedService_BetControl_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FeedServiceServer is the server API for FeedService service.
// All implementations must embed UnimplementedFeedServiceServer
// for forward compatibility
type FeedServiceServer interface {
	GetSports(context.Context, *emptypb.Empty) (*SportResponse, error)
	GetCountries(context.Context, *emptypb.Empty) (*CountryResponse, error)
	GetTournaments(context.Context, *TournamentRequest) (*TournamentResponse, error)
	GetMarketDefinitions(context.Context, *emptypb.Empty) (*MarketDefinitionsResponse, error)
	GetMatches(context.Context, *MatchesRequest) (*MatchesResponse, error)
	GetMatchSnapshots(context.Context, *MatchSnapshotsRequest) (*MatchSnapshotsResponse, error)
	SubscribeToFeed(*emptypb.Empty, FeedService_SubscribeToFeedServer) error
	SyncData(context.Context, *emptypb.Empty) (*SyncDataResponse, error)
	BetControl(context.Context, *BetControlRequest) (*BetControlResponse, error)
	mustEmbedUnimplementedFeedServiceServer()
}

// UnimplementedFeedServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFeedServiceServer struct {
}

func (UnimplementedFeedServiceServer) GetSports(context.Context, *emptypb.Empty) (*SportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSports not implemented")
}
func (UnimplementedFeedServiceServer) GetCountries(context.Context, *emptypb.Empty) (*CountryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCountries not implemented")
}
func (UnimplementedFeedServiceServer) GetTournaments(context.Context, *TournamentRequest) (*TournamentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTournaments not implemented")
}
func (UnimplementedFeedServiceServer) GetMarketDefinitions(context.Context, *emptypb.Empty) (*MarketDefinitionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMarketDefinitions not implemented")
}
func (UnimplementedFeedServiceServer) GetMatches(context.Context, *MatchesRequest) (*MatchesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMatches not implemented")
}
func (UnimplementedFeedServiceServer) GetMatchSnapshots(context.Context, *MatchSnapshotsRequest) (*MatchSnapshotsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMatchSnapshots not implemented")
}
func (UnimplementedFeedServiceServer) SubscribeToFeed(*emptypb.Empty, FeedService_SubscribeToFeedServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeToFeed not implemented")
}
func (UnimplementedFeedServiceServer) SyncData(context.Context, *emptypb.Empty) (*SyncDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncData not implemented")
}
func (UnimplementedFeedServiceServer) BetControl(context.Context, *BetControlRequest) (*BetControlResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BetControl not implemented")
}
func (UnimplementedFeedServiceServer) mustEmbedUnimplementedFeedServiceServer() {}

// UnsafeFeedServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FeedServiceServer will
// result in compilation errors.
type UnsafeFeedServiceServer interface {
	mustEmbedUnimplementedFeedServiceServer()
}

func RegisterFeedServiceServer(s grpc.ServiceRegistrar, srv FeedServiceServer) {
	s.RegisterService(&FeedService_ServiceDesc, srv)
}

func _FeedService_GetSports_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedServiceServer).GetSports(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FeedService_GetSports_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedServiceServer).GetSports(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _FeedService_GetCountries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedServiceServer).GetCountries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FeedService_GetCountries_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedServiceServer).GetCountries(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _FeedService_GetTournaments_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TournamentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedServiceServer).GetTournaments(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FeedService_GetTournaments_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedServiceServer).GetTournaments(ctx, req.(*TournamentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FeedService_GetMarketDefinitions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedServiceServer).GetMarketDefinitions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FeedService_GetMarketDefinitions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedServiceServer).GetMarketDefinitions(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _FeedService_GetMatches_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MatchesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedServiceServer).GetMatches(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FeedService_GetMatches_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedServiceServer).GetMatches(ctx, req.(*MatchesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FeedService_GetMatchSnapshots_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MatchSnapshotsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedServiceServer).GetMatchSnapshots(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FeedService_GetMatchSnapshots_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedServiceServer).GetMatchSnapshots(ctx, req.(*MatchSnapshotsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FeedService_SubscribeToFeed_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FeedServiceServer).SubscribeToFeed(m, &feedServiceSubscribeToFeedServer{stream})
}

type FeedService_SubscribeToFeedServer interface {
	Send(*FeedUpdateData) error
	grpc.ServerStream
}

type feedServiceSubscribeToFeedServer struct {
	grpc.ServerStream
}

func (x *feedServiceSubscribeToFeedServer) Send(m *FeedUpdateData) error {
	return x.ServerStream.SendMsg(m)
}

func _FeedService_SyncData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedServiceServer).SyncData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FeedService_SyncData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedServiceServer).SyncData(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _FeedService_BetControl_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BetControlRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedServiceServer).BetControl(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FeedService_BetControl_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedServiceServer).BetControl(ctx, req.(*BetControlRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FeedService_ServiceDesc is the grpc.ServiceDesc for FeedService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FeedService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "datafeed.FeedService",
	HandlerType: (*FeedServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSports",
			Handler:    _FeedService_GetSports_Handler,
		},
		{
			MethodName: "GetCountries",
			Handler:    _FeedService_GetCountries_Handler,
		},
		{
			MethodName: "GetTournaments",
			Handler:    _FeedService_GetTournaments_Handler,
		},
		{
			MethodName: "GetMarketDefinitions",
			Handler:    _FeedService_GetMarketDefinitions_Handler,
		},
		{
			MethodName: "GetMatches",
			Handler:    _FeedService_GetMatches_Handler,
		},
		{
			MethodName: "GetMatchSnapshots",
			Handler:    _FeedService_GetMatchSnapshots_Handler,
		},
		{
			MethodName: "SyncData",
			Handler:    _FeedService_SyncData_Handler,
		},
		{
			MethodName: "BetControl",
			Handler:    _FeedService_BetControl_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeToFeed",
			Handler:       _FeedService_SubscribeToFeed_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "internal/datafeed/data-feed.proto",
}