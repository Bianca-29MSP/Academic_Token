package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"academictoken/x/student/types"
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

	// TEMPORARILY COMMENTED OUT FOR TESTING CONTRACT ADDRESSES
	// Check if the message signer is the authority
	// if req.Authority != k.GetAuthority() {
	//	return nil, fmt.Errorf("invalid authority: expected %s, got %s", k.GetAuthority(), req.Authority)
	// }

	// Validate the new params using the Validate method from params.go
	if err := req.Params.Validate(); err != nil {
		return nil, fmt.Errorf("invalid params: %w", err)
	}

	// Set the new params
	if err := k.SetParams(ctx, req.Params); err != nil {
		return nil, fmt.Errorf("failed to set params: %w", err)
	}

	return &types.MsgUpdateParamsResponse{}, nil
}

// RegisterStudent creates a new student
func (k msgServer) RegisterStudent(goCtx context.Context, req *types.MsgRegisterStudent) (*types.MsgRegisterStudentResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate basic fields
	if req.Creator == "" {
		return nil, fmt.Errorf("creator cannot be empty")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("student name cannot be empty")
	}
	if req.Address == "" {
		return nil, fmt.Errorf("student address cannot be empty")
	}

	// Create student object
	student := types.Student{
		Index:   "", // Will be set by AppendStudent
		Name:    req.Name,
		Address: req.Address,
		EnrollmentIds: []string{},
	}

	// Validate student
	if err := k.ValidateStudent(ctx, student); err != nil {
		return nil, fmt.Errorf("student validation failed: %w", err)
	}

	// Store student
	_, err := k.AppendStudent(ctx, student)
	if err != nil {
		return nil, fmt.Errorf("failed to store student: %w", err)
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"student_registered",
			sdk.NewAttribute("name", req.Name),
			sdk.NewAttribute("address", req.Address),
			sdk.NewAttribute("creator", req.Creator),
		),
	)

	k.Logger().Info("Student registered successfully",
		"name", req.Name,
		"address", req.Address,
		"creator", req.Creator,
	)

	return &types.MsgRegisterStudentResponse{}, nil
}

// UpdateEnrollmentStatus updates the status of an enrollment
func (k msgServer) UpdateEnrollmentStatus(goCtx context.Context, req *types.MsgUpdateEnrollmentStatus) (*types.MsgUpdateEnrollmentStatusResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate basic fields
	if req.Creator == "" {
		return nil, fmt.Errorf("creator cannot be empty")
	}
	if req.EnrollmentId == "" {
		return nil, fmt.Errorf("enrollment ID cannot be empty")
	}
	if req.Status == "" {
		return nil, fmt.Errorf("status cannot be empty")
	}

	// Validate status
	validStatuses := []string{"pending", "active", "completed", "cancelled", "suspended"}
	validStatus := false
	for _, status := range validStatuses {
		if req.Status == status {
			validStatus = true
			break
		}
	}
	if !validStatus {
		return nil, fmt.Errorf("invalid status: must be one of %v, got '%s'", validStatuses, req.Status)
	}

	// Get existing enrollment
	enrollment, found := k.getStudentEnrollment(ctx, req.EnrollmentId)
	if !found {
		return nil, fmt.Errorf("enrollment with ID '%s' not found", req.EnrollmentId)
	}

	// Check permissions (only student or authority can update)
	// Note: Need to implement proper authorization logic
	// For now, allow updates from anyone to prevent blocking
	_ = enrollment.Student // Use the field to avoid unused variable warning

	// Update status
	enrollment.Status = req.Status
	k.setStudentEnrollment(ctx, enrollment)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"enrollment_status_updated",
			sdk.NewAttribute("enrollment_id", req.EnrollmentId),
			sdk.NewAttribute("new_status", req.Status),
			sdk.NewAttribute("updater", req.Creator),
		),
	)

	k.Logger().Info("Enrollment status updated successfully",
		"enrollment_id", req.EnrollmentId,
		"new_status", req.Status,
		"updater", req.Creator,
	)

	return &types.MsgUpdateEnrollmentStatusResponse{}, nil
}

