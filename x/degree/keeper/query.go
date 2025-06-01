package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"academictoken/x/degree/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Degree(goCtx context.Context, req *types.QueryGetDegreeRequest) (*types.QueryGetDegreeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetDegree(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetDegreeResponse{Degree: val}, nil
}

func (k Keeper) DegreeAll(goCtx context.Context, req *types.QueryAllDegreeRequest) (*types.QueryAllDegreeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	var degrees []types.Degree

	store := ctx.KVStore(k.storeKey)
	degreeStore := prefix.NewStore(store, types.DegreePrefix)

	pageRes, err := query.Paginate(degreeStore, req.Pagination, func(key []byte, value []byte) error {
		var degree types.Degree
		if err := k.cdc.Unmarshal(value, &degree); err != nil {
			return err
		}

		degrees = append(degrees, degree)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllDegreeResponse{Degree: degrees, Pagination: pageRes}, nil
}

func (k Keeper) DegreesByStudent(goCtx context.Context, req *types.QueryDegreesByStudentRequest) (*types.QueryDegreesByStudentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	degrees := k.GetDegreesByStudent(ctx, req.StudentId)

	return &types.QueryDegreesByStudentResponse{Degrees: degrees}, nil
}

func (k Keeper) DegreesByInstitution(goCtx context.Context, req *types.QueryDegreesByInstitutionRequest) (*types.QueryDegreesByInstitutionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	degrees := k.GetDegreesByInstitution(ctx, req.InstitutionId)

	return &types.QueryDegreesByInstitutionResponse{Degrees: degrees}, nil
}

func (k Keeper) DegreeRequests(goCtx context.Context, req *types.QueryDegreeRequestsRequest) (*types.QueryDegreeRequestsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var requests []types.DegreeRequest

	if req.Status != "" {
		requests = k.GetDegreeRequestsByStatus(ctx, req.Status)
	} else {
		requests = k.GetAllDegreeRequest(ctx)
	}

	return &types.QueryDegreeRequestsResponse{Requests: requests}, nil
}

func (k Keeper) DegreeValidationStatus(goCtx context.Context, req *types.QueryDegreeValidationStatusRequest) (*types.QueryDegreeValidationStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get degree by ID
	degree, found := k.GetDegree(ctx, req.DegreeId)
	if !found {
		return nil, status.Error(codes.NotFound, "degree not found")
	}

	return &types.QueryDegreeValidationStatusResponse{
		Status:            degree.Status,
		ValidationScore:   degree.ValidationScore,
		ValidationDate:    degree.IssueDate,
		ValidationDetails: "Validation completed",
	}, nil
}
