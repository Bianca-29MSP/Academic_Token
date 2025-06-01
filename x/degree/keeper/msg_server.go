package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"academictoken/x/degree/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// CASCA: Delega requisição de diploma para contrato
func (k msgServer) RequestDegree(goCtx context.Context, req *types.MsgRequestDegree) (*types.MsgRequestDegreeResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Basic validation only
	if req.Creator == "" || req.StudentId == "" || req.InstitutionId == "" || req.CurriculumId == "" {
		return nil, fmt.Errorf("required fields cannot be empty")
	}

	// Prepare contract message
	contractMsg := map[string]interface{}{
		"request_degree": map[string]interface{}{
			"student_id":               req.StudentId,
			"institution_id":           req.InstitutionId,
			"curriculum_id":            req.CurriculumId,
			"expected_graduation_date": req.ExpectedGraduationDate,
			"creator":                  req.Creator,
		},
	}

	msgBytes, err := json.Marshal(contractMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal contract message: %w", err)
	}

	// Get contract address
	contractAddr := k.GetDegreeContractAddress(ctx)

	// Convert addresses
	senderAddr, err := sdk.AccAddressFromBech32(req.Creator)
	if err != nil {
		return nil, fmt.Errorf("invalid creator address: %w", err)
	}

	contractAccAddr, err := sdk.AccAddressFromBech32(contractAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid contract address: %w", err)
	}

	// Execute contract via wasm keeper with correct signature
	if k.wasmKeeper == nil {
		return nil, fmt.Errorf("wasm keeper not available")
	}

	execResp, err := k.wasmKeeper.Execute(ctx, senderAddr, contractAccAddr, msgBytes, []sdk.Coin{})
	if err != nil {
		return nil, fmt.Errorf("contract execution failed: %w", err)
	}

	// Parse contract response - execResp is []byte, not struct with Data field
	var contractResp struct {
		DegreeRequestId string `json:"degree_request_id"`
		Status          string `json:"status"`
	}

	if err := json.Unmarshal(execResp, &contractResp); err != nil {
		return nil, fmt.Errorf("failed to parse contract response: %w", err)
	}

	// Store minimal reference on-chain (contract handles the full logic)
	degreeRequest := types.DegreeRequest{
		Id:                     contractResp.DegreeRequestId,
		StudentId:              req.StudentId,
		InstitutionId:          req.InstitutionId,
		CurriculumId:           req.CurriculumId,
		ExpectedGraduationDate: req.ExpectedGraduationDate,
		Status:                 contractResp.Status,
		RequestDate:            time.Now().UTC().Format(time.RFC3339),
		Creator:                req.Creator,
	}

	k.SetDegreeRequest(ctx, degreeRequest)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDegreeRequested,
			sdk.NewAttribute(types.AttributeKeyDegreeRequestID, contractResp.DegreeRequestId),
			sdk.NewAttribute(types.AttributeKeyStudentID, req.StudentId),
			sdk.NewAttribute(types.AttributeKeyDegreeStatus, contractResp.Status),
			sdk.NewAttribute(types.AttributeKeyContractAddress, contractAddr),
		),
	)

	k.Logger(ctx).Info("Degree request delegated to contract",
		"degree_request_id", contractResp.DegreeRequestId,
		"contract_address", contractAddr,
		"student_id", req.StudentId,
	)

	return &types.MsgRequestDegreeResponse{
		DegreeRequestId: contractResp.DegreeRequestId,
		Status:          contractResp.Status,
	}, nil
}

