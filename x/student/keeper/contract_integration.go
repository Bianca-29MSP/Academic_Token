package keeper

import (
	"encoding/json"
	"fmt"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"academictoken/x/student/types"
)

// ContractIntegration handles all CosmWasm contract interactions for the Student module
type ContractIntegration struct {
	keeper        *Keeper
	wasmMsgServer types.WasmMsgServer
	wasmQuerier   types.WasmQuerier
}

// NewContractIntegration creates a new contract integration instance
func NewContractIntegration(keeper *Keeper, wasmMsgServer types.WasmMsgServer, wasmQuerier types.WasmQuerier) *ContractIntegration {
	return &ContractIntegration{
		keeper:        keeper,
		wasmMsgServer: wasmMsgServer,
		wasmQuerier:   wasmQuerier,
	}
}

// ============================================================================
// PREREQUISITES CONTRACT INTEGRATION
// ============================================================================

// CheckPrerequisites calls the prerequisites contract to verify if student can enroll
func (ci *ContractIntegration) CheckPrerequisites(ctx sdk.Context, studentId string, subjectId string) (bool, []string, error) {
	params := ci.keeper.GetParams(ctx)
	contractAddr := params.PrerequisitesContractAddr

	if contractAddr == "" {
		return false, nil, fmt.Errorf("prerequisites contract address not configured")
	}

	// Create query message
	queryMsg := map[string]interface{}{
		"check_eligibility": map[string]interface{}{
			"student_id": studentId,
			"subject_id": subjectId,
		},
	}

	queryData, err := json.Marshal(queryMsg)
	if err != nil {
		return false, nil, fmt.Errorf("failed to marshal query message: %w", err)
	}

	// Execute contract query
	req := &wasmtypes.QuerySmartContractStateRequest{
		Address:   contractAddr,
		QueryData: queryData,
	}

	res, err := ci.wasmQuerier.SmartContractState(sdk.WrapSDKContext(ctx), req)
	if err != nil {
		return false, nil, fmt.Errorf("failed to query prerequisites contract: %w", err)
	}

	// Parse response
	var response struct {
		CanEnroll            bool     `json:"can_enroll"`
		MissingPrerequisites []string `json:"missing_prerequisites"`
		Details              string   `json:"details"`
	}

	if err := json.Unmarshal(res.Data, &response); err != nil {
		return false, nil, fmt.Errorf("failed to unmarshal contract response: %w", err)
	}

	return response.CanEnroll, response.MissingPrerequisites, nil
}

// UpdateStudentRecord updates student completion record in prerequisites contract
func (ci *ContractIntegration) UpdateStudentRecord(ctx sdk.Context, studentId string, completedSubject types.CompletedSubject) error {
	params := ci.keeper.GetParams(ctx)
	contractAddr := params.PrerequisitesContractAddr

	if contractAddr == "" {
		return fmt.Errorf("prerequisites contract address not configured")
	}

	// Create execute message
	executeMsg := map[string]interface{}{
		"update_student_record": map[string]interface{}{
			"student_id": studentId,
			"completed_subject": map[string]interface{}{
				"subject_id":      completedSubject.SubjectId,
				"credits":         completedSubject.Credits,
				"completion_date": completedSubject.CompletionDate,
				"grade":           completedSubject.Grade,
				"nft_token_id":    completedSubject.NftTokenId,
			},
		},
	}

	execData, err := json.Marshal(executeMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal execute message: %w", err)
	}

	// Get module address for contract execution
	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)

	// Execute contract
	req := &wasmtypes.MsgExecuteContract{
		Sender:   moduleAddr.String(),
		Contract: contractAddr,
		Msg:      execData,
		Funds:    nil,
	}

	_, err = ci.wasmMsgServer.ExecuteContract(sdk.WrapSDKContext(ctx), req)
	if err != nil {
		return fmt.Errorf("failed to execute prerequisites contract: %w", err)
	}

	return nil
}

