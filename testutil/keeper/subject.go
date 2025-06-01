package keeper

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"

	"academictoken/x/subject/keeper"
	"academictoken/x/subject/types"
)

func SubjectKeeper(t testing.TB) (keeper.Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	paramStoreKey := storetypes.NewKVStoreKey("params")
	paramTransientStoreKey := storetypes.NewTransientStoreKey("transient_params")

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(paramStoreKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(paramTransientStoreKey, storetypes.StoreTypeTransient, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)

	// Create a params subspace for the module
	paramSubspace := paramtypes.NewSubspace(cdc,
		codec.NewLegacyAmino(),
		paramStoreKey,
		paramTransientStoreKey,
		"subject",
	)
	paramSubspace = paramSubspace.WithKeyTable(types.ParamKeyTable())

	// Use correct mocks for subject module
	mockWasmQuerier := MockWasmQuerier{}
	mockInstitutionKeeper := MockSubjectInstitutionKeeper{} // Use subject-specific mock
	mockCourseKeeper := MockSubjectCourseKeeper{}           // Use subject-specific mock

	k := keeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(storeKey),
		log.NewNopLogger(),
		paramSubspace,
		mockWasmQuerier,
		mockInstitutionKeeper,
		mockCourseKeeper,
		authority.String(),
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return k, ctx
}
