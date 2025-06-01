package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"academictoken/x/subject/types"
)

// QueryServer represents the gRPC query service for the subject module
type QueryServer struct {
	Keeper
}

// NewQueryServerImpl returns an implementation of the QueryServer interface
func NewQueryServerImpl(keeper Keeper) types.QueryServer {
	return &QueryServer{Keeper: keeper}
}

var _ types.QueryServer = &QueryServer{}

// Params implements the Query/Params gRPC method
func (q QueryServer) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	params := q.Keeper.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

// GetSubject implements the Query/GetSubject gRPC method
func (q QueryServer) GetSubject(c context.Context, req *types.QueryGetSubjectRequest) (*types.QueryGetSubjectResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	subject, found := q.Keeper.GetSubject(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.NotFound, "subject not found")
	}

	return &types.QueryGetSubjectResponse{Subject: subject}, nil
}

// GetSubjectFull implements the Query/GetSubjectFull gRPC method
func (q QueryServer) GetSubjectFull(c context.Context, req *types.QueryGetSubjectFullRequest) (*types.QueryGetSubjectFullResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	subject, extendedContent, err := q.Keeper.GetSubjectFull(ctx, req.Index)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert extended content to JSON
	extendedContentJSON, err := json.Marshal(extendedContent)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("error serializing extended content: %v", err))
	}

	return &types.QueryGetSubjectFullResponse{
		Subject:             subject,
		ExtendedContentJson: string(extendedContentJSON),
	}, nil
}

// GetSubjectWithPrerequisites implements the Query/GetSubjectWithPrerequisites gRPC method
func (q QueryServer) GetSubjectWithPrerequisites(c context.Context, req *types.QueryGetSubjectWithPrerequisitesRequest) (*types.QueryGetSubjectWithPrerequisitesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	subjectWithPrereqs, err := q.Keeper.GetSubjectWithPrerequisites(ctx, req.SubjectId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Create SubjectWithPrerequisites structure for the response
	var protoSubjectWithPrereqs types.SubjectWithPrerequisites

	// Copy the subject - mantém como ponteiro
	protoSubjectWithPrereqs.Subject = subjectWithPrereqs.Subject

	// Copy prerequisite groups - mantém como slice de ponteiros
	protoSubjectWithPrereqs.PrerequisiteGroups = subjectWithPrereqs.PrerequisiteGroups

	return &types.QueryGetSubjectWithPrerequisitesResponse{
		SubjectWithPrerequisites: protoSubjectWithPrereqs,
	}, nil
}

// ListSubjects implements the Query/ListSubjects gRPC method
func (q QueryServer) ListSubjects(c context.Context, req *types.QueryListSubjectsRequest) (*types.QueryListSubjectsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	var subjects []types.SubjectContent

	// Get store and create iterator
	store := q.Keeper.storeService.OpenKVStore(ctx)
	prefixKey := types.SubjectPrefix

	// Use Iterator with proper error handling - capture both return values
	iterator, err := store.Iterator(nil, nil)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create iterator: %v", err))
	}
	defer iterator.Close()

	// Collect all subjects first
	var allKeys [][]byte
	var allValues [][]byte

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		if len(key) > len(prefixKey) && string(key[:len(prefixKey)]) == string(prefixKey) {
			allKeys = append(allKeys, key)
			allValues = append(allValues, iterator.Value())
		}
	}

	// Apply pagination manually
	start := 0
	limit := 100 // default limit
	if req.Pagination != nil {
		if req.Pagination.Offset > 0 {
			start = int(req.Pagination.Offset)
		}
		if req.Pagination.Limit > 0 {
			limit = int(req.Pagination.Limit)
		}
	}

	end := start + limit
	if end > len(allKeys) {
		end = len(allKeys)
	}

	// Process paginated results
	for i := start; i < end; i++ {
		var subject types.SubjectContent
		err := q.Keeper.cdc.Unmarshal(allValues[i], &subject)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		subjects = append(subjects, subject)
	}

	// Create pagination response
	var nextKey []byte
	if end < len(allKeys) {
		nextKey = allKeys[end]
	}

	pageRes := &query.PageResponse{
		NextKey: nextKey,
		Total:   uint64(len(allKeys)),
	}

	return &types.QueryListSubjectsResponse{
		Subjects:   subjects,
		Pagination: pageRes,
	}, nil
}

