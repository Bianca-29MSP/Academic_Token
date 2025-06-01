package student

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"academictoken/testutil/sample"
	studentsimulation "academictoken/x/student/simulation"
	"academictoken/x/student/types"
)

// avoid unused import issue
var (
	_ = studentsimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgRegisterStudent = "op_weight_msg_register_student"
	// TODO: Determine the simulation weight value
	defaultWeightMsgRegisterStudent int = 100

	opWeightMsgCreateEnrollment = "op_weight_msg_create_enrollment"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateEnrollment int = 100

	opWeightMsgUpdateEnrollmentStatus = "op_weight_msg_update_enrollment_status"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateEnrollmentStatus int = 100

	opWeightMsgRequestSubjectEnrollment = "op_weight_msg_request_subject_enrollment"
	// TODO: Determine the simulation weight value
	defaultWeightMsgRequestSubjectEnrollment int = 100

	opWeightMsgUpdateAcademicTree = "op_weight_msg_update_academic_tree"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateAcademicTree int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	studentGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&studentGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgRegisterStudent int
	simState.AppParams.GetOrGenerate(opWeightMsgRegisterStudent, &weightMsgRegisterStudent, nil,
		func(_ *rand.Rand) {
			weightMsgRegisterStudent = defaultWeightMsgRegisterStudent
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRegisterStudent,
		studentsimulation.SimulateMsgRegisterStudent(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCreateEnrollment int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateEnrollment, &weightMsgCreateEnrollment, nil,
		func(_ *rand.Rand) {
			weightMsgCreateEnrollment = defaultWeightMsgCreateEnrollment
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateEnrollment,
		studentsimulation.SimulateMsgCreateEnrollment(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateEnrollmentStatus int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateEnrollmentStatus, &weightMsgUpdateEnrollmentStatus, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateEnrollmentStatus = defaultWeightMsgUpdateEnrollmentStatus
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateEnrollmentStatus,
		studentsimulation.SimulateMsgUpdateEnrollmentStatus(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgRequestSubjectEnrollment int
	simState.AppParams.GetOrGenerate(opWeightMsgRequestSubjectEnrollment, &weightMsgRequestSubjectEnrollment, nil,
		func(_ *rand.Rand) {
			weightMsgRequestSubjectEnrollment = defaultWeightMsgRequestSubjectEnrollment
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRequestSubjectEnrollment,
		studentsimulation.SimulateMsgRequestSubjectEnrollment(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateAcademicTree int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateAcademicTree, &weightMsgUpdateAcademicTree, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateAcademicTree = defaultWeightMsgUpdateAcademicTree
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateAcademicTree,
		studentsimulation.SimulateMsgUpdateAcademicTree(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgRegisterStudent,
			defaultWeightMsgRegisterStudent,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				studentsimulation.SimulateMsgRegisterStudent(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreateEnrollment,
			defaultWeightMsgCreateEnrollment,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				studentsimulation.SimulateMsgCreateEnrollment(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateEnrollmentStatus,
			defaultWeightMsgUpdateEnrollmentStatus,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				studentsimulation.SimulateMsgUpdateEnrollmentStatus(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgRequestSubjectEnrollment,
			defaultWeightMsgRequestSubjectEnrollment,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				studentsimulation.SimulateMsgRequestSubjectEnrollment(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateAcademicTree,
			defaultWeightMsgUpdateAcademicTree,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				studentsimulation.SimulateMsgUpdateAcademicTree(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
