package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "academictoken/testutil/keeper"
	"academictoken/x/tokendef/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.TokendefKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
