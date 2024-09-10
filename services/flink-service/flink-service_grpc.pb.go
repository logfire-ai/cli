// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.27.0
// source: flink-service.proto

package logfire

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	FilterService_GetFilteredData_FullMethodName    = "/ai.logfire.FilterService/GetFilteredData"
	FilterService_GetStreamData_FullMethodName      = "/ai.logfire.FilterService/GetStreamData"
	FilterService_SubmitSQL_FullMethodName          = "/ai.logfire.FilterService/SubmitSQL"
	FilterService_SubmitAlertRequest_FullMethodName = "/ai.logfire.FilterService/SubmitAlertRequest"
	FilterService_DeleteAlertRequest_FullMethodName = "/ai.logfire.FilterService/DeleteAlertRequest"
	FilterService_GetOffsetData_FullMethodName      = "/ai.logfire.FilterService/GetOffsetData"
)

// FilterServiceClient is the client API for FilterService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FilterServiceClient interface {
	GetFilteredData(ctx context.Context, in *FilterRequest, opts ...grpc.CallOption) (*FilteredRecords, error)
	GetStreamData(ctx context.Context, in *FilterRequest, opts ...grpc.CallOption) (FilterService_GetStreamDataClient, error)
	SubmitSQL(ctx context.Context, in *SQLRequest, opts ...grpc.CallOption) (*SQLResponse, error)
	SubmitAlertRequest(ctx context.Context, in *AlertRequest, opts ...grpc.CallOption) (*RegisteredAlert, error)
	DeleteAlertRequest(ctx context.Context, in *RegisteredAlert, opts ...grpc.CallOption) (*Empty, error)
	GetOffsetData(ctx context.Context, in *OffsetRequest, opts ...grpc.CallOption) (*OffsetResponse, error)
}

type filterServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFilterServiceClient(cc grpc.ClientConnInterface) FilterServiceClient {
	return &filterServiceClient{cc}
}