// CASCA: Delega validação para contrato
func (k msgServer) ValidateDegreeRequirements(goCtx context.Context, req *types.MsgValidateDegreeRequirements) (*types.MsgValidateDegreeRequirementsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Basic validation only
	if req.Creator == "" || req.DegreeRequestId == "" {
		return nil, fmt.Errorf("required fields cannot be empty")
	}

	// Get degree request to verify it exists
	degreeRequest, found := k.GetDegreeRequest(ctx, req.DegreeRequestId)
	if !found {
		return nil, fmt.Errorf("degree request with ID '%s' not found", req.DegreeRequestId)
	}

	// Prepare contract message
	contractMsg := map[string]interface{}{
		"validate_degree_requirements": map[string]interface{}{
			"degree_request_id":     req.DegreeRequestId,
			"validation_parameters": req.ValidationParameters,
			"creator":               req.Creator,
		},
	}

	msgBytes, err := json.Marshal(contractMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal contract message: %w", err)
	}

	// Use contract address from request or from hardcode
	contractAddr := req.ContractAddress
	if contractAddr == "" {
		contractAddr = k.GetDegreeContractAddress(ctx)
	}

	// Convert addresses
	senderAddr, err := sdk.AccAddressFromBech32(req.Creator)
	if err != nil {
		return nil, fmt.Errorf("invalid creator address: %w", err)
	}

	contractAccAddr, err := sdk.AccAddressFromBech32(contractAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid contract address: %w", err)
	}

	// Call CosmWasm contract with correct signature
	execResp, err := k.wasmKeeper.Execute(ctx, senderAddr, contractAccAddr, msgBytes, []sdk.Coin{})
	if err != nil {
		return nil, fmt.Errorf("contract execution failed: %w", err)
	}

	// Parse contract response - execResp is []byte
	var contractResp struct {
		ValidationPassed    bool     `json:"validation_passed"`
		ValidationScore     string   `json:"validation_score"`
		ValidationDetails   string   `json:"validation_details"`
		MissingRequirements []string `json:"missing_requirements"`
	}

	if err := json.Unmarshal(execResp, &contractResp); err != nil {
		return nil, fmt.Errorf("failed to parse contract response: %w", err)
	}

	// Update degree request status based on contract response
	if contractResp.ValidationPassed {
		degreeRequest.Status = "validated"
	} else {
		degreeRequest.Status = "validation_failed"
	}

	k.SetDegreeRequest(ctx, degreeRequest)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDegreeValidated,
			sdk.NewAttribute(types.AttributeKeyDegreeRequestID, req.DegreeRequestId),
			sdk.NewAttribute(types.AttributeKeyValidationPassed, fmt.Sprintf("%t", contractResp.ValidationPassed)),
			sdk.NewAttribute(types.AttributeKeyValidationScore, contractResp.ValidationScore),
			sdk.NewAttribute(types.AttributeKeyContractAddress, contractAddr),
		),
	)

	k.Logger(ctx).Info("Degree validation delegated to contract",
		"degree_request_id", req.DegreeRequestId,
		"validation_passed", contractResp.ValidationPassed,
		"contract_address", contractAddr,
	)

	return &types.MsgValidateDegreeRequirementsResponse{
		ValidationPassed:    contractResp.ValidationPassed,
		ValidationScore:     contractResp.ValidationScore,
		ValidationDetails:   contractResp.ValidationDetails,
		MissingRequirements: contractResp.MissingRequirements,
	}, nil
}

// CASCA: Delega emissão de diploma para contrato
func (k msgServer) IssueDegree(goCtx context.Context, req *types.MsgIssueDegree) (*types.MsgIssueDegreeResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Basic validation only
	if req.Creator == "" || req.DegreeRequestId == "" || req.FinalGpa == "" || req.TotalCredits == 0 {
		return nil, fmt.Errorf("required fields cannot be empty or zero")
	}

	// Get degree request to verify it exists
	degreeRequest, found := k.GetDegreeRequest(ctx, req.DegreeRequestId)
	if !found {
		return nil, fmt.Errorf("degree request with ID '%s' not found", req.DegreeRequestId)
	}

	// Prepare contract message
	contractMsg := map[string]interface{}{
		"issue_degree": map[string]interface{}{
			"degree_request_id": req.DegreeRequestId,
			"final_gpa":         req.FinalGpa,
			"total_credits":     req.TotalCredits,
			"signatures":        req.Signatures,
			"creator":           req.Creator,
		},
	}

	msgBytes, err := json.Marshal(contractMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal contract message: %w", err)
	}

	// Get contract address
	contractAddr := k.GetDegreeContractAddress(ctx)

	// Convert addresses
	senderAddr, err := sdk.AccAddressFromBech32(req.Creator)
	if err != nil {
		return nil, fmt.Errorf("invalid creator address: %w", err)
	}

	contractAccAddr, err := sdk.AccAddressFromBech32(contractAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid contract address: %w", err)
	}

	// Call CosmWasm contract with correct signature
	execResp, err := k.wasmKeeper.Execute(ctx, senderAddr, contractAccAddr, msgBytes, []sdk.Coin{})
	if err != nil {
		return nil, fmt.Errorf("contract execution failed: %w", err)
	}

	// Parse contract response - execResp is []byte
	var contractResp struct {
		DegreeId   string `json:"degree_id"`
		NftTokenId string `json:"nft_token_id"`
		IpfsHash   string `json:"ipfs_hash"`
		IssueDate  string `json:"issue_date"`
	}

	if err := json.Unmarshal(execResp, &contractResp); err != nil {
		return nil, fmt.Errorf("failed to parse contract response: %w", err)
	}

	// Store minimal degree reference on-chain (contract handles full diploma data)
	degree := types.Degree{
		Index:        contractResp.DegreeId,
		DegreeId:     contractResp.DegreeId,
		Student:      degreeRequest.StudentId,
		Institution:  degreeRequest.InstitutionId,
		CourseId:     degreeRequest.CurriculumId,
		IssueDate:    contractResp.IssueDate,
		Status:       "issued",
		NftTokenId:   contractResp.NftTokenId,
		IpfsLink:     contractResp.IpfsHash,
		FinalGrade:   req.FinalGpa,
		TotalCredits: req.TotalCredits,
		Signatures:   req.Signatures,
	}

	k.SetDegree(ctx, degree)

	// Update degree request status
	degreeRequest.Status = "approved"
	k.SetDegreeRequest(ctx, degreeRequest)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDegreeIssued,
			sdk.NewAttribute(types.AttributeKeyDegreeID, contractResp.DegreeId),
			sdk.NewAttribute(types.AttributeKeyNFTTokenID, contractResp.NftTokenId),
			sdk.NewAttribute(types.AttributeKeyIPFSHash, contractResp.IpfsHash),
			sdk.NewAttribute(types.AttributeKeyContractAddress, contractAddr),
		),
	)

	k.Logger(ctx).Info("Degree issuance delegated to contract",
		"degree_id", contractResp.DegreeId,
		"nft_token_id", contractResp.NftTokenId,
		"contract_address", contractAddr,
	)

	return &types.MsgIssueDegreeResponse{
		DegreeId:   contractResp.DegreeId,
		NftTokenId: contractResp.NftTokenId,
		IpfsHash:   contractResp.IpfsHash,
		IssueDate:  contractResp.IssueDate,
	}, nil
}

