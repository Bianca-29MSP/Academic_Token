package tokendef

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"academictoken/testutil/sample"
	tokendefsimulation "academictoken/x/tokendef/simulation"
	"academictoken/x/tokendef/types"
)

// avoid unused import issue
var (
	_ = tokendefsimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgCreateTokenDefinition = "op_weight_msg_create_token_definition"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateTokenDefinition int = 100

	opWeightMsgUpdateTokenDefinition = "op_weight_msg_update_token_definition"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateTokenDefinition int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	tokendefGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&tokendefGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateTokenDefinition int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateTokenDefinition, &weightMsgCreateTokenDefinition, nil,
		func(_ *rand.Rand) {
			weightMsgCreateTokenDefinition = defaultWeightMsgCreateTokenDefinition
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateTokenDefinition,
		tokendefsimulation.SimulateMsgCreateTokenDefinition(am.keeper), // Only keeper
	))

	var weightMsgUpdateTokenDefinition int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateTokenDefinition, &weightMsgUpdateTokenDefinition, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateTokenDefinition = defaultWeightMsgUpdateTokenDefinition
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateTokenDefinition,
		tokendefsimulation.SimulateMsgUpdateTokenDefinition(am.accountKeeper, am.bankKeeper, am.keeper), // All three keepers
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreateTokenDefinition,
			defaultWeightMsgCreateTokenDefinition,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				tokendefsimulation.SimulateMsgCreateTokenDefinition(am.keeper) // Only keeper
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateTokenDefinition,
			defaultWeightMsgUpdateTokenDefinition,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				tokendefsimulation.SimulateMsgUpdateTokenDefinition(am.accountKeeper, am.bankKeeper, am.keeper) // All three keepers
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