func (c *filterServiceClient) GetFilteredData(ctx context.Context, in *FilterRequest, opts ...grpc.CallOption) (*FilteredRecords, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FilteredRecords)
	err := c.cc.Invoke(ctx, FilterService_GetFilteredData_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *filterServiceClient) GetStreamData(ctx context.Context, in *FilterRequest, opts ...grpc.CallOption) (FilterService_GetStreamDataClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &FilterService_ServiceDesc.Streams[0], FilterService_GetStreamData_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &filterServiceGetStreamDataClient{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type FilterService_GetStreamDataClient interface {
	Recv() (*FilteredRecords, error)
	grpc.ClientStream
}

type filterServiceGetStreamDataClient struct {
	grpc.ClientStream
}

func (x *filterServiceGetStreamDataClient) Recv() (*FilteredRecords, error) {
	m := new(FilteredRecords)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *filterServiceClient) SubmitSQL(ctx context.Context, in *SQLRequest, opts ...grpc.CallOption) (*SQLResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SQLResponse)
	err := c.cc.Invoke(ctx, FilterService_SubmitSQL_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *filterServiceClient) SubmitAlertRequest(ctx context.Context, in *AlertRequest, opts ...grpc.CallOption) (*RegisteredAlert, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisteredAlert)
	err := c.cc.Invoke(ctx, FilterService_SubmitAlertRequest_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *filterServiceClient) DeleteAlertRequest(ctx context.Context, in *RegisteredAlert, opts ...grpc.CallOption) (*Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Empty)
	err := c.cc.Invoke(ctx, FilterService_DeleteAlertRequest_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *filterServiceClient) GetOffsetData(ctx context.Context, in *OffsetRequest, opts ...grpc.CallOption) (*OffsetResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OffsetResponse)
	err := c.cc.Invoke(ctx, FilterService_GetOffsetData_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FilterServiceServer is the server API for FilterService service.
// All implementations must embed UnimplementedFilterServiceServer
// for forward compatibility
type FilterServiceServer interface {
	GetFilteredData(context.Context, *FilterRequest) (*FilteredRecords, error)
	GetStreamData(*FilterRequest, FilterService_GetStreamDataServer) error
	SubmitSQL(context.Context, *SQLRequest) (*SQLResponse, error)
	SubmitAlertRequest(context.Context, *AlertRequest) (*RegisteredAlert, error)
	DeleteAlertRequest(context.Context, *RegisteredAlert) (*Empty, error)
	GetOffsetData(context.Context, *OffsetRequest) (*OffsetResponse, error)
	mustEmbedUnimplementedFilterServiceServer()
}

// UnimplementedFilterServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFilterServiceServer struct {
}

func (UnimplementedFilterServiceServer) GetFilteredData(context.Context, *FilterRequest) (*FilteredRecords, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFilteredData not implemented")
}
func (UnimplementedFilterServiceServer) GetStreamData(*FilterRequest, FilterService_GetStreamDataServer) error {
	return status.Errorf(codes.Unimplemented, "method GetStreamData not implemented")
}
func (UnimplementedFilterServiceServer) SubmitSQL(context.Context, *SQLRequest) (*SQLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitSQL not implemented")
}
func (UnimplementedFilterServiceServer) SubmitAlertRequest(context.Context, *AlertRequest) (*RegisteredAlert, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitAlertRequest not implemented")
}
func (UnimplementedFilterServiceServer) DeleteAlertRequest(context.Context, *RegisteredAlert) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAlertRequest not implemented")
}
func (UnimplementedFilterServiceServer) GetOffsetData(context.Context, *OffsetRequest) (*OffsetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOffsetData not implemented")
}
func (UnimplementedFilterServiceServer) mustEmbedUnimplementedFilterServiceServer() {}

// UnsafeFilterServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FilterServiceServer will
// result in compilation errors.
type UnsafeFilterServiceServer interface {
	mustEmbedUnimplementedFilterServiceServer()
}

func RegisterFilterServiceServer(s grpc.ServiceRegistrar, srv FilterServiceServer) {
	s.RegisterService(&FilterService_ServiceDesc, srv)
}

func _FilterService_GetFilteredData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FilterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilterServiceServer).GetFilteredData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FilterService_GetFilteredData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServiceServer).GetFilteredData(ctx, req.(*FilterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FilterService_GetStreamData_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(FilterRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FilterServiceServer).GetStreamData(m, &filterServiceGetStreamDataServer{ServerStream: stream})
}

type FilterService_GetStreamDataServer interface {
	Send(*FilteredRecords) error
	grpc.ServerStream
}

type filterServiceGetStreamDataServer struct {
	grpc.ServerStream
}

func (x *filterServiceGetStreamDataServer) Send(m *FilteredRecords) error {
	return x.ServerStream.SendMsg(m)
}

func _FilterService_SubmitSQL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SQLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilterServiceServer).SubmitSQL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FilterService_SubmitSQL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServiceServer).SubmitSQL(ctx, req.(*SQLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FilterService_SubmitAlertRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AlertRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilterServiceServer).SubmitAlertRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FilterService_SubmitAlertRequest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServiceServer).SubmitAlertRequest(ctx, req.(*AlertRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FilterService_DeleteAlertRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisteredAlert)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilterServiceServer).DeleteAlertRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FilterService_DeleteAlertRequest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServiceServer).DeleteAlertRequest(ctx, req.(*RegisteredAlert))
	}
	return interceptor(ctx, in, info, handler)
}

