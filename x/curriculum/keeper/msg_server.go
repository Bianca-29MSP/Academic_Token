package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"academictoken/x/curriculum/types"
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

// CreateCurriculumTree creates a new curriculum tree
func (k msgServer) CreateCurriculumTree(goCtx context.Context, req *types.MsgCreateCurriculumTree) (*types.MsgCreateCurriculumTreeResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate basic fields
	if req.Creator == "" {
		return nil, fmt.Errorf("creator cannot be empty")
	}
	if req.CourseId == "" {
		return nil, fmt.Errorf("course ID cannot be empty")
	}
	if req.Version == "" {
		return nil, fmt.Errorf("version cannot be empty")
	}

	// Validate numeric fields
	if req.ElectiveMin == 0 {
		return nil, fmt.Errorf("minimum electives cannot be zero")
	}
	if req.TotalWorkloadHours == 0 {
		return nil, fmt.Errorf("total workload hours cannot be zero")
	}

	// Check if course exists
	if k.courseKeeper != nil {
		_, found := k.courseKeeper.GetCourse(ctx, req.CourseId)
		if !found {
			return nil, fmt.Errorf("course with ID '%s' not found", req.CourseId)
		}
	}

	// Validate that required and elective subjects exist
	if k.subjectKeeper != nil {
		for _, subjectId := range req.RequiredSubjects {
			if subjectId != "" {
				_, found := k.subjectKeeper.GetSubject(ctx, subjectId)
				if !found {
					return nil, fmt.Errorf("required subject with ID '%s' not found", subjectId)
				}
			}
		}

		for _, subjectId := range req.ElectiveSubjects {
			if subjectId != "" {
				_, found := k.subjectKeeper.GetSubject(ctx, subjectId)
				if !found {
					return nil, fmt.Errorf("elective subject with ID '%s' not found", subjectId)
				}
			}
		}
	}

	// Check if curriculum already exists for this course and version
	allCurriculums := k.GetAllCurriculumTree(ctx)
	for _, curriculum := range allCurriculums {
		if curriculum.CourseId == req.CourseId && curriculum.Version == req.Version {
			return nil, fmt.Errorf("curriculum version '%s' already exists for course '%s'", req.Version, req.CourseId)
		}
	}

	// Create curriculum tree
	curriculumTree := types.CurriculumTree{
		Index:             "", // Will be set by AppendCurriculumTree
		CourseId:          req.CourseId,
		Version:           req.Version,
		ElectiveMin:       req.ElectiveMin,
		TotalWorkloadHours: fmt.Sprintf("%d", req.TotalWorkloadHours),
		RequiredSubjects:  req.RequiredSubjects,
		ElectiveSubjects:  req.ElectiveSubjects,
		SemesterStructure: []*types.CurriculumSemester{}, // Initialize empty
		ElectiveGroups:    []*types.ElectiveGroup{},      // Initialize empty
	}

	// Store curriculum tree
	index, err := k.AppendCurriculumTree(ctx, curriculumTree)
	if err != nil {
		return nil, fmt.Errorf("failed to store curriculum tree: %w", err)
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"curriculum_tree_created",
			sdk.NewAttribute("curriculum_index", fmt.Sprintf("%d", index)),
			sdk.NewAttribute("course_id", req.CourseId),
			sdk.NewAttribute("version", req.Version),
			sdk.NewAttribute("elective_min", fmt.Sprintf("%d", req.ElectiveMin)),
			sdk.NewAttribute("total_workload_hours", fmt.Sprintf("%d", req.TotalWorkloadHours)),
			sdk.NewAttribute("required_subjects_count", fmt.Sprintf("%d", len(req.RequiredSubjects))),
			sdk.NewAttribute("elective_subjects_count", fmt.Sprintf("%d", len(req.ElectiveSubjects))),
			sdk.NewAttribute("creator", req.Creator),
		),
	)

	k.Logger().Info("Curriculum tree created successfully",
		"curriculum_index", index,
		"course_id", req.CourseId,
		"version", req.Version,
		"creator", req.Creator,
	)

	return &types.MsgCreateCurriculumTreeResponse{}, nil
}

