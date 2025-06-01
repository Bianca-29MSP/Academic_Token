package keeper

import (
	"crypto/sha256"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"academictoken/x/academicnft/types"
)

// ============================================================================
// PASSIVE MODULE AUTHORIZATION METHODS
// ============================================================================

// isAuthorizedContractCall verifies if the caller is an authorized contract
func (k Keeper) isAuthorizedContractCall(ctx sdk.Context, caller string) bool {
	// Get module parameters to check for authorized contract addresses
	params := k.GetParams(ctx)
	
	// For development/testing, we'll allow module addresses and known contract addresses
	// In production, this should be more restrictive
	
	// Check if caller is the module account (for testing)
	// Note: Using authtypes since the accountKeeper interface is limited
	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)
	if caller == moduleAddr.String() {
		k.Logger().Info("Authorized module account call", "caller", caller)
		return true
	}
	
	// Check against authorized contract addresses
	// Since AuthorizedContracts is not in Params, we use a hardcoded list for now
	// TODO: Add AuthorizedContracts to Params in future update
	authorizedContracts := []string{
		// Add actual contract addresses here
		// For now, we'll use mock addresses for testing
		"cosmos1mockcontract1address", // Mock NFT Minting Contract
		"cosmos1mockcontract2address", // Mock Degree Contract  
		"cosmos1mockcontract3address", // Mock Academic Progress Contract
	}
	
	for _, contractAddr := range authorizedContracts {
		if caller == contractAddr {
			k.Logger().Info("Authorized contract call", "caller", caller, "contract", contractAddr)
			return true
		}
	}
	
	// Check if caller is admin (from params)
	if params.Admin != "" && caller == params.Admin {
		k.Logger().Info("Authorized admin call", "caller", caller)
		return true
	}
	
	k.Logger().Error("Unauthorized contract call attempt", 
		"caller", caller,
		"module_addr", moduleAddr.String(),
		"authorized_contracts", authorizedContracts,
		"admin", params.Admin,
	)
	
	return false
}

// verifyContractAuthorization performs cryptographic verification of contract authorization
func (k Keeper) verifyContractAuthorization(ctx sdk.Context, extMsg *types.ExtendedMsgMintSubjectToken) bool {
	// STEP 1: Basic validation - ensure we have the required fields
	if !extMsg.IsContractAuthorized() {
		k.Logger().Error("Missing contract authorization")
		return false
	}
	
	if extMsg.SubjectId == "" {
		k.Logger().Error("Missing subject ID for contract authorization")
		return false
	}
	
	// STEP 2: Construct the data that should have been signed by the contract
	// This should match exactly what the contract used to generate the hash
	authData := k.constructContractAuthorizationData(extMsg)
	
	// STEP 3: Calculate expected hash
	expectedHash := k.calculateContractAuthorizationHash(authData)
	
	// STEP 4: Compare hashes
	isValid := expectedHash == extMsg.ContractAuthorizationHash
	
	if isValid {
		k.Logger().Info("Contract authorization verified successfully",
			"token_def_id", extMsg.TokenDefId,
			"student", extMsg.Student,
			"subject_id", extMsg.SubjectId,
			"authorization_hash", extMsg.ContractAuthorizationHash,
		)
	} else {
		k.Logger().Error("Contract authorization verification failed",
			"token_def_id", extMsg.TokenDefId,
			"student", extMsg.Student,
			"subject_id", extMsg.SubjectId,
			"provided_hash", extMsg.ContractAuthorizationHash,
			"expected_hash", expectedHash,
		)
	}
	
	return isValid
}

// constructContractAuthorizationData creates the canonical data representation for hashing
func (k Keeper) constructContractAuthorizationData(extMsg *types.ExtendedMsgMintSubjectToken) string {
	// Create a deterministic string representation of the authorization data
	// This must match exactly what the contract uses
	return fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s:%s",
		extMsg.TokenDefId,
		extMsg.Student,
		extMsg.SubjectId,
		extMsg.CompletionDate,
		extMsg.Grade,
		extMsg.IssuerInstitution,
		extMsg.Semester,
		extMsg.Creator, // Contract address that authorized this
	)
}

// calculateContractAuthorizationHash calculates the expected authorization hash
func (k Keeper) calculateContractAuthorizationHash(data string) string {
	// Use SHA-256 to create a deterministic hash
	// This should match what the contract uses
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// AddAuthorizedContract adds a contract address to the authorized list (governance)
// Note: Since AuthorizedContracts is not in current Params, this is a placeholder
func (k Keeper) AddAuthorizedContract(ctx sdk.Context, contractAddress string) error {
	// TODO: Implement when AuthorizedContracts is added to Params
	
	k.Logger().Info("Contract authorization request (not implemented)", 
		"contract_address", contractAddress)
	
	// Emit event for tracking
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"contract_authorization_requested",
			sdk.NewAttribute("contract_address", contractAddress),
			sdk.NewAttribute("requested_by", k.authority),
			sdk.NewAttribute("status", "pending_params_update"),
		),
	)
	
	return fmt.Errorf("contract authorization not implemented - AuthorizedContracts field needed in Params")
}

