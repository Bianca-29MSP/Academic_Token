package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/runtime"

	"academictoken/x/equivalence/types"
)

// GetParams get all parameters as types.Params with HARDCODED fallback
func (k Keeper) GetParams(ctx context.Context) (params types.Params) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := store.Get(types.ParamsKey)
	
	// If no params in store, return HARDCODED defaults
	if bz == nil {
		k.logger.Info("No params found in store, using hardcoded defaults")
		return types.DefaultParams()
	}

	// Try to unmarshal, fallback to defaults on error
	err := k.cdc.Unmarshal(bz, &params)
	if err != nil {
		k.logger.Error("Failed to unmarshal params, using hardcoded defaults", "error", err)
		return types.DefaultParams()
	}
	
	return params
}

// SetParams set the params with error handling
func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		k.logger.Error("Failed to marshal params", "error", err)
		return err
	}
	
	store.Set(types.ParamsKey, bz)
	k.logger.Info("Parameters updated successfully")
	return nil
}

// GetHardcodedParams returns hardcoded parameters (for emergency use)
func (k Keeper) GetHardcodedParams() types.Params {
	return types.DefaultParams()
}

// ============================================================================
// HARDCODED PARAMETER ACCESS METHODS
// ============================================================================

// GetEquivalenceContractAddress returns the equivalence contract address
func (k Keeper) GetEquivalenceContractAddress(ctx context.Context) string {
	// Since protobuf Params is empty, always return hardcoded value
	address := types.GetHardcodedEquivalenceContractAddress()
	k.logger.Debug("Using hardcoded equivalence contract address", "address", address)
	return address
}

// GetIPFSGateway returns the IPFS gateway URL
func (k Keeper) GetIPFSGateway(ctx context.Context) string {
	// Since protobuf Params is empty, always return hardcoded value
	gateway := types.GetHardcodedIPFSGateway()
	k.logger.Debug("Using hardcoded IPFS gateway", "gateway", gateway)
	return gateway
}

// IsIPFSEnabled returns if IPFS is enabled
func (k Keeper) IsIPFSEnabled(ctx context.Context) bool {
	// Since protobuf Params is empty, always return hardcoded value
	enabled := types.IsHardcodedIPFSEnabled()
	k.logger.Debug("Using hardcoded IPFS enabled status", "enabled", enabled)
	return enabled
}

// GetMinApprovalThreshold returns minimum approval threshold
func (k Keeper) GetMinApprovalThreshold(ctx context.Context) string {
	// Since protobuf Params is empty, always return hardcoded value
	threshold := types.GetHardcodedMinApprovalThreshold()
	k.logger.Debug("Using hardcoded min approval threshold", "threshold", threshold)
	return threshold
}

// GetMaxAnalysisRetries returns max analysis retries
func (k Keeper) GetMaxAnalysisRetries(ctx context.Context) uint64 {
	// Since protobuf Params is empty, always return hardcoded value
	retries := types.GetHardcodedMaxAnalysisRetries()
	k.logger.Debug("Using hardcoded max analysis retries", "retries", retries)
	return retries
}

// GetAnalysisTimeoutSeconds returns analysis timeout in seconds
func (k Keeper) GetAnalysisTimeoutSeconds(ctx context.Context) uint64 {
	// Since protobuf Params is empty, always return hardcoded value
	timeout := types.GetHardcodedAnalysisTimeoutSeconds()
	k.logger.Debug("Using hardcoded analysis timeout", "timeout_seconds", timeout)
	return timeout
}

// IsContractAuthRequired returns if contract authorization is required
func (k Keeper) IsContractAuthRequired(ctx context.Context) bool {
	// Since protobuf Params is empty, always return hardcoded value
	required := types.IsHardcodedContractAuthRequired()
	k.logger.Debug("Using hardcoded contract auth requirement", "required", required)
	return required
}

// GetAdminAddress returns the admin address
func (k Keeper) GetAdminAddress(ctx context.Context) string {
	// Since protobuf Params is empty, always return hardcoded value
	admin := types.GetHardcodedAdmin()
	k.logger.Debug("Using hardcoded admin address", "admin", admin)
	return admin
}

// ============================================================================
// PARAMETER CONFIGURATION HELPERS (FIXED SIGNATURES)
// ============================================================================

// GetAllHardcodedParams returns all hardcoded parameters as a map
func (k Keeper) GetAllHardcodedParams(ctx context.Context) map[string]interface{} {
	config := types.GetDefaultContractConfig()
	k.logger.Debug("Retrieved all hardcoded parameters", "config", config)
	return config
}

// ValidateHardcodedParams validates that all hardcoded parameters are valid - FIXED
func (k Keeper) ValidateHardcodedParams(ctx context.Context) error {
	// Validate contract address
	contractAddr := k.GetEquivalenceContractAddress(ctx)
	if !types.IsValidContractAddress(contractAddr) {
		return types.ErrInvalidContractAddress
	}
	
	// Validate IPFS gateway
	gateway := k.GetIPFSGateway(ctx)
	if gateway == "" {
		return types.ErrIPFSAccessFailed
	}
	
	// Validate threshold
	threshold := k.GetMinApprovalThreshold(ctx)
	if !types.IsValidEquivalencePercent(threshold) {
		return types.ErrInvalidEquivalencePercent
	}
	
	k.logger.Info("All hardcoded parameters validated successfully")
	return nil
}
