package simulation

import (
	"fmt"
	"math/rand"

	"academictoken/x/schedule/keeper"
	"academictoken/x/schedule/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

func SimulateMsgUpdateStudyPlanStatus(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		
		// Select random account for creator
		creatorAcc, _ := simtypes.RandomAcc(r, accs)
		
		// Generate random study plan ID (in a real scenario, this would come from existing study plans)
		studyPlanId := fmt.Sprintf("sp_simulation_%d", r.Intn(1000))
		
		// Random status
		statuses := []string{
			types.StudyPlanStatusDraft,
			types.StudyPlanStatusActive,
			types.StudyPlanStatusCompleted,
			types.StudyPlanStatusArchived,
		}
		status := statuses[r.Intn(len(statuses))]

		msg := &types.MsgUpdateStudyPlanStatus{
			Creator:     creatorAcc.Address.String(),
			StudyPlanId: studyPlanId,
			Status:      status,
		}

		// Validate the message
		if err := msg.ValidateBasic(); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "invalid message"), nil, err
		}

		// Check account balance
		creatorAccount := ak.GetAccount(ctx, creatorAcc.Address)
		spendable := bk.SpendableCoins(ctx, creatorAccount.GetAddress())

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           nil, // Will be provided by simulation framework
			Cdc:             nil,
			Msg:             msg,
			Context:         ctx,
			SimAccount:      creatorAcc,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		// Deliver the transaction
		opMsg, fops, err := simulation.GenAndDeliverTxWithRandFees(txCtx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "failed to deliver tx"), nil, err
		}

		return opMsg, fops, nil
	}
}
