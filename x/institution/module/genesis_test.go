package institution_test

import (
	"testing"

	keepertest "academictoken/testutil/keeper"
	"academictoken/testutil/nullify"
	institution "academictoken/x/institution/module"
	"academictoken/x/institution/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.InstitutionKeeper(t)
	institution.InitGenesis(ctx, k, genesisState)
	got := institution.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