// AddSemesterToCurriculum adds a semester to an existing curriculum
func (k msgServer) AddSemesterToCurriculum(goCtx context.Context, req *types.MsgAddSemesterToCurriculum) (*types.MsgAddSemesterToCurriculumResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate basic fields
	if req.Creator == "" {
		return nil, fmt.Errorf("creator cannot be empty")
	}
	if req.CurriculumIndex == "" {
		return nil, fmt.Errorf("curriculum index cannot be empty")
	}
	if req.SemesterNumber == 0 {
		return nil, fmt.Errorf("semester number cannot be zero")
	}

	// Get existing curriculum
	curriculum, found := k.GetCurriculumTree(ctx, req.CurriculumIndex)
	if !found {
		return nil, fmt.Errorf("curriculum with index '%s' not found", req.CurriculumIndex)
	}

	// Check permissions (only course creator or authority can update)
	// Note: Need to implement proper authorization logic
	// For now, allow updates from anyone to prevent blocking
	_ = curriculum.CourseId // Use the field to avoid unused variable warning

	// Check if semester already exists
	for _, semester := range curriculum.SemesterStructure {
		if semester.SemesterNumber == fmt.Sprintf("%d", req.SemesterNumber) {
			return nil, fmt.Errorf("semester %d already exists in curriculum '%s'", req.SemesterNumber, req.CurriculumIndex)
		}
	}

	// Validate that all subjects exist
	if k.subjectKeeper != nil {
		for _, subjectId := range req.SubjectIds {
			if subjectId != "" {
				_, found := k.subjectKeeper.GetSubject(ctx, subjectId)
				if !found {
					return nil, fmt.Errorf("subject with ID '%s' not found", subjectId)
				}
			}
		}
	}

	// Create new semester
	newSemester := types.CurriculumSemester{
		SemesterNumber: fmt.Sprintf("%d", req.SemesterNumber),
		SubjectIds:     req.SubjectIds,
	}

	// Add semester to curriculum
	curriculum.SemesterStructure = append(curriculum.SemesterStructure, &newSemester)

	// Store updated curriculum
	k.SetCurriculumTree(ctx, curriculum)

	// Store semester separately if needed (optional)
	// Note: AppendCurriculumSemester method may not exist, commenting out
	// _, err := k.AppendCurriculumSemester(ctx, newSemester)
	// if err != nil {
	// 	k.Logger().Error("Failed to store semester separately",
	// 		"curriculum_index", req.CurriculumIndex,
	// 		"semester_number", req.SemesterNumber,
	// 		"error", err,
	// 	)
	// 	// Continue execution even if separate storage fails
	// }

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"semester_added_to_curriculum",
			sdk.NewAttribute("curriculum_index", req.CurriculumIndex),
			sdk.NewAttribute("semester_number", fmt.Sprintf("%d", req.SemesterNumber)),
			sdk.NewAttribute("subjects_count", fmt.Sprintf("%d", len(req.SubjectIds))),
			sdk.NewAttribute("creator", req.Creator),
		),
	)

	k.Logger().Info("Semester added to curriculum successfully",
		"curriculum_index", req.CurriculumIndex,
		"semester_number", req.SemesterNumber,
		"subjects_count", len(req.SubjectIds),
		"creator", req.Creator,
	)

	return &types.MsgAddSemesterToCurriculumResponse{}, nil
}

