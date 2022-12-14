// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: api/read.proto

package api

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

// MsgReaderClient is the client API for MsgReader service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MsgReaderClient interface {
	ReadMsg(ctx context.Context, in *ReadMsgRequest, opts ...grpc.CallOption) (*ReadMsgResponse, error)
}

type msgReaderClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgReaderClient(cc grpc.ClientConnInterface) MsgReaderClient {
	return &msgReaderClient{cc}
}

func (c *msgReaderClient) ReadMsg(ctx context.Context, in *ReadMsgRequest, opts ...grpc.CallOption) (*ReadMsgResponse, error) {
	out := new(ReadMsgResponse)
	err := c.cc.Invoke(ctx, "/api.MsgReader/ReadMsg", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgReaderServer is the server API for MsgReader service.
// All implementations must embed UnimplementedMsgReaderServer
// for forward compatibility
type MsgReaderServer interface {
	ReadMsg(context.Context, *ReadMsgRequest) (*ReadMsgResponse, error)
	mustEmbedUnimplementedMsgReaderServer()
}

// UnimplementedMsgReaderServer must be embedded to have forward compatible implementations.
type UnimplementedMsgReaderServer struct {
}

func (UnimplementedMsgReaderServer) ReadMsg(context.Context, *ReadMsgRequest) (*ReadMsgResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadMsg not implemented")
}
func (UnimplementedMsgReaderServer) mustEmbedUnimplementedMsgReaderServer() {}

// UnsafeMsgReaderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MsgReaderServer will
// result in compilation errors.
type UnsafeMsgReaderServer interface {
	mustEmbedUnimplementedMsgReaderServer()
}

func RegisterMsgReaderServer(s grpc.ServiceRegistrar, srv MsgReaderServer) {
	s.RegisterService(&MsgReader_ServiceDesc, srv)
}

func _MsgReader_ReadMsg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadMsgRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgReaderServer).ReadMsg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.MsgReader/ReadMsg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgReaderServer).ReadMsg(ctx, req.(*ReadMsgRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MsgReader_ServiceDesc is the grpc.ServiceDesc for MsgReader service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MsgReader_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.MsgReader",
	HandlerType: (*MsgReaderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ReadMsg",
			Handler:    _MsgReader_ReadMsg_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/read.proto",
}
