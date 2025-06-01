package keeper_test

import (
	"testing"

	keepertest "academictoken/testutil/keeper"
	"academictoken/x/subject/types"

	"github.com/stretchr/testify/require"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx := keepertest.SubjectKeeper(t)
	params := types.DefaultParams()

	keeper.SetParams(ctx, params)

	response, err := keeper.Params(ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
