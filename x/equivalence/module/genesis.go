package equivalence

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"academictoken/x/equivalence/keeper"
	"academictoken/x/equivalence/types"
)

// InitGenesis initializes the module's state from a provided genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Note: SubjectEquivalenceList and SubjectEquivalenceCount would be initialized here
	// when they are properly defined in the protobuf genesis.proto file
	
	// For now, we only initialize params
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// Note: SubjectEquivalenceList and SubjectEquivalenceCount would be exported here
	// when they are properly defined in the protobuf genesis.proto file
	// genesis.SubjectEquivalenceList = k.GetAllSubjectEquivalences(ctx)
	// genesis.SubjectEquivalenceCount = k.GetSubjectEquivalenceCount(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