// RegisterPrerequisites registers prerequisite groups for a subject in the contract
func (ci *ContractIntegration) RegisterPrerequisites(ctx sdk.Context, subjectId string, prerequisites []types.PrerequisiteGroup) error {
	params := ci.keeper.GetParams(ctx)
	contractAddr := params.PrerequisitesContractAddr

	if contractAddr == "" {
		return fmt.Errorf("prerequisites contract address not configured")
	}

	// Convert to contract format
	contractPrereqs := make([]map[string]interface{}, len(prerequisites))
	for i, prereq := range prerequisites {
		contractPrereqs[i] = map[string]interface{}{
			"id":                         prereq.Id,
			"subject_id":                 prereq.SubjectId,
			"group_type":                 prereq.GroupType,
			"minimum_credits":            prereq.MinimumCredits,
			"minimum_completed_subjects": prereq.MinimumCompletedSubjects,
			"subject_ids":                prereq.SubjectIds,
			"logic":                      prereq.Logic,
			"priority":                   prereq.Priority,
			"confidence":                 prereq.Confidence,
		}
	}

	// Create execute message
	executeMsg := map[string]interface{}{
		"register_prerequisites": map[string]interface{}{
			"subject_id":    subjectId,
			"prerequisites": contractPrereqs,
		},
	}

	execData, err := json.Marshal(executeMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal execute message: %w", err)
	}

	// Get module address for contract execution
	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)

	// Execute contract
	req := &wasmtypes.MsgExecuteContract{
		Sender:   moduleAddr.String(),
		Contract: contractAddr,
		Msg:      execData,
		Funds:    nil,
	}

	_, err = ci.wasmMsgServer.ExecuteContract(sdk.WrapSDKContext(ctx), req)
	if err != nil {
		return fmt.Errorf("failed to execute prerequisites contract: %w", err)
	}

	return nil
}

// ============================================================================
// EQUIVALENCE CONTRACT INTEGRATION
// ============================================================================

// RequestEquivalence requests equivalence analysis between two subjects
func (ci *ContractIntegration) RequestEquivalence(ctx sdk.Context, sourceSubjectId string, targetSubjectId string) (string, error) {
	params := ci.keeper.GetParams(ctx)
	contractAddr := params.EquivalenceContractAddr

	if contractAddr == "" {
		return "", fmt.Errorf("equivalence contract address not configured")
	}

	// Get source and target subject info
	sourceSubject, found := ci.keeper.subjectKeeper.GetSubject(ctx, sourceSubjectId)
	if !found {
		return "", fmt.Errorf("source subject not found: %s", sourceSubjectId)
	}

	targetSubject, found := ci.keeper.subjectKeeper.GetSubject(ctx, targetSubjectId)
	if !found {
		return "", fmt.Errorf("target subject not found: %s", targetSubjectId)
	}

	// Convert to contract format
	sourceInfo := map[string]interface{}{
		"subject_id":     sourceSubject.Index,
		"institution_id": sourceSubject.Institution,
		"subject_name":   sourceSubject.Title,
		"credits":        sourceSubject.Credits,
		"course_code":    sourceSubject.Code,
		"content_hash":   sourceSubject.ContentHash,
		"ipfs_link":      sourceSubject.IPFSLink,
		"metadata": map[string]interface{}{
			"department":          sourceSubject.KnowledgeArea,
			"level":               "undergraduate", // Default for now
			"duration_weeks":      16,              // Default for now
			"workload_hours":      sourceSubject.WorkloadHours,
			"prerequisites_count": 0,           // Could be calculated
			"language":            "português", // Default for now
		},
	}

	targetInfo := map[string]interface{}{
		"subject_id":     targetSubject.Index,
		"institution_id": targetSubject.Institution,
		"subject_name":   targetSubject.Title,
		"credits":        targetSubject.Credits,
		"course_code":    targetSubject.Code,
		"content_hash":   targetSubject.ContentHash,
		"ipfs_link":      targetSubject.IPFSLink,
		"metadata": map[string]interface{}{
			"department":          targetSubject.KnowledgeArea,
			"level":               "undergraduate",
			"duration_weeks":      16,
			"workload_hours":      targetSubject.WorkloadHours,
			"prerequisites_count": 0,
			"language":            "português",
		},
	}

	// Create execute message
	executeMsg := map[string]interface{}{
		"register_equivalence": map[string]interface{}{
			"source_subject":  sourceInfo,
			"target_subject":  targetInfo,
			"analysis_method": "automatic",
			"notes":           nil,
		},
	}

	execData, err := json.Marshal(executeMsg)
	if err != nil {
		return "", fmt.Errorf("failed to marshal execute message: %w", err)
	}

	// Get module address for contract execution
	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)

	// Execute contract
	req := &wasmtypes.MsgExecuteContract{
		Sender:   moduleAddr.String(),
		Contract: contractAddr,
		Msg:      execData,
		Funds:    nil,
	}

	res, err := ci.wasmMsgServer.ExecuteContract(sdk.WrapSDKContext(ctx), req)
	if err != nil {
		return "", fmt.Errorf("failed to execute equivalence contract: %w", err)
	}

	// Parse response to get equivalence ID
	var equivalenceId string

	// Try to find equivalence_id in the response data
	if res.Data != nil {
		var responseData map[string]interface{}
		if err := json.Unmarshal(res.Data, &responseData); err == nil {
			if id, exists := responseData["equivalence_id"]; exists {
				if idStr, ok := id.(string); ok {
					equivalenceId = idStr
				}
			}
		}
	}

	if equivalenceId == "" {
		return "", fmt.Errorf("equivalence ID not found in contract response")
	}

	return equivalenceId, nil
}

