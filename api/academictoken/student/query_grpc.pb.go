// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: academictoken/student/query.proto

package student

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
	Query_Params_FullMethodName                     = "/academictoken.student.Query/Params"
	Query_ListStudents_FullMethodName               = "/academictoken.student.Query/ListStudents"
	Query_GetStudent_FullMethodName                 = "/academictoken.student.Query/GetStudent"
	Query_ListEnrollments_FullMethodName            = "/academictoken.student.Query/ListEnrollments"
	Query_GetEnrollment_FullMethodName              = "/academictoken.student.Query/GetEnrollment"
	Query_GetEnrollmentsByStudent_FullMethodName    = "/academictoken.student.Query/GetEnrollmentsByStudent"
	Query_GetStudentProgress_FullMethodName         = "/academictoken.student.Query/GetStudentProgress"
	Query_GetStudentsByInstitution_FullMethodName   = "/academictoken.student.Query/GetStudentsByInstitution"
	Query_GetStudentsByCourse_FullMethodName        = "/academictoken.student.Query/GetStudentsByCourse"
	Query_GetStudentAcademicTree_FullMethodName     = "/academictoken.student.Query/GetStudentAcademicTree"
	Query_CheckGraduationEligibility_FullMethodName = "/academictoken.student.Query/CheckGraduationEligibility"
)

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueryClient interface {
	// Parameters queries the parameters of the module.
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
	// ListStudents queries all students with pagination
	ListStudents(ctx context.Context, in *QueryListStudentsRequest, opts ...grpc.CallOption) (*QueryListStudentsResponse, error)
	// GetStudent queries a student by index
	GetStudent(ctx context.Context, in *QueryGetStudentRequest, opts ...grpc.CallOption) (*QueryGetStudentResponse, error)
	// ListEnrollments queries all enrollments with pagination
	ListEnrollments(ctx context.Context, in *QueryListEnrollmentsRequest, opts ...grpc.CallOption) (*QueryListEnrollmentsResponse, error)
	// GetEnrollment queries an enrollment by index
	GetEnrollment(ctx context.Context, in *QueryGetEnrollmentRequest, opts ...grpc.CallOption) (*QueryGetEnrollmentResponse, error)
	// GetEnrollmentsByStudent queries all enrollments for a specific student
	GetEnrollmentsByStudent(ctx context.Context, in *QueryGetEnrollmentsByStudentRequest, opts ...grpc.CallOption) (*QueryGetEnrollmentsByStudentResponse, error)
	// GetStudentProgress queries academic progress for a specific student
	GetStudentProgress(ctx context.Context, in *QueryGetStudentProgressRequest, opts ...grpc.CallOption) (*QueryGetStudentProgressResponse, error)
	// GetStudentsByInstitution queries students enrolled in a specific institution
	GetStudentsByInstitution(ctx context.Context, in *QueryGetStudentsByInstitutionRequest, opts ...grpc.CallOption) (*QueryGetStudentsByInstitutionResponse, error)
	// GetStudentsByCourse queries students enrolled in a specific course
	GetStudentsByCourse(ctx context.Context, in *QueryGetStudentsByCourseRequest, opts ...grpc.CallOption) (*QueryGetStudentsByCourseResponse, error)
	// GetStudentAcademicTree queries the academic tree for a specific student
	GetStudentAcademicTree(ctx context.Context, in *QueryGetStudentAcademicTreeRequest, opts ...grpc.CallOption) (*QueryGetStudentAcademicTreeResponse, error)
	// CheckGraduationEligibility checks if a student is eligible for graduation
	CheckGraduationEligibility(ctx context.Context, in *QueryCheckGraduationEligibilityRequest, opts ...grpc.CallOption) (*QueryCheckGraduationEligibilityResponse, error)
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

func (c *queryClient) ListStudents(ctx context.Context, in *QueryListStudentsRequest, opts ...grpc.CallOption) (*QueryListStudentsResponse, error) {
	out := new(QueryListStudentsResponse)
	err := c.cc.Invoke(ctx, Query_ListStudents_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetStudent(ctx context.Context, in *QueryGetStudentRequest, opts ...grpc.CallOption) (*QueryGetStudentResponse, error) {
	out := new(QueryGetStudentResponse)
	err := c.cc.Invoke(ctx, Query_GetStudent_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ListEnrollments(ctx context.Context, in *QueryListEnrollmentsRequest, opts ...grpc.CallOption) (*QueryListEnrollmentsResponse, error) {
	out := new(QueryListEnrollmentsResponse)
	err := c.cc.Invoke(ctx, Query_ListEnrollments_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetEnrollment(ctx context.Context, in *QueryGetEnrollmentRequest, opts ...grpc.CallOption) (*QueryGetEnrollmentResponse, error) {
	out := new(QueryGetEnrollmentResponse)
	err := c.cc.Invoke(ctx, Query_GetEnrollment_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetEnrollmentsByStudent(ctx context.Context, in *QueryGetEnrollmentsByStudentRequest, opts ...grpc.CallOption) (*QueryGetEnrollmentsByStudentResponse, error) {
	out := new(QueryGetEnrollmentsByStudentResponse)
	err := c.cc.Invoke(ctx, Query_GetEnrollmentsByStudent_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetStudentProgress(ctx context.Context, in *QueryGetStudentProgressRequest, opts ...grpc.CallOption) (*QueryGetStudentProgressResponse, error) {
	out := new(QueryGetStudentProgressResponse)
	err := c.cc.Invoke(ctx, Query_GetStudentProgress_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetStudentsByInstitution(ctx context.Context, in *QueryGetStudentsByInstitutionRequest, opts ...grpc.CallOption) (*QueryGetStudentsByInstitutionResponse, error) {
	out := new(QueryGetStudentsByInstitutionResponse)
	err := c.cc.Invoke(ctx, Query_GetStudentsByInstitution_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetStudentsByCourse(ctx context.Context, in *QueryGetStudentsByCourseRequest, opts ...grpc.CallOption) (*QueryGetStudentsByCourseResponse, error) {
	out := new(QueryGetStudentsByCourseResponse)
	err := c.cc.Invoke(ctx, Query_GetStudentsByCourse_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetStudentAcademicTree(ctx context.Context, in *QueryGetStudentAcademicTreeRequest, opts ...grpc.CallOption) (*QueryGetStudentAcademicTreeResponse, error) {
	out := new(QueryGetStudentAcademicTreeResponse)
	err := c.cc.Invoke(ctx, Query_GetStudentAcademicTree_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) CheckGraduationEligibility(ctx context.Context, in *QueryCheckGraduationEligibilityRequest, opts ...grpc.CallOption) (*QueryCheckGraduationEligibilityResponse, error) {
	out := new(QueryCheckGraduationEligibilityResponse)
	err := c.cc.Invoke(ctx, Query_CheckGraduationEligibility_FullMethodName, in, out, opts...)
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
	// ListStudents queries all students with pagination
	ListStudents(context.Context, *QueryListStudentsRequest) (*QueryListStudentsResponse, error)
	// GetStudent queries a student by index
	GetStudent(context.Context, *QueryGetStudentRequest) (*QueryGetStudentResponse, error)
	// ListEnrollments queries all enrollments with pagination
	ListEnrollments(context.Context, *QueryListEnrollmentsRequest) (*QueryListEnrollmentsResponse, error)
	// GetEnrollment queries an enrollment by index
	GetEnrollment(context.Context, *QueryGetEnrollmentRequest) (*QueryGetEnrollmentResponse, error)
	// GetEnrollmentsByStudent queries all enrollments for a specific student
	GetEnrollmentsByStudent(context.Context, *QueryGetEnrollmentsByStudentRequest) (*QueryGetEnrollmentsByStudentResponse, error)
	// GetStudentProgress queries academic progress for a specific student
	GetStudentProgress(context.Context, *QueryGetStudentProgressRequest) (*QueryGetStudentProgressResponse, error)
	// GetStudentsByInstitution queries students enrolled in a specific institution
	GetStudentsByInstitution(context.Context, *QueryGetStudentsByInstitutionRequest) (*QueryGetStudentsByInstitutionResponse, error)
	// GetStudentsByCourse queries students enrolled in a specific course
	GetStudentsByCourse(context.Context, *QueryGetStudentsByCourseRequest) (*QueryGetStudentsByCourseResponse, error)
	// GetStudentAcademicTree queries the academic tree for a specific student
	GetStudentAcademicTree(context.Context, *QueryGetStudentAcademicTreeRequest) (*QueryGetStudentAcademicTreeResponse, error)
	// CheckGraduationEligibility checks if a student is eligible for graduation
	CheckGraduationEligibility(context.Context, *QueryCheckGraduationEligibilityRequest) (*QueryCheckGraduationEligibilityResponse, error)
	mustEmbedUnimplementedQueryServer()
}

// UnimplementedQueryServer must be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (UnimplementedQueryServer) Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (UnimplementedQueryServer) ListStudents(context.Context, *QueryListStudentsRequest) (*QueryListStudentsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListStudents not implemented")
}
func (UnimplementedQueryServer) GetStudent(context.Context, *QueryGetStudentRequest) (*QueryGetStudentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStudent not implemented")
}
func (UnimplementedQueryServer) ListEnrollments(context.Context, *QueryListEnrollmentsRequest) (*QueryListEnrollmentsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEnrollments not implemented")
}
func (UnimplementedQueryServer) GetEnrollment(context.Context, *QueryGetEnrollmentRequest) (*QueryGetEnrollmentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEnrollment not implemented")
}
func (UnimplementedQueryServer) GetEnrollmentsByStudent(context.Context, *QueryGetEnrollmentsByStudentRequest) (*QueryGetEnrollmentsByStudentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEnrollmentsByStudent not implemented")
}
func (UnimplementedQueryServer) GetStudentProgress(context.Context, *QueryGetStudentProgressRequest) (*QueryGetStudentProgressResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStudentProgress not implemented")
}
func (UnimplementedQueryServer) GetStudentsByInstitution(context.Context, *QueryGetStudentsByInstitutionRequest) (*QueryGetStudentsByInstitutionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStudentsByInstitution not implemented")
}
func (UnimplementedQueryServer) GetStudentsByCourse(context.Context, *QueryGetStudentsByCourseRequest) (*QueryGetStudentsByCourseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStudentsByCourse not implemented")
}
func (UnimplementedQueryServer) GetStudentAcademicTree(context.Context, *QueryGetStudentAcademicTreeRequest) (*QueryGetStudentAcademicTreeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStudentAcademicTree not implemented")
}
func (UnimplementedQueryServer) CheckGraduationEligibility(context.Context, *QueryCheckGraduationEligibilityRequest) (*QueryCheckGraduationEligibilityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckGraduationEligibility not implemented")
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

func _Query_ListStudents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryListStudentsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ListStudents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_ListStudents_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ListStudents(ctx, req.(*QueryListStudentsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetStudent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetStudentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetStudent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetStudent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetStudent(ctx, req.(*QueryGetStudentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ListEnrollments_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryListEnrollmentsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ListEnrollments(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_ListEnrollments_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ListEnrollments(ctx, req.(*QueryListEnrollmentsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetEnrollment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetEnrollmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetEnrollment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetEnrollment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetEnrollment(ctx, req.(*QueryGetEnrollmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetEnrollmentsByStudent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetEnrollmentsByStudentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetEnrollmentsByStudent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetEnrollmentsByStudent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetEnrollmentsByStudent(ctx, req.(*QueryGetEnrollmentsByStudentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetStudentProgress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetStudentProgressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetStudentProgress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetStudentProgress_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetStudentProgress(ctx, req.(*QueryGetStudentProgressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetStudentsByInstitution_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetStudentsByInstitutionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetStudentsByInstitution(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetStudentsByInstitution_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetStudentsByInstitution(ctx, req.(*QueryGetStudentsByInstitutionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetStudentsByCourse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetStudentsByCourseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetStudentsByCourse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetStudentsByCourse_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetStudentsByCourse(ctx, req.(*QueryGetStudentsByCourseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetStudentAcademicTree_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetStudentAcademicTreeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetStudentAcademicTree(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_GetStudentAcademicTree_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetStudentAcademicTree(ctx, req.(*QueryGetStudentAcademicTreeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_CheckGraduationEligibility_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryCheckGraduationEligibilityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).CheckGraduationEligibility(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_CheckGraduationEligibility_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).CheckGraduationEligibility(ctx, req.(*QueryCheckGraduationEligibilityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Query_ServiceDesc is the grpc.ServiceDesc for Query service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "academictoken.student.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "ListStudents",
			Handler:    _Query_ListStudents_Handler,
		},
		{
			MethodName: "GetStudent",
			Handler:    _Query_GetStudent_Handler,
		},
		{
			MethodName: "ListEnrollments",
			Handler:    _Query_ListEnrollments_Handler,
		},
		{
			MethodName: "GetEnrollment",
			Handler:    _Query_GetEnrollment_Handler,
		},
		{
			MethodName: "GetEnrollmentsByStudent",
			Handler:    _Query_GetEnrollmentsByStudent_Handler,
		},
		{
			MethodName: "GetStudentProgress",
			Handler:    _Query_GetStudentProgress_Handler,
		},
		{
			MethodName: "GetStudentsByInstitution",
			Handler:    _Query_GetStudentsByInstitution_Handler,
		},
		{
			MethodName: "GetStudentsByCourse",
			Handler:    _Query_GetStudentsByCourse_Handler,
		},
		{
			MethodName: "GetStudentAcademicTree",
			Handler:    _Query_GetStudentAcademicTree_Handler,
		},
		{
			MethodName: "CheckGraduationEligibility",
			Handler:    _Query_CheckGraduationEligibility_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "academictoken/student/query.proto",
}