// CASCA: Delega atualização de contrato
func (k msgServer) UpdateDegreeContract(goCtx context.Context, msg *types.MsgUpdateDegreeContract) (*types.MsgUpdateDegreeContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check authority
	if k.GetAuthority() != msg.Authority {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.GetAuthority(), msg.Authority)
	}

	// For hardcoded params, this would require a code update
	// In production, this could trigger a governance proposal
	oldAddress := k.GetDegreeContractAddress(ctx)

	// Emit event for tracking
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeContractUpdated,
			sdk.NewAttribute(types.AttributeKeyContractAddress, msg.NewContractAddress),
			sdk.NewAttribute("old_contract_address", oldAddress),
		),
	)

	k.Logger(ctx).Info("Degree contract update requested",
		"old_address", oldAddress,
		"new_address", msg.NewContractAddress,
		"authority", msg.Authority,
	)

	return &types.MsgUpdateDegreeContractResponse{
		OldContractAddress: oldAddress,
		NewContractAddress: msg.NewContractAddress,
		UpdateDate:         time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// CASCA: Delega cancelamento para contrato
func (k msgServer) CancelDegreeRequest(goCtx context.Context, msg *types.MsgCancelDegreeRequest) (*types.MsgCancelDegreeRequestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Basic validation
	if msg.Creator == "" || msg.DegreeRequestId == "" {
		return nil, fmt.Errorf("required fields cannot be empty")
	}

	// Get degree request
	degreeRequest, found := k.GetDegreeRequest(ctx, msg.DegreeRequestId)
	if !found {
		return nil, fmt.Errorf("degree request not found")
	}

	// Prepare contract message
	contractMsg := map[string]interface{}{
		"cancel_degree_request": map[string]interface{}{
			"degree_request_id":   msg.DegreeRequestId,
			"cancellation_reason": msg.CancellationReason,
			"creator":             msg.Creator,
		},
	}

	msgBytes, err := json.Marshal(contractMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal contract message: %w", err)
	}

	// Use hardcoded contract address since DegreeRequest doesn't have ContractAddress field
	contractAddr := k.GetDegreeContractAddress(ctx)

	// Convert addresses
	senderAddr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, fmt.Errorf("invalid creator address: %w", err)
	}

	contractAccAddr, err := sdk.AccAddressFromBech32(contractAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid contract address: %w", err)
	}

	// Call CosmWasm contract with correct signature
	_, err = k.wasmKeeper.Execute(ctx, senderAddr, contractAccAddr, msgBytes, []sdk.Coin{})
	if err != nil {
		return nil, fmt.Errorf("contract execution failed: %w", err)
	}

	// Update local state
	degreeRequest.Status = "cancelled"
	k.SetDegreeRequest(ctx, degreeRequest)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDegreeRejected,
			sdk.NewAttribute(types.AttributeKeyDegreeRequestID, msg.DegreeRequestId),
			sdk.NewAttribute(types.AttributeKeyDegreeStatus, "cancelled"),
			sdk.NewAttribute(types.AttributeKeyRejectionReason, msg.CancellationReason),
		),
	)

	return &types.MsgCancelDegreeRequestResponse{
		DegreeRequestId:  msg.DegreeRequestId,
		Status:           "cancelled",
		CancellationDate: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// UpdateParams is deprecated with hardcoded params
func (k msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	return nil, fmt.Errorf("UpdateParams is deprecated - contract addresses are hardcoded")
}