// RemoveAuthorizedContract removes a contract address from the authorized list (governance)
// Note: Since AuthorizedContracts is not in current Params, this is a placeholder
func (k Keeper) RemoveAuthorizedContract(ctx sdk.Context, contractAddress string) error {
	// TODO: Implement when AuthorizedContracts is added to Params
	
	k.Logger().Info("Contract deauthorization request (not implemented)", 
		"contract_address", contractAddress)
	
	// Emit event for tracking
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"contract_deauthorization_requested",
			sdk.NewAttribute("contract_address", contractAddress),
			sdk.NewAttribute("requested_by", k.authority),
			sdk.NewAttribute("status", "pending_params_update"),
		),
	)
	
	return fmt.Errorf("contract deauthorization not implemented - AuthorizedContracts field needed in Params")
}

// GetAuthorizedContracts returns the list of authorized contract addresses
func (k Keeper) GetAuthorizedContracts(ctx sdk.Context) []string {
	// TODO: Return from params when AuthorizedContracts is added
	
	// Return hardcoded list for now
	return []string{
		"cosmos1mockcontract1address",
		"cosmos1mockcontract2address", 
		"cosmos1mockcontract3address",
	}
}

// IsContractAuthorized checks if a specific contract address is authorized
func (k Keeper) IsContractAuthorized(ctx sdk.Context, contractAddress string) bool {
	authorizedContracts := k.GetAuthorizedContracts(ctx)
	
	for _, addr := range authorizedContracts {
		if addr == contractAddress {
			return true
		}
	}
	
	// Also check if it's the admin
	params := k.GetParams(ctx)
	if params.Admin != "" && contractAddress == params.Admin {
		return true
	}
	
	return false
}

// ============================================================================
// PASSIVE MODULE HELPER METHODS
// ============================================================================

// ValidatePassiveModeOperation ensures that the operation is properly authorized for passive mode
func (k Keeper) ValidatePassiveModeOperation(ctx sdk.Context, operation string, caller string) error {
	// Check if the caller is authorized
	if !k.isAuthorizedContractCall(ctx, caller) {
		return types.ErrInvalidContractCaller
	}
	
	k.Logger().Info("Passive mode operation validated",
		"operation", operation,
		"caller", caller,
		"block_height", ctx.BlockHeight(),
	)
	
	return nil
}

// RecordPassiveOperation logs passive mode operations for audit trail
func (k Keeper) RecordPassiveOperation(ctx sdk.Context, operation string, details map[string]string) {
	// Create audit event
	eventAttributes := []sdk.Attribute{
		sdk.NewAttribute("operation", operation),
		sdk.NewAttribute("mode", "passive"),
		sdk.NewAttribute("block_height", fmt.Sprintf("%d", ctx.BlockHeight())),
	}
	
	// Add custom details
	for key, value := range details {
		eventAttributes = append(eventAttributes, sdk.NewAttribute(key, value))
	}
	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"passive_operation_executed",
			eventAttributes...,
		),
	)
	
	k.Logger().Info("Passive operation recorded",
		"operation", operation,
		"details", details,
	)
}

// ExtendTokenInstanceForPassiveMode adds passive mode fields to a token instance
func (k Keeper) ExtendTokenInstanceForPassiveMode(
	tokenInstance types.SubjectTokenInstance,
	subjectId string,
	contractAuthorizationHash string,
	mintedByContract bool,
) types.ExtendedSubjectTokenInstance {
	return types.ExtendedSubjectTokenInstance{
		SubjectTokenInstance:      tokenInstance,
		SubjectId:                subjectId,
		ContractAuthorizationHash: contractAuthorizationHash,
		MintedByContract:         mintedByContract,
		PassiveModeEnabled:       true,
	}
}

// ProcessPassiveModeMessage processes a standard message for passive mode
func (k Keeper) ProcessPassiveModeMessage(ctx sdk.Context, msg *types.MsgMintSubjectToken) (*types.ExtendedMsgMintSubjectToken, error) {
	// Extract passive mode data from context or message metadata
	// This is where contract-provided data would be processed
	
	subjectId, authHash, isPassive := types.ExtractPassiveModeData(ctx, msg)
	
	if !isPassive {
		return nil, fmt.Errorf("message does not contain passive mode data")
	}
	
	// Create extended message
	extMsg := types.ProcessPassiveModeMessage(msg, subjectId, authHash)
	
	// Validate extended message
	if err := types.ValidateContractAuthorization(extMsg); err != nil {
		return nil, fmt.Errorf("invalid contract authorization: %w", err)
	}
	
	// Verify authorization
	if !k.verifyContractAuthorization(ctx, extMsg) {
		return nil, types.ErrInvalidContractAuthorization
	}
	
	return extMsg, nil
}