func _FilterService_GetOffsetData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OffsetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilterServiceServer).GetOffsetData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FilterService_GetOffsetData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServiceServer).GetOffsetData(ctx, req.(*OffsetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FilterService_ServiceDesc is the grpc.ServiceDesc for FilterService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FilterService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ai.logfire.FilterService",
	HandlerType: (*FilterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetFilteredData",
			Handler:    _FilterService_GetFilteredData_Handler,
		},
		{
			MethodName: "SubmitSQL",
			Handler:    _FilterService_SubmitSQL_Handler,
		},
		{
			MethodName: "SubmitAlertRequest",
			Handler:    _FilterService_SubmitAlertRequest_Handler,
		},
		{
			MethodName: "DeleteAlertRequest",
			Handler:    _FilterService_DeleteAlertRequest_Handler,
		},
		{
			MethodName: "GetOffsetData",
			Handler:    _FilterService_GetOffsetData_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetStreamData",
			Handler:       _FilterService_GetStreamData_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "flink-service.proto",
}

const (
	MetaService_GetBarGraph_FullMethodName = "/ai.logfire.MetaService/GetBarGraph"
	MetaService_GetStatus_FullMethodName   = "/ai.logfire.MetaService/GetStatus"
)

// MetaServiceClient is the client API for MetaService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetaServiceClient interface {
	GetBarGraph(ctx context.Context, in *GraphRequest, opts ...grpc.CallOption) (MetaService_GetBarGraphClient, error)
	GetStatus(ctx context.Context, in *GraphRequest, opts ...grpc.CallOption) (MetaService_GetStatusClient, error)
}

type metaServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMetaServiceClient(cc grpc.ClientConnInterface) MetaServiceClient {
	return &metaServiceClient{cc}
}

func (c *metaServiceClient) GetBarGraph(ctx context.Context, in *GraphRequest, opts ...grpc.CallOption) (MetaService_GetBarGraphClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &MetaService_ServiceDesc.Streams[0], MetaService_GetBarGraph_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &metaServiceGetBarGraphClient{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MetaService_GetBarGraphClient interface {
	Recv() (*GraphResponse, error)
	grpc.ClientStream
}

type metaServiceGetBarGraphClient struct {
	grpc.ClientStream
}

func (x *metaServiceGetBarGraphClient) Recv() (*GraphResponse, error) {
	m := new(GraphResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *metaServiceClient) GetStatus(ctx context.Context, in *GraphRequest, opts ...grpc.CallOption) (MetaService_GetStatusClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &MetaService_ServiceDesc.Streams[1], MetaService_GetStatus_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &metaServiceGetStatusClient{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MetaService_GetStatusClient interface {
	Recv() (*GraphResponse, error)
	grpc.ClientStream
}

type metaServiceGetStatusClient struct {
	grpc.ClientStream
}

func (x *metaServiceGetStatusClient) Recv() (*GraphResponse, error) {
	m := new(GraphResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MetaServiceServer is the server API for MetaService service.
// All implementations must embed UnimplementedMetaServiceServer
// for forward compatibility
type MetaServiceServer interface {
	GetBarGraph(*GraphRequest, MetaService_GetBarGraphServer) error
	GetStatus(*GraphRequest, MetaService_GetStatusServer) error
	mustEmbedUnimplementedMetaServiceServer()
}

// UnimplementedMetaServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMetaServiceServer struct {
}

func (UnimplementedMetaServiceServer) GetBarGraph(*GraphRequest, MetaService_GetBarGraphServer) error {
	return status.Errorf(codes.Unimplemented, "method GetBarGraph not implemented")
}
func (UnimplementedMetaServiceServer) GetStatus(*GraphRequest, MetaService_GetStatusServer) error {
	return status.Errorf(codes.Unimplemented, "method GetStatus not implemented")
}
func (UnimplementedMetaServiceServer) mustEmbedUnimplementedMetaServiceServer() {}

// UnsafeMetaServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetaServiceServer will
// result in compilation errors.
type UnsafeMetaServiceServer interface {
	mustEmbedUnimplementedMetaServiceServer()
}

func RegisterMetaServiceServer(s grpc.ServiceRegistrar, srv MetaServiceServer) {
	s.RegisterService(&MetaService_ServiceDesc, srv)
}

func _MetaService_GetBarGraph_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GraphRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MetaServiceServer).GetBarGraph(m, &metaServiceGetBarGraphServer{ServerStream: stream})
}

type MetaService_GetBarGraphServer interface {
	Send(*GraphResponse) error
	grpc.ServerStream
}

type metaServiceGetBarGraphServer struct {
	grpc.ServerStream
}

func (x *metaServiceGetBarGraphServer) Send(m *GraphResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _MetaService_GetStatus_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GraphRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MetaServiceServer).GetStatus(m, &metaServiceGetStatusServer{ServerStream: stream})
}

type MetaService_GetStatusServer interface {
	Send(*GraphResponse) error
	grpc.ServerStream
}

type metaServiceGetStatusServer struct {
	grpc.ServerStream
}

func (x *metaServiceGetStatusServer) Send(m *GraphResponse) error {
	return x.ServerStream.SendMsg(m)
}

// MetaService_ServiceDesc is the grpc.ServiceDesc for MetaService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MetaService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ai.logfire.MetaService",
	HandlerType: (*MetaServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetBarGraph",
			Handler:       _MetaService_GetBarGraph_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetStatus",
			Handler:       _MetaService_GetStatus_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "flink-service.proto",
}

const (
	NotificationService_ReceiveNotification_FullMethodName = "/ai.logfire.NotificationService/ReceiveNotification"
	NotificationService_SendNotification_FullMethodName    = "/ai.logfire.NotificationService/SendNotification"
)

// NotificationServiceClient is the client API for NotificationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NotificationServiceClient interface {
	ReceiveNotification(ctx context.Context, in *ReceiveNotificationRequest, opts ...grpc.CallOption) (NotificationService_ReceiveNotificationClient, error)
	SendNotification(ctx context.Context, in *SendNotificationRequest, opts ...grpc.CallOption) (*SendNotificationResponse, error)
}

type notificationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNotificationServiceClient(cc grpc.ClientConnInterface) NotificationServiceClient {
	return &notificationServiceClient{cc}
}

func (c *notificationServiceClient) ReceiveNotification(ctx context.Context, in *ReceiveNotificationRequest, opts ...grpc.CallOption) (NotificationService_ReceiveNotificationClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &NotificationService_ServiceDesc.Streams[0], NotificationService_ReceiveNotification_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &notificationServiceReceiveNotificationClient{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type NotificationService_ReceiveNotificationClient interface {
	Recv() (*ReceiveNotificationResponse, error)
	grpc.ClientStream
}

type notificationServiceReceiveNotificationClient struct {
	grpc.ClientStream
}

func (x *notificationServiceReceiveNotificationClient) Recv() (*ReceiveNotificationResponse, error) {
	m := new(ReceiveNotificationResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *notificationServiceClient) SendNotification(ctx context.Context, in *SendNotificationRequest, opts ...grpc.CallOption) (*SendNotificationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SendNotificationResponse)
	err := c.cc.Invoke(ctx, NotificationService_SendNotification_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NotificationServiceServer is the server API for NotificationService service.
// All implementations must embed UnimplementedNotificationServiceServer
// for forward compatibility
type NotificationServiceServer interface {
	ReceiveNotification(*ReceiveNotificationRequest, NotificationService_ReceiveNotificationServer) error
	SendNotification(context.Context, *SendNotificationRequest) (*SendNotificationResponse, error)
	mustEmbedUnimplementedNotificationServiceServer()
}

// UnimplementedNotificationServiceServer must be embedded to have forward compatible implementations.
type UnimplementedNotificationServiceServer struct {
}

func (UnimplementedNotificationServiceServer) ReceiveNotification(*ReceiveNotificationRequest, NotificationService_ReceiveNotificationServer) error {
	return status.Errorf(codes.Unimplemented, "method ReceiveNotification not implemented")
}
func (UnimplementedNotificationServiceServer) SendNotification(context.Context, *SendNotificationRequest) (*SendNotificationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendNotification not implemented")
}
func (UnimplementedNotificationServiceServer) mustEmbedUnimplementedNotificationServiceServer() {}

// UnsafeNotificationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NotificationServiceServer will
// result in compilation errors.
type UnsafeNotificationServiceServer interface {
	mustEmbedUnimplementedNotificationServiceServer()
}

func RegisterNotificationServiceServer(s grpc.ServiceRegistrar, srv NotificationServiceServer) {
	s.RegisterService(&NotificationService_ServiceDesc, srv)
}

func _NotificationService_ReceiveNotification_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ReceiveNotificationRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(NotificationServiceServer).ReceiveNotification(m, &notificationServiceReceiveNotificationServer{ServerStream: stream})
}

type NotificationService_ReceiveNotificationServer interface {
	Send(*ReceiveNotificationResponse) error
	grpc.ServerStream
}

type notificationServiceReceiveNotificationServer struct {
	grpc.ServerStream
}

func (x *notificationServiceReceiveNotificationServer) Send(m *ReceiveNotificationResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _NotificationService_SendNotification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendNotificationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).SendNotification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_SendNotification_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).SendNotification(ctx, req.(*SendNotificationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NotificationService_ServiceDesc is the grpc.ServiceDesc for NotificationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NotificationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ai.logfire.NotificationService",
	HandlerType: (*NotificationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendNotification",
			Handler:    _NotificationService_SendNotification_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ReceiveNotification",
			Handler:       _NotificationService_ReceiveNotification_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "flink-service.proto",
}
