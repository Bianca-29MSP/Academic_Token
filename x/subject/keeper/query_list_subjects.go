package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"academictoken/x/subject/types"
)

func (k Keeper) ListSubjects(goCtx context.Context, req *types.QueryListSubjectsRequest) (*types.QueryListSubjectsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	subjects, err := k.GetAllSubjects(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryListSubjectsResponse{
		Subjects: subjects,
	}, nil
}
