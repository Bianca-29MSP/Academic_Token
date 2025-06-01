package keeper

import (
	"academictoken/x/degree/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HARDCODED CONTRACT CONFIGURATION
const (
	// Degree validation contract address
	DEGREE_CONTRACT_ADDRESS = "cosmos1degreecontract123456789abcdef"

	// Contract version for compatibility checks
	DEGREE_CONTRACT_VERSION = "1.0.0"

	// Contract execution gas limit
	DEGREE_CONTRACT_GAS_LIMIT = uint64(500000)
)

// GetDegreeContractAddress returns the hardcoded contract address
func (k Keeper) GetDegreeContractAddress(ctx sdk.Context) string {
	return DEGREE_CONTRACT_ADDRESS
}

// GetDegreeContractVersion returns the hardcoded contract version
func (k Keeper) GetDegreeContractVersion(ctx sdk.Context) string {
	return DEGREE_CONTRACT_VERSION
}

// GetDegreeContractGasLimit returns the hardcoded gas limit for contract execution
func (k Keeper) GetDegreeContractGasLimit(ctx sdk.Context) uint64 {
	return DEGREE_CONTRACT_GAS_LIMIT
}

// DEPRECATED: Remove after migration
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.Params{
		ContractAddress: DEGREE_CONTRACT_ADDRESS,
		ContractVersion: DEGREE_CONTRACT_VERSION,
	}
}

// DEPRECATED: Remove after migration
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	// No-op for hardcoded params
	return nil
}
