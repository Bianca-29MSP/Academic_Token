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
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"

	coursekeeper "academictoken/x/course/keeper"
	coursetypes "academictoken/x/course/types"
	"academictoken/x/curriculum/keeper"
	"academictoken/x/curriculum/types"
	subjectkeeper "academictoken/x/subject/keeper"
	subjecttypes "academictoken/x/subject/types"
)

func CurriculumKeeper(t testing.TB) (keeper.Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	courseStoreKey := storetypes.NewKVStoreKey(coursetypes.StoreKey)
	subjectStoreKey := storetypes.NewKVStoreKey(subjecttypes.StoreKey)
	paramsStoreKey := storetypes.NewKVStoreKey(paramstypes.StoreKey)
	tparamsStoreKey := storetypes.NewTransientStoreKey(paramstypes.TStoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(courseStoreKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(subjectStoreKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(paramsStoreKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(tparamsStoreKey, storetypes.StoreTypeTransient, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())

	// Create params keeper
	paramsKeeper := paramskeeper.NewKeeper(cdc, codec.NewLegacyAmino(), paramsStoreKey, tparamsStoreKey)

	// Use module-specific mocks
	mockCourseInstitutionKeeper := MockCourseInstitutionKeeper{}   // For course module
	mockSubjectInstitutionKeeper := MockSubjectInstitutionKeeper{} // For subject module
	mockSubjectCourseKeeper := MockSubjectCourseKeeper{}           // For subject module
	mockWasmQuerier := MockWasmQuerier{}

	// Create course keeper using course-specific mock
	courseKeeper := coursekeeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(courseStoreKey),
		log.NewNopLogger(),
		authority.String(),
		mockCourseInstitutionKeeper, // Use course-specific mock
		nil,                         // accountKeeper - using nil for tests
		nil,                         // bankKeeper - using nil for tests
	)

	// Create subject keeper with subject-specific mocks
	subjectParamSpace := paramsKeeper.Subspace(subjecttypes.ModuleName)
	if !subjectParamSpace.HasKeyTable() {
		subjectParamSpace = subjectParamSpace.WithKeyTable(subjecttypes.ParamKeyTable())
	}

	subjectKeeper := subjectkeeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(subjectStoreKey),
		log.NewNopLogger(),
		subjectParamSpace,
		mockWasmQuerier,              // WasmQuerier from mocks.go
		mockSubjectInstitutionKeeper, // Use subject-specific mock
		mockSubjectCourseKeeper,      // Use subject-specific mock
		authority.String(),
	)

	// Create curriculum keeper param space
	curriculumParamSpace := paramsKeeper.Subspace(types.ModuleName)
	if !curriculumParamSpace.HasKeyTable() {
		curriculumParamSpace = curriculumParamSpace.WithKeyTable(types.ParamKeyTable())
	}

	k := keeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(storeKey),
		storeKey,
		log.NewNopLogger(),
		curriculumParamSpace,
		&courseKeeper, // Use real keeper for curriculum
		&subjectKeeper,
		authority.String(),
	)

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return k, ctx
}
