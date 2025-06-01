package keeper

import (
	"context"

	"academictoken/x/schedule/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

// Params returns the module parameters
func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	params := k.GetParams(c)

	return &types.QueryParamsResponse{Params: params}, nil
}

// StudyPlan returns a specific study plan by ID
func (k Keeper) StudyPlan(c context.Context, req *types.QueryGetStudyPlanRequest) (*types.QueryGetStudyPlanResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.StudyPlanId == "" {
		return nil, status.Error(codes.InvalidArgument, "study plan ID cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	studyPlan, found := k.GetStudyPlan(ctx, req.StudyPlanId)
	if !found {
		return nil, status.Error(codes.NotFound, "study plan not found")
	}

	return &types.QueryGetStudyPlanResponse{StudyPlan: studyPlan}, nil
}

// StudyPlanAll returns all study plans with pagination
func (k Keeper) StudyPlanAll(c context.Context, req *types.QueryAllStudyPlanRequest) (*types.QueryAllStudyPlanResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	studyPlans, pageRes, err := k.GetAllStudyPlans(ctx, req.Pagination)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllStudyPlanResponse{
		StudyPlan:  studyPlans,
		Pagination: pageRes,
	}, nil
}

// StudyPlansByStudent returns all study plans for a specific student
func (k Keeper) StudyPlansByStudent(c context.Context, req *types.QueryStudyPlansByStudentRequest) (*types.QueryStudyPlansByStudentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.StudentId == "" {
		return nil, status.Error(codes.InvalidArgument, "student ID cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	studyPlans := k.GetStudyPlansByStudent(ctx, req.StudentId)

	return &types.QueryStudyPlansByStudentResponse{
		StudyPlans: studyPlans,
		StudentId:  req.StudentId,
	}, nil
}

// PlannedSemester returns a specific planned semester by ID
func (k Keeper) PlannedSemester(c context.Context, req *types.QueryGetPlannedSemesterRequest) (*types.QueryGetPlannedSemesterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.PlannedSemesterId == "" {
		return nil, status.Error(codes.InvalidArgument, "planned semester ID cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	plannedSemester, found := k.GetPlannedSemester(ctx, req.PlannedSemesterId)
	if !found {
		return nil, status.Error(codes.NotFound, "planned semester not found")
	}

	return &types.QueryGetPlannedSemesterResponse{PlannedSemester: plannedSemester}, nil
}

// SubjectRecommendation returns a specific subject recommendation by ID
func (k Keeper) SubjectRecommendation(c context.Context, req *types.QueryGetSubjectRecommendationRequest) (*types.QueryGetSubjectRecommendationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.RecommendationId == "" {
		return nil, status.Error(codes.InvalidArgument, "recommendation ID cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	recommendation, found := k.GetSubjectRecommendation(ctx, req.RecommendationId)
	if !found {
		return nil, status.Error(codes.NotFound, "subject recommendation not found")
	}

	return &types.QueryGetSubjectRecommendationResponse{SubjectRecommendation: recommendation}, nil
}

// SubjectRecommendationsByStudent returns all recommendations for a student
func (k Keeper) SubjectRecommendationsByStudent(c context.Context, req *types.QuerySubjectRecommendationsByStudentRequest) (*types.QuerySubjectRecommendationsByStudentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.StudentId == "" {
		return nil, status.Error(codes.InvalidArgument, "student ID cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	recommendations := k.GetSubjectRecommendationsByStudent(ctx, req.StudentId)

	return &types.QuerySubjectRecommendationsByStudentResponse{
		SubjectRecommendations: recommendations,
		StudentId:              req.StudentId,
	}, nil
}

// GenerateRecommendations generates automatic recommendations for a student
func (k Keeper) GenerateRecommendations(c context.Context, req *types.QueryGenerateRecommendationsRequest) (*types.QueryGenerateRecommendationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.StudentId == "" {
		return nil, status.Error(codes.InvalidArgument, "student ID cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	recommendations, err := k.BuildRecommendations(ctx, req.StudentId, req.SemesterCode)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGenerateRecommendationsResponse{
		Recommendations: recommendations,
		StudentId:       req.StudentId,
		SemesterCode:    req.SemesterCode,
	}, nil
}

// CheckStudentProgress checks the academic progress of a student (simplified)
func (k Keeper) CheckStudentProgress(c context.Context, req *types.QueryCheckStudentProgressRequest) (*types.QueryCheckStudentProgressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.StudentId == "" {
		return nil, status.Error(codes.InvalidArgument, "student ID cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)

	// Simplified implementation - just return basic data
	// TODO: Implement actual progress checking logic using ctx
	_ = ctx // Suppress unused variable warning
	return &types.QueryCheckStudentProgressResponse{
		StudentId:            req.StudentId,
		CourseId:             req.CourseId,
		CompletedCredits:     0,
		TotalRequiredCredits: 0,
		CompletionPercentage: 0.0,
		RemainingSemesters:   0,
		CurrentGPA:           0.0,
		AcademicTree:         types.AcademicTree{},
	}, nil
}

// OptimizeSchedule optimizes a student's schedule (simplified)
func (k Keeper) OptimizeSchedule(c context.Context, req *types.QueryOptimizeScheduleRequest) (*types.QueryOptimizeScheduleResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.StudyPlanId == "" {
		return nil, status.Error(codes.InvalidArgument, "study plan ID cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)

	// Get the study plan
	studyPlan, found := k.GetStudyPlan(ctx, req.StudyPlanId)
	if !found {
		return nil, status.Error(codes.NotFound, "study plan not found")
	}

	// Simplified optimization - just return the existing plan
	return &types.QueryOptimizeScheduleResponse{
		StudyPlanId:      req.StudyPlanId,
		OptimizedPlan:    studyPlan,
		Suggestions:      []string{"Consider taking prerequisite subjects earlier"},
		EstimatedSavings: "0 semesters",
	}, nil
}
