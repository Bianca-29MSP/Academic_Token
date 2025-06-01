// x/student/keeper/params.go
package keeper

import (
	"academictoken/x/student/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetHardcodedParams returns fixed parameters for development
func (k Keeper) GetHardcodedParams() types.Params {
	return types.Params{
		IpfsGateway:                  "http://localhost:5001",
		IpfsEnabled:                  true,
		Admin:                        "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",                     // gov module address
		PrerequisitesContractAddr:    "cosmos14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s4hmalr", // Prerequisites Contract
		EquivalenceContractAddr:      "cosmos1nc5tatafv6eyq7llkr2gv50ff9e22mnf70qgjlv737ktmt4eswrqrr2r7y", // Equivalence Contract
		AcademicProgressContractAddr: "cosmos1a7x8aj7k38vnj9edrlymkerhrl5d4ud3makmqhx6vt3dhu0d824swy3kus", // Academic Progress Contract
		DegreeContractAddr:           "cosmos15f3n26jmjyh3qfwd7rtmpnr0c6n9qhc9w3j6e2jqz9x8h7f2k6hqeqvhzm", // Degree Contract
		NftMintingContractAddr:       "cosmos1z6a9x2xs5n8j7k3l4f8m9c2v1b6h4g8r7e3w5q9t8y7u6i5o4p3a2s1d0f", // NFT Minting Contract
	}
}

// GetParams now returns hardcoded params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return k.GetHardcodedParams()
}

// SetParams is kept for compatibility but doesn't persist
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	// No-op in hardcode mode - always return success
	return nil
}