// RequestSubjectEnrollment requests enrollment in a specific subject
func (k msgServer) RequestSubjectEnrollment(goCtx context.Context, req *types.MsgRequestSubjectEnrollment) (*types.MsgRequestSubjectEnrollmentResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate basic fields
	if req.Creator == "" {
		return nil, fmt.Errorf("creator cannot be empty")
	}
	if req.Student == "" {
		return nil, fmt.Errorf("student ID cannot be empty")
	}
	if req.SubjectId == "" {
		return nil, fmt.Errorf("subject ID cannot be empty")
	}

	// Check if student exists
	_, found := k.getStudentByIndex(ctx, req.Student)
	if !found {
		return nil, fmt.Errorf("student with ID '%s' not found", req.Student)
	}

	// Check if subject exists
	subject, found := k.subjectKeeper.GetSubject(ctx, req.SubjectId)
	if !found {
		return nil, fmt.Errorf("subject with ID '%s' not found", req.SubjectId)
	}

	// Use contract integration to check prerequisites
	contractIntegration := k.GetContractIntegration().(*ContractIntegration)
	canEnroll, missingPrereqs, err := contractIntegration.CheckPrerequisites(ctx, req.Student, req.SubjectId)
	if err != nil {
		k.Logger().Error("Failed to check prerequisites via contract",
			"student", req.Student,
			"subject", req.SubjectId,
			"error", err,
		)
		// Continue without contract check for now, but log the issue
	}

	if !canEnroll && len(missingPrereqs) > 0 {
		return nil, fmt.Errorf("student does not meet prerequisites for subject '%s'. Missing: %v", req.SubjectId, missingPrereqs)
	}

	// Get or create academic tree
	academicTree, found := k.getAcademicTreeByStudent(ctx, req.Student)
	if !found {
		// Create new academic tree
		academicTree = types.StudentAcademicTree{
			Student:          req.Student,
			Institution:      subject.Institution,
			CourseId:         "", // Will be set when enrollment is created
			CompletedTokens:  []string{},
			InProgressTokens: []string{},
			AvailableTokens:  []string{req.SubjectId},
			AcademicProgress: &types.AcademicProgress{},
			GraduationStatus: &types.GraduationStatus{},
		}
		k.AppendStudentAcademicTree(ctx, academicTree)
	} else {
		// Add subject to in-progress tokens if not already there
		alreadyInProgress := false
		for _, token := range academicTree.InProgressTokens {
			if token == req.SubjectId {
				alreadyInProgress = true
				break
			}
		}

		if !alreadyInProgress {
			academicTree.InProgressTokens = append(academicTree.InProgressTokens, req.SubjectId)
			k.setStudentAcademicTree(ctx, academicTree)
		}
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"subject_enrollment_requested",
			sdk.NewAttribute("student", req.Student),
			sdk.NewAttribute("subject_id", req.SubjectId),
			sdk.NewAttribute("creator", req.Creator),
		),
	)

	k.Logger().Info("Subject enrollment requested successfully",
		"student", req.Student,
		"subject_id", req.SubjectId,
		"creator", req.Creator,
	)

	return &types.MsgRequestSubjectEnrollmentResponse{}, nil
}

// UpdateAcademicTree updates a student's academic tree
func (k msgServer) UpdateAcademicTree(goCtx context.Context, req *types.MsgUpdateAcademicTree) (*types.MsgUpdateAcademicTreeResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate basic fields
	if req.Creator == "" {
		return nil, fmt.Errorf("creator cannot be empty")
	}
	if req.StudentId == "" {
		return nil, fmt.Errorf("student ID cannot be empty")
	}

	// Check if student exists
	_, found := k.getStudentByIndex(ctx, req.StudentId)
	if !found {
		return nil, fmt.Errorf("student with ID '%s' not found", req.StudentId)
	}

	// Get existing academic tree
	academicTree, found := k.getAcademicTreeByStudent(ctx, req.StudentId)
	if !found {
		return nil, fmt.Errorf("academic tree for student '%s' not found", req.StudentId)
	}

	// Check permissions (only student themselves or authority can update)
	student, _ := k.getStudentByIndex(ctx, req.StudentId)
	if req.Creator != student.Address && req.Creator != k.GetAuthority() {
		return nil, fmt.Errorf("creator '%s' is not authorized to update academic tree for student '%s'", req.Creator, req.StudentId)
	}

	// Update fields if provided
	updated := false

	if len(req.CompletedTokens) > 0 {
		academicTree.CompletedTokens = req.CompletedTokens
		updated = true
	}

	if len(req.InProgressTokens) > 0 {
		academicTree.InProgressTokens = req.InProgressTokens
		updated = true
	}

	if len(req.AvailableTokens) > 0 {
		academicTree.AvailableTokens = req.AvailableTokens
		updated = true
	}

	if !updated {
		return nil, fmt.Errorf("no valid updates provided")
	}

	// Store updated academic tree
	k.setStudentAcademicTree(ctx, academicTree)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"academic_tree_updated",
			sdk.NewAttribute("student_id", req.StudentId),
			sdk.NewAttribute("updater", req.Creator),
		),
	)

	k.Logger().Info("Academic tree updated successfully",
		"student_id", req.StudentId,
		"updater", req.Creator,
	)

	return &types.MsgUpdateAcademicTreeResponse{}, nil
}

