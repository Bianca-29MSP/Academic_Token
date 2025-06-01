package keeper

import (
	"context"
	"fmt"

	"academictoken/x/schedule/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// CASCA: Delega criação de recomendação para módulo simples (sem contrato por enquanto)
func (k msgServer) CreateSubjectRecommendation(goCtx context.Context, req *types.MsgCreateSubjectRecommendation) (*types.MsgCreateSubjectRecommendationResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Basic validation only
	if req.Creator == "" || req.Student == "" {
		return nil, fmt.Errorf("required fields cannot be empty")
	}

	// Use unique variable name to avoid conflicts
	recId, err := k.Keeper.CreateSubjectRecommendation(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create subject recommendation: %w", err)
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSubjectRecommendationCreated,
			sdk.NewAttribute(types.AttributeKeyStudentID, req.Student),
			sdk.NewAttribute("recommendation_id", recId), // recId is string
			sdk.NewAttribute("semester", req.RecommendationSemester),
		),
	)

	k.Logger().Info("Subject recommendation created",
		"recommendation_id", recId,
		"student_id", req.Student,
		"creator", req.Creator,
	)

	return &types.MsgCreateSubjectRecommendationResponse{}, nil
}

// CASCA: Delega criação de plano de estudos para módulo simples (sem contrato por enquanto)
func (k msgServer) CreateStudyPlan(goCtx context.Context, req *types.MsgCreateStudyPlan) (*types.MsgCreateStudyPlanResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Basic validation only
	if req.Creator == "" || req.Student == "" {
		return nil, fmt.Errorf("required fields cannot be empty")
	}

	// Check max study plans per student
	existingPlans := k.GetStudyPlansByStudent(ctx, req.Student)
	if uint64(len(existingPlans)) >= k.GetMaxStudyPlansPerStudent(ctx) {
		return nil, fmt.Errorf("student has reached maximum number of study plans (%d)", k.GetMaxStudyPlansPerStudent(ctx))
	}

	// Use unique variable name to avoid conflicts
	planId, err := k.Keeper.CreateStudyPlan(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create study plan: %w", err)
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeStudyPlanCreated,
			sdk.NewAttribute(types.AttributeKeyStudentID, req.Student),
			sdk.NewAttribute(types.AttributeKeyStudyPlanID, planId), // planId is string
			sdk.NewAttribute("completion_target", req.CompletionTarget),
		),
	)

	k.Logger().Info("Study plan created",
		"study_plan_id", planId,
		"student_id", req.Student,
		"creator", req.Creator,
	)

	return &types.MsgCreateStudyPlanResponse{}, nil
}

// CASCA: Delega adição de semestre planejado para módulo simples (sem contrato por enquanto)
func (k msgServer) AddPlannedSemester(goCtx context.Context, req *types.MsgAddPlannedSemester) (*types.MsgAddPlannedSemesterResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Basic validation only
	if req.Creator == "" || req.StudyPlanId == "" || req.SemesterCode == "" {
		return nil, fmt.Errorf("required fields cannot be empty")
	}

	// Check if study plan exists
	_, found := k.GetStudyPlan(ctx, req.StudyPlanId)
	if !found {
		return nil, fmt.Errorf("study plan with ID '%s' not found", req.StudyPlanId)
	}

	// Validate max credits per semester
	if req.TotalCredits > k.GetMaxCreditsPerSemester(ctx) {
		return nil, fmt.Errorf("total credits (%d) exceeds maximum allowed per semester (%d)",
			req.TotalCredits, k.GetMaxCreditsPerSemester(ctx))
	}

	// Call keeper method directly - AddPlannedSemester returns only error
	err := k.Keeper.AddPlannedSemester(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to add planned semester: %w", err)
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePlannedSemesterAdded,
			sdk.NewAttribute(types.AttributeKeyStudyPlanID, req.StudyPlanId),
			sdk.NewAttribute(types.AttributeKeySemesterNumber, req.SemesterCode),
			sdk.NewAttribute("total_credits", fmt.Sprintf("%d", req.TotalCredits)),
		),
	)

	k.Logger().Info("Planned semester added",
		"study_plan_id", req.StudyPlanId,
		"semester_code", req.SemesterCode,
		"total_credits", req.TotalCredits,
		"creator", req.Creator,
	)

	return &types.MsgAddPlannedSemesterResponse{}, nil
}

// CASCA: Delega atualização de status para módulo simples (sem contrato por enquanto)
func (k msgServer) UpdateStudyPlanStatus(goCtx context.Context, req *types.MsgUpdateStudyPlanStatus) (*types.MsgUpdateStudyPlanStatusResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Basic validation only
	if req.Creator == "" || req.StudyPlanId == "" || req.Status == "" {
		return nil, fmt.Errorf("required fields cannot be empty")
	}

	// Validate status
	validStatuses := []string{
		types.StudyPlanStatusDraft,
		types.StudyPlanStatusActive,
		types.StudyPlanStatusCompleted,
		types.StudyPlanStatusArchived,
	}

	validStatus := false
	for _, status := range validStatuses {
		if req.Status == status {
			validStatus = true
			break
		}
	}
	if !validStatus {
		return nil, fmt.Errorf("invalid study plan status: %s", req.Status)
	}

	// Check if study plan exists
	_, found := k.GetStudyPlan(ctx, req.StudyPlanId)
	if !found {
		return nil, fmt.Errorf("study plan with ID '%s' not found", req.StudyPlanId)
	}

	// Call keeper method directly - UpdateStudyPlanStatus returns only error
	err := k.Keeper.UpdateStudyPlanStatus(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update study plan status: %w", err)
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeStudyPlanStatusChanged,
			sdk.NewAttribute(types.AttributeKeyStudyPlanID, req.StudyPlanId),
			sdk.NewAttribute(types.AttributeKeyStatus, req.Status),
		),
	)

	k.Logger().Info("Study plan status updated",
		"study_plan_id", req.StudyPlanId,
		"new_status", req.Status,
		"creator", req.Creator,
	)

	return &types.MsgUpdateStudyPlanStatusResponse{}, nil
}

// UpdateParams is deprecated with hardcoded params
func (k msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	return nil, fmt.Errorf("UpdateParams is deprecated - schedule configuration is hardcoded")
}
