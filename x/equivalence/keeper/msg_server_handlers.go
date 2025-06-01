package keeper

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"academictoken/x/equivalence/types"
)

// RequestEquivalence handles equivalence request
func (k msgServer) RequestEquivalence(goCtx context.Context, req *types.MsgRequestEquivalence) (*types.MsgRequestEquivalenceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate basic request
	if req.SourceSubjectId == "" || req.TargetSubjectId == "" || req.TargetInstitution == "" {
		return nil, errors.Wrap(types.ErrInvalidRequest, "missing required fields")
	}

	// Validate that source and target subjects are different
	if req.SourceSubjectId == req.TargetSubjectId {
		return nil, errors.Wrap(types.ErrInvalidRequest, "source and target subjects cannot be the same")
	}

	// TODO: Validate subjects exist using subject keeper
	// if ms.subjectKeeper != nil {
	//     _, found := ms.subjectKeeper.GetSubject(ctx, req.SourceSubjectId)
	//     if !found {
	//         return nil, errors.Wrapf(types.ErrSubjectNotFound, "source subject %s", req.SourceSubjectId)
	//     }
	//     _, found = ms.subjectKeeper.GetSubject(ctx, req.TargetSubjectId)
	//     if !found {
	//         return nil, errors.Wrapf(types.ErrSubjectNotFound, "target subject %s", req.TargetSubjectId)
	//     }
	// }

	// TODO: Validate institution exists and is authorized
	// if ms.institutionKeeper != nil {
	//     _, found := ms.institutionKeeper.GetInstitution(ctx, req.TargetInstitution)
	//     if !found {
	//         return nil, errors.Wrapf(types.ErrInstitutionNotFound, "institution %s", req.TargetInstitution)
	//     }
	//     if !ms.institutionKeeper.IsInstitutionAuthorized(ctx, req.TargetInstitution) {
	//         return nil, errors.Wrapf(types.ErrInstitutionNotAuthorized, "institution %s", req.TargetInstitution)
	//     }
	// }

	// Create equivalence request
	equivalenceId, err := k.Keeper.CreateEquivalenceRequest(goCtx, req.SourceSubjectId, req.TargetInstitution, req.TargetSubjectId, req.ForceRecalculation)
	if err != nil {
		return nil, errors.Wrap(types.ErrEquivalenceUpdateFailed, err.Error())
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"equivalence_requested",
			sdk.NewAttribute("equivalence_id", equivalenceId),
			sdk.NewAttribute("source_subject_id", req.SourceSubjectId),
			sdk.NewAttribute("target_subject_id", req.TargetSubjectId),
			sdk.NewAttribute("target_institution", req.TargetInstitution),
			sdk.NewAttribute("creator", req.Creator),
		),
	)

	return &types.MsgRequestEquivalenceResponse{
		EquivalenceId:     equivalenceId,
		Status:            types.EquivalenceStatusPending,
		AnalysisTriggered: true, // In real implementation, this would depend on contract availability
	}, nil
}

