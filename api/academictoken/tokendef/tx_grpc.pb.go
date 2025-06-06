// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: academictoken/tokendef/tx.proto

package tokendef

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
	Msg_UpdateParams_FullMethodName          = "/academictoken.tokendef.Msg/UpdateParams"
	Msg_CreateTokenDefinition_FullMethodName = "/academictoken.tokendef.Msg/CreateTokenDefinition"
	Msg_UpdateTokenDefinition_FullMethodName = "/academictoken.tokendef.Msg/UpdateTokenDefinition"
)

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MsgClient interface {
	// UpdateParams defines a (governance) operation for updating the module
	// parameters. The authority defaults to the x/gov module account.
	UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error)
	// CreateTokenDefinition creates a new token definition for a subject
	CreateTokenDefinition(ctx context.Context, in *MsgCreateTokenDefinition, opts ...grpc.CallOption) (*MsgCreateTokenDefinitionResponse, error)
	// UpdateTokenDefinition updates an existing token definition
	UpdateTokenDefinition(ctx context.Context, in *MsgUpdateTokenDefinition, opts ...grpc.CallOption) (*MsgUpdateTokenDefinitionResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error) {
	out := new(MsgUpdateParamsResponse)
	err := c.cc.Invoke(ctx, Msg_UpdateParams_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) CreateTokenDefinition(ctx context.Context, in *MsgCreateTokenDefinition, opts ...grpc.CallOption) (*MsgCreateTokenDefinitionResponse, error) {
	out := new(MsgCreateTokenDefinitionResponse)
	err := c.cc.Invoke(ctx, Msg_CreateTokenDefinition_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateTokenDefinition(ctx context.Context, in *MsgUpdateTokenDefinition, opts ...grpc.CallOption) (*MsgUpdateTokenDefinitionResponse, error) {
	out := new(MsgUpdateTokenDefinitionResponse)
	err := c.cc.Invoke(ctx, Msg_UpdateTokenDefinition_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations must embed UnimplementedMsgServer
// for forward compatibility
type MsgServer interface {
	// UpdateParams defines a (governance) operation for updating the module
	// parameters. The authority defaults to the x/gov module account.
	UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error)
	// CreateTokenDefinition creates a new token definition for a subject
	CreateTokenDefinition(context.Context, *MsgCreateTokenDefinition) (*MsgCreateTokenDefinitionResponse, error)
	// UpdateTokenDefinition updates an existing token definition
	UpdateTokenDefinition(context.Context, *MsgUpdateTokenDefinition) (*MsgUpdateTokenDefinitionResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (UnimplementedMsgServer) UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateParams not implemented")
}
func (UnimplementedMsgServer) CreateTokenDefinition(context.Context, *MsgCreateTokenDefinition) (*MsgCreateTokenDefinitionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTokenDefinition not implemented")
}
func (UnimplementedMsgServer) UpdateTokenDefinition(context.Context, *MsgUpdateTokenDefinition) (*MsgUpdateTokenDefinitionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTokenDefinition not implemented")
}
func (UnimplementedMsgServer) mustEmbedUnimplementedMsgServer() {}

// UnsafeMsgServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MsgServer will
// result in compilation errors.
type UnsafeMsgServer interface {
	mustEmbedUnimplementedMsgServer()
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&Msg_ServiceDesc, srv)
}

func _Msg_UpdateParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_UpdateParams_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateParams(ctx, req.(*MsgUpdateParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CreateTokenDefinition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateTokenDefinition)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateTokenDefinition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_CreateTokenDefinition_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateTokenDefinition(ctx, req.(*MsgCreateTokenDefinition))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateTokenDefinition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateTokenDefinition)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateTokenDefinition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_UpdateTokenDefinition_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateTokenDefinition(ctx, req.(*MsgUpdateTokenDefinition))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "academictoken.tokendef.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateParams",
			Handler:    _Msg_UpdateParams_Handler,
		},
		{
			MethodName: "CreateTokenDefinition",
			Handler:    _Msg_CreateTokenDefinition_Handler,
		},
		{
			MethodName: "UpdateTokenDefinition",
			Handler:    _Msg_UpdateTokenDefinition_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "academictoken/tokendef/tx.proto",
}