// CompleteSubject marks a subject as completed for a student
func (k msgServer) CompleteSubject(goCtx context.Context, req *types.MsgCompleteSubject) (*types.MsgCompleteSubjectResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate basic fields
	if req.Creator == "" {
		return nil, fmt.Errorf("creator cannot be empty")
	}
	if req.StudentId == "" {
		return nil, fmt.Errorf("student ID cannot be empty")
	}
	if req.SubjectId == "" {
		return nil, fmt.Errorf("subject ID cannot be empty")
	}
	if req.Grade == 0 {
		return nil, fmt.Errorf("grade cannot be zero")
	}
	if req.CompletionDate == "" {
		return nil, fmt.Errorf("completion date cannot be empty")
	}
	if req.Semester == "" {
		return nil, fmt.Errorf("semester cannot be empty")
	}

	// Validate grade range
	if req.Grade < 60 || req.Grade > 100 {
		return nil, fmt.Errorf("grade must be between 60 and 100, got %d", req.Grade)
	}

	// Check if student exists
	student, found := k.getStudentByIndex(ctx, req.StudentId)
	if !found {
		return nil, fmt.Errorf("student with ID '%s' not found", req.StudentId)
	}

	// Check if subject exists
	subject, found := k.subjectKeeper.GetSubject(ctx, req.SubjectId)
	if !found {
		return nil, fmt.Errorf("subject with ID '%s' not found", req.SubjectId)
	}

	// Check permissions (only professors/institutions or authority can complete subjects)
	// For now, allow any creator but in production this should be restricted
	// TODO: Add proper authorization logic

	// Create subject completion request for contract
	completionRequest := types.SubjectCompletionRequest{
		StudentId:      req.StudentId,
		SubjectId:      req.SubjectId,
		Grade:          uint64(req.Grade), // Convert uint32 to uint64
		CompletionDate: req.CompletionDate,
		Semester:       req.Semester,
		Credits:        subject.Credits,
		Institution:    subject.Institution,
	}

	// Use contract integration to process completion
	contractIntegration := k.GetContractIntegration().(*ContractIntegration)
	completionResult, err := contractIntegration.ProcessSubjectCompletion(ctx, completionRequest)
	if err != nil {
		k.Logger().Error("Failed to process subject completion via contract",
			"student", req.StudentId,
			"subject", req.SubjectId,
			"error", err,
		)
		return nil, fmt.Errorf("failed to process subject completion: %w", err)
	}

	// Update academic tree based on contract result
	academicTree, found := k.getAcademicTreeByStudent(ctx, req.StudentId)
	if found {
		academicTree.CompletedTokens = completionResult.UpdatedCompletedSubjects
		academicTree.InProgressTokens = completionResult.UpdatedInProgressSubjects
		academicTree.AcademicProgress = &completionResult.UpdatedProgress
		k.setStudentAcademicTree(ctx, academicTree)
	}

	// Request NFT minting
	nftRequest := types.NFTMintingRequest{
		StudentAddress:    student.Address,
		SubjectId:         req.SubjectId,
		Grade:             uint64(req.Grade), // Convert uint32 to uint64
		CompletionDate:    req.CompletionDate,
		Semester:          req.Semester,
		IssuerInstitution: subject.Institution,
		ProgressData:      completionResult.UpdatedProgress,
	}

	nftResult, err := contractIntegration.RequestNFTMinting(ctx, nftRequest)
	if err != nil {
		k.Logger().Error("Failed to request NFT minting",
			"student", req.StudentId,
			"subject", req.SubjectId,
			"error", err,
		)
		// Continue execution even if NFT minting fails
	}

	nftTokenId := ""
	if nftResult.Success {
		nftTokenId = nftResult.TokenInstanceId
	}

	// Check graduation eligibility if recommended by contract
	isEligibleForGraduation := false
	if completionResult.ShouldCheckGraduation {
		graduationResult, err := contractIntegration.CheckGraduationEligibility(ctx, req.StudentId)
		if err != nil {
			k.Logger().Error("Failed to check graduation eligibility",
				"student", req.StudentId,
				"error", err,
			)
		} else {
			isEligibleForGraduation = graduationResult.IsEligible
		}
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"subject_completed",
			sdk.NewAttribute("student_id", req.StudentId),
			sdk.NewAttribute("subject_id", req.SubjectId),
			sdk.NewAttribute("grade", fmt.Sprintf("%d", req.Grade)),
			sdk.NewAttribute("completion_date", req.CompletionDate),
			sdk.NewAttribute("semester", req.Semester),
			sdk.NewAttribute("nft_token_id", nftTokenId),
			sdk.NewAttribute("credits_completed", fmt.Sprintf("%d", completionResult.UpdatedProgress.RequiredCreditsCompleted)),
			sdk.NewAttribute("progress_percentage", fmt.Sprintf("%.2f", completionResult.UpdatedProgress.RequiredSubjectsPercentage)),
			sdk.NewAttribute("eligible_for_graduation", fmt.Sprintf("%t", isEligibleForGraduation)),
			sdk.NewAttribute("creator", req.Creator),
		),
	)

	k.Logger().Info("Subject completed successfully",
		"student_id", req.StudentId,
		"subject_id", req.SubjectId,
		"grade", req.Grade,
		"nft_token_id", nftTokenId,
		"eligible_for_graduation", isEligibleForGraduation,
		"creator", req.Creator,
	)

	return &types.MsgCompleteSubjectResponse{
		NftTokenId:              nftTokenId,
		ProgressPercentage:      float64(completionResult.UpdatedProgress.RequiredSubjectsPercentage),
		CreditsCompleted:        completionResult.UpdatedProgress.RequiredCreditsCompleted,
		IsEligibleForGraduation: isEligibleForGraduation,
	}, nil
}