// ExecuteEquivalenceAnalysis handles contract execution for analysis
func (k msgServer) ExecuteEquivalenceAnalysis(goCtx context.Context, req *types.MsgExecuteEquivalenceAnalysis) (*types.MsgExecuteEquivalenceAnalysisResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate request
	if req.EquivalenceId == "" || req.ContractAddress == "" {
		return nil, errors.Wrap(types.ErrInvalidRequest, "missing required fields")
	}

	// Get equivalence
	equivalence, found := k.Keeper.GetSubjectEquivalence(goCtx, req.EquivalenceId)
	if !found {
		return nil, errors.Wrapf(types.ErrEquivalenceNotFound, "equivalence %s", req.EquivalenceId)
	}

	// TODO: Execute CosmWasm contract analysis
	// For now, simulate contract execution with mock data
	mockEquivalencePercent := "85.50"
	mockAnalysisMetadata := fmt.Sprintf(`{
		"analysis_date": "%s",
		"contract_address": "%s",
		"source_subject": "%s",
		"target_subject": "%s",
		"similarity_score": %s,
		"analysis_method": "content_comparison",
		"confidence_level": "high"
	}`, time.Now().Format(time.RFC3339), req.ContractAddress, equivalence.SourceSubjectId, equivalence.TargetSubjectId, mockEquivalencePercent)

	// Update equivalence with analysis results
	err := k.Keeper.UpdateEquivalenceAnalysis(
		goCtx,
		req.EquivalenceId,
		req.ContractAddress,
		mockEquivalencePercent,
		mockAnalysisMetadata,
		"v1.0.0", // Mock contract version
	)
	if err != nil {
		return nil, errors.Wrap(types.ErrContractAnalysisFailed, err.Error())
	}

	// Get updated equivalence
	updatedEquivalence, _ := k.Keeper.GetSubjectEquivalence(goCtx, req.EquivalenceId)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"equivalence_analyzed",
			sdk.NewAttribute("equivalence_id", req.EquivalenceId),
			sdk.NewAttribute("contract_address", req.ContractAddress),
			sdk.NewAttribute("equivalence_percent", mockEquivalencePercent),
			sdk.NewAttribute("status", updatedEquivalence.EquivalenceStatus),
			sdk.NewAttribute("analysis_hash", updatedEquivalence.AnalysisHash),
		),
	)

	return &types.MsgExecuteEquivalenceAnalysisResponse{
		Success:            true,
		EquivalencePercent: mockEquivalencePercent,
		AnalysisMetadata:   mockAnalysisMetadata,
		AnalysisHash:       updatedEquivalence.AnalysisHash,
		UpdatedStatus:      updatedEquivalence.EquivalenceStatus,
	}, nil
}

// BatchRequestEquivalence handles batch equivalence requests
func (k msgServer) BatchRequestEquivalence(goCtx context.Context, req *types.MsgBatchRequestEquivalence) (*types.MsgBatchRequestEquivalenceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if len(req.Requests) == 0 {
		return nil, errors.Wrap(types.ErrInvalidRequest, "no requests provided")
	}

	var results []*types.BatchEquivalenceResult
	var successfulRequests uint64
	var failedRequests uint64

	// Process each request
	for _, equivalenceReq := range req.Requests {
		result := &types.BatchEquivalenceResult{
			SourceSubjectId: equivalenceReq.SourceSubjectId,
			TargetSubjectId: equivalenceReq.TargetSubjectId,
		}

		// Create individual equivalence request
		equivalenceId, err := k.Keeper.CreateEquivalenceRequest(
			goCtx,
			equivalenceReq.SourceSubjectId,
			equivalenceReq.TargetInstitution,
			equivalenceReq.TargetSubjectId,
			req.ForceRecalculation,
		)

		if err != nil {
			result.Status = "failed"
			result.Error = err.Error()
			failedRequests++
		} else {
			result.EquivalenceId = equivalenceId
			result.Status = types.EquivalenceStatusPending
			successfulRequests++
		}

		results = append(results, result)
	}

	// Emit batch event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"batch_equivalence_requested",
			sdk.NewAttribute("total_requests", strconv.Itoa(len(req.Requests))),
			sdk.NewAttribute("successful_requests", strconv.FormatUint(successfulRequests, 10)),
			sdk.NewAttribute("failed_requests", strconv.FormatUint(failedRequests, 10)),
			sdk.NewAttribute("creator", req.Creator),
		),
	)

	return &types.MsgBatchRequestEquivalenceResponse{
		Results:            results,
		SuccessfulRequests: successfulRequests,
		FailedRequests:     failedRequests,
	}, nil
}

