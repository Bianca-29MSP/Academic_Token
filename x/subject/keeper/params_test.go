package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "academictoken/testutil/keeper"
	"academictoken/x/subject/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.SubjectKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)
	require.EqualValues(t, params, k.GetParams(ctx))
}
