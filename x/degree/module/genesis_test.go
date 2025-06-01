package degree_test

import (
	"testing"

	keepertest "academictoken/testutil/keeper"
	"academictoken/testutil/nullify"
	degree "academictoken/x/degree/module"
	"academictoken/x/degree/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.DegreeKeeper(t)
	degree.InitGenesis(ctx, *k, genesisState)
	got := degree.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
