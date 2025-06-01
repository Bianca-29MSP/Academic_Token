package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"academictoken/x/equivalence/keeper"
	"academictoken/x/equivalence/types"
)

func SimulateMsgRequestEquivalence(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		msg := &types.MsgRequestEquivalence{
			Creator:            simAccount.Address.String(),
			SourceSubjectId:    simtypes.RandStringOfLength(r, 10),
			TargetInstitution:  simtypes.RandStringOfLength(r, 8),
			TargetSubjectId:    simtypes.RandStringOfLength(r, 10),
			ForceRecalculation: r.Intn(2) == 1,
			ContractAddress:    simtypes.RandStringOfLength(r, 42),
		}

		// TODO: Handling the RequestEquivalence simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "RequestEquivalence simulation not implemented"), nil, nil
	}
}

func SimulateMsgExecuteEquivalenceAnalysis(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		msg := &types.MsgExecuteEquivalenceAnalysis{
			Creator:            simAccount.Address.String(),
			EquivalenceId:      simtypes.RandStringOfLength(r, 10),
			ContractAddress:    simtypes.RandStringOfLength(r, 42),
			AnalysisParameters: simtypes.RandStringOfLength(r, 100),
		}

		// TODO: Handling the ExecuteEquivalenceAnalysis simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "ExecuteEquivalenceAnalysis simulation not implemented"), nil, nil
	}
}

func SimulateMsgBatchRequestEquivalence(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		// Create 1-3 random requests
		numRequests := r.Intn(3) + 1
		requests := make([]*types.EquivalenceRequest, numRequests)

		for i := 0; i < numRequests; i++ {
			requests[i] = &types.EquivalenceRequest{
				SourceSubjectId:   simtypes.RandStringOfLength(r, 10),
				TargetInstitution: simtypes.RandStringOfLength(r, 8),
				TargetSubjectId:   simtypes.RandStringOfLength(r, 10),
			}
		}

		msg := &types.MsgBatchRequestEquivalence{
			Creator:            simAccount.Address.String(),
			Requests:           requests,
			ForceRecalculation: r.Intn(2) == 1,
			ContractAddress:    simtypes.RandStringOfLength(r, 42),
		}

		// TODO: Handling the BatchRequestEquivalence simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "BatchRequestEquivalence simulation not implemented"), nil, nil
	}
}

func SimulateMsgUpdateContractAddress(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		msg := &types.MsgUpdateContractAddress{
			Authority:          simAccount.Address.String(),
			NewContractAddress: simtypes.RandStringOfLength(r, 42),
			ContractVersion:    simtypes.RandStringOfLength(r, 6),
		}

		// TODO: Handling the UpdateContractAddress simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "UpdateContractAddress simulation not implemented"), nil, nil
	}
}

func SimulateMsgReanalyzeEquivalence(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		msg := &types.MsgReanalyzeEquivalence{
			Creator:          simAccount.Address.String(),
			EquivalenceId:    simtypes.RandStringOfLength(r, 10),
			ContractAddress:  simtypes.RandStringOfLength(r, 42),
			ReanalysisReason: simtypes.RandStringOfLength(r, 50),
		}

		// TODO: Handling the ReanalyzeEquivalence simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "ReanalyzeEquivalence simulation not implemented"), nil, nil
	}
}
