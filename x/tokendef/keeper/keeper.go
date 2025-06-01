package keeper

import (
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"academictoken/x/tokendef/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message.
		// Typically, this should be the x/gov module account.
		authority string

		// Cross-module keepers
		subjectKeeper     types.SubjectKeeper
		institutionKeeper types.InstitutionKeeper
		courseKeeper      types.CourseKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	subjectKeeper types.SubjectKeeper,
	institutionKeeper types.InstitutionKeeper,
	courseKeeper types.CourseKeeper,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:               cdc,
		storeService:      storeService,
		authority:         authority,
		logger:            logger,
		subjectKeeper:     subjectKeeper,
		institutionKeeper: institutionKeeper,
		courseKeeper:      courseKeeper,
	}
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetParams get all parameters as types.Params
/*func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	store := k.storeService.OpenKVStore(ctx)
	var params types.Params
	bz, err := store.Get(types.ParamsKey)
	if err != nil {
		return params
	}
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}*/

// SetParams set the params
/*func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	store := k.storeService.OpenKVStore(ctx)
	bz := k.cdc.MustMarshal(&params)
	return store.Set(types.ParamsKey, bz)
}*/

// GetTokenDefinitionCount get the total number of token definitions
func (k Keeper) GetTokenDefinitionCount(ctx sdk.Context) uint64 {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.KeyPrefix(types.TokenDefinitionCountKey))
	if err != nil {
		return 0
	}
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

// SetTokenDefinitionCount set the total number of token definitions
func (k Keeper) SetTokenDefinitionCount(ctx sdk.Context, count uint64) {
	store := k.storeService.OpenKVStore(ctx)
	bz := sdk.Uint64ToBigEndian(count)
	store.Set(types.KeyPrefix(types.TokenDefinitionCountKey), bz)
}

// GetTokenDefinitionByIndex returns a token definition from its index (internal method)
func (k Keeper) GetTokenDefinitionByIndex(ctx sdk.Context, index string) (val types.TokenDefinition, found bool) {
	store := k.storeService.OpenKVStore(ctx)

	b, err := store.Get(types.KeyPrefix(types.TokenDefinitionKeyPrefix + index))
	if err != nil || b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// SetTokenDefinition set a specific token definition in the store from its index
func (k Keeper) SetTokenDefinition(ctx sdk.Context, tokenDefinition types.TokenDefinition) {
	store := k.storeService.OpenKVStore(ctx)
	b := k.cdc.MustMarshal(&tokenDefinition)
	store.Set(types.KeyPrefix(types.TokenDefinitionKeyPrefix+tokenDefinition.Index), b)
}

// RemoveTokenDefinition removes a token definition from the store
func (k Keeper) RemoveTokenDefinition(ctx sdk.Context, index string) {
	store := k.storeService.OpenKVStore(ctx)
	store.Delete(types.KeyPrefix(types.TokenDefinitionKeyPrefix + index))
}

// GetAllTokenDefinition returns all token definitions
func (k Keeper) GetAllTokenDefinition(ctx sdk.Context) (list []types.TokenDefinition) {
	store := k.storeService.OpenKVStore(ctx)
	iterator, err := store.Iterator(types.KeyPrefix(types.TokenDefinitionKeyPrefix), nil)
	if err != nil {
		return list
	}

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.TokenDefinition
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetTokenDefinitionsBySubject returns all token definitions for a specific subject
func (k Keeper) GetTokenDefinitionsBySubject(ctx sdk.Context, subjectId string) []types.TokenDefinition {
	var tokenDefs []types.TokenDefinition
	allTokenDefs := k.GetAllTokenDefinition(ctx)

	for _, tokenDef := range allTokenDefs {
		if tokenDef.SubjectId == subjectId {
			tokenDefs = append(tokenDefs, tokenDef)
		}
	}

	return tokenDefs
}

// GetNextTokenDefinitionIndex generates the next token definition index
func (k Keeper) GetNextTokenDefinitionIndex(ctx sdk.Context) string {
	count := k.GetTokenDefinitionCount(ctx)
	k.SetTokenDefinitionCount(ctx, count+1)
	return fmt.Sprintf("tokendef-%d", count+1)
}

// SetTokenDefinitionBySubjectIndex creates an index for token definitions by subject
func (k Keeper) SetTokenDefinitionBySubjectIndex(ctx sdk.Context, subjectId, tokenDefIndex string) error {
	store := k.storeService.OpenKVStore(ctx)
	key := types.KeyPrefix(types.TokenDefinitionBySubjectKeyPrefix + subjectId + "/" + tokenDefIndex)
	return store.Set(key, []byte(tokenDefIndex))
}

// SetTokenDefinitionByCourseIndex creates an index for token definitions by course
func (k Keeper) SetTokenDefinitionByCourseIndex(ctx sdk.Context, courseId, tokenDefIndex string) error {
	store := k.storeService.OpenKVStore(ctx)
	key := types.KeyPrefix(types.TokenDefinitionByCourseKeyPrefix + courseId + "/" + tokenDefIndex)
	return store.Set(key, []byte(tokenDefIndex))
}

// SetTokenDefinitionByInstitutionIndex creates an index for token definitions by institution
func (k Keeper) SetTokenDefinitionByInstitutionIndex(ctx sdk.Context, institutionId, tokenDefIndex string) error {
	store := k.storeService.OpenKVStore(ctx)
	key := types.KeyPrefix(types.TokenDefinitionByInstitutionKeyPrefix + institutionId + "/" + tokenDefIndex)
	return store.Set(key, []byte(tokenDefIndex))
}