// UpdateContractAddress handles contract address updates
func (k msgServer) UpdateContractAddress(goCtx context.Context, req *types.MsgUpdateContractAddress) (*types.MsgUpdateContractAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate authority
	if k.Keeper.GetAuthority() != req.Authority {
		return nil, errors.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.Keeper.GetAuthority(), req.Authority)
	}

	// TODO: Store contract address in module params or separate store
	// For now, just emit event as placeholder
	previousContractAddress := "previous_contract_placeholder" // Would get from params

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"contract_address_updated",
			sdk.NewAttribute("previous_contract_address", previousContractAddress),
			sdk.NewAttribute("new_contract_address", req.NewContractAddress),
			sdk.NewAttribute("contract_version", req.ContractVersion),
			sdk.NewAttribute("authority", req.Authority),
		),
	)

	return &types.MsgUpdateContractAddressResponse{
		Success:                 true,
		PreviousContractAddress: previousContractAddress,
		NewContractAddress:      req.NewContractAddress,
	}, nil
}

// ReanalyzeEquivalence handles re-analysis requests
func (k msgServer) ReanalyzeEquivalence(goCtx context.Context, req *types.MsgReanalyzeEquivalence) (*types.MsgReanalyzeEquivalenceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get existing equivalence
	equivalence, found := k.Keeper.GetSubjectEquivalence(goCtx, req.EquivalenceId)
	if !found {
		return nil, errors.Wrapf(types.ErrEquivalenceNotFound, "equivalence %s", req.EquivalenceId)
	}

	previousStatus := equivalence.EquivalenceStatus

	// Use provided contract address or existing one
	contractAddress := req.ContractAddress
	if contractAddress == "" {
		contractAddress = equivalence.ContractAddress
	}

	if contractAddress == "" {
		return nil, errors.Wrap(types.ErrInvalidContractAddress, "no contract address available for analysis")
	}

	// TODO: Execute contract re-analysis
	// For now, simulate with mock data
	newEquivalencePercent := "78.25" // Mock new percentage
	newAnalysisMetadata := fmt.Sprintf(`{
		"reanalysis_date": "%s",
		"contract_address": "%s",
		"source_subject": "%s",
		"target_subject": "%s",
		"similarity_score": %s,
		"analysis_method": "enhanced_content_comparison",
		"confidence_level": "high",
		"reanalysis_reason": "%s"
	}`, time.Now().Format(time.RFC3339), contractAddress, equivalence.SourceSubjectId, equivalence.TargetSubjectId, newEquivalencePercent, req.ReanalysisReason)

	// Update equivalence
	err := k.Keeper.UpdateEquivalenceAnalysis(
		goCtx,
		req.EquivalenceId,
		contractAddress,
		newEquivalencePercent,
		newAnalysisMetadata,
		"v1.0.1", // Mock updated contract version
	)
	if err != nil {
		return nil, errors.Wrap(types.ErrContractAnalysisFailed, err.Error())
	}

	// Get updated equivalence
	updatedEquivalence, _ := k.Keeper.GetSubjectEquivalence(goCtx, req.EquivalenceId)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"equivalence_reanalyzed",
			sdk.NewAttribute("equivalence_id", req.EquivalenceId),
			sdk.NewAttribute("previous_status", previousStatus),
			sdk.NewAttribute("new_status", updatedEquivalence.EquivalenceStatus),
			sdk.NewAttribute("previous_percent", equivalence.EquivalencePercent),
			sdk.NewAttribute("new_percent", newEquivalencePercent),
			sdk.NewAttribute("reanalysis_reason", req.ReanalysisReason),
			sdk.NewAttribute("creator", req.Creator),
		),
	)

	return &types.MsgReanalyzeEquivalenceResponse{
		Success:            true,
		PreviousStatus:     previousStatus,
		NewStatus:          updatedEquivalence.EquivalenceStatus,
		EquivalencePercent: newEquivalencePercent,
		AnalysisMetadata:   newAnalysisMetadata,
	}, nil
}
