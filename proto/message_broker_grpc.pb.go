// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.6
// source: message_broker.proto

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

// MessageBrokerClient is the client API for MessageBroker service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessageBrokerClient interface {
	CreateChannel(ctx context.Context, in *CreateChannelRequest, opts ...grpc.CallOption) (*CreateChannelResponse, error)
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (MessageBroker_SubscribeClient, error)
	CloseChannel(ctx context.Context, in *CloseChannelRequest, opts ...grpc.CallOption) (*CloseChannelResponse, error)
	Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*PublishResponse, error)
}

type messageBrokerClient struct {
	cc grpc.ClientConnInterface
}

func NewMessageBrokerClient(cc grpc.ClientConnInterface) MessageBrokerClient {
	return &messageBrokerClient{cc}
}

func (c *messageBrokerClient) CreateChannel(ctx context.Context, in *CreateChannelRequest, opts ...grpc.CallOption) (*CreateChannelResponse, error) {
	out := new(CreateChannelResponse)
	err := c.cc.Invoke(ctx, "/MessageBroker/CreateChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messageBrokerClient) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (MessageBroker_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &MessageBroker_ServiceDesc.Streams[0], "/MessageBroker/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &messageBrokerSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MessageBroker_SubscribeClient interface {
	Recv() (*SubscribeResponse, error)
	grpc.ClientStream
}

type messageBrokerSubscribeClient struct {
	grpc.ClientStream
}

func (x *messageBrokerSubscribeClient) Recv() (*SubscribeResponse, error) {
	m := new(SubscribeResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *messageBrokerClient) CloseChannel(ctx context.Context, in *CloseChannelRequest, opts ...grpc.CallOption) (*CloseChannelResponse, error) {
	out := new(CloseChannelResponse)
	err := c.cc.Invoke(ctx, "/MessageBroker/CloseChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messageBrokerClient) Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*PublishResponse, error) {
	out := new(PublishResponse)
	err := c.cc.Invoke(ctx, "/MessageBroker/Publish", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MessageBrokerServer is the server API for MessageBroker service.
// All implementations must embed UnimplementedMessageBrokerServer
// for forward compatibility
type MessageBrokerServer interface {
	CreateChannel(context.Context, *CreateChannelRequest) (*CreateChannelResponse, error)
	Subscribe(*SubscribeRequest, MessageBroker_SubscribeServer) error
	CloseChannel(context.Context, *CloseChannelRequest) (*CloseChannelResponse, error)
	Publish(context.Context, *PublishRequest) (*PublishResponse, error)
	mustEmbedUnimplementedMessageBrokerServer()
}

// UnimplementedMessageBrokerServer must be embedded to have forward compatible implementations.
type UnimplementedMessageBrokerServer struct {
}

func (UnimplementedMessageBrokerServer) CreateChannel(context.Context, *CreateChannelRequest) (*CreateChannelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateChannel not implemented")
}
func (UnimplementedMessageBrokerServer) Subscribe(*SubscribeRequest, MessageBroker_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedMessageBrokerServer) CloseChannel(context.Context, *CloseChannelRequest) (*CloseChannelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CloseChannel not implemented")
}
func (UnimplementedMessageBrokerServer) Publish(context.Context, *PublishRequest) (*PublishResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedMessageBrokerServer) mustEmbedUnimplementedMessageBrokerServer() {}

// UnsafeMessageBrokerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessageBrokerServer will
// result in compilation errors.
type UnsafeMessageBrokerServer interface {
	mustEmbedUnimplementedMessageBrokerServer()
}

func RegisterMessageBrokerServer(s grpc.ServiceRegistrar, srv MessageBrokerServer) {
	s.RegisterService(&MessageBroker_ServiceDesc, srv)
}

func _MessageBroker_CreateChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageBrokerServer).CreateChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MessageBroker/CreateChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageBrokerServer).CreateChannel(ctx, req.(*CreateChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessageBroker_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MessageBrokerServer).Subscribe(m, &messageBrokerSubscribeServer{stream})
}

type MessageBroker_SubscribeServer interface {
	Send(*SubscribeResponse) error
	grpc.ServerStream
}

type messageBrokerSubscribeServer struct {
	grpc.ServerStream
}

func (x *messageBrokerSubscribeServer) Send(m *SubscribeResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _MessageBroker_CloseChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CloseChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageBrokerServer).CloseChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MessageBroker/CloseChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageBrokerServer).CloseChannel(ctx, req.(*CloseChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessageBroker_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageBrokerServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MessageBroker/Publish",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageBrokerServer).Publish(ctx, req.(*PublishRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MessageBroker_ServiceDesc is the grpc.ServiceDesc for MessageBroker service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MessageBroker_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "MessageBroker",
	HandlerType: (*MessageBrokerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateChannel",
			Handler:    _MessageBroker_CreateChannel_Handler,
		},
		{
			MethodName: "CloseChannel",
			Handler:    _MessageBroker_CloseChannel_Handler,
		},
		{
			MethodName: "Publish",
			Handler:    _MessageBroker_Publish_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _MessageBroker_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "message_broker.proto",
}