// CheckEquivalenceStatus checks the status of an equivalence analysis
func (ci *ContractIntegration) CheckEquivalenceStatus(ctx sdk.Context, equivalenceId string) (types.EquivalenceResult, error) {
	params := ci.keeper.GetParams(ctx)
	contractAddr := params.EquivalenceContractAddr

	if contractAddr == "" {
		return types.EquivalenceResult{}, fmt.Errorf("equivalence contract address not configured")
	}

	// Create query message
	queryMsg := map[string]interface{}{
		"get_equivalence": map[string]interface{}{
			"equivalence_id": equivalenceId,
		},
	}

	queryData, err := json.Marshal(queryMsg)
	if err != nil {
		return types.EquivalenceResult{}, fmt.Errorf("failed to marshal query message: %w", err)
	}

	// Execute contract query
	req := &wasmtypes.QuerySmartContractStateRequest{
		Address:   contractAddr,
		QueryData: queryData,
	}

	res, err := ci.wasmQuerier.SmartContractState(sdk.WrapSDKContext(ctx), req)
	if err != nil {
		return types.EquivalenceResult{}, fmt.Errorf("failed to query equivalence contract: %w", err)
	}

	// Parse response
	var response struct {
		Equivalence struct {
			Id                   string `json:"id"`
			SimilarityPercentage uint32 `json:"similarity_percentage"`
			Status               string `json:"status"`
			EquivalenceType      string `json:"equivalence_type"`
			Notes                string `json:"notes"`
		} `json:"equivalence"`
	}

	if err := json.Unmarshal(res.Data, &response); err != nil {
		return types.EquivalenceResult{}, fmt.Errorf("failed to unmarshal contract response: %w", err)
	}

	return types.EquivalenceResult{
		EquivalenceId:        response.Equivalence.Id,
		SimilarityPercentage: response.Equivalence.SimilarityPercentage,
		Status:               response.Equivalence.Status,
		EquivalenceType:      response.Equivalence.EquivalenceType,
		Notes:                response.Equivalence.Notes,
	}, nil
}

// ============================================================================
// ACADEMIC PROGRESS CONTRACT INTEGRATION
// ============================================================================

