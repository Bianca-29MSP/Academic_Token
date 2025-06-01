package keeper

import (
	"context"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"

	"academictoken/x/degree/keeper"
	"academictoken/x/degree/types"
)

// Mock keepers for testing
type mockStudentKeeper struct{}

func (m mockStudentKeeper) GetStudent(ctx sdk.Context, id string) (types.Student, bool) {
	return types.Student{}, true
}
func (m mockStudentKeeper) GetStudentAcademicRecord(ctx sdk.Context, studentId string) (types.AcademicRecord, bool) {
	return types.AcademicRecord{}, true
}
func (m mockStudentKeeper) ValidateStudentExists(ctx sdk.Context, studentId string) error { return nil }
func (m mockStudentKeeper) GetStudentsByInstitution(ctx sdk.Context, institutionId string) []types.Student {
	return []types.Student{}
}
func (m mockStudentKeeper) GetStudentGPA(ctx sdk.Context, studentId string) (string, error) {
	return "3.5", nil
}
func (m mockStudentKeeper) GetStudentTotalCredits(ctx sdk.Context, studentId string) (uint64, error) {
	return 120, nil
}
func (m mockStudentKeeper) GetCompletedSubjects(ctx sdk.Context, studentId string) ([]string, error) {
	return []string{}, nil
}
func (m mockStudentKeeper) GetStudentByIndex(ctx sdk.Context, index string) (types.Student, bool) {
	return types.Student{}, true
}
func (m mockStudentKeeper) GetAcademicTreeByStudent(ctx sdk.Context, studentId string) (interface{}, bool) {
	return struct{}{}, true
}
func (m mockStudentKeeper) GetContractIntegration() interface{} {
	return struct{}{}
}

type mockCurriculumKeeper struct{}

func (m mockCurriculumKeeper) GetCurriculum(ctx sdk.Context, id string) (types.Curriculum, bool) {
	return types.Curriculum{}, true
}
func (m mockCurriculumKeeper) ValidateCurriculumRequirements(ctx sdk.Context, curriculumId string, completedSubjects []string) error {
	return nil
}
func (m mockCurriculumKeeper) GetCurriculumRequiredCredits(ctx sdk.Context, curriculumId string) (uint64, error) {
	return 120, nil
}
func (m mockCurriculumKeeper) GetCurriculumRequiredSubjects(ctx sdk.Context, curriculumId string) ([]string, error) {
	return []string{}, nil
}
func (m mockCurriculumKeeper) GetCurriculumsByInstitution(ctx sdk.Context, institutionId string) []types.Curriculum {
	return []types.Curriculum{}
}
func (m mockCurriculumKeeper) GetCurriculumTree(ctx sdk.Context, id string) (interface{}, bool) {
	return struct{}{}, true
}

type mockAcademicNFTKeeper struct{}

func (m mockAcademicNFTKeeper) MintDegreeNFT(ctx sdk.Context, recipient string, degreeData types.DegreeNFTData) (string, error) {
	return "nft123", nil
}
func (m mockAcademicNFTKeeper) GetNFTByTokenID(ctx sdk.Context, tokenId string) (types.AcademicNFT, bool) {
	return types.AcademicNFT{}, true
}
func (m mockAcademicNFTKeeper) TransferNFT(ctx sdk.Context, from string, to string, tokenId string) error {
	return nil
}
func (m mockAcademicNFTKeeper) BurnNFT(ctx sdk.Context, tokenId string) error { return nil }
func (m mockAcademicNFTKeeper) ValidateNFTOwnership(ctx sdk.Context, owner string, tokenId string) error {
	return nil
}

type mockWasmKeeper struct{}

func (m mockWasmKeeper) Sudo(ctx sdk.Context, contractAddress sdk.AccAddress, msg []byte) ([]byte, error) {
	return []byte("{}"), nil
}
func (m mockWasmKeeper) QuerySmart(ctx context.Context, contractAddr sdk.AccAddress, req []byte) ([]byte, error) {
	return []byte("{}"), nil
}
func (m mockWasmKeeper) Execute(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress, msg []byte, coins sdk.Coins) ([]byte, error) {
	return []byte("{}"), nil
}

type mockInstitutionKeeper struct{}

func (m mockInstitutionKeeper) GetInstitution(ctx sdk.Context, id string) (interface{}, bool) {
	return struct{}{}, true
}
func (m mockInstitutionKeeper) ValidateInstitutionAuthorization(ctx sdk.Context, institutionId string) error {
	return nil
}
func (m mockInstitutionKeeper) IsAuthorizedToIssueDegrees(ctx sdk.Context, institutionId string) bool {
	return true
}

func DegreeKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	types.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)

	paramsSubspace := paramtypes.NewSubspace(cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"DegreeParams",
	)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		mockStudentKeeper{},
		mockCurriculumKeeper{},
		mockAcademicNFTKeeper{},
		mockWasmKeeper{},
		mockInstitutionKeeper{},
		authority.String(),
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	if err := k.SetParams(ctx, types.DefaultParams()); err != nil {
		panic(err)
	}

	return k, ctx
}
