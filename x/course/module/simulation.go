package course

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"academictoken/testutil/sample"
	coursesimulation "academictoken/x/course/simulation"
	"academictoken/x/course/types"
)

// avoid unused import issue
var (
	_ = coursesimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgCreateCourse = "op_weight_msg_create_course"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateCourse int = 100

	opWeightMsgUpdateCourse = "op_weight_msg_update_course"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateCourse int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	courseGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&courseGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateCourse int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateCourse, &weightMsgCreateCourse, nil,
		func(_ *rand.Rand) {
			weightMsgCreateCourse = defaultWeightMsgCreateCourse
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateCourse,
		coursesimulation.SimulateMsgCreateCourse(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateCourse int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateCourse, &weightMsgUpdateCourse, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateCourse = defaultWeightMsgUpdateCourse
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateCourse,
		coursesimulation.SimulateMsgUpdateCourse(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreateCourse,
			defaultWeightMsgCreateCourse,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				coursesimulation.SimulateMsgCreateCourse(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateCourse,
			defaultWeightMsgUpdateCourse,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				coursesimulation.SimulateMsgUpdateCourse(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
