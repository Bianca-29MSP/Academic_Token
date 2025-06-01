package curriculum

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"academictoken/testutil/sample"
	curriculumsimulation "academictoken/x/curriculum/simulation"
	"academictoken/x/curriculum/types"
)

// avoid unused import issue
var (
	_ = curriculumsimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgCreateCurriculumTree = "op_weight_msg_create_curriculum_tree"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateCurriculumTree int = 100

	opWeightMsgAddSemesterToCurriculum = "op_weight_msg_add_semester_to_curriculum"
	// TODO: Determine the simulation weight value
	defaultWeightMsgAddSemesterToCurriculum int = 100

	opWeightMsgAddElectiveGroup = "op_weight_msg_add_elective_group"
	// TODO: Determine the simulation weight value
	defaultWeightMsgAddElectiveGroup int = 100

	opWeightMsgSetGraduationRequirements = "op_weight_msg_set_graduation_requirements"
	// TODO: Determine the simulation weight value
	defaultWeightMsgSetGraduationRequirements int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	curriculumGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&curriculumGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateCurriculumTree int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateCurriculumTree, &weightMsgCreateCurriculumTree, nil,
		func(_ *rand.Rand) {
			weightMsgCreateCurriculumTree = defaultWeightMsgCreateCurriculumTree
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateCurriculumTree,
		curriculumsimulation.SimulateMsgCreateCurriculumTree(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgAddSemesterToCurriculum int
	simState.AppParams.GetOrGenerate(opWeightMsgAddSemesterToCurriculum, &weightMsgAddSemesterToCurriculum, nil,
		func(_ *rand.Rand) {
			weightMsgAddSemesterToCurriculum = defaultWeightMsgAddSemesterToCurriculum
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddSemesterToCurriculum,
		curriculumsimulation.SimulateMsgAddSemesterToCurriculum(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgAddElectiveGroup int
	simState.AppParams.GetOrGenerate(opWeightMsgAddElectiveGroup, &weightMsgAddElectiveGroup, nil,
		func(_ *rand.Rand) {
			weightMsgAddElectiveGroup = defaultWeightMsgAddElectiveGroup
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddElectiveGroup,
		curriculumsimulation.SimulateMsgAddElectiveGroup(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgSetGraduationRequirements int
	simState.AppParams.GetOrGenerate(opWeightMsgSetGraduationRequirements, &weightMsgSetGraduationRequirements, nil,
		func(_ *rand.Rand) {
			weightMsgSetGraduationRequirements = defaultWeightMsgSetGraduationRequirements
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSetGraduationRequirements,
		curriculumsimulation.SimulateMsgSetGraduationRequirements(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreateCurriculumTree,
			defaultWeightMsgCreateCurriculumTree,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				curriculumsimulation.SimulateMsgCreateCurriculumTree(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgAddSemesterToCurriculum,
			defaultWeightMsgAddSemesterToCurriculum,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				curriculumsimulation.SimulateMsgAddSemesterToCurriculum(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgAddElectiveGroup,
			defaultWeightMsgAddElectiveGroup,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				curriculumsimulation.SimulateMsgAddElectiveGroup(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgSetGraduationRequirements,
			defaultWeightMsgSetGraduationRequirements,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				curriculumsimulation.SimulateMsgSetGraduationRequirements(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
