package simulation

import (
	"math/rand"
	//"strconv"
	"time"

	"academictoken/x/schedule/keeper"
	"academictoken/x/schedule/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

func SimulateMsgCreateStudyPlan(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		// Select random accounts for creator and student
		creatorAcc, _ := simtypes.RandomAcc(r, accs)
		studentAcc, _ := simtypes.RandomAcc(r, accs)

		// Generate random semester codes
		semesterCodes := []string{
			"2024-1", "2024-2", "2025-1", "2025-2", "2026-1", "2026-2",
		}

		// Select random number of semesters (1-4)
		numSemesters := r.Intn(4) + 1
		selectedSemesters := make([]string, numSemesters)
		for i := 0; i < numSemesters; i++ {
			selectedSemesters[i] = semesterCodes[r.Intn(len(semesterCodes))]
		}

		// Generate completion target (random date 2-6 years from now)
		yearsToComplete := r.Intn(4) + 2
		completionTarget := time.Now().AddDate(yearsToComplete, 0, 0).Format("2006-01-02")

		// Generate random additional notes
		notes := []string{
			"Fast track completion",
			"Part-time study plan",
			"Regular academic plan",
			"Extended study period",
			"Accelerated program",
		}
		additionalNotes := notes[r.Intn(len(notes))]

		// Create message with correct fields from protobuf
		msg := &types.MsgCreateStudyPlan{
			Creator:          creatorAcc.Address.String(),
			Student:          studentAcc.Address.String(),
			CompletionTarget: completionTarget,
			SemesterCodes:    selectedSemesters,
			AdditionalNotes:  additionalNotes,
			Status:           "draft",
		}

		// Validate the message
		if err := msg.ValidateBasic(); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "invalid message"), nil, err
		}

		// Check account balance
		creatorAccount := ak.GetAccount(ctx, creatorAcc.Address)
		spendable := bk.SpendableCoins(ctx, creatorAccount.GetAddress())

		// Simulate transaction fees
		_, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "unable to generate fees"), nil, err
		}

		// Create operation input with correct fields for newer SDK versions
		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           nil, // TxGen is handled internally in newer versions
			Cdc:             nil,
			Msg:             msg,
			Context:         ctx,
			SimAccount:      creatorAcc,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		// Deliver the transaction with correct API
		opMsg, fops, err := simulation.GenAndDeliverTxWithRandFees(txCtx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "failed to deliver tx"), nil, err
		}

		return opMsg, fops, nil
	}
}