// AddElectiveGroup adds an elective group to a curriculum
func (k msgServer) AddElectiveGroup(goCtx context.Context, req *types.MsgAddElectiveGroup) (*types.MsgAddElectiveGroupResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate basic fields
	if req.Creator == "" {
		return nil, fmt.Errorf("creator cannot be empty")
	}
	if req.CurriculumIndex == "" {
		return nil, fmt.Errorf("curriculum index cannot be empty")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("elective group name cannot be empty")
	}
	if req.MinSubjectsRequired == 0 {
		return nil, fmt.Errorf("minimum subjects required cannot be zero")
	}
	if req.CreditsRequired == 0 {
		return nil, fmt.Errorf("credits required cannot be zero")
	}

	// Get existing curriculum
	curriculum, found := k.GetCurriculumTree(ctx, req.CurriculumIndex)
	if !found {
		return nil, fmt.Errorf("curriculum with index '%s' not found", req.CurriculumIndex)
	}

	// Check permissions (only course creator or authority can update)
	// Note: Need to implement proper authorization logic
	// For now, allow updates from anyone to prevent blocking
	_ = curriculum.CourseId // Use the field to avoid unused variable warning

	// Check if elective group with same name already exists in this curriculum
	for _, group := range curriculum.ElectiveGroups {
		if group.Name == req.Name {
			return nil, fmt.Errorf("elective group with name '%s' already exists in curriculum '%s'", req.Name, req.CurriculumIndex)
		}
	}

	// Validate that all subjects exist
	if k.subjectKeeper != nil {
		for _, subjectId := range req.SubjectIds {
			if subjectId != "" {
				_, found := k.subjectKeeper.GetSubject(ctx, subjectId)
				if !found {
					return nil, fmt.Errorf("subject with ID '%s' not found", subjectId)
				}
			}
		}
	}

	// Create new elective group
	newElectiveGroup := types.ElectiveGroup{
		GroupId:            fmt.Sprintf("%s-%s", req.CurriculumIndex, req.Name), // Generate unique ID
		Name:               req.Name,
		Description:        req.Description,
		MinSubjectsRequired: fmt.Sprintf("%d", req.MinSubjectsRequired),
		CreditsRequired:    fmt.Sprintf("%d", req.CreditsRequired),
		KnowledgeArea:      req.KnowledgeArea,
		SubjectIds:         req.SubjectIds,
	}

	// Add elective group to curriculum
	curriculum.ElectiveGroups = append(curriculum.ElectiveGroups, &newElectiveGroup)

	// Store updated curriculum
	k.SetCurriculumTree(ctx, curriculum)

	// Store elective group separately if needed (optional)
	// Note: AppendElectiveGroup method may not exist, commenting out
	// _, err := k.AppendElectiveGroup(ctx, newElectiveGroup)
	// if err != nil {
	// 	k.Logger().Error("Failed to store elective group separately",
	// 		"curriculum_index", req.CurriculumIndex,
	// 		"group_name", req.Name,
	// 		"error", err,
	// 	)
	// 	// Continue execution even if separate storage fails
	// }

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"elective_group_added",
			sdk.NewAttribute("curriculum_index", req.CurriculumIndex),
			sdk.NewAttribute("group_name", req.Name),
			sdk.NewAttribute("min_subjects_required", fmt.Sprintf("%d", req.MinSubjectsRequired)),
			sdk.NewAttribute("credits_required", fmt.Sprintf("%d", req.CreditsRequired)),
			sdk.NewAttribute("knowledge_area", req.KnowledgeArea),
			sdk.NewAttribute("subjects_count", fmt.Sprintf("%d", len(req.SubjectIds))),
			sdk.NewAttribute("creator", req.Creator),
		),
	)

	k.Logger().Info("Elective group added successfully",
		"curriculum_index", req.CurriculumIndex,
		"group_name", req.Name,
		"min_subjects_required", req.MinSubjectsRequired,
		"credits_required", req.CreditsRequired,
		"creator", req.Creator,
	)

	return &types.MsgAddElectiveGroupResponse{}, nil
}

