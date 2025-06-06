// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: academictoken/tokendef/query.proto

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
	Query_Params_FullMethodName                        = "/academictoken.tokendef.Query/Params"
	Query_GetTokenDefinition_FullMethodName            = "/academictoken.tokendef.Query/GetTokenDefinition"
	Query_GetTokenDefinitionFull_FullMethodName        = "/academictoken.tokendef.Query/GetTokenDefinitionFull"
	Query_ListTokenDefinitions_FullMethodName          = "/academictoken.tokendef.Query/ListTokenDefinitions"
	Query_ListTokenDefinitionsBySubject_FullMethodName = "/academictoken.tokendef.Query/ListTokenDefinitionsBySubject"
)

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueryClient interface {
	// Parameters queries the parameters of the module.
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
	// GetTokenDefinition queries a token definition by index
	GetTokenDefinition(ctx context.Context, in *QueryGetTokenDefinitionRequest, opts ...grpc.CallOption) (*QueryGetTokenDefinitionResponse, error)
	// GetTokenDefinitionFull queries a token definition with full content (including IPFS data)
	GetTokenDefinitionFull(ctx context.Context, in *QueryGetTokenDefinitionFullRequest, opts ...grpc.CallOption) (*QueryGetTokenDefinitionFullResponse, error)
	// ListTokenDefinitions queries all token definitions with pagination
	ListTokenDefinitions(ctx context.Context, in *QueryListTokenDefinitionsRequest, opts ...grpc.CallOption) (*QueryListTokenDefinitionsResponse, error)
	// ListTokenDefinitionsBySubject queries token definitions by subject
	ListTokenDefinitionsBySubject(ctx context.Context, in *QueryListTokenDefinitionsBySubjectRequest, opts ...grpc.CallOption) (*QueryListTokenDefinitionsBySubjectResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, Query_Params_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetTokenDefinition(ctx context.Context, in *QueryGetTokenDefinitionRequest, opts ...grpc.CallOption) (*QueryGetTokenDefinitionResponse, error) {
	out := new(QueryGetTokenDefinitionResponse)
	err := c.cc.Invoke(ctx, Query_GetTokenDefinition_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetTokenDefinitionFull(ctx context.Context, in *QueryGetTokenDefinitionFullRequest, opts ...grpc.CallOption) (*QueryGetTokenDefinitionFullResponse, error) {
	out := new(QueryGetTokenDefinitionFullResponse)
	err := c.cc.Invoke(ctx, Query_GetTokenDefinitionFull_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ListTokenDefinitions(ctx context.Context, in *QueryListTokenDefinitionsRequest, opts ...grpc.CallOption) (*QueryListTokenDefinitionsResponse, error) {
	out := new(QueryListTokenDefinitionsResponse)
	err := c.cc.Invoke(ctx, Query_ListTokenDefinitions_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ListTokenDefinitionsBySubject(ctx context.Context, in *QueryListTokenDefinitionsBySubjectRequest, opts ...grpc.CallOption) (*QueryListTokenDefinitionsBySubjectResponse, error) {
	out := new(QueryListTokenDefinitionsBySubjectResponse)
	err := c.cc.Invoke(ctx, Query_ListTokenDefinitionsBySubject_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
// All implementations must embed UnimplementedQueryServer
// for forward compatibility
type QueryServer interface {
	// Parameters queries the parameters of the module.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	// GetTokenDefinition queries a token definition by index
	GetTokenDefinition(context.Context, *QueryGetTokenDefinitionRequest) (*QueryGetTokenDefinitionResponse, error)
	// GetTokenDefinitionFull queries a token definition with full content (including IPFS data)
	GetTokenDefinitionFull(context.Context, *QueryGetTokenDefinitionFullRequest) (*QueryGetTokenDefinitionFullResponse, error)
	// ListTokenDefinitions queries all token definitions with pagination
	ListTokenDefinitions(context.Context, *QueryListTokenDefinitionsRequest) (*QueryListTokenDefinitionsResponse, error)
	// ListTokenDefinitionsBySubject queries token definitions by subject
	ListTokenDefinitionsBySubject(context.Context, *QueryListTokenDefinitionsBySubjectRequest) (*QueryListTokenDefinitionsBySubjectResponse, error)
	mustEmbedUnimplementedQueryServer()
}

// UnimplementedQueryServer must be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (UnimplementedQueryServer) Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (UnimplementedQueryServer) GetTokenDefinition(context.Context, *QueryGetTokenDefinitionRequest) (*QueryGetTokenDefinitionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTokenDefinition not implemented")
}
func (UnimplementedQueryServer) GetTokenDefinitionFull(context.Context, *QueryGetTokenDefinitionFullRequest) (*QueryGetTokenDefinitionFullResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTokenDefinitionFull not implemented")
}
func (UnimplementedQueryServer) ListTokenDefinitions(context.Context, *QueryListTokenDefinitionsRequest) (*QueryListTokenDefinitionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTokenDefinitions not implemented")
}
func (UnimplementedQueryServer) ListTokenDefinitionsBySubject(context.Context, *QueryListTokenDefinitionsBySubjectRequest) (*QueryListTokenDefinitionsBySubjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTokenDefinitionsBySubject not implemented")
}
func (UnimplementedQueryServer) mustEmbedUnimplementedQueryServer() {}

// UnsafeQueryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueryServer will
// result in compilation errors.
type UnsafeQueryServer interface {
	mustEmbedUnimplementedQueryServer()
}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&Query_ServiceDesc, srv)
}

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Params_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetTokenDefinition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetTokenDefinitionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetTokenDefinition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetTokenDefinition_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetTokenDefinition(ctx, req.(*QueryGetTokenDefinitionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetTokenDefinitionFull_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetTokenDefinitionFullRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetTokenDefinitionFull(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetTokenDefinitionFull_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetTokenDefinitionFull(ctx, req.(*QueryGetTokenDefinitionFullRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ListTokenDefinitions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryListTokenDefinitionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ListTokenDefinitions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_ListTokenDefinitions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ListTokenDefinitions(ctx, req.(*QueryListTokenDefinitionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ListTokenDefinitionsBySubject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryListTokenDefinitionsBySubjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ListTokenDefinitionsBySubject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_ListTokenDefinitionsBySubject_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ListTokenDefinitionsBySubject(ctx, req.(*QueryListTokenDefinitionsBySubjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Query_ServiceDesc is the grpc.ServiceDesc for Query service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "academictoken.tokendef.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "GetTokenDefinition",
			Handler:    _Query_GetTokenDefinition_Handler,
		},
		{
			MethodName: "GetTokenDefinitionFull",
			Handler:    _Query_GetTokenDefinitionFull_Handler,
		},
		{
			MethodName: "ListTokenDefinitions",
			Handler:    _Query_ListTokenDefinitions_Handler,
		},
		{
			MethodName: "ListTokenDefinitionsBySubject",
			Handler:    _Query_ListTokenDefinitionsBySubject_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "academictoken/tokendef/query.proto",
}
