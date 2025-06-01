package simulation

import (
	"math/rand"

	"academictoken/x/tokendef/keeper"
	"academictoken/x/tokendef/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgCreateTokenDefinition(
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgCreateTokenDefinition{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the CreateTokenDefinition simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "CreateTokenDefinition simulation not implemented"), nil, nil
	}
}
