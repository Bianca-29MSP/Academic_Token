package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"academictoken/x/subject/types"
)

// PrerequisitesByCourse implements the Query/PrerequisitesByCourse gRPC method
// Returns all prerequisites for subjects in a specific course
func (q QueryServer) PrerequisitesByCourse(c context.Context, req *types.QueryPrerequisitesByCourseRequest) (*types.QueryPrerequisitesByCourseResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.CourseId == "" {
		return nil, status.Error(codes.InvalidArgument, "course ID cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)

	// First, get all subjects for the course
	subjects, err := q.Keeper.GetSubjectsByCourse(ctx, req.CourseId)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get subjects for course: %v", err))
	}

	// Create the response map
	prerequisites := make(map[string]*types.PrerequisiteGroups)

	// For each subject, get its prerequisite groups
	for _, subject := range subjects {
		// Get prerequisite groups for this subject
		prereqGroups := q.Keeper.GetPrerequisiteGroupsBySubject(ctx, subject.SubjectId)

		if len(prereqGroups) > 0 {
			// Convert []PrerequisiteGroup to []*PrerequisiteGroup
			groupPtrs := make([]*types.PrerequisiteGroup, len(prereqGroups))
			for i := range prereqGroups {
				groupPtrs[i] = &prereqGroups[i] // Create pointers to each element
			}

			prerequisites[subject.SubjectId] = &types.PrerequisiteGroups{
				Groups: groupPtrs, // Now using the correct pointer type
			}
		}
	}

	return &types.QueryPrerequisitesByCourseResponse{
		Prerequisites: prerequisites,
	}, nil
}
