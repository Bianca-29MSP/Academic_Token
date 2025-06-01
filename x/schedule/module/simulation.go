package schedule

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"academictoken/testutil/sample"
	schedulesimulation "academictoken/x/schedule/simulation"
	"academictoken/x/schedule/types"
)

// avoid unused import issue
var (
	_ = schedulesimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgCreateSubjectRecommendation = "op_weight_msg_create_subject_recommendation"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateSubjectRecommendation int = 100

	opWeightMsgCreateStudyPlan = "op_weight_msg_create_study_plan"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateStudyPlan int = 100

	opWeightMsgAddPlannedSemester = "op_weight_msg_add_planned_semester"
	// TODO: Determine the simulation weight value
	defaultWeightMsgAddPlannedSemester int = 100

	opWeightMsgUpdateStudyPlanStatus = "op_weight_msg_update_study_plan_status"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateStudyPlanStatus int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	scheduleGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&scheduleGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateSubjectRecommendation int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateSubjectRecommendation, &weightMsgCreateSubjectRecommendation, nil,
		func(_ *rand.Rand) {
			weightMsgCreateSubjectRecommendation = defaultWeightMsgCreateSubjectRecommendation
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateSubjectRecommendation,
		schedulesimulation.SimulateMsgCreateSubjectRecommendation(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCreateStudyPlan int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateStudyPlan, &weightMsgCreateStudyPlan, nil,
		func(_ *rand.Rand) {
			weightMsgCreateStudyPlan = defaultWeightMsgCreateStudyPlan
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateStudyPlan,
		schedulesimulation.SimulateMsgCreateStudyPlan(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgAddPlannedSemester int
	simState.AppParams.GetOrGenerate(opWeightMsgAddPlannedSemester, &weightMsgAddPlannedSemester, nil,
		func(_ *rand.Rand) {
			weightMsgAddPlannedSemester = defaultWeightMsgAddPlannedSemester
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddPlannedSemester,
		schedulesimulation.SimulateMsgAddPlannedSemester(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateStudyPlanStatus int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateStudyPlanStatus, &weightMsgUpdateStudyPlanStatus, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateStudyPlanStatus = defaultWeightMsgUpdateStudyPlanStatus
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateStudyPlanStatus,
		schedulesimulation.SimulateMsgUpdateStudyPlanStatus(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreateSubjectRecommendation,
			defaultWeightMsgCreateSubjectRecommendation,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				schedulesimulation.SimulateMsgCreateSubjectRecommendation(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreateStudyPlan,
			defaultWeightMsgCreateStudyPlan,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				schedulesimulation.SimulateMsgCreateStudyPlan(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgAddPlannedSemester,
			defaultWeightMsgAddPlannedSemester,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				schedulesimulation.SimulateMsgAddPlannedSemester(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateStudyPlanStatus,
			defaultWeightMsgUpdateStudyPlanStatus,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				schedulesimulation.SimulateMsgUpdateStudyPlanStatus(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
