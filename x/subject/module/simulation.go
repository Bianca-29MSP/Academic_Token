package subject

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"academictoken/testutil/sample"
	subjectsimulation "academictoken/x/subject/simulation"
	"academictoken/x/subject/types"
)

// avoid unused import issue
var (
	_ = subjectsimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgCreateSubject = "op_weight_msg_create_subject"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateSubject int = 100

	opWeightMsgAddPrerequisiteGroup = "op_weight_msg_add_prerequisite_group"
	// TODO: Determine the simulation weight value
	defaultWeightMsgAddPrerequisiteGroup int = 100

	opWeightMsgUpdateSubjectContent = "op_weight_msg_update_subject_content"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateSubjectContent int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	subjectGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&subjectGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateSubject int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateSubject, &weightMsgCreateSubject, nil,
		func(_ *rand.Rand) {
			weightMsgCreateSubject = defaultWeightMsgCreateSubject
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateSubject,
		subjectsimulation.SimulateMsgCreateSubject(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgAddPrerequisiteGroup int
	simState.AppParams.GetOrGenerate(opWeightMsgAddPrerequisiteGroup, &weightMsgAddPrerequisiteGroup, nil,
		func(_ *rand.Rand) {
			weightMsgAddPrerequisiteGroup = defaultWeightMsgAddPrerequisiteGroup
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddPrerequisiteGroup,
		subjectsimulation.SimulateMsgAddPrerequisiteGroup(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateSubjectContent int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateSubjectContent, &weightMsgUpdateSubjectContent, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateSubjectContent = defaultWeightMsgUpdateSubjectContent
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateSubjectContent,
		subjectsimulation.SimulateMsgUpdateSubjectContent(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreateSubject,
			defaultWeightMsgCreateSubject,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				subjectsimulation.SimulateMsgCreateSubject(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgAddPrerequisiteGroup,
			defaultWeightMsgAddPrerequisiteGroup,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				subjectsimulation.SimulateMsgAddPrerequisiteGroup(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateSubjectContent,
			defaultWeightMsgUpdateSubjectContent,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				subjectsimulation.SimulateMsgUpdateSubjectContent(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
