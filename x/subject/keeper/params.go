// x/subject/keeper/params.go
package keeper

import (
	"academictoken/x/subject/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetHardcodedParams returns fixed parameters for development
func (k Keeper) GetHardcodedParams() types.Params {
	return types.Params{
		IpfsGateway:                   "http://localhost:5001",
		IpfsEnabled:                   true,
		PrerequisiteValidatorContract: "cosmos14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s4hmalr", // hardcode
		EquivalenceValidatorContract:  "cosmos1nc5tatafv6eyq7llkr2gv50ff9e22mnf70qgjlv737ktmt4eswrqrr2r7y",  // hardcode
		Admin:                         "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",                        // gov module address
	}
}

// GetParams now returns hardcoded params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return k.GetHardcodedParams()
}

// SetParams is kept for compatibility but doesn't persist
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	// No-op in hardcode mode
	k.UpdateIPFSClientFromParams(ctx)
}

func (k Keeper) SetIPFSParams(ctx sdk.Context, gateway string, enabled bool) {
	// Update the IPFS client directly
	k.UpdateIPFSClientFromParams(ctx)
}

func (k Keeper) IsIPFSEnabled(ctx sdk.Context) bool {
	return true // always enabled
}

func (k Keeper) GetIPFSGateway(ctx sdk.Context) string {
	return "http://localhost:5001" // hardcoded
}