// ProcessSubjectCompletion calls Academic Progress Contract to validate and update progress
func (ci *ContractIntegration) ProcessSubjectCompletion(ctx sdk.Context, request types.SubjectCompletionRequest) (types.SubjectCompletionResult, error) {
	params := ci.keeper.GetParams(ctx)
	contractAddr := params.AcademicProgressContractAddr

	if contractAddr == "" {
		ci.LogContractCall(ctx, "ProcessSubjectCompletion", "none", false, fmt.Errorf("contract not configured"))
		return ci.mockProcessSubjectCompletion(ctx, request)
	}

	// Create execute message for contract
	executeMsg := map[string]interface{}{
		"process_completion": map[string]interface{}{
			"student_id":      request.StudentId,
			"subject_id":      request.SubjectId,
			"grade":           request.Grade,
			"completion_date": request.CompletionDate,
			"semester":        request.Semester,
			"credits":         request.Credits,
			"institution":     request.Institution,
		},
	}

	execData, err := json.Marshal(executeMsg)
	if err != nil {
		ci.LogContractCall(ctx, "ProcessSubjectCompletion", contractAddr, false, err)
		return types.SubjectCompletionResult{}, fmt.Errorf("failed to marshal execute message: %w", err)
	}

	// Execute contract
	req := &wasmtypes.MsgExecuteContract{
		Sender:   ci.GetModuleAddress().String(),
		Contract: contractAddr,
		Msg:      execData,
		Funds:    nil,
	}

	res, err := ci.wasmMsgServer.ExecuteContract(sdk.WrapSDKContext(ctx), req)
	if err != nil {
		ci.LogContractCall(ctx, "ProcessSubjectCompletion", contractAddr, false, err)
		return types.SubjectCompletionResult{}, fmt.Errorf("failed to execute academic progress contract: %w", err)
	}

	// Parse response from data field
	var response types.SubjectCompletionResult
	if res.Data != nil {
		if err := json.Unmarshal(res.Data, &response); err != nil {
			ci.LogContractCall(ctx, "ProcessSubjectCompletion", contractAddr, false, err)
			return types.SubjectCompletionResult{}, fmt.Errorf("failed to parse contract response: %w", err)
		}
		ci.LogContractCall(ctx, "ProcessSubjectCompletion", contractAddr, true, nil)
		return response, nil
	}

	ci.LogContractCall(ctx, "ProcessSubjectCompletion", contractAddr, false, fmt.Errorf("no completion result in response"))
	return types.SubjectCompletionResult{}, fmt.Errorf("no completion result found in contract response")
}

// CheckGraduationEligibility calls Degree Contract to check graduation status
func (ci *ContractIntegration) CheckGraduationEligibility(ctx sdk.Context, studentId string) (types.GraduationEligibilityResult, error) {
	params := ci.keeper.GetParams(ctx)
	contractAddr := params.DegreeContractAddr

	if contractAddr == "" {
		ci.LogContractCall(ctx, "CheckGraduationEligibility", "none", false, fmt.Errorf("contract not configured"))
		return ci.mockCheckGraduationEligibility(ctx, studentId)
	}

	// Create query message
	queryMsg := map[string]interface{}{
		"check_graduation_eligibility": map[string]interface{}{
			"student_id": studentId,
		},
	}

	queryData, err := json.Marshal(queryMsg)
	if err != nil {
		ci.LogContractCall(ctx, "CheckGraduationEligibility", contractAddr, false, err)
		return types.GraduationEligibilityResult{}, fmt.Errorf("failed to marshal query message: %w", err)
	}

	// Query contract
	req := &wasmtypes.QuerySmartContractStateRequest{
		Address:   contractAddr,
		QueryData: queryData,
	}

	res, err := ci.wasmQuerier.SmartContractState(sdk.WrapSDKContext(ctx), req)
	if err != nil {
		ci.LogContractCall(ctx, "CheckGraduationEligibility", contractAddr, false, err)
		return types.GraduationEligibilityResult{}, fmt.Errorf("failed to query degree contract: %w", err)
	}

	// Parse response
	var response types.GraduationEligibilityResult
	if err := json.Unmarshal(res.Data, &response); err != nil {
		ci.LogContractCall(ctx, "CheckGraduationEligibility", contractAddr, false, err)
		return types.GraduationEligibilityResult{}, fmt.Errorf("failed to unmarshal contract response: %w", err)
	}

	ci.LogContractCall(ctx, "CheckGraduationEligibility", contractAddr, true, nil)
	return response, nil
}

