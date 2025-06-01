package keeper

import (
	"context"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"academictoken/x/course/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// UpdateParams updates the module parameters
func (k msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the message signer is the authority
	if req.Authority != k.GetAuthority() {
		return nil, fmt.Errorf("invalid authority: expected %s, got %s", k.GetAuthority(), req.Authority)
	}

	// Validate the new params
	if err := req.Params.Validate(); err != nil {
		return nil, fmt.Errorf("invalid params: %w", err)
	}

	// Set the new params
	k.SetParams(ctx, req.Params)

	return &types.MsgUpdateParamsResponse{}, nil
}

// CreateCourse creates a new course
func (k msgServer) CreateCourse(goCtx context.Context, req *types.MsgCreateCourse) (*types.MsgCreateCourseResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate basic fields
	if req.Creator == "" {
		return nil, fmt.Errorf("creator cannot be empty")
	}
	if req.Institution == "" {
		return nil, fmt.Errorf("institution cannot be empty")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("course name cannot be empty")
	}
	if req.Code == "" {
		return nil, fmt.Errorf("course code cannot be empty")
	}
	if req.TotalCredits == "" {
		return nil, fmt.Errorf("total credits cannot be empty")
	}
	if req.DegreeLevel == "" {
		return nil, fmt.Errorf("degree level cannot be empty")
	}

	// Validate total credits is a valid number
	_, err := strconv.ParseUint(req.TotalCredits, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid total credits format: %w", err)
	}

	// Validate degree level
	validDegreeLevels := []string{"undergraduate", "graduate", "postgraduate", "doctorate", "technical"}
	validDegreeLevel := false
	for _, level := range validDegreeLevels {
		if req.DegreeLevel == level {
			validDegreeLevel = true
			break
		}
	}
	if !validDegreeLevel {
		return nil, fmt.Errorf("invalid degree level: must be one of %v, got '%s'", validDegreeLevels, req.DegreeLevel)
	}

	// Check if institution exists (assuming we have access to institution keeper)
	if k.institutionKeeper != nil {
		_, found := k.institutionKeeper.GetInstitution(ctx, req.Institution)
		if !found {
			return nil, fmt.Errorf("institution with ID '%s' not found", req.Institution)
		}
	}

	// Check if course with same code already exists for this institution
	allCourses := k.GetAllCourse(ctx)
	for _, course := range allCourses {
		if course.Institution == req.Institution && course.Code == req.Code {
			return nil, fmt.Errorf("course with code '%s' already exists in institution '%s'", req.Code, req.Institution)
		}
	}

	// Generate a new course index
	courseIndex := k.GetNextCourseIndex(ctx)

	// Create course object
	course := types.Course{
		Index:        courseIndex,
		Institution:  req.Institution,
		Name:         req.Name,
		Code:         req.Code,
		Description:  req.Description,
		TotalCredits: req.TotalCredits,
		DegreeLevel:  req.DegreeLevel,
	}

	// Store the course
	err = k.SetCourse(ctx, course)
	if err != nil {
		return nil, fmt.Errorf("failed to store course: %w", err)
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"course_created",
			sdk.NewAttribute("course_index", courseIndex),
			sdk.NewAttribute("institution", req.Institution),
			sdk.NewAttribute("name", req.Name),
			sdk.NewAttribute("code", req.Code),
			sdk.NewAttribute("degree_level", req.DegreeLevel),
			sdk.NewAttribute("total_credits", req.TotalCredits),
			sdk.NewAttribute("creator", req.Creator),
		),
	)

	k.Logger().Info("Course created successfully",
		"course_index", courseIndex,
		"institution", req.Institution,
		"name", req.Name,
		"code", req.Code,
		"creator", req.Creator,
	)

	return &types.MsgCreateCourseResponse{}, nil
}

// UpdateCourse updates an existing course
func (k msgServer) UpdateCourse(goCtx context.Context, req *types.MsgUpdateCourse) (*types.MsgUpdateCourseResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate basic fields
	if req.Creator == "" {
		return nil, fmt.Errorf("creator cannot be empty")
	}
	if req.Index == "" {
		return nil, fmt.Errorf("course index cannot be empty")
	}

	// Get existing course
	course, found := k.GetCourse(ctx, req.Index)
	if !found {
		return nil, fmt.Errorf("course with index '%s' not found", req.Index)
	}

	// Check permissions using the CanUpdateCourse method
	if !k.CanUpdateCourse(ctx, req.Index, req.Creator) {
		return nil, fmt.Errorf("creator '%s' is not authorized to update course '%s'", req.Creator, req.Index)
	}

	// Update fields if provided
	updated := false

	if req.Name != "" && req.Name != course.Name {
		course.Name = req.Name
		updated = true
	}

	if req.Description != "" && req.Description != course.Description {
		course.Description = req.Description
		updated = true
	}

	if req.TotalCredits != "" && req.TotalCredits != course.TotalCredits {
		// Validate that total credits is a valid number
		_, err := strconv.ParseUint(req.TotalCredits, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid total credits format: %w", err)
		}
		course.TotalCredits = req.TotalCredits
		updated = true
	}

	if !updated {
		return nil, fmt.Errorf("no valid updates provided")
	}

	// Store updated course
	k.SetCourse(ctx, course)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"course_updated",
			sdk.NewAttribute("course_index", req.Index),
			sdk.NewAttribute("name", course.Name),
			sdk.NewAttribute("total_credits", course.TotalCredits),
			sdk.NewAttribute("updater", req.Creator),
		),
	)

	k.Logger().Info("Course updated successfully",
		"course_index", req.Index,
		"name", course.Name,
		"updater", req.Creator,
	)

	return &types.MsgUpdateCourseResponse{}, nil
}