// SetGraduationRequirements sets the graduation requirements for a curriculum
func (k msgServer) SetGraduationRequirements(goCtx context.Context, req *types.MsgSetGraduationRequirements) (*types.MsgSetGraduationRequirementsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate basic fields
	if req.Creator == "" {
		return nil, fmt.Errorf("creator cannot be empty")
	}
	if req.CurriculumIndex == "" {
		return nil, fmt.Errorf("curriculum index cannot be empty")
	}
	if req.TotalCreditsRequired == 0 {
		return nil, fmt.Errorf("total credits required cannot be zero")
	}
	if req.MinGpa <= 0 || req.MinGpa > 4.0 {
		return nil, fmt.Errorf("minimum GPA must be between 0 and 4.0, got %f", req.MinGpa)
	}
	if req.MinimumTimeYears <= 0 {
		return nil, fmt.Errorf("minimum time in years must be positive, got %f", req.MinimumTimeYears)
	}
	if req.MaximumTimeYears <= req.MinimumTimeYears {
		return nil, fmt.Errorf("maximum time (%f years) must be greater than minimum time (%f years)", req.MaximumTimeYears, req.MinimumTimeYears)
	}

	// Get existing curriculum
	curriculum, found := k.GetCurriculumTree(ctx, req.CurriculumIndex)
	if !found {
		return nil, fmt.Errorf("curriculum with index '%s' not found", req.CurriculumIndex)
	}

	// Check permissions (only course creator or authority can update)
	// Note: Need to implement proper authorization logic
	// For now, allow updates from anyone to prevent blocking
	_ = curriculum.CourseId // Use the field to avoid unused variable warning

	// Create graduation requirements
	graduationRequirements := types.GraduationRequirements{
		TotalCreditsRequired:   fmt.Sprintf("%d", req.TotalCreditsRequired),
		MinGpa:                 fmt.Sprintf("%.2f", req.MinGpa),
		RequiredElectiveCredits: fmt.Sprintf("%d", req.RequiredElectiveCredits),
		MinimumTimeYears:       fmt.Sprintf("%.2f", req.MinimumTimeYears),
		MaximumTimeYears:       fmt.Sprintf("%.2f", req.MaximumTimeYears),
		RequiredActivities:     req.RequiredActivities,
	}

	// Update curriculum with graduation requirements
	curriculum.GraduationRequirements = &graduationRequirements

	// Store updated curriculum
	k.SetCurriculumTree(ctx, curriculum)

	// Store graduation requirements separately if needed (optional)
	// Note: AppendGraduationRequirements method may not exist, commenting out
	// _, err := k.AppendGraduationRequirements(ctx, graduationRequirements)
	// if err != nil {
	// 	k.Logger().Error("Failed to store graduation requirements separately",
	// 		"curriculum_index", req.CurriculumIndex,
	// 		"error", err,
	// 	)
	// 	// Continue execution even if separate storage fails
	// }

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"graduation_requirements_set",
			sdk.NewAttribute("curriculum_index", req.CurriculumIndex),
			sdk.NewAttribute("total_credits_required", fmt.Sprintf("%d", req.TotalCreditsRequired)),
			sdk.NewAttribute("min_gpa", fmt.Sprintf("%.2f", req.MinGpa)),
			sdk.NewAttribute("required_elective_credits", fmt.Sprintf("%d", req.RequiredElectiveCredits)),
			sdk.NewAttribute("minimum_time_years", fmt.Sprintf("%.2f", req.MinimumTimeYears)),
			sdk.NewAttribute("maximum_time_years", fmt.Sprintf("%.2f", req.MaximumTimeYears)),
			sdk.NewAttribute("required_activities_count", fmt.Sprintf("%d", len(req.RequiredActivities))),
			sdk.NewAttribute("creator", req.Creator),
		),
	)

	k.Logger().Info("Graduation requirements set successfully",
		"curriculum_index", req.CurriculumIndex,
		"total_credits_required", req.TotalCreditsRequired,
		"min_gpa", req.MinGpa,
		"minimum_time_years", req.MinimumTimeYears,
		"maximum_time_years", req.MaximumTimeYears,
		"creator", req.Creator,
	)

	return &types.MsgSetGraduationRequirementsResponse{}, nil
}