// RequestNFTMinting calls NFT Minting Contract to authorize NFT creation
func (ci *ContractIntegration) RequestNFTMinting(ctx sdk.Context, request types.NFTMintingRequest) (types.NFTMintingResult, error) {
	params := ci.keeper.GetParams(ctx)
	contractAddr := params.NftMintingContractAddr

	if contractAddr == "" {
		ci.LogContractCall(ctx, "RequestNFTMinting", "none", false, fmt.Errorf("contract not configured"))
		return ci.mockRequestNFTMinting(ctx, request)
	}

	// Create execute message
	executeMsg := map[string]interface{}{
		"authorize_minting": map[string]interface{}{
			"student_address":    request.StudentAddress,
			"subject_id":         request.SubjectId,
			"grade":              request.Grade,
			"completion_date":    request.CompletionDate,
			"semester":           request.Semester,
			"issuer_institution": request.IssuerInstitution,
			"progress_data":      request.ProgressData,
		},
	}

	execData, err := json.Marshal(executeMsg)
	if err != nil {
		ci.LogContractCall(ctx, "RequestNFTMinting", contractAddr, false, err)
		return types.NFTMintingResult{}, fmt.Errorf("failed to marshal execute message: %w", err)
	}

	// Execute contract
	req := &wasmtypes.MsgExecuteContract{
		Sender:   ci.GetModuleAddress().String(),
		Contract: contractAddr,
		Msg:      execData,
		Funds:    nil,
	}

	res, err := ci.wasmMsgServer.ExecuteContract(sdk.WrapSDKContext(ctx), req)
	if err != nil {
		ci.LogContractCall(ctx, "RequestNFTMinting", contractAddr, false, err)
		return types.NFTMintingResult{}, fmt.Errorf("failed to execute NFT minting contract: %w", err)
	}

	// Parse response from data field
	var response types.NFTMintingResult
	if res.Data != nil {
		if err := json.Unmarshal(res.Data, &response); err != nil {
			ci.LogContractCall(ctx, "RequestNFTMinting", contractAddr, false, err)
			return types.NFTMintingResult{}, fmt.Errorf("failed to parse contract response: %w", err)
		}
		ci.LogContractCall(ctx, "RequestNFTMinting", contractAddr, true, nil)
		return response, nil
	}

	ci.LogContractCall(ctx, "RequestNFTMinting", contractAddr, false, fmt.Errorf("no minting result in response"))
	return types.NFTMintingResult{}, fmt.Errorf("no minting result found in contract response")
}

// ============================================================================
// DEGREE CONTRACT INTEGRATION
// ============================================================================

// ValidateDegreeRequirements calls Degree Contract to validate graduation requirements
func (ci *ContractIntegration) ValidateDegreeRequirements(ctx sdk.Context, request types.DegreeValidationRequest) (types.DegreeValidationResult, error) {
	params := ci.keeper.GetParams(ctx)
	contractAddr := params.DegreeContractAddr

	if contractAddr == "" {
		ci.LogContractCall(ctx, "ValidateDegreeRequirements", "none", false, fmt.Errorf("contract not configured"))
		return ci.mockValidateDegreeRequirements(ctx, request)
	}

	// Create execute message for degree validation
	executeMsg := map[string]interface{}{
		"validate_degree_requirements": map[string]interface{}{
			"student_id":     request.StudentId,
			"curriculum_id":  request.CurriculumId,
			"institution_id": request.InstitutionId,
			"final_gpa":      request.FinalGPA,
			"total_credits":  request.TotalCredits,
			"signatures":     request.Signatures,
			"requested_date": request.RequestedDate,
		},
	}

	execData, err := json.Marshal(executeMsg)
	if err != nil {
		ci.LogContractCall(ctx, "ValidateDegreeRequirements", contractAddr, false, err)
		return types.DegreeValidationResult{}, fmt.Errorf("failed to marshal execute message: %w", err)
	}

	// Execute contract
	req := &wasmtypes.MsgExecuteContract{
		Sender:   ci.GetModuleAddress().String(),
		Contract: contractAddr,
		Msg:      execData,
		Funds:    nil,
	}

	res, err := ci.wasmMsgServer.ExecuteContract(sdk.WrapSDKContext(ctx), req)
	if err != nil {
		ci.LogContractCall(ctx, "ValidateDegreeRequirements", contractAddr, false, err)
		return types.DegreeValidationResult{}, fmt.Errorf("failed to execute degree validation contract: %w", err)
	}

	// Parse response from data field
	var response types.DegreeValidationResult
	if res.Data != nil {
		if err := json.Unmarshal(res.Data, &response); err != nil {
			ci.LogContractCall(ctx, "ValidateDegreeRequirements", contractAddr, false, err)
			return types.DegreeValidationResult{}, fmt.Errorf("failed to parse contract response: %w", err)
		}
		ci.LogContractCall(ctx, "ValidateDegreeRequirements", contractAddr, true, nil)
		return response, nil
	}

	ci.LogContractCall(ctx, "ValidateDegreeRequirements", contractAddr, false, fmt.Errorf("no validation result in response"))
	return types.DegreeValidationResult{}, fmt.Errorf("no validation result found in contract response")
}

