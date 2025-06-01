package institution

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"academictoken/testutil/sample"
	institutionsimulation "academictoken/x/institution/simulation"
	"academictoken/x/institution/types"
)

// avoid unused import issue
var (
	_ = institutionsimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgRegisterInstitution = "op_weight_msg_register_institution"
	// TODO: Determine the simulation weight value
	defaultWeightMsgRegisterInstitution int = 100

	opWeightMsgUpdateInstitution = "op_weight_msg_update_institution"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateInstitution int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	institutionGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&institutionGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgRegisterInstitution int
	simState.AppParams.GetOrGenerate(opWeightMsgRegisterInstitution, &weightMsgRegisterInstitution, nil,
		func(_ *rand.Rand) {
			weightMsgRegisterInstitution = defaultWeightMsgRegisterInstitution
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRegisterInstitution,
		institutionsimulation.SimulateMsgRegisterInstitution(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateInstitution int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateInstitution, &weightMsgUpdateInstitution, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateInstitution = defaultWeightMsgUpdateInstitution
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateInstitution,
		institutionsimulation.SimulateMsgUpdateInstitution(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgRegisterInstitution,
			defaultWeightMsgRegisterInstitution,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				institutionsimulation.SimulateMsgRegisterInstitution(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateInstitution,
			defaultWeightMsgUpdateInstitution,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				institutionsimulation.SimulateMsgUpdateInstitution(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
