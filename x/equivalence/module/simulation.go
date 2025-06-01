package equivalence

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"academictoken/testutil/sample"
	equivalencesimulation "academictoken/x/equivalence/simulation"
	"academictoken/x/equivalence/types"
)

// avoid unused import issue
var (
	_ = equivalencesimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgRequestEquivalence = "op_weight_msg_request_equivalence"
	// TODO: Determine the simulation weight value
	defaultWeightMsgRequestEquivalence int = 100

	opWeightMsgExecuteEquivalenceAnalysis = "op_weight_msg_execute_equivalence_analysis"
	// TODO: Determine the simulation weight value
	defaultWeightMsgExecuteEquivalenceAnalysis int = 50

	opWeightMsgBatchRequestEquivalence = "op_weight_msg_batch_request_equivalence"
	// TODO: Determine the simulation weight value
	defaultWeightMsgBatchRequestEquivalence int = 30

	opWeightMsgUpdateContractAddress = "op_weight_msg_update_contract_address"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateContractAddress int = 10

	opWeightMsgReanalyzeEquivalence = "op_weight_msg_reanalyze_equivalence"
	// TODO: Determine the simulation weight value
	defaultWeightMsgReanalyzeEquivalence int = 40

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	equivalenceGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&equivalenceGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgRequestEquivalence int
	simState.AppParams.GetOrGenerate(opWeightMsgRequestEquivalence, &weightMsgRequestEquivalence, nil,
		func(_ *rand.Rand) {
			weightMsgRequestEquivalence = defaultWeightMsgRequestEquivalence
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRequestEquivalence,
		equivalencesimulation.SimulateMsgRequestEquivalence(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgExecuteEquivalenceAnalysis int
	simState.AppParams.GetOrGenerate(opWeightMsgExecuteEquivalenceAnalysis, &weightMsgExecuteEquivalenceAnalysis, nil,
		func(_ *rand.Rand) {
			weightMsgExecuteEquivalenceAnalysis = defaultWeightMsgExecuteEquivalenceAnalysis
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgExecuteEquivalenceAnalysis,
		equivalencesimulation.SimulateMsgExecuteEquivalenceAnalysis(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgBatchRequestEquivalence int
	simState.AppParams.GetOrGenerate(opWeightMsgBatchRequestEquivalence, &weightMsgBatchRequestEquivalence, nil,
		func(_ *rand.Rand) {
			weightMsgBatchRequestEquivalence = defaultWeightMsgBatchRequestEquivalence
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgBatchRequestEquivalence,
		equivalencesimulation.SimulateMsgBatchRequestEquivalence(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateContractAddress int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateContractAddress, &weightMsgUpdateContractAddress, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateContractAddress = defaultWeightMsgUpdateContractAddress
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateContractAddress,
		equivalencesimulation.SimulateMsgUpdateContractAddress(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgReanalyzeEquivalence int
	simState.AppParams.GetOrGenerate(opWeightMsgReanalyzeEquivalence, &weightMsgReanalyzeEquivalence, nil,
		func(_ *rand.Rand) {
			weightMsgReanalyzeEquivalence = defaultWeightMsgReanalyzeEquivalence
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgReanalyzeEquivalence,
		equivalencesimulation.SimulateMsgReanalyzeEquivalence(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgRequestEquivalence,
			defaultWeightMsgRequestEquivalence,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				equivalencesimulation.SimulateMsgRequestEquivalence(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgExecuteEquivalenceAnalysis,
			defaultWeightMsgExecuteEquivalenceAnalysis,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				equivalencesimulation.SimulateMsgExecuteEquivalenceAnalysis(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgBatchRequestEquivalence,
			defaultWeightMsgBatchRequestEquivalence,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				equivalencesimulation.SimulateMsgBatchRequestEquivalence(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgReanalyzeEquivalence,
			defaultWeightMsgReanalyzeEquivalence,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				equivalencesimulation.SimulateMsgReanalyzeEquivalence(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