// AuthorizeDegreeNFTMinting calls Degree Contract to authorize degree NFT minting
func (ci *ContractIntegration) AuthorizeDegreeNFTMinting(ctx sdk.Context, request types.DegreeNFTMintingRequest) (types.DegreeNFTMintingResult, error) {
	params := ci.keeper.GetParams(ctx)
	contractAddr := params.DegreeContractAddr

	if contractAddr == "" {
		ci.LogContractCall(ctx, "AuthorizeDegreeNFTMinting", "none", false, fmt.Errorf("contract not configured"))
		return ci.mockAuthorizeDegreeNFTMinting(ctx, request)
	}

	// Create execute message for degree NFT authorization
	executeMsg := map[string]interface{}{
		"authorize_degree_nft": map[string]interface{}{
			"student_id":      request.StudentId,
			"curriculum_id":   request.CurriculumId,
			"institution_id":  request.InstitutionId,
			"degree_type":     request.DegreeType,
			"final_gpa":       request.FinalGPA,
			"total_credits":   request.TotalCredits,
			"validation_data": request.ValidationData,
			"issue_date":      request.IssueDate,
		},
	}

	execData, err := json.Marshal(executeMsg)
	if err != nil {
		ci.LogContractCall(ctx, "AuthorizeDegreeNFTMinting", contractAddr, false, err)
		return types.DegreeNFTMintingResult{}, fmt.Errorf("failed to marshal execute message: %w", err)
	}

	// Execute contract
	req := &wasmtypes.MsgExecuteContract{
		Sender:   ci.GetModuleAddress().String(),
		Contract: contractAddr,
		Msg:      execData,
		Funds:    nil,
	}

	res, err := ci.wasmMsgServer.ExecuteContract(sdk.WrapSDKContext(ctx), req)
	if err != nil {
		ci.LogContractCall(ctx, "AuthorizeDegreeNFTMinting", contractAddr, false, err)
		return types.DegreeNFTMintingResult{}, fmt.Errorf("failed to execute degree NFT authorization contract: %w", err)
	}

	// Parse response from data field
	var response types.DegreeNFTMintingResult
	if res.Data != nil {
		if err := json.Unmarshal(res.Data, &response); err != nil {
			ci.LogContractCall(ctx, "AuthorizeDegreeNFTMinting", contractAddr, false, err)
			return types.DegreeNFTMintingResult{}, fmt.Errorf("failed to parse contract response: %w", err)
		}
		ci.LogContractCall(ctx, "AuthorizeDegreeNFTMinting", contractAddr, true, nil)
		return response, nil
	}

	ci.LogContractCall(ctx, "AuthorizeDegreeNFTMinting", contractAddr, false, fmt.Errorf("no degree NFT result in response"))
	return types.DegreeNFTMintingResult{}, fmt.Errorf("no degree NFT result found in contract response")
}

