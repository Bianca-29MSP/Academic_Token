// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: academictoken/equivalence/query.proto

package equivalence

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
	Query_Params_FullMethodName                           = "/academictoken.equivalence.Query/Params"
	Query_ListEquivalences_FullMethodName                 = "/academictoken.equivalence.Query/ListEquivalences"
	Query_GetEquivalence_FullMethodName                   = "/academictoken.equivalence.Query/GetEquivalence"
	Query_GetEquivalencesBySourceSubject_FullMethodName   = "/academictoken.equivalence.Query/GetEquivalencesBySourceSubject"
	Query_GetEquivalencesByTargetSubject_FullMethodName   = "/academictoken.equivalence.Query/GetEquivalencesByTargetSubject"
	Query_GetEquivalencesByInstitution_FullMethodName     = "/academictoken.equivalence.Query/GetEquivalencesByInstitution"
	Query_CheckEquivalenceStatus_FullMethodName           = "/academictoken.equivalence.Query/CheckEquivalenceStatus"
	Query_GetPendingAnalysis_FullMethodName               = "/academictoken.equivalence.Query/GetPendingAnalysis"
	Query_GetApprovedEquivalences_FullMethodName          = "/academictoken.equivalence.Query/GetApprovedEquivalences"
	Query_GetRejectedEquivalences_FullMethodName          = "/academictoken.equivalence.Query/GetRejectedEquivalences"
	Query_GetEquivalencesByContract_FullMethodName        = "/academictoken.equivalence.Query/GetEquivalencesByContract"
	Query_GetEquivalencesByContractVersion_FullMethodName = "/academictoken.equivalence.Query/GetEquivalencesByContractVersion"
	Query_GetEquivalenceHistory_FullMethodName            = "/academictoken.equivalence.Query/GetEquivalenceHistory"
	Query_GetEquivalenceStats_FullMethodName              = "/academictoken.equivalence.Query/GetEquivalenceStats"
	Query_GetAnalysisMetadata_FullMethodName              = "/academictoken.equivalence.Query/GetAnalysisMetadata"
	Query_VerifyAnalysisIntegrity_FullMethodName          = "/academictoken.equivalence.Query/VerifyAnalysisIntegrity"
)

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueryClient interface {
	// Parameters queries the parameters of the module.
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
	// ListEquivalences queries all subject equivalences with pagination
	ListEquivalences(ctx context.Context, in *QueryListEquivalencesRequest, opts ...grpc.CallOption) (*QueryListEquivalencesResponse, error)
	// GetEquivalence queries a specific subject equivalence by index
	GetEquivalence(ctx context.Context, in *QueryGetEquivalenceRequest, opts ...grpc.CallOption) (*QueryGetEquivalenceResponse, error)
	// GetEquivalencesBySourceSubject queries equivalences by source subject ID
	GetEquivalencesBySourceSubject(ctx context.Context, in *QueryGetEquivalencesBySourceSubjectRequest, opts ...grpc.CallOption) (*QueryGetEquivalencesBySourceSubjectResponse, error)
	// GetEquivalencesByTargetSubject queries equivalences by target subject ID
	GetEquivalencesByTargetSubject(ctx context.Context, in *QueryGetEquivalencesByTargetSubjectRequest, opts ...grpc.CallOption) (*QueryGetEquivalencesByTargetSubjectResponse, error)
	// GetEquivalencesByInstitution queries equivalences by target institution
	GetEquivalencesByInstitution(ctx context.Context, in *QueryGetEquivalencesByInstitutionRequest, opts ...grpc.CallOption) (*QueryGetEquivalencesByInstitutionResponse, error)
	// CheckEquivalenceStatus checks if two subjects have an established equivalence
	CheckEquivalenceStatus(ctx context.Context, in *QueryCheckEquivalenceStatusRequest, opts ...grpc.CallOption) (*QueryCheckEquivalenceStatusResponse, error)
	// GetPendingAnalysis queries equivalences awaiting contract analysis
	GetPendingAnalysis(ctx context.Context, in *QueryGetPendingAnalysisRequest, opts ...grpc.CallOption) (*QueryGetPendingAnalysisResponse, error)
	// GetApprovedEquivalences queries equivalences with approved status (by contract)
	GetApprovedEquivalences(ctx context.Context, in *QueryGetApprovedEquivalencesRequest, opts ...grpc.CallOption) (*QueryGetApprovedEquivalencesResponse, error)
	// GetRejectedEquivalences queries equivalences rejected by contract analysis
	GetRejectedEquivalences(ctx context.Context, in *QueryGetRejectedEquivalencesRequest, opts ...grpc.CallOption) (*QueryGetRejectedEquivalencesResponse, error)
	// GetEquivalencesByContract queries equivalences analyzed by a specific contract
	GetEquivalencesByContract(ctx context.Context, in *QueryGetEquivalencesByContractRequest, opts ...grpc.CallOption) (*QueryGetEquivalencesByContractResponse, error)
	// GetEquivalencesByContractVersion queries equivalences by contract version
	GetEquivalencesByContractVersion(ctx context.Context, in *QueryGetEquivalencesByContractVersionRequest, opts ...grpc.CallOption) (*QueryGetEquivalencesByContractVersionResponse, error)
	// GetEquivalenceHistory queries the analysis history of equivalence requests for a subject
	GetEquivalenceHistory(ctx context.Context, in *QueryGetEquivalenceHistoryRequest, opts ...grpc.CallOption) (*QueryGetEquivalenceHistoryResponse, error)
	// GetEquivalenceStats queries statistics about automated equivalence analysis
	GetEquivalenceStats(ctx context.Context, in *QueryGetEquivalenceStatsRequest, opts ...grpc.CallOption) (*QueryGetEquivalenceStatsResponse, error)
	// GetAnalysisMetadata queries detailed analysis metadata for an equivalence
	GetAnalysisMetadata(ctx context.Context, in *QueryGetAnalysisMetadataRequest, opts ...grpc.CallOption) (*QueryGetAnalysisMetadataResponse, error)
	// VerifyAnalysisIntegrity verifies the integrity of an equivalence analysis
	VerifyAnalysisIntegrity(ctx context.Context, in *QueryVerifyAnalysisIntegrityRequest, opts ...grpc.CallOption) (*QueryVerifyAnalysisIntegrityResponse, error)
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

func (c *queryClient) ListEquivalences(ctx context.Context, in *QueryListEquivalencesRequest, opts ...grpc.CallOption) (*QueryListEquivalencesResponse, error) {
	out := new(QueryListEquivalencesResponse)
	err := c.cc.Invoke(ctx, Query_ListEquivalences_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetEquivalence(ctx context.Context, in *QueryGetEquivalenceRequest, opts ...grpc.CallOption) (*QueryGetEquivalenceResponse, error) {
	out := new(QueryGetEquivalenceResponse)
	err := c.cc.Invoke(ctx, Query_GetEquivalence_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetEquivalencesBySourceSubject(ctx context.Context, in *QueryGetEquivalencesBySourceSubjectRequest, opts ...grpc.CallOption) (*QueryGetEquivalencesBySourceSubjectResponse, error) {
	out := new(QueryGetEquivalencesBySourceSubjectResponse)
	err := c.cc.Invoke(ctx, Query_GetEquivalencesBySourceSubject_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetEquivalencesByTargetSubject(ctx context.Context, in *QueryGetEquivalencesByTargetSubjectRequest, opts ...grpc.CallOption) (*QueryGetEquivalencesByTargetSubjectResponse, error) {
	out := new(QueryGetEquivalencesByTargetSubjectResponse)
	err := c.cc.Invoke(ctx, Query_GetEquivalencesByTargetSubject_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetEquivalencesByInstitution(ctx context.Context, in *QueryGetEquivalencesByInstitutionRequest, opts ...grpc.CallOption) (*QueryGetEquivalencesByInstitutionResponse, error) {
	out := new(QueryGetEquivalencesByInstitutionResponse)
	err := c.cc.Invoke(ctx, Query_GetEquivalencesByInstitution_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) CheckEquivalenceStatus(ctx context.Context, in *QueryCheckEquivalenceStatusRequest, opts ...grpc.CallOption) (*QueryCheckEquivalenceStatusResponse, error) {
	out := new(QueryCheckEquivalenceStatusResponse)
	err := c.cc.Invoke(ctx, Query_CheckEquivalenceStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetPendingAnalysis(ctx context.Context, in *QueryGetPendingAnalysisRequest, opts ...grpc.CallOption) (*QueryGetPendingAnalysisResponse, error) {
	out := new(QueryGetPendingAnalysisResponse)
	err := c.cc.Invoke(ctx, Query_GetPendingAnalysis_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetApprovedEquivalences(ctx context.Context, in *QueryGetApprovedEquivalencesRequest, opts ...grpc.CallOption) (*QueryGetApprovedEquivalencesResponse, error) {
	out := new(QueryGetApprovedEquivalencesResponse)
	err := c.cc.Invoke(ctx, Query_GetApprovedEquivalences_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetRejectedEquivalences(ctx context.Context, in *QueryGetRejectedEquivalencesRequest, opts ...grpc.CallOption) (*QueryGetRejectedEquivalencesResponse, error) {
	out := new(QueryGetRejectedEquivalencesResponse)
	err := c.cc.Invoke(ctx, Query_GetRejectedEquivalences_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetEquivalencesByContract(ctx context.Context, in *QueryGetEquivalencesByContractRequest, opts ...grpc.CallOption) (*QueryGetEquivalencesByContractResponse, error) {
	out := new(QueryGetEquivalencesByContractResponse)
	err := c.cc.Invoke(ctx, Query_GetEquivalencesByContract_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetEquivalencesByContractVersion(ctx context.Context, in *QueryGetEquivalencesByContractVersionRequest, opts ...grpc.CallOption) (*QueryGetEquivalencesByContractVersionResponse, error) {
	out := new(QueryGetEquivalencesByContractVersionResponse)
	err := c.cc.Invoke(ctx, Query_GetEquivalencesByContractVersion_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetEquivalenceHistory(ctx context.Context, in *QueryGetEquivalenceHistoryRequest, opts ...grpc.CallOption) (*QueryGetEquivalenceHistoryResponse, error) {
	out := new(QueryGetEquivalenceHistoryResponse)
	err := c.cc.Invoke(ctx, Query_GetEquivalenceHistory_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetEquivalenceStats(ctx context.Context, in *QueryGetEquivalenceStatsRequest, opts ...grpc.CallOption) (*QueryGetEquivalenceStatsResponse, error) {
	out := new(QueryGetEquivalenceStatsResponse)
	err := c.cc.Invoke(ctx, Query_GetEquivalenceStats_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetAnalysisMetadata(ctx context.Context, in *QueryGetAnalysisMetadataRequest, opts ...grpc.CallOption) (*QueryGetAnalysisMetadataResponse, error) {
	out := new(QueryGetAnalysisMetadataResponse)
	err := c.cc.Invoke(ctx, Query_GetAnalysisMetadata_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) VerifyAnalysisIntegrity(ctx context.Context, in *QueryVerifyAnalysisIntegrityRequest, opts ...grpc.CallOption) (*QueryVerifyAnalysisIntegrityResponse, error) {
	out := new(QueryVerifyAnalysisIntegrityResponse)
	err := c.cc.Invoke(ctx, Query_VerifyAnalysisIntegrity_FullMethodName, in, out, opts...)
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
	// ListEquivalences queries all subject equivalences with pagination
	ListEquivalences(context.Context, *QueryListEquivalencesRequest) (*QueryListEquivalencesResponse, error)
	// GetEquivalence queries a specific subject equivalence by index
	GetEquivalence(context.Context, *QueryGetEquivalenceRequest) (*QueryGetEquivalenceResponse, error)
	// GetEquivalencesBySourceSubject queries equivalences by source subject ID
	GetEquivalencesBySourceSubject(context.Context, *QueryGetEquivalencesBySourceSubjectRequest) (*QueryGetEquivalencesBySourceSubjectResponse, error)
	// GetEquivalencesByTargetSubject queries equivalences by target subject ID
	GetEquivalencesByTargetSubject(context.Context, *QueryGetEquivalencesByTargetSubjectRequest) (*QueryGetEquivalencesByTargetSubjectResponse, error)
	// GetEquivalencesByInstitution queries equivalences by target institution
	GetEquivalencesByInstitution(context.Context, *QueryGetEquivalencesByInstitutionRequest) (*QueryGetEquivalencesByInstitutionResponse, error)
	// CheckEquivalenceStatus checks if two subjects have an established equivalence
	CheckEquivalenceStatus(context.Context, *QueryCheckEquivalenceStatusRequest) (*QueryCheckEquivalenceStatusResponse, error)
	// GetPendingAnalysis queries equivalences awaiting contract analysis
	GetPendingAnalysis(context.Context, *QueryGetPendingAnalysisRequest) (*QueryGetPendingAnalysisResponse, error)
	// GetApprovedEquivalences queries equivalences with approved status (by contract)
	GetApprovedEquivalences(context.Context, *QueryGetApprovedEquivalencesRequest) (*QueryGetApprovedEquivalencesResponse, error)
	// GetRejectedEquivalences queries equivalences rejected by contract analysis
	GetRejectedEquivalences(context.Context, *QueryGetRejectedEquivalencesRequest) (*QueryGetRejectedEquivalencesResponse, error)
	// GetEquivalencesByContract queries equivalences analyzed by a specific contract
	GetEquivalencesByContract(context.Context, *QueryGetEquivalencesByContractRequest) (*QueryGetEquivalencesByContractResponse, error)
	// GetEquivalencesByContractVersion queries equivalences by contract version
	GetEquivalencesByContractVersion(context.Context, *QueryGetEquivalencesByContractVersionRequest) (*QueryGetEquivalencesByContractVersionResponse, error)
	// GetEquivalenceHistory queries the analysis history of equivalence requests for a subject
	GetEquivalenceHistory(context.Context, *QueryGetEquivalenceHistoryRequest) (*QueryGetEquivalenceHistoryResponse, error)
	// GetEquivalenceStats queries statistics about automated equivalence analysis
	GetEquivalenceStats(context.Context, *QueryGetEquivalenceStatsRequest) (*QueryGetEquivalenceStatsResponse, error)
	// GetAnalysisMetadata queries detailed analysis metadata for an equivalence
	GetAnalysisMetadata(context.Context, *QueryGetAnalysisMetadataRequest) (*QueryGetAnalysisMetadataResponse, error)
	// VerifyAnalysisIntegrity verifies the integrity of an equivalence analysis
	VerifyAnalysisIntegrity(context.Context, *QueryVerifyAnalysisIntegrityRequest) (*QueryVerifyAnalysisIntegrityResponse, error)
	mustEmbedUnimplementedQueryServer()
}

// UnimplementedQueryServer must be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (UnimplementedQueryServer) Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (UnimplementedQueryServer) ListEquivalences(context.Context, *QueryListEquivalencesRequest) (*QueryListEquivalencesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEquivalences not implemented")
}
func (UnimplementedQueryServer) GetEquivalence(context.Context, *QueryGetEquivalenceRequest) (*QueryGetEquivalenceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEquivalence not implemented")
}
func (UnimplementedQueryServer) GetEquivalencesBySourceSubject(context.Context, *QueryGetEquivalencesBySourceSubjectRequest) (*QueryGetEquivalencesBySourceSubjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEquivalencesBySourceSubject not implemented")
}
func (UnimplementedQueryServer) GetEquivalencesByTargetSubject(context.Context, *QueryGetEquivalencesByTargetSubjectRequest) (*QueryGetEquivalencesByTargetSubjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEquivalencesByTargetSubject not implemented")
}
func (UnimplementedQueryServer) GetEquivalencesByInstitution(context.Context, *QueryGetEquivalencesByInstitutionRequest) (*QueryGetEquivalencesByInstitutionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEquivalencesByInstitution not implemented")
}
func (UnimplementedQueryServer) CheckEquivalenceStatus(context.Context, *QueryCheckEquivalenceStatusRequest) (*QueryCheckEquivalenceStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckEquivalenceStatus not implemented")
}
func (UnimplementedQueryServer) GetPendingAnalysis(context.Context, *QueryGetPendingAnalysisRequest) (*QueryGetPendingAnalysisResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPendingAnalysis not implemented")
}
func (UnimplementedQueryServer) GetApprovedEquivalences(context.Context, *QueryGetApprovedEquivalencesRequest) (*QueryGetApprovedEquivalencesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetApprovedEquivalences not implemented")
}
func (UnimplementedQueryServer) GetRejectedEquivalences(context.Context, *QueryGetRejectedEquivalencesRequest) (*QueryGetRejectedEquivalencesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRejectedEquivalences not implemented")
}
func (UnimplementedQueryServer) GetEquivalencesByContract(context.Context, *QueryGetEquivalencesByContractRequest) (*QueryGetEquivalencesByContractResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEquivalencesByContract not implemented")
}
func (UnimplementedQueryServer) GetEquivalencesByContractVersion(context.Context, *QueryGetEquivalencesByContractVersionRequest) (*QueryGetEquivalencesByContractVersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEquivalencesByContractVersion not implemented")
}
func (UnimplementedQueryServer) GetEquivalenceHistory(context.Context, *QueryGetEquivalenceHistoryRequest) (*QueryGetEquivalenceHistoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEquivalenceHistory not implemented")
}
func (UnimplementedQueryServer) GetEquivalenceStats(context.Context, *QueryGetEquivalenceStatsRequest) (*QueryGetEquivalenceStatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEquivalenceStats not implemented")
}
func (UnimplementedQueryServer) GetAnalysisMetadata(context.Context, *QueryGetAnalysisMetadataRequest) (*QueryGetAnalysisMetadataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAnalysisMetadata not implemented")
}
func (UnimplementedQueryServer) VerifyAnalysisIntegrity(context.Context, *QueryVerifyAnalysisIntegrityRequest) (*QueryVerifyAnalysisIntegrityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyAnalysisIntegrity not implemented")
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

func _Query_ListEquivalences_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryListEquivalencesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ListEquivalences(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_ListEquivalences_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ListEquivalences(ctx, req.(*QueryListEquivalencesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetEquivalence_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetEquivalenceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetEquivalence(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetEquivalence_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetEquivalence(ctx, req.(*QueryGetEquivalenceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetEquivalencesBySourceSubject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetEquivalencesBySourceSubjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetEquivalencesBySourceSubject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetEquivalencesBySourceSubject_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetEquivalencesBySourceSubject(ctx, req.(*QueryGetEquivalencesBySourceSubjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetEquivalencesByTargetSubject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetEquivalencesByTargetSubjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetEquivalencesByTargetSubject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetEquivalencesByTargetSubject_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetEquivalencesByTargetSubject(ctx, req.(*QueryGetEquivalencesByTargetSubjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetEquivalencesByInstitution_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetEquivalencesByInstitutionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetEquivalencesByInstitution(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetEquivalencesByInstitution_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetEquivalencesByInstitution(ctx, req.(*QueryGetEquivalencesByInstitutionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_CheckEquivalenceStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryCheckEquivalenceStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).CheckEquivalenceStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_CheckEquivalenceStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).CheckEquivalenceStatus(ctx, req.(*QueryCheckEquivalenceStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetPendingAnalysis_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetPendingAnalysisRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetPendingAnalysis(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetPendingAnalysis_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetPendingAnalysis(ctx, req.(*QueryGetPendingAnalysisRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetApprovedEquivalences_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetApprovedEquivalencesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetApprovedEquivalences(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetApprovedEquivalences_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetApprovedEquivalences(ctx, req.(*QueryGetApprovedEquivalencesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetRejectedEquivalences_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetRejectedEquivalencesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetRejectedEquivalences(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetRejectedEquivalences_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetRejectedEquivalences(ctx, req.(*QueryGetRejectedEquivalencesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetEquivalencesByContract_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetEquivalencesByContractRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetEquivalencesByContract(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetEquivalencesByContract_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetEquivalencesByContract(ctx, req.(*QueryGetEquivalencesByContractRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetEquivalencesByContractVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetEquivalencesByContractVersionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetEquivalencesByContractVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetEquivalencesByContractVersion_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetEquivalencesByContractVersion(ctx, req.(*QueryGetEquivalencesByContractVersionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetEquivalenceHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetEquivalenceHistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetEquivalenceHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetEquivalenceHistory_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetEquivalenceHistory(ctx, req.(*QueryGetEquivalenceHistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetEquivalenceStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetEquivalenceStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetEquivalenceStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetEquivalenceStats_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetEquivalenceStats(ctx, req.(*QueryGetEquivalenceStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetAnalysisMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetAnalysisMetadataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetAnalysisMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetAnalysisMetadata_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetAnalysisMetadata(ctx, req.(*QueryGetAnalysisMetadataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_VerifyAnalysisIntegrity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryVerifyAnalysisIntegrityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).VerifyAnalysisIntegrity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_VerifyAnalysisIntegrity_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).VerifyAnalysisIntegrity(ctx, req.(*QueryVerifyAnalysisIntegrityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Query_ServiceDesc is the grpc.ServiceDesc for Query service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "academictoken.equivalence.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "ListEquivalences",
			Handler:    _Query_ListEquivalences_Handler,
		},
		{
			MethodName: "GetEquivalence",
			Handler:    _Query_GetEquivalence_Handler,
		},
		{
			MethodName: "GetEquivalencesBySourceSubject",
			Handler:    _Query_GetEquivalencesBySourceSubject_Handler,
		},
		{
			MethodName: "GetEquivalencesByTargetSubject",
			Handler:    _Query_GetEquivalencesByTargetSubject_Handler,
		},
		{
			MethodName: "GetEquivalencesByInstitution",
			Handler:    _Query_GetEquivalencesByInstitution_Handler,
		},
		{
			MethodName: "CheckEquivalenceStatus",
			Handler:    _Query_CheckEquivalenceStatus_Handler,
		},
		{
			MethodName: "GetPendingAnalysis",
			Handler:    _Query_GetPendingAnalysis_Handler,
		},
		{
			MethodName: "GetApprovedEquivalences",
			Handler:    _Query_GetApprovedEquivalences_Handler,
		},
		{
			MethodName: "GetRejectedEquivalences",
			Handler:    _Query_GetRejectedEquivalences_Handler,
		},
		{
			MethodName: "GetEquivalencesByContract",
			Handler:    _Query_GetEquivalencesByContract_Handler,
		},
		{
			MethodName: "GetEquivalencesByContractVersion",
			Handler:    _Query_GetEquivalencesByContractVersion_Handler,
		},
		{
			MethodName: "GetEquivalenceHistory",
			Handler:    _Query_GetEquivalenceHistory_Handler,
		},
		{
			MethodName: "GetEquivalenceStats",
			Handler:    _Query_GetEquivalenceStats_Handler,
		},
		{
			MethodName: "GetAnalysisMetadata",
			Handler:    _Query_GetAnalysisMetadata_Handler,
		},
		{
			MethodName: "VerifyAnalysisIntegrity",
			Handler:    _Query_VerifyAnalysisIntegrity_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "academictoken/equivalence/query.proto",
}
