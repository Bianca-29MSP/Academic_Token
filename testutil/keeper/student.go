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
	"github.com/stretchr/testify/require"

	"academictoken/x/student/keeper"
	"academictoken/x/student/types"
)

func StudentKeeper(t testing.TB) (keeper.Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)

	// Use mocks from mocks.go
	accountKeeper := MockAccountKeeper{}
	bankKeeper := MockBankKeeper{}
	institutionKeeper := MockStudentInstitutionKeeper{}
	courseKeeper := MockStudentCourseKeeper{}
	curriculumKeeper := MockStudentCurriculumKeeper{}
	subjectKeeper := MockStudentSubjectKeeper{}
	tokenDefKeeper := MockStudentTokenDefKeeper{}
	academicNFTKeeper := MockStudentAcademicNFTKeeper{}
	wasmMsgServer := MockWasmMsgServer{}
	wasmQuerier := MockWasmQuerier{}

	k := keeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(storeKey),
		log.NewNopLogger(),
		authority.String(),
		accountKeeper,
		bankKeeper,
		institutionKeeper,
		courseKeeper,
		curriculumKeeper,
		subjectKeeper,
		tokenDefKeeper,
		academicNFTKeeper,
		wasmMsgServer,
		wasmQuerier,
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	if err := k.SetParams(ctx, types.DefaultParams()); err != nil {
		panic(err)
	}

	return k, ctx
}