// ============================================================================
// MOCK IMPLEMENTATIONS FOR DEVELOPMENT
// ============================================================================

// mockProcessSubjectCompletion provides mock data when contract is not available
func (ci *ContractIntegration) mockProcessSubjectCompletion(ctx sdk.Context, request types.SubjectCompletionRequest) (types.SubjectCompletionResult, error) {
	// Get current academic tree to simulate progress calculation
	academicTree, found := ci.keeper.getAcademicTreeByStudent(ctx, request.StudentId)
	if !found {
		// Return basic progress data
		return types.SubjectCompletionResult{
			Success:                   true,
			Message:                   "Subject completion validated (mock)",
			UpdatedCompletedSubjects:  []string{request.SubjectId},
			UpdatedInProgressSubjects: []string{},
			UpdatedProgress: types.AcademicProgress{
				RequiredCreditsCompleted:   request.Credits,
				ElectiveCreditsCompleted:   0,
				RequiredSubjectsPercentage: 10.0, // Mock 10% progress
				CurrentSemester:            1,
				CurrentYear:                1,
				EnrollmentYears:            0.5,
				ElectivesByAreaCompleted:   make(map[string]uint64),
			},
			ShouldCheckGraduation: false,
		}, nil
	}

	// Simulate progress update
	updatedCompleted := append(academicTree.CompletedTokens, request.SubjectId)
	updatedInProgress := ci.removeTokenFromList(academicTree.InProgressTokens, request.SubjectId)

	// Calculate mock progress
	totalCredits := request.Credits
	if academicTree.AcademicProgress != nil {
		totalCredits += academicTree.AcademicProgress.RequiredCreditsCompleted
	}

	progressPercentage := float32(float64(len(updatedCompleted)) * 5.0) // Mock: each subject = 5% progress
	shouldCheckGraduation := len(updatedCompleted) >= 8                 // Mock: 8 subjects = eligible for graduation

	return types.SubjectCompletionResult{
		Success:                   true,
		Message:                   "Subject completion validated (mock)",
		UpdatedCompletedSubjects:  updatedCompleted,
		UpdatedInProgressSubjects: updatedInProgress,
		UpdatedProgress: types.AcademicProgress{
			RequiredCreditsCompleted:   totalCredits,
			ElectiveCreditsCompleted:   0,
			RequiredSubjectsPercentage: progressPercentage,
			CurrentSemester:            1,
			CurrentYear:                1,
			EnrollmentYears:            0.5,
			ElectivesByAreaCompleted:   make(map[string]uint64),
		},
		ShouldCheckGraduation: shouldCheckGraduation,
	}, nil
}

// mockCheckGraduationEligibility provides mock graduation eligibility data
func (ci *ContractIntegration) mockCheckGraduationEligibility(ctx sdk.Context, studentId string) (types.GraduationEligibilityResult, error) {
	academicTree, found := ci.keeper.getAcademicTreeByStudent(ctx, studentId)
	if !found {
		return types.GraduationEligibilityResult{
			IsEligible: false,
			Message:    "No academic tree found",
		}, nil
	}

	// Mock graduation criteria: 8+ completed subjects
	isEligible := len(academicTree.CompletedTokens) >= 8
	message := "Not eligible yet"
	estimatedDate := "2025-12-01"

	if isEligible {
		message = "Eligible for graduation"
		estimatedDate = "2025-06-01"
	}

	return types.GraduationEligibilityResult{
		IsEligible:                isEligible,
		Message:                   message,
		EstimatedGraduationDate:   estimatedDate,
		RequiredCreditsRemaining:  0,
		RequiredSubjectsRemaining: []string{},
	}, nil
}

// mockRequestNFTMinting provides mock NFT minting authorization
func (ci *ContractIntegration) mockRequestNFTMinting(ctx sdk.Context, request types.NFTMintingRequest) (types.NFTMintingResult, error) {
	// Generate mock token instance ID
	tokenInstanceId := fmt.Sprintf("mock-token-%s-%s-%d", request.StudentAddress[:8], request.SubjectId, ctx.BlockHeight())

	return types.NFTMintingResult{
		Success:         true,
		TokenInstanceId: tokenInstanceId,
		Message:         "NFT minting authorized (mock)",
		MetadataHash:    "mock-metadata-hash",
	}, nil
}

