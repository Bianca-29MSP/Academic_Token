package academicnft

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"academictoken/testutil/sample"
	academicnftsimulation "academictoken/x/academicnft/simulation"
	"academictoken/x/academicnft/types"
)

// avoid unused import issue
var (
	_ = academicnftsimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgMintSubjectToken = "op_weight_msg_mint_subject_token"
	// TODO: Determine the simulation weight value
	defaultWeightMsgMintSubjectToken int = 100

	opWeightMsgVerifyTokenInstance = "op_weight_msg_verify_token_instance"
	// TODO: Determine the simulation weight value
	defaultWeightMsgVerifyTokenInstance int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	academicnftGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&academicnftGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgMintSubjectToken int
	simState.AppParams.GetOrGenerate(opWeightMsgMintSubjectToken, &weightMsgMintSubjectToken, nil,
		func(_ *rand.Rand) {
			weightMsgMintSubjectToken = defaultWeightMsgMintSubjectToken
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgMintSubjectToken,
		academicnftsimulation.SimulateMsgMintSubjectToken(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgVerifyTokenInstance int
	simState.AppParams.GetOrGenerate(opWeightMsgVerifyTokenInstance, &weightMsgVerifyTokenInstance, nil,
		func(_ *rand.Rand) {
			weightMsgVerifyTokenInstance = defaultWeightMsgVerifyTokenInstance
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgVerifyTokenInstance,
		academicnftsimulation.SimulateMsgVerifyTokenInstance(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgMintSubjectToken,
			defaultWeightMsgMintSubjectToken,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				academicnftsimulation.SimulateMsgMintSubjectToken(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgVerifyTokenInstance,
			defaultWeightMsgVerifyTokenInstance,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				academicnftsimulation.SimulateMsgVerifyTokenInstance(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
