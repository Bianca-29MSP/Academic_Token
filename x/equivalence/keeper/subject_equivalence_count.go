package keeper

import (
	"context"
	"strconv"

	"academictoken/x/equivalence/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AppendSubjectEquivalence appends a subject equivalence in the store with a new id and update the count
func (k Keeper) AppendSubjectEquivalence(ctx context.Context, equivalence types.SubjectEquivalence) string {
	// Create the equivalence
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	count := k.GetSubjectEquivalenceCount(sdkCtx)

	// Set the ID of the appended value
	equivalence.Index = strconv.FormatUint(count, 10)

	store := k.GetStore(ctx)
	appendedValue := k.cdc.MustMarshal(&equivalence)
	store.Set(types.KeyPrefix(types.SubjectEquivalenceKeyPrefix+equivalence.Index), appendedValue)

	// Update subject equivalence count
	k.SetSubjectEquivalenceCount(sdkCtx, count+1)

	// Set secondary indexes
	k.setEquivalenceIndexes(ctx, equivalence)

	return equivalence.Index
}