// mockValidateDegreeRequirements provides mock degree validation
func (ci *ContractIntegration) mockValidateDegreeRequirements(ctx sdk.Context, request types.DegreeValidationRequest) (types.DegreeValidationResult, error) {
	// Get student's academic record for mock validation
	academicTree, found := ci.keeper.getAcademicTreeByStudent(ctx, request.StudentId)
	if !found {
		return types.DegreeValidationResult{
			IsValid:             false,
			Message:             "No academic record found for student",
			MissingRequirements: []string{"academic_record"},
		}, nil
	}

	// Mock graduation criteria
	requiredCredits := uint64(120)
	requiredSubjects := 8
	currentCredits := uint64(0)
	if academicTree.AcademicProgress != nil {
		currentCredits = academicTree.AcademicProgress.RequiredCreditsCompleted
	}

	// Check if requirements are met
	isValid := len(academicTree.CompletedTokens) >= requiredSubjects && currentCredits >= requiredCredits
	message := "Requirements not met"
	requirementsMet := []string{}
	missingRequirements := []string{}

	if len(academicTree.CompletedTokens) >= requiredSubjects {
		requirementsMet = append(requirementsMet, "required_subjects")
	} else {
		missingRequirements = append(missingRequirements, "required_subjects")
	}

	if currentCredits >= requiredCredits {
		requirementsMet = append(requirementsMet, "required_credits")
	} else {
		missingRequirements = append(missingRequirements, "required_credits")
	}

	if isValid {
		message = "All graduation requirements met"
	}

	return types.DegreeValidationResult{
		IsValid:             isValid,
		Message:             message,
		DegreeType:          "Bachelor",
		CurriculumVersion:   "1.0",
		ValidationHash:      fmt.Sprintf("mock-validation-%s-%d", request.StudentId, ctx.BlockHeight()),
		RequirementsMet:     requirementsMet,
		MissingRequirements: missingRequirements,
	}, nil
}

// mockAuthorizeDegreeNFTMinting provides mock degree NFT authorization
func (ci *ContractIntegration) mockAuthorizeDegreeNFTMinting(ctx sdk.Context, request types.DegreeNFTMintingRequest) (types.DegreeNFTMintingResult, error) {
	// Generate mock degree token ID
	tokenId := fmt.Sprintf("degree-nft-%s-%s-%d", request.StudentId, request.InstitutionId, ctx.BlockHeight())
	ipfsHash := fmt.Sprintf("QmMockDegreeHash%s", request.StudentId[:8])

	return types.DegreeNFTMintingResult{
		Success:          true,
		TokenId:          tokenId,
		MetadataIPFSHash: ipfsHash,
		Message:          "Degree NFT authorization granted (mock)",
	}, nil
}

// removeTokenFromList helper function
func (ci *ContractIntegration) removeTokenFromList(tokens []string, tokenToRemove string) []string {
	result := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if token != tokenToRemove {
			result = append(result, token)
		}
	}
	return result
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// GetModuleAddress returns the module's address for contract execution
func (ci *ContractIntegration) GetModuleAddress() sdk.AccAddress {
	return authtypes.NewModuleAddress(types.ModuleName)
}

// LogContractCall logs contract interactions for debugging
func (ci *ContractIntegration) LogContractCall(ctx sdk.Context, operation string, contractAddr string, success bool, err error) {
	if success {
		ci.keeper.Logger().Info("Contract call successful",
			"operation", operation,
			"contract", contractAddr,
			"block_height", ctx.BlockHeight(),
		)
	} else {
		ci.keeper.Logger().Error("Contract call failed",
			"operation", operation,
			"contract", contractAddr,
			"error", err,
			"block_height", ctx.BlockHeight(),
		)
	}
}
