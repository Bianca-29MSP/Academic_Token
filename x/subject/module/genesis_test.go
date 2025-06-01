package subject_test

import (
	"testing"

	keepertest "academictoken/testutil/keeper"
	"academictoken/testutil/nullify"
	subject "academictoken/x/subject/module"
	"academictoken/x/subject/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.SubjectKeeper(t)
	subject.InitGenesis(ctx, k, genesisState)
	got := subject.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