// CheckPrerequisites implements the Query/CheckPrerequisites gRPC method
func (q QueryServer) CheckPrerequisites(c context.Context, req *types.QueryCheckPrerequisitesRequest) (*types.QueryCheckPrerequisitesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	isEligible, missingPrereqs, err := q.Keeper.CheckPrerequisitesViaContract(ctx, req.StudentId, req.SubjectId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCheckPrerequisitesResponse{
		IsEligible:           isEligible,
		MissingPrerequisites: missingPrereqs,
	}, nil
}

// CheckEquivalence implements the Query/CheckEquivalence gRPC method
func (q QueryServer) CheckEquivalence(c context.Context, req *types.QueryCheckEquivalenceRequest) (*types.QueryCheckEquivalenceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	equivalencePercent, statusStr, err := q.Keeper.CheckEquivalenceViaContract(
		ctx,
		req.SourceSubjectId,
		req.TargetSubjectId,
		req.ForceRecalculate,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCheckEquivalenceResponse{
		EquivalencePercent: equivalencePercent,
		Status:             statusStr,
	}, nil
}

// SubjectsByCourse implements the Query/SubjectsByCourse gRPC method
func (q QueryServer) SubjectsByCourse(c context.Context, req *types.QuerySubjectsByCourseRequest) (*types.QuerySubjectsByCourseResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	subjects, err := q.Keeper.GetSubjectsByCourse(ctx, req.CourseId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Apply pagination manually if needed
	start := 0
	limit := len(subjects)
	if req.Pagination != nil {
		if req.Pagination.Offset > 0 {
			start = int(req.Pagination.Offset)
		}
		if req.Pagination.Limit > 0 && int(req.Pagination.Limit) < len(subjects)-start {
			limit = start + int(req.Pagination.Limit)
		} else {
			limit = len(subjects)
		}
	}

	if start > len(subjects) {
		start = len(subjects)
	}
	if limit > len(subjects) {
		limit = len(subjects)
	}

	paginatedSubjects := subjects[start:limit]

	// Create pagination response
	var nextKey []byte
	if limit < len(subjects) {
		nextKey = []byte(fmt.Sprintf("%d", limit))
	}

	pageRes := &query.PageResponse{
		NextKey: nextKey,
		Total:   uint64(len(subjects)),
	}

	return &types.QuerySubjectsByCourseResponse{
		Subjects:   paginatedSubjects,
		Pagination: pageRes,
	}, nil
}

// SubjectsByInstitution implements the Query/SubjectsByInstitution gRPC method
func (q QueryServer) SubjectsByInstitution(c context.Context, req *types.QuerySubjectsByInstitutionRequest) (*types.QuerySubjectsByInstitutionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	subjects, err := q.Keeper.GetSubjectsByInstitution(ctx, req.InstitutionId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Apply pagination manually if needed
	start := 0
	limit := len(subjects)
	if req.Pagination != nil {
		if req.Pagination.Offset > 0 {
			start = int(req.Pagination.Offset)
		}
		if req.Pagination.Limit > 0 && int(req.Pagination.Limit) < len(subjects)-start {
			limit = start + int(req.Pagination.Limit)
		} else {
			limit = len(subjects)
		}
	}

	if start > len(subjects) {
		start = len(subjects)
	}
	if limit > len(subjects) {
		limit = len(subjects)
	}

	paginatedSubjects := subjects[start:limit]

	// Create pagination response
	var nextKey []byte
	if limit < len(subjects) {
		nextKey = []byte(fmt.Sprintf("%d", limit))
	}

	pageRes := &query.PageResponse{
		NextKey: nextKey,
		Total:   uint64(len(subjects)),
	}

	return &types.QuerySubjectsByInstitutionResponse{
		Subjects:   paginatedSubjects,
		Pagination: pageRes,
	}, nil
}
