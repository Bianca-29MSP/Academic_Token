package degree

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"academictoken/testutil/sample"
	degreesimulation "academictoken/x/degree/simulation"
	"academictoken/x/degree/types"
)

// avoid unused import issue
var (
	_ = degreesimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgRequestDegree = "op_weight_msg_request_degree_verification"
	// TODO: Determine the simulation weight value
	defaultWeightMsgRequestDegree int = 100

	opWeightMsgIssueDegree = "op_weight_msg_issue_degree"
	// TODO: Determine the simulation weight value
	defaultWeightMsgIssueDegree int = 100

	opWeightMsgValidateDegreeRequirements = "op_weight_msg_verify_degree"
	// TODO: Determine the simulation weight value
	defaultWeightMsgValidateDegreeRequirements int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	degreeGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&degreeGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgRequestDegree int
	simState.AppParams.GetOrGenerate(opWeightMsgRequestDegree, &weightMsgRequestDegree, nil,
		func(_ *rand.Rand) {
			weightMsgRequestDegree = defaultWeightMsgRequestDegree
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRequestDegree,
		degreesimulation.SimulateMsgRequestDegree(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgIssueDegree int
	simState.AppParams.GetOrGenerate(opWeightMsgIssueDegree, &weightMsgIssueDegree, nil,
		func(_ *rand.Rand) {
			weightMsgIssueDegree = defaultWeightMsgIssueDegree
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgIssueDegree,
		degreesimulation.SimulateMsgIssueDegree(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgValidateDegreeRequirements int
	simState.AppParams.GetOrGenerate(opWeightMsgValidateDegreeRequirements, &weightMsgValidateDegreeRequirements, nil,
		func(_ *rand.Rand) {
			weightMsgValidateDegreeRequirements = defaultWeightMsgValidateDegreeRequirements
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgValidateDegreeRequirements,
		degreesimulation.SimulateMsgValidateDegreeRequirements(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgRequestDegree,
			defaultWeightMsgRequestDegree,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				degreesimulation.SimulateMsgRequestDegree(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgIssueDegree,
			defaultWeightMsgIssueDegree,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				degreesimulation.SimulateMsgIssueDegree(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgValidateDegreeRequirements,
			defaultWeightMsgValidateDegreeRequirements,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				degreesimulation.SimulateMsgValidateDegreeRequirements(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
