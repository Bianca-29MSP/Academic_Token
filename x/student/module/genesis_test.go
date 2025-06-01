package student_test

import (
	"testing"

	keepertest "academictoken/testutil/keeper"
	"academictoken/testutil/nullify"
	student "academictoken/x/student/module"
	"academictoken/x/student/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.StudentKeeper(t)
	student.InitGenesis(ctx, k, genesisState)
	got := student.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