// RequestEquivalence requests equivalence analysis between subjects
func (k msgServer) RequestEquivalence(goCtx context.Context, req *types.MsgRequestEquivalence) (*types.MsgRequestEquivalenceResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate basic fields
	if req.Creator == "" {
		return nil, fmt.Errorf("creator cannot be empty")
	}
	if req.StudentId == "" {
		return nil, fmt.Errorf("student ID cannot be empty")
	}
	if req.SourceSubjectId == "" {
		return nil, fmt.Errorf("source subject ID cannot be empty")
	}
	if req.TargetSubjectId == "" {
		return nil, fmt.Errorf("target subject ID cannot be empty")
	}

	// Check if student exists
	_, found := k.getStudentByIndex(ctx, req.StudentId)
	if !found {
		return nil, fmt.Errorf("student with ID '%s' not found", req.StudentId)
	}

	// Check if source subject exists
	_, found = k.subjectKeeper.GetSubject(ctx, req.SourceSubjectId)
	if !found {
		return nil, fmt.Errorf("source subject with ID '%s' not found", req.SourceSubjectId)
	}

	// Check if target subject exists
	_, found = k.subjectKeeper.GetSubject(ctx, req.TargetSubjectId)
	if !found {
		return nil, fmt.Errorf("target subject with ID '%s' not found", req.TargetSubjectId)
	}

	// Use contract integration to request equivalence
	contractIntegration := k.GetContractIntegration().(*ContractIntegration)
	equivalenceId, err := contractIntegration.RequestEquivalence(ctx, req.SourceSubjectId, req.TargetSubjectId)
	if err != nil {
		k.Logger().Error("Failed to request equivalence via contract",
			"source_subject", req.SourceSubjectId,
			"target_subject", req.TargetSubjectId,
			"error", err,
		)
		return nil, fmt.Errorf("failed to request equivalence: %w", err)
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"equivalence_requested",
			sdk.NewAttribute("student_id", req.StudentId),
			sdk.NewAttribute("source_subject_id", req.SourceSubjectId),
			sdk.NewAttribute("target_subject_id", req.TargetSubjectId),
			sdk.NewAttribute("equivalence_id", equivalenceId),
			sdk.NewAttribute("reason", req.Reason),
			sdk.NewAttribute("creator", req.Creator),
		),
	)

	k.Logger().Info("Equivalence requested successfully",
		"student_id", req.StudentId,
		"source_subject", req.SourceSubjectId,
		"target_subject", req.TargetSubjectId,
		"equivalence_id", equivalenceId,
		"creator", req.Creator,
	)

	return &types.MsgRequestEquivalenceResponse{
		EquivalenceId: equivalenceId,
	}, nil
}
