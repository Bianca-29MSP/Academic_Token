// x/tokendef/keeper/params.go
package keeper

import (
	"academictoken/x/tokendef/types"
	"context"
)

// GetHardcodedParams returns fixed parameters for development
func (k Keeper) GetHardcodedParams() types.Params {
	return types.Params{
		IpfsGateway: "http://localhost:5001",
		IpfsEnabled: true,
		Admin:       "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn", // gov module address
	}
}

// GetParams now returns hardcoded params
func (k Keeper) GetParams(ctx context.Context) types.Params {
	return k.GetHardcodedParams()
}

func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	// No-op in hardcode mode - always return success
	return nil
}
