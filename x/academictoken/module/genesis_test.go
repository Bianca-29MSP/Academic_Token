package academictoken_test

import (
	"testing"

	keepertest "academictoken/testutil/keeper"
	"academictoken/testutil/nullify"
	academictoken "academictoken/x/academictoken/module"
	"academictoken/x/academictoken/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.AcademictokenKeeper(t)
	academictoken.InitGenesis(ctx, k, genesisState)
	got := academictoken.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
