package app

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	_ "cosmossdk.io/api/cosmos/tx/config/v1" // import for side-effects
	clienthelpers "cosmossdk.io/client/v2/helpers"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	_ "cosmossdk.io/x/circuit" // import for side-effects
	circuitkeeper "cosmossdk.io/x/circuit/keeper"
	_ "cosmossdk.io/x/evidence" // import for side-effects
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"
	_ "cosmossdk.io/x/feegrant/module" // import for side-effects
	nftkeeper "cosmossdk.io/x/nft/keeper"
	_ "cosmossdk.io/x/nft/module" // import for side-effects
	_ "cosmossdk.io/x/upgrade"    // import for side-effects
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	abci "github.com/cometbft/cometbft/abci/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	_ "github.com/cosmos/cosmos-sdk/x/auth/tx/config" // import for side-effects
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	_ "github.com/cosmos/cosmos-sdk/x/auth/vesting" // import for side-effects
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	_ "github.com/cosmos/cosmos-sdk/x/authz/module" // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/bank"         // import for side-effects
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	_ "github.com/cosmos/cosmos-sdk/x/consensus" // import for side-effects
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	_ "github.com/cosmos/cosmos-sdk/x/crisis" // import for side-effects
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	_ "github.com/cosmos/cosmos-sdk/x/distribution" // import for side-effects
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	groupkeeper "github.com/cosmos/cosmos-sdk/x/group/keeper"
	_ "github.com/cosmos/cosmos-sdk/x/group/module" // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/mint"         // import for side-effects
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	_ "github.com/cosmos/cosmos-sdk/x/params" // import for side-effects
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	_ "github.com/cosmos/cosmos-sdk/x/slashing" // import for side-effects
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	_ "github.com/cosmos/cosmos-sdk/x/staking" // import for side-effects
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	_ "github.com/cosmos/ibc-go/modules/capability" // import for side-effects
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	_ "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts" // import for side-effects
	icacontrollerkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/keeper"
	icahostkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/keeper"
	_ "github.com/cosmos/ibc-go/v8/modules/apps/29-fee" // import for side-effects
	ibcfeekeeper "github.com/cosmos/ibc-go/v8/modules/apps/29-fee/keeper"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"

	academicnftmodulekeeper "academictoken/x/academicnft/keeper"
	academictokenmodulekeeper "academictoken/x/academictoken/keeper"
	coursemodulekeeper "academictoken/x/course/keeper"
	curriculummodulekeeper "academictoken/x/curriculum/keeper"
	degreemodulekeeper "academictoken/x/degree/keeper"
	equivalencemodulekeeper "academictoken/x/equivalence/keeper"
	institutionmodulekeeper "academictoken/x/institution/keeper"
	schedulemodulekeeper "academictoken/x/schedule/keeper"
	studentmodulekeeper "academictoken/x/student/keeper"
	subjectmodulekeeper "academictoken/x/subject/keeper"
	tokendefmodulekeeper "academictoken/x/tokendef/keeper"

	academicnftmoduletypes "academictoken/x/academicnft/types"
	coursemoduletypes "academictoken/x/course/types"
	curriculummoduletypes "academictoken/x/curriculum/types"
	degreemoduletypes "academictoken/x/degree/types"
	equivalencemoduletypes "academictoken/x/equivalence/types"
	institutionmoduletypes "academictoken/x/institution/types"
	schedulemoduletypes "academictoken/x/schedule/types"

	// Manual modules
	academicnft "academictoken/x/academicnft/module"
	course "academictoken/x/course/module"
	curriculum "academictoken/x/curriculum/module"
	degree "academictoken/x/degree/module"
	equivalence "academictoken/x/equivalence/module"
	institution "academictoken/x/institution/module"
	schedule "academictoken/x/schedule/module"
	student "academictoken/x/student/module"
	studentmoduletypes "academictoken/x/student/types"
	subject "academictoken/x/subject/module"
	subjectmoduletypes "academictoken/x/subject/types"
	tokendef "academictoken/x/tokendef/module"
	tokendefmoduletypes "academictoken/x/tokendef/types"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	// Name is the name of the application.
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	// this line is used by starport scaffolding # stargate/app/moduleImport

	"academictoken/docs"
)

const (
	Name = "academictoken"
	// AccountAddressPrefix is the prefix for accounts addresses.
	AccountAddressPrefix = "academic"
	// ChainCoinType is the coin type of the chain.
	ChainCoinType = 118
)

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string
)

var (
	_ runtime.AppI            = (*App)(nil)
	_ servertypes.Application = (*App)(nil)
)

// App extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type App struct {
	*runtime.App
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry codectypes.InterfaceRegistry

	// keepers
	AccountKeeper         authkeeper.AccountKeeper
	BankKeeper            bankkeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	DistrKeeper           distrkeeper.Keeper
	ConsensusParamsKeeper consensuskeeper.Keeper

	SlashingKeeper       slashingkeeper.Keeper
	MintKeeper           mintkeeper.Keeper
	GovKeeper            *govkeeper.Keeper
	CrisisKeeper         *crisiskeeper.Keeper
	UpgradeKeeper        *upgradekeeper.Keeper
	ParamsKeeper         paramskeeper.Keeper
	AuthzKeeper          authzkeeper.Keeper
	EvidenceKeeper       evidencekeeper.Keeper
	FeeGrantKeeper       feegrantkeeper.Keeper
	GroupKeeper          groupkeeper.Keeper
	NFTKeeper            nftkeeper.Keeper
	CircuitBreakerKeeper circuitkeeper.Keeper

	// IBC
	IBCKeeper           *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	CapabilityKeeper    *capabilitykeeper.Keeper
	IBCFeeKeeper        ibcfeekeeper.Keeper
	ICAControllerKeeper icacontrollerkeeper.Keeper
	ICAHostKeeper       icahostkeeper.Keeper
	TransferKeeper      ibctransferkeeper.Keeper

	// Scoped IBC
	ScopedIBCKeeper           capabilitykeeper.ScopedKeeper
	ScopedIBCTransferKeeper   capabilitykeeper.ScopedKeeper
	ScopedICAControllerKeeper capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper       capabilitykeeper.ScopedKeeper
	ScopedKeepers             map[string]capabilitykeeper.ScopedKeeper

	AcademictokenKeeper academictokenmodulekeeper.Keeper

	// CosmWasm
	WasmKeeper       wasmkeeper.Keeper
	ScopedWasmKeeper capabilitykeeper.ScopedKeeper

	// Manual keepers
	InstitutionKeeper institutionmodulekeeper.Keeper
	CourseKeeper      coursemodulekeeper.Keeper
	SubjectKeeper     subjectmodulekeeper.Keeper
	CurriculumKeeper  curriculummodulekeeper.Keeper
	TokendefKeeper    tokendefmodulekeeper.Keeper
	AcademicnftKeeper academicnftmodulekeeper.Keeper
	StudentKeeper     studentmodulekeeper.Keeper
	EquivalenceKeeper equivalencemodulekeeper.Keeper
	DegreeKeeper      degreemodulekeeper.Keeper
	ScheduleKeeper    schedulemodulekeeper.Keeper

	// Manual store keys
	keys map[string]*storetypes.KVStoreKey
	// this line is used by starport scaffolding # stargate/app/keeperDeclaration

	// simulation manager
	sm *module.SimulationManager
}

func init() {
	var err error
	clienthelpers.EnvPrefix = Name
	DefaultNodeHome, err = clienthelpers.GetNodeHomeDirectory("." + Name)
	if err != nil {
		panic(err)
	}
}

// getGovProposalHandlers return the chain proposal handlers.
func getGovProposalHandlers() []govclient.ProposalHandler {
	var govProposalHandlers []govclient.ProposalHandler
	// this line is used by starport scaffolding # stargate/app/govProposalHandlers

	govProposalHandlers = append(govProposalHandlers,
		paramsclient.ProposalHandler,
		// this line is used by starport scaffolding # stargate/app/govProposalHandler
	)

	return govProposalHandlers
}

// WasmQuerier implements the interface required by the subject keeper
type WasmQuerier struct {
	wasmKeeper *wasmkeeper.Keeper
}

// WasmMsgServerAdapter implements the WasmMsgServer interface for Student module
type WasmMsgServerAdapter struct {
	wasmKeeper *wasmkeeper.Keeper
}

// NewWasmQuerier creates a new WasmQuerier
func NewWasmQuerier(wasmKeeper *wasmkeeper.Keeper) *WasmQuerier {
	return &WasmQuerier{
		wasmKeeper: wasmKeeper,
	}
}

// SmartContractState queries a smart contract
func (q *WasmQuerier) SmartContractState(ctx context.Context, req *wasmtypes.QuerySmartContractStateRequest) (*wasmtypes.QuerySmartContractStateResponse, error) {
	if q.wasmKeeper == nil {
		// Return a mock response for development/testing
		return &wasmtypes.QuerySmartContractStateResponse{
			Data: []byte(`{"is_eligible":true,"missing_prerequisites":[]}`),
		}, nil
	}

	// In a real implementation, this would forward to the actual WasmKeeper
	// For now, return a mock response
	return &wasmtypes.QuerySmartContractStateResponse{
		Data: []byte(`{"is_eligible":true,"missing_prerequisites":[]}`),
	}, nil
}

// NewWasmMsgServerAdapter creates a new WasmMsgServerAdapter
func NewWasmMsgServerAdapter(wasmKeeper *wasmkeeper.Keeper) *WasmMsgServerAdapter {
	return &WasmMsgServerAdapter{
		wasmKeeper: wasmKeeper,
	}
}

// ExecuteContract implements the WasmMsgServer interface
func (a *WasmMsgServerAdapter) ExecuteContract(ctx context.Context, req *wasmtypes.MsgExecuteContract) (*wasmtypes.MsgExecuteContractResponse, error) {
	if a.wasmKeeper == nil {
		// Return a mock response for development/testing
		return &wasmtypes.MsgExecuteContractResponse{
			Data: []byte(`{"success":true,"result":"mock_execution"}`),
		}, nil
	}

	// In a real implementation, this would forward to the actual WasmKeeper
	// For now, return a mock response
	return &wasmtypes.MsgExecuteContractResponse{
		Data: []byte(`{"success":true,"result":"mock_execution"}`),
	}, nil
}

// AppConfig returns the default app config.
func AppConfig() depinject.Config {
	return depinject.Configs(
		appConfig,
		// Alternatively, load the app config from a YAML file.
		// appconfig.LoadYAML(AppConfigYAML),
		depinject.Supply(
			// supply custom module basics
			map[string]module.AppModuleBasic{
				genutiltypes.ModuleName:           genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
				govtypes.ModuleName:               gov.NewAppModuleBasic(getGovProposalHandlers()),
				institutionmoduletypes.ModuleName: institution.AppModuleBasic{},
				coursemoduletypes.ModuleName:      course.AppModuleBasic{},
				subjectmoduletypes.ModuleName:     subject.NewAppModuleBasic(nil), // Subject uses different pattern
				curriculummoduletypes.ModuleName:  curriculum.AppModuleBasic{},
				tokendefmoduletypes.ModuleName:    tokendef.NewAppModuleBasic(nil),
				academicnftmoduletypes.ModuleName: academicnft.NewAppModuleBasic(nil),
				studentmoduletypes.ModuleName:     student.NewAppModuleBasic(nil),
				equivalencemoduletypes.ModuleName: equivalence.NewAppModuleBasic(nil),
				degreemoduletypes.ModuleName:      degree.NewAppModuleBasic(nil),
				schedulemoduletypes.ModuleName:    schedule.NewAppModuleBasic(nil),
				// this line is used by starport scaffolding # stargate/appConfig/moduleBasic
			},
		),
	)
}

// ============================================================================
// HELPER FUNCTIONS FOR TYPE CONVERSIONS
// ============================================================================

// convertStringToUint64 safely converts string to uint64
func convertStringToUint64(s string) uint64 {
	if s == "" {
		return 0
	}
	val, err := strconv.ParseUint(strings.TrimSpace(s), 10, 64)
	if err != nil {
		return 0
	}
	return val
}

// convertStringToFloat64 safely converts string to float64
func convertStringToFloat64(s string) float64 {
	if s == "" {
		return 0.0
	}
	val, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0.0
	}
	return val
}

// ============================================================================
// ADAPTERS FOR INTERFACE COMPATIBILITY ISSUES
// ============================================================================

// AccountKeeperAdapter adapts SDK AccountKeeper to degree module interface
type AccountKeeperAdapter struct {
	keeper authkeeper.AccountKeeper
}

func (a AccountKeeperAdapter) GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI {
	// Convert sdk.Context to context.Context for the keeper call
	return a.keeper.GetAccount(sdk.WrapSDKContext(ctx), addr)
}

func (a AccountKeeperAdapter) SetAccount(ctx sdk.Context, acc authtypes.AccountI) {
	a.keeper.SetAccount(sdk.WrapSDKContext(ctx), acc)
}

func (a AccountKeeperAdapter) NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI {
	return a.keeper.NewAccountWithAddress(sdk.WrapSDKContext(ctx), addr)
}

func (a AccountKeeperAdapter) HasAccount(ctx sdk.Context, addr sdk.AccAddress) bool {
	return a.keeper.HasAccount(sdk.WrapSDKContext(ctx), addr)
}

func (a AccountKeeperAdapter) GetModuleAddress(moduleName string) sdk.AccAddress {
	return a.keeper.GetModuleAddress(moduleName)
}

func (a AccountKeeperAdapter) GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI {
	return a.keeper.GetModuleAccount(sdk.WrapSDKContext(ctx), moduleName)
}

// BankKeeperAdapter adapts SDK BankKeeper to degree module interface
type BankKeeperAdapter struct {
	keeper bankkeeper.Keeper
}

func (a BankKeeperAdapter) SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return a.keeper.SpendableCoins(sdk.WrapSDKContext(ctx), addr)
}

func (a BankKeeperAdapter) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	return a.keeper.SendCoins(sdk.WrapSDKContext(ctx), fromAddr, toAddr, amt)
}

func (a BankKeeperAdapter) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	return a.keeper.SendCoinsFromModuleToAccount(sdk.WrapSDKContext(ctx), senderModule, recipientAddr, amt)
}

func (a BankKeeperAdapter) SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	return a.keeper.SendCoinsFromAccountToModule(sdk.WrapSDKContext(ctx), senderAddr, recipientModule, amt)
}

func (a BankKeeperAdapter) SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error {
	return a.keeper.SendCoinsFromModuleToModule(sdk.WrapSDKContext(ctx), senderModule, recipientModule, amt)
}

func (a BankKeeperAdapter) MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	return a.keeper.MintCoins(sdk.WrapSDKContext(ctx), moduleName, amt)
}

func (a BankKeeperAdapter) BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	return a.keeper.BurnCoins(sdk.WrapSDKContext(ctx), moduleName, amt)
}

// ============================================================================
// ADAPTERS FOR TOKENDEF MODULE INTERFACES
// ============================================================================

// SubjectKeeperAdapterForTokenDef adapts subject keeper to tokendef interface
type SubjectKeeperAdapterForTokenDef struct {
	keeper *subjectmodulekeeper.Keeper
}

func (a SubjectKeeperAdapterForTokenDef) HasSubject(ctx sdk.Context, subjectId string) bool {
	return a.keeper.SubjectExists(ctx, subjectId)
}

func (a SubjectKeeperAdapterForTokenDef) GetSubject(ctx sdk.Context, subjectId string) (tokendefmoduletypes.SubjectContent, bool) {
	subject, found := a.keeper.GetSubject(ctx, subjectId)
	if !found {
		return tokendefmoduletypes.SubjectContent{}, false
	}

	// Convert from subject.types to tokendef.types
	return tokendefmoduletypes.SubjectContent{
		Index:         subject.Index,
		SubjectId:     subject.SubjectId,
		Institution:   subject.Institution,
		CourseId:      subject.CourseId,
		Title:         subject.Title,
		Code:          subject.Code,
		WorkloadHours: subject.WorkloadHours,
		Credits:       subject.Credits,
		Description:   subject.Description,
		ContentHash:   subject.ContentHash,
		SubjectType:   subject.SubjectType,
		KnowledgeArea: subject.KnowledgeArea,
		IpfsLink:      subject.IpfsLink,
		Creator:       subject.Creator,
	}, true
}

// InstitutionKeeperAdapterForTokenDef adapts institution keeper to tokendef interface
type InstitutionKeeperAdapterForTokenDef struct {
	keeper *institutionmodulekeeper.Keeper
}

func (a InstitutionKeeperAdapterForTokenDef) GetInstitution(ctx sdk.Context, index string) (tokendefmoduletypes.Institution, bool) {
	institution, found := a.keeper.GetInstitution(ctx, index)
	if !found {
		return tokendefmoduletypes.Institution{}, false
	}

	// Convert from institution.types to tokendef.types
	return tokendefmoduletypes.Institution{
		Index:        institution.Index,
		Address:      institution.Address,
		Name:         institution.Name,
		IsAuthorized: institution.IsAuthorized,
		Creator:      institution.Creator,
	}, true
}

func (a InstitutionKeeperAdapterForTokenDef) IsInstitutionAuthorized(ctx sdk.Context, index string) bool {
	return a.keeper.IsInstitutionAuthorized(ctx, index)
}

// CourseKeeperAdapterForTokenDef adapts course keeper to tokendef interface
type CourseKeeperAdapterForTokenDef struct {
	keeper *coursemodulekeeper.Keeper
}

func (a CourseKeeperAdapterForTokenDef) GetCourse(ctx sdk.Context, index string) (tokendefmoduletypes.Course, bool) {
	course, found := a.keeper.GetCourse(ctx, index)
	if !found {
		return tokendefmoduletypes.Course{}, false
	}

	// Convert from course.types to tokendef.types
	return tokendefmoduletypes.Course{
		Index:        course.Index,
		Institution:  course.Institution,
		Name:         course.Name,
		Code:         course.Code,
		Description:  course.Description,
		TotalCredits: course.TotalCredits,
		DegreeLevel:  course.DegreeLevel,
	}, true
}

func (a CourseKeeperAdapterForTokenDef) HasCourse(ctx sdk.Context, index string) bool {
	return a.keeper.HasCourse(ctx, index)
}

// ============================================================================
// ADAPTERS FOR STUDENT MODULE INTERFACES
// ============================================================================

// InstitutionKeeperAdapterForStudent adapts institution keeper to student interface
type InstitutionKeeperAdapterForStudent struct {
	keeper *institutionmodulekeeper.Keeper
}

func (a InstitutionKeeperAdapterForStudent) GetInstitution(ctx sdk.Context, index string) (studentmoduletypes.Institution, bool) {
	institution, found := a.keeper.GetInstitution(ctx, index)
	if !found {
		return studentmoduletypes.Institution{}, false
	}

	// Convert IsAuthorized from string to bool
	isAuthorized := institution.IsAuthorized == "true"

	// Convert from institution.types to student.types
	return studentmoduletypes.Institution{
		Index:        institution.Index,
		Name:         institution.Name,
		Address:      institution.Address,
		IsAuthorized: isAuthorized,
	}, true
}

func (a InstitutionKeeperAdapterForStudent) IsAuthorized(ctx sdk.Context, institutionIndex string) bool {
	return a.keeper.IsInstitutionAuthorized(ctx, institutionIndex)
}

func (a InstitutionKeeperAdapterForStudent) IsInstitutionAuthorized(ctx sdk.Context, institutionIndex string) bool {
	return a.keeper.IsInstitutionAuthorized(ctx, institutionIndex)
}

// CourseKeeperAdapterForStudent adapts course keeper to student interface
type CourseKeeperAdapterForStudent struct {
	keeper *coursemodulekeeper.Keeper
}

func (a CourseKeeperAdapterForStudent) GetCourse(ctx sdk.Context, index string) (studentmoduletypes.Course, bool) {
	course, found := a.keeper.GetCourse(ctx, index)
	if !found {
		return studentmoduletypes.Course{}, false
	}

	// Convert from course.types to student.types
	return studentmoduletypes.Course{
		Index:        course.Index,
		Institution:  course.Institution,
		Name:         course.Name,
		Code:         course.Code,
		Description:  course.Description,
		TotalCredits: course.TotalCredits,
		DegreeLevel:  course.DegreeLevel,
	}, true
}

func (a CourseKeeperAdapterForStudent) GetCoursesByInstitution(ctx sdk.Context, institutionIndex string) []studentmoduletypes.Course {
	courses := a.keeper.GetCoursesByInstitution(ctx, institutionIndex)

	// Convert to student types
	var result []studentmoduletypes.Course
	for _, course := range courses {
		result = append(result, studentmoduletypes.Course{
			Index:        course.Index,
			Institution:  course.Institution,
			Name:         course.Name,
			Code:         course.Code,
			Description:  course.Description,
			TotalCredits: course.TotalCredits,
			DegreeLevel:  course.DegreeLevel,
		})
	}

	return result
}

// CurriculumKeeperAdapterForStudent adapts curriculum keeper to student interface
type CurriculumKeeperAdapterForStudent struct {
	keeper *curriculummodulekeeper.Keeper
}

func (a CurriculumKeeperAdapterForStudent) GetCurriculumTree(ctx context.Context, index string) (studentmoduletypes.CurriculumTree, bool) {
	// Convert context.Context to sdk.Context
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	curriculum, found := a.keeper.GetCurriculumTree(sdkCtx, index)
	if !found {
		return studentmoduletypes.CurriculumTree{}, false
	}

	// Convert from curriculum.types to student.types with proper conversions
	return studentmoduletypes.CurriculumTree{
		Index:              curriculum.Index,
		CourseId:           curriculum.CourseId,
		Version:            curriculum.Version,
		RequiredSubjects:   curriculum.RequiredSubjects,
		ElectiveMin:        curriculum.ElectiveMin,
		ElectiveSubjects:   curriculum.ElectiveSubjects,
		TotalWorkloadHours: 0, // Default value
		GraduationRequirements: studentmoduletypes.GraduationRequirements{
			TotalCreditsRequired:    convertStringToUint64(curriculum.GraduationRequirements.TotalCreditsRequired),
			MinGPA:                  convertStringToFloat64(curriculum.GraduationRequirements.MinGpa),
			RequiredElectiveCredits: convertStringToUint64(curriculum.GraduationRequirements.RequiredElectiveCredits),
			RequiredActivities:      curriculum.GraduationRequirements.RequiredActivities,
			MinimumTimeYears:        convertStringToFloat64(curriculum.GraduationRequirements.MinimumTimeYears),
			MaximumTimeYears:        convertStringToFloat64(curriculum.GraduationRequirements.MaximumTimeYears),
		},
	}, true
}

func (a CurriculumKeeperAdapterForStudent) GetGraduationRequirements(ctx sdk.Context, curriculumIndex string) (studentmoduletypes.GraduationRequirements, bool) {
	curriculum, found := a.keeper.GetCurriculumTree(ctx, curriculumIndex)
	if !found {
		return studentmoduletypes.GraduationRequirements{}, false
	}

	return studentmoduletypes.GraduationRequirements{
		TotalCreditsRequired:    convertStringToUint64(curriculum.GraduationRequirements.TotalCreditsRequired),
		MinGPA:                  convertStringToFloat64(curriculum.GraduationRequirements.MinGpa), // Convert MinGpa string to float64
		RequiredElectiveCredits: convertStringToUint64(curriculum.GraduationRequirements.RequiredElectiveCredits),
		RequiredActivities:      curriculum.GraduationRequirements.RequiredActivities,
		MinimumTimeYears:        convertStringToFloat64(curriculum.GraduationRequirements.MinimumTimeYears),
		MaximumTimeYears:        convertStringToFloat64(curriculum.GraduationRequirements.MaximumTimeYears),
	}, true
}

func (a CurriculumKeeperAdapterForStudent) GetCurriculumTreesByCourse(ctx context.Context, courseId string) []studentmoduletypes.CurriculumTree {
	// Convert context.Context to sdk.Context
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	curriculums := a.keeper.GetCurriculumTreesByCourse(sdkCtx, courseId)

	// Convert to student types
	var result []studentmoduletypes.CurriculumTree
	for _, curriculum := range curriculums {
		result = append(result, studentmoduletypes.CurriculumTree{
			Index:              curriculum.Index,
			CourseId:           curriculum.CourseId,
			Version:            curriculum.Version,
			RequiredSubjects:   curriculum.RequiredSubjects,
			ElectiveMin:        curriculum.ElectiveMin,
			ElectiveSubjects:   curriculum.ElectiveSubjects,
			TotalWorkloadHours: 0,
			GraduationRequirements: studentmoduletypes.GraduationRequirements{
				TotalCreditsRequired:    convertStringToUint64(curriculum.GraduationRequirements.TotalCreditsRequired),
				MinGPA:                  convertStringToFloat64(curriculum.GraduationRequirements.MinGpa),
				RequiredElectiveCredits: convertStringToUint64(curriculum.GraduationRequirements.RequiredElectiveCredits),
				RequiredActivities:      curriculum.GraduationRequirements.RequiredActivities,
				MinimumTimeYears:        convertStringToFloat64(curriculum.GraduationRequirements.MinimumTimeYears),
				MaximumTimeYears:        convertStringToFloat64(curriculum.GraduationRequirements.MaximumTimeYears),
			},
		})
	}

	return result
}

// SubjectKeeperAdapterForStudent adapts subject keeper to student interface
type SubjectKeeperAdapterForStudent struct {
	keeper *subjectmodulekeeper.Keeper
}

func (a SubjectKeeperAdapterForStudent) GetSubject(ctx sdk.Context, subjectId string) (studentmoduletypes.SubjectContent, bool) {
	subject, found := a.keeper.GetSubject(ctx, subjectId)
	if !found {
		return studentmoduletypes.SubjectContent{}, false
	}

	// Convert from subject.types to student.types
	return studentmoduletypes.SubjectContent{
		Index:         subject.Index,
		SubjectId:     subject.SubjectId,
		Institution:   subject.Institution,
		Title:         subject.Title,
		Code:          subject.Code,
		WorkloadHours: subject.WorkloadHours,
		Credits:       subject.Credits,
		Description:   subject.Description,
		ContentHash:   subject.ContentHash,
		SubjectType:   subject.SubjectType,
		KnowledgeArea: subject.KnowledgeArea,
		IPFSLink:      subject.IpfsLink, // Note: IPFSLink in student types
	}, true
}

func (a SubjectKeeperAdapterForStudent) SubjectExists(ctx sdk.Context, subjectId string) bool {
	return a.keeper.SubjectExists(ctx, subjectId)
}

func (a SubjectKeeperAdapterForStudent) CheckPrerequisites(ctx sdk.Context, studentId string, subjectId string) (bool, []string, error) {
	return a.keeper.CheckPrerequisitesViaContract(ctx, studentId, subjectId)
}

func (a SubjectKeeperAdapterForStudent) GetSubjectsByCourse(ctx sdk.Context, courseId string) []studentmoduletypes.SubjectContent {
	subjects, err := a.keeper.GetSubjectsByCourse(ctx, courseId)
	if err != nil {
		return []studentmoduletypes.SubjectContent{}
	}

	// Convert to student types
	var result []studentmoduletypes.SubjectContent
	for _, subject := range subjects {
		result = append(result, studentmoduletypes.SubjectContent{
			Index:         subject.Index,
			SubjectId:     subject.SubjectId,
			Institution:   subject.Institution,
			Title:         subject.Title,
			Code:          subject.Code,
			WorkloadHours: subject.WorkloadHours,
			Credits:       subject.Credits,
			Description:   subject.Description,
			ContentHash:   subject.ContentHash,
			SubjectType:   subject.SubjectType,
			KnowledgeArea: subject.KnowledgeArea,
			IPFSLink:      subject.IpfsLink,
		})
	}

	return result
}

func (a SubjectKeeperAdapterForStudent) CheckEquivalence(ctx sdk.Context, sourceSubjectId string, targetSubjectId string) (float64, bool) {
	percentage, _, err := a.keeper.CheckEquivalenceViaContract(ctx, sourceSubjectId, targetSubjectId, false)
	if err != nil {
		return 0.0, false
	}
	return float64(percentage) / 100.0, true
}

func (a SubjectKeeperAdapterForStudent) CheckPrerequisitesViaContract(ctx sdk.Context, studentID string, subjectID string) (bool, []string, error) {
	return a.keeper.CheckPrerequisitesViaContract(ctx, studentID, subjectID)
}

func (a SubjectKeeperAdapterForStudent) CheckEquivalenceViaContract(ctx sdk.Context, sourceSubjectID string, targetSubjectID string, forceRecalculate bool) (uint64, string, error) {
	return a.keeper.CheckEquivalenceViaContract(ctx, sourceSubjectID, targetSubjectID, forceRecalculate)
}

// TokenDefKeeperAdapterForStudent adapts tokendef keeper to student interface
type TokenDefKeeperAdapterForStudent struct {
	keeper *tokendefmodulekeeper.Keeper
}

func (a TokenDefKeeperAdapterForStudent) GetTokenDefinition(ctx sdk.Context, index string) (studentmoduletypes.TokenDefinition, bool) {
	tokenDef, found := a.keeper.GetTokenDefinitionByIndex(ctx, index)
	if !found {
		return studentmoduletypes.TokenDefinition{}, false
	}

	// Convert from tokendef.types to student.types
	return studentmoduletypes.TokenDefinition{
		Index:          tokenDef.Index,
		SubjectId:      tokenDef.SubjectId,
		TokenName:      tokenDef.TokenName,
		TokenSymbol:    tokenDef.TokenSymbol,
		Description:    tokenDef.Metadata.Description, // Use Description from Metadata
		TokenType:      tokenDef.TokenType,
		IsTransferable: tokenDef.IsTransferable,
		IsBurnable:     tokenDef.IsBurnable,
		MaxSupply:      tokenDef.MaxSupply,
		ImageUri:       tokenDef.Metadata.ImageUri,
		ContentHash:    tokenDef.ContentHash,
		IPFSLink:       tokenDef.IpfsLink,
	}, true
}

func (a TokenDefKeeperAdapterForStudent) GetTokenDefinitionsBySubject(ctx sdk.Context, subjectId string) []studentmoduletypes.TokenDefinition {
	tokenDefs := a.keeper.GetTokenDefinitionsBySubject(ctx, subjectId)

	// Convert to student types
	var result []studentmoduletypes.TokenDefinition
	for _, tokenDef := range tokenDefs {
		result = append(result, studentmoduletypes.TokenDefinition{
			Index:          tokenDef.Index,
			SubjectId:      tokenDef.SubjectId,
			TokenName:      tokenDef.TokenName,
			TokenSymbol:    tokenDef.TokenSymbol,
			Description:    tokenDef.Metadata.Description, // Use Description from Metadata
			TokenType:      tokenDef.TokenType,
			IsTransferable: tokenDef.IsTransferable,
			IsBurnable:     tokenDef.IsBurnable,
			MaxSupply:      tokenDef.MaxSupply,
			ImageUri:       tokenDef.Metadata.ImageUri,
			ContentHash:    tokenDef.ContentHash,
			IPFSLink:       tokenDef.IpfsLink,
		})
	}

	return result
}

func (a TokenDefKeeperAdapterForStudent) GetTokenDefinitionByIndex(ctx sdk.Context, index string) (studentmoduletypes.TokenDefinition, bool) {
	return a.GetTokenDefinition(ctx, index)
}

// ============================================================================
// ADAPTERS FOR ACADEMICNFT MODULE INTERFACES
// ============================================================================

// StudentKeeperAdapterForAcademicNFT adapts student keeper for AcademicNFT interface
type StudentKeeperAdapterForAcademicNFT struct {
	keeper *studentmodulekeeper.Keeper
}

func (a StudentKeeperAdapterForAcademicNFT) GetStudentByAddress(ctx sdk.Context, address sdk.AccAddress) (academicnftmoduletypes.Student, bool) {
	student, found := a.keeper.GetStudentByAddress(ctx, address.String())
	if !found {
		return academicnftmoduletypes.Student{}, false
	}

	// Convert from student.types to academicnft.types
	// Note: Student in Expected Keepers has no Institution field, so use empty string
	return academicnftmoduletypes.Student{
		Index:       student.Index,
		Address:     student.Address,
		Name:        student.Name,
		Institution: "", // Not present in actual Student type
	}, true
}

func (a StudentKeeperAdapterForAcademicNFT) AddCompletedSubject(ctx sdk.Context, studentAddress string, tokenInstanceID string) error {
	// Get the student first
	student, found := a.keeper.GetStudentByAddress(ctx, studentAddress)
	if !found {
		return fmt.Errorf("student not found: %s", studentAddress)
	}

	// Get the academic tree using the typed method
	academicTree, found := a.keeper.GetAcademicTreeByStudentTyped(ctx, student.Index)
	if !found {
		return fmt.Errorf("academic tree not found for student: %s", student.Index)
	}

	// Add token to completed list if not already there
	for _, completedToken := range academicTree.CompletedTokens {
		if completedToken == tokenInstanceID {
			return nil // Already completed
		}
	}

	// Remove from in-progress if it exists
	var newInProgressTokens []string
	for _, token := range academicTree.InProgressTokens {
		if token != tokenInstanceID {
			newInProgressTokens = append(newInProgressTokens, token)
		}
	}

	// Update the academic tree using the public method
	academicTree.CompletedTokens = append(academicTree.CompletedTokens, tokenInstanceID)
	academicTree.InProgressTokens = newInProgressTokens
	a.keeper.SetStudentAcademicTree(ctx, academicTree)

	return nil
}

func (a StudentKeeperAdapterForAcademicNFT) ValidateStudentEligibility(ctx sdk.Context, studentAddress string, institutionID string) error {
	// Get student using the public method
	student, found := a.keeper.GetStudentByAddress(ctx, studentAddress)
	if !found {
		return fmt.Errorf("student not found: %s", studentAddress)
	}

	// Use the renamed method GetStudentEnrollments
	enrollments, err := a.keeper.GetStudentEnrollments(ctx, student.Index)
	if err != nil {
		return fmt.Errorf("failed to get enrollments: %w", err)
	}

	// Check if student has enrollment in the institution
	for _, enrollment := range enrollments {
		if enrollment.Institution == institutionID && enrollment.Status == "active" {
			return nil // Student is eligible
		}
	}

	return fmt.Errorf("student %s is not eligible to receive tokens from institution %s", studentAddress, institutionID)
}

// AcademicNFTKeeperAdapterForStudent adapts AcademicNFT keeper for Student interface
type AcademicNFTKeeperAdapterForStudent struct {
	keeper *academicnftmodulekeeper.Keeper
}

func (a AcademicNFTKeeperAdapterForStudent) GetTokenInstance(ctx sdk.Context, tokenInstanceId string) (studentmoduletypes.SubjectTokenInstance, bool) {
	tokenInstance, found := a.keeper.GetSubjectTokenInstance(ctx, tokenInstanceId)
	if !found {
		return studentmoduletypes.SubjectTokenInstance{}, false
	}

	// Convert from academicnft.types to student.types
	// Based on proto definition: index, tokenDefId, student, completionDate, grade, issuerInstitution, semester, professorSignature
	return studentmoduletypes.SubjectTokenInstance{
		TokenInstanceId:    tokenInstance.Index, // Use Index as TokenInstanceId
		TokenDefId:         tokenInstance.TokenDefId,
		Student:            tokenInstance.Student,
		CompletionDate:     tokenInstance.CompletionDate,
		Grade:              tokenInstance.Grade,
		IssuerInstitution:  tokenInstance.IssuerInstitution,
		Semester:           tokenInstance.Semester,
		ProfessorSignature: tokenInstance.ProfessorSignature,
		MintedAt:           "",   // Set default empty string since CreatedAt doesn't exist
		IsValid:            true, // Default value
	}, true
}

func (a AcademicNFTKeeperAdapterForStudent) GetStudentTokens(ctx sdk.Context, studentAddress string) []studentmoduletypes.SubjectTokenInstance {
	tokens, err := a.keeper.GetStudentTokenInstances(ctx, studentAddress)
	if err != nil {
		return []studentmoduletypes.SubjectTokenInstance{}
	}

	// Convert to student types
	var result []studentmoduletypes.SubjectTokenInstance
	for _, token := range tokens {
		result = append(result, studentmoduletypes.SubjectTokenInstance{
			TokenInstanceId:    token.Index, // Use Index as TokenInstanceId
			TokenDefId:         token.TokenDefId,
			Student:            token.Student,
			CompletionDate:     token.CompletionDate,
			Grade:              token.Grade,
			IssuerInstitution:  token.IssuerInstitution,
			Semester:           token.Semester,
			ProfessorSignature: token.ProfessorSignature,
			MintedAt:           "",   // Set default empty string since CreatedAt doesn't exist
			IsValid:            true, // Default value
		})
	}

	return result
}

func (a AcademicNFTKeeperAdapterForStudent) MintSubjectToken(ctx sdk.Context, tokenDefId string, student string, completionDate string, grade string, issuerInstitution string, semester string, professorSignature string) (string, error) {
	// Create the token instance directly using keeper functions
	tokenInstanceId := a.keeper.GenerateTokenInstanceID(ctx)

	// Create token instance using only the fields that exist in the proto definition
	tokenInstance := academicnftmoduletypes.SubjectTokenInstance{
		Index:              tokenInstanceId,    // Use Index (field 1 in proto)
		TokenDefId:         tokenDefId,         // field 2
		Student:            student,            // field 3
		CompletionDate:     completionDate,     // field 4
		Grade:              grade,              // field 5
		IssuerInstitution:  issuerInstitution,  // field 6
		Semester:           semester,           // field 7
		ProfessorSignature: professorSignature, // field 8
		// No CreatedAt, Creator, MintedAt, IsValid fields in proto
	}

	// Store the token instance using the existing keeper function
	a.keeper.SetSubjectTokenInstance(ctx, tokenInstance)

	return tokenInstanceId, nil
}

// ============================================================================
// ADAPTERS FOR EQUIVALENCE MODULE INTERFACES
// ============================================================================

// SubjectKeeperAdapterForEquivalence adapts subject keeper to equivalence interface
type SubjectKeeperAdapterForEquivalence struct {
	keeper *subjectmodulekeeper.Keeper
}

func (a SubjectKeeperAdapterForEquivalence) GetSubject(ctx sdk.Context, subjectID string) (equivalencemoduletypes.SubjectContent, bool) {
	subject, found := a.keeper.GetSubject(ctx, subjectID)
	if !found {
		return equivalencemoduletypes.SubjectContent{}, false
	}

	// Convert from subject.types to equivalence.types
	return equivalencemoduletypes.SubjectContent{
		Index:         subject.Index,
		SubjectId:     subject.SubjectId,
		Institution:   subject.Institution,
		Title:         subject.Title,
		Code:          subject.Code,
		WorkloadHours: subject.WorkloadHours,
		Credits:       subject.Credits,
		ContentHash:   subject.ContentHash,
		TopicUnits:    []string{}, // Initialize empty slice since subject doesn't have this field
		Keywords:      []string{}, // Initialize empty slice since subject doesn't have this field
		KnowledgeArea: subject.KnowledgeArea,
		IPFSLink:      subject.IpfsLink,
	}, true
}

func (a SubjectKeeperAdapterForEquivalence) GetSubjectsByInstitution(ctx sdk.Context, institutionID string) []equivalencemoduletypes.SubjectContent {
	subjects, err := a.keeper.GetSubjectsByInstitution(ctx, institutionID)
	if err != nil {
		return []equivalencemoduletypes.SubjectContent{}
	}

	// Convert to equivalence types
	var result []equivalencemoduletypes.SubjectContent
	for _, subject := range subjects {
		result = append(result, equivalencemoduletypes.SubjectContent{
			Index:         subject.Index,
			SubjectId:     subject.SubjectId,
			Institution:   subject.Institution,
			Title:         subject.Title,
			Code:          subject.Code,
			WorkloadHours: subject.WorkloadHours,
			Credits:       subject.Credits,
			ContentHash:   subject.ContentHash,
			TopicUnits:    []string{}, // Initialize empty slice since subject doesn't have this field
			Keywords:      []string{}, // Initialize empty slice since subject doesn't have this field
			KnowledgeArea: subject.KnowledgeArea,
			IPFSLink:      subject.IpfsLink,
		})
	}

	return result
}

func (a SubjectKeeperAdapterForEquivalence) CompareSubjectContent(ctx sdk.Context, sourceSubjectID string, targetSubjectID string) (float64, error) {
	// Use the existing equivalence checking via contract
	percentage, _, err := a.keeper.CheckEquivalenceViaContract(ctx, sourceSubjectID, targetSubjectID, false)
	if err != nil {
		return 0.0, err
	}

	// Convert percentage from uint64 to float64 (divide by 100 since contract returns percentage as integer)
	return float64(percentage) / 100.0, nil
}

// ============================================================================
// ADAPTERS FOR DEGREE MODULE INTERFACES
// ============================================================================

// StudentKeeperAdapterForDegree adapts student keeper to degree interface
type StudentKeeperAdapterForDegree struct {
	keeper *studentmodulekeeper.Keeper
}

func (a StudentKeeperAdapterForDegree) GetStudent(ctx sdk.Context, id string) (degreemoduletypes.Student, bool) {
	student, found := a.keeper.GetStudentByAddress(ctx, id)
	if !found {
		return degreemoduletypes.Student{}, false
	}
	return degreemoduletypes.Student{
		Id:            student.Index,
		InstitutionId: "", // Will need to get from enrollment
		Status:        "active",
	}, true
}

func (a StudentKeeperAdapterForDegree) GetStudentAcademicRecord(ctx sdk.Context, studentId string) (degreemoduletypes.AcademicRecord, bool) {
	// Get academic tree to build record
	academicTree, found := a.keeper.GetAcademicTreeByStudentTyped(ctx, studentId)
	if !found {
		return degreemoduletypes.AcademicRecord{}, false
	}

	return degreemoduletypes.AcademicRecord{
		StudentId:         studentId,
		CompletedCredits:  uint64(len(academicTree.CompletedTokens) * 4), // Estimate 4 credits per token
		GPA:               "8.0",                                         // Default GPA
		CompletedSubjects: academicTree.CompletedTokens,
	}, true
}

func (a StudentKeeperAdapterForDegree) GetStudentsByInstitution(ctx sdk.Context, institutionId string) []degreemoduletypes.Student {
	// This would require querying by institution - simplified implementation
	return []degreemoduletypes.Student{}
}

func (a StudentKeeperAdapterForDegree) ValidateStudentExists(ctx sdk.Context, studentId string) error {
	_, found := a.keeper.GetStudentByAddress(ctx, studentId)
	if !found {
		return fmt.Errorf("student not found: %s", studentId)
	}
	return nil
}

func (a StudentKeeperAdapterForDegree) GetStudentGPA(ctx sdk.Context, studentId string) (string, error) {
	_, found := a.keeper.GetStudentByAddress(ctx, studentId)
	if !found {
		return "", fmt.Errorf("student not found: %s", studentId)
	}
	return "8.0", nil // Default GPA
}

func (a StudentKeeperAdapterForDegree) GetStudentTotalCredits(ctx sdk.Context, studentId string) (uint64, error) {
	academicTree, found := a.keeper.GetAcademicTreeByStudentTyped(ctx, studentId)
	if !found {
		return 0, fmt.Errorf("academic tree not found for student: %s", studentId)
	}
	return uint64(len(academicTree.CompletedTokens) * 4), nil // Estimate 4 credits per token
}

func (a StudentKeeperAdapterForDegree) GetCompletedSubjects(ctx sdk.Context, studentId string) ([]string, error) {
	academicTree, found := a.keeper.GetAcademicTreeByStudentTyped(ctx, studentId)
	if !found {
		return nil, fmt.Errorf("academic tree not found for student: %s", studentId)
	}
	return academicTree.CompletedTokens, nil
}

func (a StudentKeeperAdapterForDegree) GetAcademicTreeByStudent(ctx sdk.Context, studentId string) (interface{}, bool) {
	return a.keeper.GetAcademicTreeByStudent(ctx, studentId)
}

// GetContractIntegration needed by degree module interface
func (a StudentKeeperAdapterForDegree) GetContractIntegration() interface{} {
	return a.keeper.GetContractIntegration()
}

// GetStudentByIndex returns a student by index
func (a StudentKeeperAdapterForDegree) GetStudentByIndex(ctx sdk.Context, index string) (degreemoduletypes.Student, bool) {
	student, found := a.keeper.GetStudentByIndex(ctx, index)
	if !found {
		return degreemoduletypes.Student{}, false
	}
	return degreemoduletypes.Student{
		Id:            student.Index,
		InstitutionId: "", // Will need to get from enrollment
		Status:        "active",
	}, true
}

// CurriculumKeeperAdapterForDegree adapts curriculum keeper to degree interface
type CurriculumKeeperAdapterForDegree struct {
	keeper *curriculummodulekeeper.Keeper
}

func (a CurriculumKeeperAdapterForDegree) GetCurriculum(ctx sdk.Context, id string) (degreemoduletypes.Curriculum, bool) {
	curriculum, found := a.keeper.GetCurriculumTreeSDK(ctx, id)
	if !found {
		return degreemoduletypes.Curriculum{}, false
	}

	totalCredits := uint64(120) // Default
	if curriculum.GraduationRequirements != nil {
		totalCredits = convertStringToUint64(curriculum.GraduationRequirements.TotalCreditsRequired)
	}

	return degreemoduletypes.Curriculum{
		Id:               curriculum.Index,
		InstitutionId:    "", // Get from course
		RequiredCredits:  totalCredits,
		RequiredSubjects: curriculum.RequiredSubjects,
	}, true
}

func (a CurriculumKeeperAdapterForDegree) ValidateCurriculumRequirements(ctx sdk.Context, curriculumId string, completedSubjects []string) error {
	curriculum, found := a.keeper.GetCurriculumTreeSDK(ctx, curriculumId)
	if !found {
		return fmt.Errorf("curriculum not found: %s", curriculumId)
	}

	// Simple validation - check if all required subjects are completed
	requiredMap := make(map[string]bool)
	for _, req := range curriculum.RequiredSubjects {
		requiredMap[req] = false
	}

	for _, completed := range completedSubjects {
		if _, exists := requiredMap[completed]; exists {
			requiredMap[completed] = true
		}
	}

	for subject, completed := range requiredMap {
		if !completed {
			return fmt.Errorf("required subject not completed: %s", subject)
		}
	}

	return nil
}

func (a CurriculumKeeperAdapterForDegree) GetCurriculumRequiredCredits(ctx sdk.Context, curriculumId string) (uint64, error) {
	curriculum, found := a.keeper.GetCurriculumTreeSDK(ctx, curriculumId)
	if !found {
		return 0, fmt.Errorf("curriculum not found: %s", curriculumId)
	}

	if curriculum.GraduationRequirements != nil {
		return convertStringToUint64(curriculum.GraduationRequirements.TotalCreditsRequired), nil
	}

	return 120, nil // Default
}

func (a CurriculumKeeperAdapterForDegree) GetCurriculumRequiredSubjects(ctx sdk.Context, curriculumId string) ([]string, error) {
	curriculum, found := a.keeper.GetCurriculumTreeSDK(ctx, curriculumId)
	if !found {
		return nil, fmt.Errorf("curriculum not found: %s", curriculumId)
	}
	return curriculum.RequiredSubjects, nil
}

func (a CurriculumKeeperAdapterForDegree) GetCurriculumsByInstitution(ctx sdk.Context, institutionId string) []degreemoduletypes.Curriculum {
	// This would require additional indexing - simplified implementation
	return []degreemoduletypes.Curriculum{}
}

func (a CurriculumKeeperAdapterForDegree) GetCurriculumTree(ctx sdk.Context, id string) (interface{}, bool) {
	return a.keeper.GetCurriculumTreeSDK(ctx, id)
}

// AcademicNFTKeeperAdapterForDegree adapts academic NFT keeper to degree interface
type AcademicNFTKeeperAdapterForDegree struct {
	keeper *academicnftmodulekeeper.Keeper
}

func (a AcademicNFTKeeperAdapterForDegree) MintDegreeNFT(ctx sdk.Context, recipient string, degreeData degreemoduletypes.DegreeNFTData) (string, error) {
	// Generate a degree NFT token ID
	tokenId := fmt.Sprintf("degree-%s-%s", degreeData.StudentId, degreeData.CurriculumId)

	// Create degree NFT instance - simplified implementation
	// In a real implementation, this would create a proper degree NFT
	return tokenId, nil
}

func (a AcademicNFTKeeperAdapterForDegree) GetNFTByTokenID(ctx sdk.Context, tokenId string) (degreemoduletypes.AcademicNFT, bool) {
	// Simplified implementation
	return degreemoduletypes.AcademicNFT{
		TokenId:   tokenId,
		Owner:     "student-address",
		Metadata:  "Degree NFT metadata",
		TokenType: "degree",
	}, true
}

func (a AcademicNFTKeeperAdapterForDegree) TransferNFT(ctx sdk.Context, from string, to string, tokenId string) error {
	// Simplified implementation
	return nil
}

func (a AcademicNFTKeeperAdapterForDegree) BurnNFT(ctx sdk.Context, tokenId string) error {
	// Simplified implementation
	return nil
}

func (a AcademicNFTKeeperAdapterForDegree) ValidateNFTOwnership(ctx sdk.Context, owner string, tokenId string) error {
	// Simplified implementation
	return nil
}

// WasmKeeperAdapterForDegree adapts wasm keeper to degree interface
type WasmKeeperAdapterForDegree struct {
	keeper *wasmkeeper.Keeper
}

func (a WasmKeeperAdapterForDegree) Sudo(ctx sdk.Context, contractAddress sdk.AccAddress, msg []byte) ([]byte, error) {
	if a.keeper == nil {
		return []byte(`{"success":true}`), nil // Mock response for development
	}
	return a.keeper.Sudo(ctx, contractAddress, msg)
}

func (a WasmKeeperAdapterForDegree) QuerySmart(ctx context.Context, contractAddr sdk.AccAddress, req []byte) ([]byte, error) {
	if a.keeper == nil {
		return []byte(`{"success":true}`), nil // Mock response for development
	}
	return a.keeper.QuerySmart(ctx, contractAddr, req)
}
func (a WasmKeeperAdapterForDegree) Execute(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress, msg []byte, coins sdk.Coins) ([]byte, error) {
	if a.keeper == nil {
		return []byte(`{"success":true}`), nil // Mock response for development
	}

	// Since the Execute method is not publicly available, we'll use a workaround
	// In a production environment, you might need to use reflection or find an alternative approach
	// For now, return a mock response
	return []byte(`{"success":true,"data":"mock_execution_result"}`), nil
}

// InstitutionKeeperAdapterForDegree adapts institution keeper to degree interface
type InstitutionKeeperAdapterForDegree struct {
	keeper *institutionmodulekeeper.Keeper
}

func (a InstitutionKeeperAdapterForDegree) GetInstitution(ctx sdk.Context, id string) (interface{}, bool) {
	institution, found := a.keeper.GetInstitution(ctx, id)
	if !found {
		return nil, false
	}

	// Return as generic interface{}
	return map[string]interface{}{
		"id":            institution.Index,
		"name":          institution.Name,
		"address":       institution.Address,
		"is_authorized": institution.IsAuthorized == "true",
	}, true
}

func (a InstitutionKeeperAdapterForDegree) ValidateInstitutionAuthorization(ctx sdk.Context, institutionId string) error {
	if !a.keeper.IsInstitutionAuthorized(ctx, institutionId) {
		return fmt.Errorf("institution not authorized: %s", institutionId)
	}
	return nil
}

func (a InstitutionKeeperAdapterForDegree) IsAuthorizedToIssueDegrees(ctx sdk.Context, institutionId string) bool {
	return a.keeper.IsInstitutionAuthorized(ctx, institutionId)
}

// ============================================================================
// ADAPTERS FOR SCHEDULE MODULE INTERFACES
// ============================================================================

// SubjectKeeperAdapterForSchedule adapts subject keeper to schedule interface
type SubjectKeeperAdapterForSchedule struct {
	keeper *subjectmodulekeeper.Keeper
}

func (a SubjectKeeperAdapterForSchedule) GetSubject(ctx sdk.Context, subjectID string) (schedulemoduletypes.SubjectContent, bool) {
	subject, found := a.keeper.GetSubject(ctx, subjectID)
	if !found {
		return schedulemoduletypes.SubjectContent{}, false
	}

	// Convert from subject.types to schedule.types
	return schedulemoduletypes.SubjectContent{
		Index:              subject.Index,
		SubjectId:          subject.SubjectId,
		Institution:        subject.Institution,
		Title:              subject.Title,
		Code:               subject.Code,
		WorkloadHours:      subject.WorkloadHours,
		Credits:            subject.Credits,
		SubjectType:        subject.SubjectType,
		KnowledgeArea:      subject.KnowledgeArea,
		PrerequisiteGroups: []schedulemoduletypes.PrerequisiteGroup{}, // Initialize empty
		DifficultyLevel:    "medium",                                  // Default difficulty
	}, true
}

func (a SubjectKeeperAdapterForSchedule) CheckPrerequisites(ctx sdk.Context, studentID string, subjectID string) (bool, []string, error) {
	return a.keeper.CheckPrerequisitesViaContract(ctx, studentID, subjectID)
}

func (a SubjectKeeperAdapterForSchedule) GetSubjectsByArea(ctx sdk.Context, area string) []schedulemoduletypes.SubjectContent {
	// Simplified implementation - get all subjects and filter by area
	subjects, err := a.keeper.GetAllSubjects(ctx)
	if err != nil {
		return []schedulemoduletypes.SubjectContent{}
	}

	var result []schedulemoduletypes.SubjectContent
	for _, subject := range subjects {
		if subject.KnowledgeArea == area {
			result = append(result, schedulemoduletypes.SubjectContent{
				Index:              subject.Index,
				SubjectId:          subject.SubjectId,
				Institution:        subject.Institution,
				Title:              subject.Title,
				Code:               subject.Code,
				WorkloadHours:      subject.WorkloadHours,
				Credits:            subject.Credits,
				SubjectType:        subject.SubjectType,
				KnowledgeArea:      subject.KnowledgeArea,
				PrerequisiteGroups: []schedulemoduletypes.PrerequisiteGroup{},
				DifficultyLevel:    "medium",
			})
		}
	}

	return result
}

// StudentKeeperAdapterForSchedule adapts student keeper to schedule interface
type StudentKeeperAdapterForSchedule struct {
	keeper *studentmodulekeeper.Keeper
}

func (a StudentKeeperAdapterForSchedule) GetAcademicTree(ctx sdk.Context, studentID string, courseID string) (schedulemoduletypes.StudentAcademicTree, bool) {
	academicTree, found := a.keeper.GetAcademicTreeByStudentTyped(ctx, studentID)
	if !found {
		return schedulemoduletypes.StudentAcademicTree{}, false
	}

	return schedulemoduletypes.StudentAcademicTree{
		Index:               academicTree.Index,
		Student:             academicTree.Student,
		Institution:         academicTree.Institution,
		CourseId:            academicTree.CourseId,
		CurriculumVersion:   academicTree.CurriculumVersion,
		CompletedTokens:     academicTree.CompletedTokens,
		InProgressTokens:    academicTree.InProgressTokens,
		AvailableTokens:     academicTree.AvailableTokens,
		TotalCredits:        uint64(len(academicTree.CompletedTokens) * 4),  // Estimate
		TotalCompletedHours: uint64(len(academicTree.CompletedTokens) * 60), // Estimate
		CoefficientGPA:      8.0,                                            // Default GPA
	}, true
}

func (a StudentKeeperAdapterForSchedule) GetCompletedSubjects(ctx sdk.Context, studentID string) []string {
	academicTree, found := a.keeper.GetAcademicTreeByStudentTyped(ctx, studentID)
	if !found {
		return []string{}
	}
	return academicTree.CompletedTokens
}

func (a StudentKeeperAdapterForSchedule) GetInProgressSubjects(ctx sdk.Context, studentID string) []string {
	academicTree, found := a.keeper.GetAcademicTreeByStudentTyped(ctx, studentID)
	if !found {
		return []string{}
	}
	return academicTree.InProgressTokens
}

// CurriculumKeeperAdapterForSchedule adapts curriculum keeper to schedule interface
type CurriculumKeeperAdapterForSchedule struct {
	keeper *curriculummodulekeeper.Keeper
}

func (a CurriculumKeeperAdapterForSchedule) GetCurriculumTree(ctx sdk.Context, courseID string, version string) (schedulemoduletypes.CurriculumTree, bool) {
	// Use courseID as curriculum index since we don't have direct course->curriculum lookup
	curriculum, found := a.keeper.GetCurriculumTree(ctx, courseID)
	if !found {
		return schedulemoduletypes.CurriculumTree{}, false
	}

	return schedulemoduletypes.CurriculumTree{
		Index:             curriculum.Index,
		CourseId:          curriculum.CourseId,
		Version:           curriculum.Version,
		RequiredSubjects:  curriculum.RequiredSubjects,
		ElectiveSubjects:  curriculum.ElectiveSubjects,
		SemesterStructure: []schedulemoduletypes.CurriculumSemester{}, // Initialize empty
		ElectiveGroups:    []schedulemoduletypes.ElectiveGroup{},      // Initialize empty
	}, true
}

func (a CurriculumKeeperAdapterForSchedule) GetCurrentCurriculumVersion(ctx sdk.Context, courseID string) string {
	curriculum, found := a.keeper.GetCurriculumTree(ctx, courseID)
	if !found {
		return "v1.0" // Default version
	}
	return curriculum.Version
}

// New returns a reference to an initialized App.
func New(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) (*App, error) {
	var (
		app        = &App{ScopedKeepers: make(map[string]capabilitykeeper.ScopedKeeper)}
		appBuilder *runtime.AppBuilder

		// merge the AppConfig and other configuration in one config
		appConfig = depinject.Configs(
			AppConfig(),
			depinject.Supply(
				appOpts, // supply app options
				logger,  // supply logger
				// Supply with IBC keeper getter for the IBC modules with App Wiring.
				// The IBC Keeper cannot be passed because it has not been initiated yet.
				// Passing the getter, the app IBC Keeper will always be accessible.
				// This needs to be removed after IBC supports App Wiring.
				app.GetIBCKeeper,
				app.GetCapabilityScopedKeeper,

				// here alternative options can be supplied to the DI container.
				// those options can be used f.e to override the default behavior of some modules.
				// for instance supplying a custom address codec for not using bech32 addresses.
				// read the depinject documentation and depinject module wiring for more information
				// on available options and how to use them.
			),
		)
	)

	if err := depinject.Inject(appConfig,
		&appBuilder,
		&app.appCodec,
		&app.legacyAmino,
		&app.txConfig,
		&app.interfaceRegistry,
		&app.AccountKeeper,
		&app.BankKeeper,
		&app.StakingKeeper,
		&app.DistrKeeper,
		&app.ConsensusParamsKeeper,
		&app.SlashingKeeper,
		&app.MintKeeper,
		&app.GovKeeper,
		&app.CrisisKeeper,
		&app.UpgradeKeeper,
		&app.ParamsKeeper,
		&app.AuthzKeeper,
		&app.EvidenceKeeper,
		&app.FeeGrantKeeper,
		&app.NFTKeeper,
		&app.GroupKeeper,
		&app.CircuitBreakerKeeper,
		&app.AcademictokenKeeper,
		// Manual keepers NOT injected via DI
		// this line is used by starport scaffolding # stargate/app/keeperDefinition
	); err != nil {
		panic(err)
	}

	// add to default baseapp options
	// enable optimistic execution
	baseAppOptions = append(baseAppOptions, baseapp.SetOptimisticExecution())

	// Create and register manual store keys FIRST
	app.keys = storetypes.NewKVStoreKeys(
		institutionmoduletypes.StoreKey,
		coursemoduletypes.StoreKey,
		subjectmoduletypes.StoreKey,
		curriculummoduletypes.StoreKey,
		tokendefmoduletypes.StoreKey,
		academicnftmoduletypes.StoreKey,
		studentmoduletypes.StoreKey,
		equivalencemoduletypes.StoreKey,
		degreemoduletypes.StoreKey,
		schedulemoduletypes.StoreKey,
	)

	// build app WITH manual store keys
	app.App = appBuilder.Build(db, traceStore, baseAppOptions...)

	// CRITICAL: Mount manual stores AFTER app is built but BEFORE keeper initialization
	for _, key := range app.keys {
		app.MountKVStores(map[string]*storetypes.KVStoreKey{key.Name(): key})
	}

	// Initialize manual keepers AFTER stores are mounted
	app.initializeCustomKeepers()

	// register legacy modules
	if err := app.registerIBCModules(appOpts); err != nil {
		return nil, err
	}

	// Register manual modules with adapters
	institutionModule := institution.NewAppModule(app.appCodec, app.InstitutionKeeper, app.AccountKeeper, app.BankKeeper)
	courseModule := course.NewAppModule(app.appCodec, app.CourseKeeper, app.AccountKeeper, app.BankKeeper)
	subjectModule := subject.NewAppModule(app.appCodec, app.SubjectKeeper, app.AccountKeeper, app.BankKeeper)
	curriculumModule := curriculum.NewAppModule(app.appCodec, app.CurriculumKeeper, app.AccountKeeper, app.BankKeeper)
	tokendefModule := tokendef.NewAppModule(app.appCodec, app.TokendefKeeper, app.AccountKeeper, app.BankKeeper)

	// Create adapters for AcademicNFT module
	studentAdapterForAcademicNFT := StudentKeeperAdapterForAcademicNFT{keeper: &app.StudentKeeper}

	academicnftModule := academicnft.NewAppModule(
		app.appCodec,
		app.AcademicnftKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		&app.TokendefKeeper,
		studentAdapterForAcademicNFT, // Using adapter
		&app.InstitutionKeeper,
	)

	studentModule := student.NewAppModule(
		app.appCodec,
		app.StudentKeeper,
		app.AccountKeeper,
		app.BankKeeper,
	)

	equivalenceModule := equivalence.NewAppModule(
		app.appCodec,
		app.EquivalenceKeeper,
		app.AccountKeeper,
		app.BankKeeper,
	)

	// Create account and bank keeper adapters for Degree module
	accountKeeperAdapter := AccountKeeperAdapter{keeper: app.AccountKeeper}
	bankKeeperAdapter := BankKeeperAdapter{keeper: app.BankKeeper}

	// Create Degree module with adapters
	degreeModule := degree.NewManualAppModule(
		app.appCodec,
		app.DegreeKeeper,
		accountKeeperAdapter, // Use adapter instead of app.AccountKeeper
		bankKeeperAdapter,    // Use adapter instead of app.BankKeeper
	)

	// Create adapters for Schedule module
	subjectAdapterForSchedule := SubjectKeeperAdapterForSchedule{keeper: &app.SubjectKeeper}
	studentAdapterForSchedule := StudentKeeperAdapterForSchedule{keeper: &app.StudentKeeper}
	curriculumAdapterForSchedule := CurriculumKeeperAdapterForSchedule{keeper: &app.CurriculumKeeper}

	scheduleModule := schedule.NewAppModule(
		app.appCodec,
		app.ScheduleKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		subjectAdapterForSchedule,
		studentAdapterForSchedule,
		curriculumAdapterForSchedule,
	)

	if err := app.RegisterModules(
		institutionModule,
		courseModule,
		subjectModule,
		curriculumModule,
		tokendefModule,
		academicnftModule,
		studentModule,
		equivalenceModule,
		degreeModule,
		scheduleModule,
	); err != nil {
		return nil, err
	}

	// register streaming services
	if err := app.RegisterStreamingServices(appOpts, app.kvStoreKeys()); err != nil {
		return nil, err
	}

	/****  Module Options ****/

	app.ModuleManager.RegisterInvariants(app.CrisisKeeper)

	// create the simulation manager and define the order of the modules for deterministic simulations
	overrideModules := map[string]module.AppModuleSimulation{
		authtypes.ModuleName: auth.NewAppModule(app.appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
	}
	app.sm = module.NewSimulationManagerFromAppModules(app.ModuleManager.Modules, overrideModules)
	app.sm.RegisterStoreDecoders()

	// A custom InitChainer sets if extra pre-init-genesis logic is required.
	// This is necessary for manually registered modules that do not support app wiring.
	// Manually set the module version map as shown below.
	// The upgrade module will automatically handle de-duplication of the module version map.
	app.SetInitChainer(func(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
		// CRITICAL: Initialize default parameters for modules FIRST
		// This prevents the "UnmarshalJSON cannot decode empty bytes" error

		// Initialize Subject module params
		defaultSubjectParams := subjectmoduletypes.DefaultParams()
		app.SubjectKeeper.SetParams(ctx, defaultSubjectParams)

		// Initialize Curriculum module params - USE ctx NOT context.Background()
		defaultCurriculumParams := curriculummoduletypes.DefaultParams()
		app.CurriculumKeeper.SetParams(ctx, defaultCurriculumParams)

		// Initialize Tokendef module params - USE ctx NOT context.Background()
		defaultTokendefParams := tokendefmoduletypes.DefaultParams()
		if err := app.TokendefKeeper.SetParams(ctx, defaultTokendefParams); err != nil {
			return nil, err
		}

		// Initialize AcademicNFT module params
		defaultAcademicNFTParams := academicnftmoduletypes.DefaultParams()
		app.AcademicnftKeeper.SetParams(ctx, defaultAcademicNFTParams)

		// Initialize Student module params - USE ctx NOT context.Background()
		defaultStudentParams := studentmoduletypes.DefaultParams()
		app.StudentKeeper.SetParams(ctx, defaultStudentParams)

		// Initialize Equivalence module params
		defaultEquivalenceParams := equivalencemoduletypes.DefaultParams()
		app.EquivalenceKeeper.SetParams(ctx, defaultEquivalenceParams)

		// Initialize Degree module params
		defaultDegreeParams := degreemoduletypes.DefaultParams()
		if err := app.DegreeKeeper.SetParams(ctx, defaultDegreeParams); err != nil {
			return nil, err
		}

		// Initialize Schedule module params
		defaultScheduleParams := schedulemoduletypes.DefaultParams()
		app.ScheduleKeeper.SetParams(ctx, defaultScheduleParams)

		if err := app.UpgradeKeeper.SetModuleVersionMap(ctx, app.ModuleManager.GetVersionMap()); err != nil {
			return nil, err
		}
		return app.App.InitChainer(ctx, req)
	})

	if err := app.Load(loadLatest); err != nil {
		return nil, err
	}
	if err := app.WasmKeeper.InitializePinnedCodes(app.NewUncachedContext(true, tmproto.Header{})); err != nil {
		panic(err)
	}

	return app, nil
}

// initializeCustomKeepers initializes keepers that are not managed by DI
func (app *App) initializeCustomKeepers() {
	// 1. Institution first (no dependencies on other custom modules)
	app.InstitutionKeeper = institutionmodulekeeper.NewKeeper(
		app.appCodec,
		runtime.NewKVStoreService(app.keys[institutionmoduletypes.StoreKey]),
		app.Logger(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		app.AccountKeeper,
		app.BankKeeper,
	)

	// 2. Course second (depends on Institution)
	app.CourseKeeper = coursemodulekeeper.NewKeeper(
		app.appCodec,
		runtime.NewKVStoreService(app.keys[coursemoduletypes.StoreKey]),
		app.Logger(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		&app.InstitutionKeeper,
		app.AccountKeeper,
		app.BankKeeper,
	)

	// 3. Subject third (depends on Course and Institution)
	subjectParamSpace := app.ParamsKeeper.Subspace(subjectmoduletypes.ModuleName)
	if !subjectParamSpace.HasKeyTable() {
		subjectParamSpace = subjectParamSpace.WithKeyTable(subjectmoduletypes.ParamKeyTable())
	}

	wasmQuerier := NewWasmQuerier(&app.WasmKeeper)

	app.SubjectKeeper = subjectmodulekeeper.NewKeeper(
		app.appCodec,
		runtime.NewKVStoreService(app.keys[subjectmoduletypes.StoreKey]),
		app.Logger(),
		subjectParamSpace,
		wasmQuerier,
		&app.InstitutionKeeper,
		&app.CourseKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// 4. Curriculum fourth (depends on Course and Subject)
	curriculumParamSpace := app.ParamsKeeper.Subspace(curriculummoduletypes.ModuleName)
	if !curriculumParamSpace.HasKeyTable() {
		curriculumParamSpace = curriculumParamSpace.WithKeyTable(curriculummoduletypes.ParamKeyTable())
	}

	app.CurriculumKeeper = curriculummodulekeeper.NewKeeper(
		app.appCodec,
		runtime.NewKVStoreService(app.keys[curriculummoduletypes.StoreKey]),
		app.keys[curriculummoduletypes.StoreKey],
		app.Logger(),
		curriculumParamSpace,
		&app.CourseKeeper,
		&app.SubjectKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// 5. TokenDef (depends on Subject, Course and Institution) - USING ADAPTERS
	subjectAdapterForTokenDef := SubjectKeeperAdapterForTokenDef{keeper: &app.SubjectKeeper}
	institutionAdapterForTokenDef := InstitutionKeeperAdapterForTokenDef{keeper: &app.InstitutionKeeper}
	courseAdapterForTokenDef := CourseKeeperAdapterForTokenDef{keeper: &app.CourseKeeper}

	app.TokendefKeeper = tokendefmodulekeeper.NewKeeper(
		app.appCodec,
		runtime.NewKVStoreService(app.keys[tokendefmoduletypes.StoreKey]),
		app.Logger(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		subjectAdapterForTokenDef,
		institutionAdapterForTokenDef,
		courseAdapterForTokenDef,
	)

	// 6. Student (depends on Institution, Course, Curriculum, Subject, TokenDef) - USING ADAPTERS
	institutionAdapterForStudent := InstitutionKeeperAdapterForStudent{keeper: &app.InstitutionKeeper}
	courseAdapterForStudent := CourseKeeperAdapterForStudent{keeper: &app.CourseKeeper}
	curriculumAdapterForStudent := CurriculumKeeperAdapterForStudent{keeper: &app.CurriculumKeeper}
	subjectAdapterForStudent := SubjectKeeperAdapterForStudent{keeper: &app.SubjectKeeper}
	tokendefAdapterForStudent := TokenDefKeeperAdapterForStudent{keeper: &app.TokendefKeeper}

	// Create WasmMsgServer adapter for Student module
	wasmMsgServerForStudent := NewWasmMsgServerAdapter(&app.WasmKeeper)
	wasmQuerierForStudent := NewWasmQuerier(&app.WasmKeeper)

	app.StudentKeeper = studentmodulekeeper.NewKeeper(
		app.appCodec,
		runtime.NewKVStoreService(app.keys[studentmoduletypes.StoreKey]),
		app.Logger(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		app.AccountKeeper,
		app.BankKeeper,
		institutionAdapterForStudent,
		courseAdapterForStudent,
		curriculumAdapterForStudent,
		subjectAdapterForStudent,
		tokendefAdapterForStudent,
		nil,                     // AcademicNFTKeeper - será setado depois para evitar dependência circular
		wasmMsgServerForStudent, // NEW: WasmMsgServer for contract integration
		wasmQuerierForStudent,   // NEW: WasmQuerier for contract integration
	)

	// 7. AcademicNFT (depends on TokenDef, Student, Institution) - USING ADAPTERS
	tokendefAdapterForAcademicNFT := &app.TokendefKeeper // Direct reference since TokenDef uses interface
	studentAdapterForAcademicNFT := StudentKeeperAdapterForAcademicNFT{keeper: &app.StudentKeeper}

	app.AcademicnftKeeper = academicnftmodulekeeper.NewKeeper(
		app.appCodec,
		runtime.NewKVStoreService(app.keys[academicnftmoduletypes.StoreKey]),
		app.Logger(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		tokendefAdapterForAcademicNFT, // TokenDefKeeper
		studentAdapterForAcademicNFT,  // StudentKeeper
		&app.InstitutionKeeper,        // InstitutionKeeper - Direct reference (returns real Institution type)
		app.AccountKeeper,             // AccountKeeper
		app.BankKeeper,                // BankKeeper
	)

	// 8. Equivalence (depends on Subject)
	subjectAdapterForEquivalence := SubjectKeeperAdapterForEquivalence{keeper: &app.SubjectKeeper}

	app.EquivalenceKeeper = equivalencemodulekeeper.NewKeeper(
		app.appCodec,
		runtime.NewKVStoreService(app.keys[equivalencemoduletypes.StoreKey]),
		app.Logger(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Set the subject keeper using the adapter
	app.EquivalenceKeeper.SetSubjectKeeper(subjectAdapterForEquivalence)

	// 9. Degree (depends on Student, Curriculum, AcademicNFT, Wasm, Institution)
	// Create adapters for Degree module
	studentAdapterForDegree := StudentKeeperAdapterForDegree{keeper: &app.StudentKeeper}
	curriculumAdapterForDegree := CurriculumKeeperAdapterForDegree{keeper: &app.CurriculumKeeper}
	academicNFTAdapterForDegree := AcademicNFTKeeperAdapterForDegree{keeper: &app.AcademicnftKeeper}
	wasmAdapterForDegree := WasmKeeperAdapterForDegree{keeper: &app.WasmKeeper}
	institutionAdapterForDegree := InstitutionKeeperAdapterForDegree{keeper: &app.InstitutionKeeper}

	// Initialize param subspace for degree
	degreeParamSpace := app.ParamsKeeper.Subspace(degreemoduletypes.ModuleName)
	if !degreeParamSpace.HasKeyTable() {
		degreeParamSpace = degreeParamSpace.WithKeyTable(degreemoduletypes.ParamKeyTable())
	}

	// Note: Dereference the pointer since NewKeeper returns *Keeper but we need Keeper
	keeperPtr := degreemodulekeeper.NewKeeper(
		app.appCodec,
		app.keys[degreemoduletypes.StoreKey],
		storetypes.NewMemoryStoreKey(degreemoduletypes.MemStoreKey),
		degreeParamSpace,
		studentAdapterForDegree,
		curriculumAdapterForDegree,
		academicNFTAdapterForDegree,
		wasmAdapterForDegree,
		institutionAdapterForDegree,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Dereference the pointer to get the value
	app.DegreeKeeper = *keeperPtr

	// 10. Schedule (depends on Subject, Student, Curriculum) - USING ADAPTERS
	subjectAdapterForSchedule := SubjectKeeperAdapterForSchedule{keeper: &app.SubjectKeeper}
	studentAdapterForSchedule := StudentKeeperAdapterForSchedule{keeper: &app.StudentKeeper}
	curriculumAdapterForSchedule := CurriculumKeeperAdapterForSchedule{keeper: &app.CurriculumKeeper}

	app.ScheduleKeeper = schedulemodulekeeper.NewKeeper(
		app.appCodec,
		runtime.NewKVStoreService(app.keys[schedulemoduletypes.StoreKey]),
		app.Logger(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		subjectAdapterForSchedule,
		studentAdapterForSchedule,
		curriculumAdapterForSchedule,
		app.AccountKeeper,
		app.BankKeeper,
	)
}

// LegacyAmino returns App's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

// AppCodec returns App's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns App's interfaceRegistry.
func (app *App) InterfaceRegistry() codectypes.InterfaceRegistry {
	return app.interfaceRegistry
}

// TxConfig returns App's tx config.
func (app *App) TxConfig() client.TxConfig {
	return app.txConfig
}

// GetKey returns the KVStoreKey for the provided store key.
func (app *App) GetKey(storeKey string) *storetypes.KVStoreKey {
	// First check manual keys
	if key, ok := app.keys[storeKey]; ok {
		return key
	}

	// Then check app store keys
	kvStoreKey, ok := app.UnsafeFindStoreKey(storeKey).(*storetypes.KVStoreKey)
	if !ok {
		return nil
	}
	return kvStoreKey
}

// GetMemKey returns the MemoryStoreKey for the provided store key.
func (app *App) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	key, ok := app.UnsafeFindStoreKey(storeKey).(*storetypes.MemoryStoreKey)
	if !ok {
		return nil
	}

	return key
}

// kvStoreKeys returns all the kv store keys registered inside App.
func (app *App) kvStoreKeys() map[string]*storetypes.KVStoreKey {
	keys := make(map[string]*storetypes.KVStoreKey)
	for _, k := range app.GetStoreKeys() {
		if kv, ok := k.(*storetypes.KVStoreKey); ok {
			keys[kv.Name()] = kv
		}
	}

	// Add manual store keys
	for name, key := range app.keys {
		keys[name] = key
	}

	return keys
}

// GetSubspace returns a param subspace for a given module name.
func (app *App) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// GetIBCKeeper returns the IBC keeper.
func (app *App) GetIBCKeeper() *ibckeeper.Keeper {
	return app.IBCKeeper
}

// GetCapabilityScopedKeeper returns the capability scoped keeper.
func (app *App) GetCapabilityScopedKeeper(moduleName string) capabilitykeeper.ScopedKeeper {
	sk, ok := app.ScopedKeepers[moduleName]
	if !ok {
		sk = app.CapabilityKeeper.ScopeToModule(moduleName)
		app.ScopedKeepers[moduleName] = sk
	}
	return sk
}

// SimulationManager implements the SimulationApp interface.
func (app *App) SimulationManager() *module.SimulationManager {
	return app.sm
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *App) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	app.App.RegisterAPIRoutes(apiSvr, apiConfig)
	// register swagger API in app.go so that other applications can override easily
	if err := server.RegisterSwaggerAPI(apiSvr.ClientCtx, apiSvr.Router, apiConfig.Swagger); err != nil {
		panic(err)
	}

	// register app's OpenAPI routes.
	docs.RegisterOpenAPIService(Name, apiSvr.Router)
}

// GetMaccPerms returns a copy of the module account permissions
//
// NOTE: This is solely to be used for testing purposes.
func GetMaccPerms() map[string][]string {
	dup := make(map[string][]string)
	for _, perms := range moduleAccPerms {
		dup[perms.Account] = perms.Permissions
	}
	return dup
}

// BlockedAddresses returns all the app's blocked account addresses.
func BlockedAddresses() map[string]bool {
	result := make(map[string]bool)
	if len(blockAccAddrs) > 0 {
		for _, addr := range blockAccAddrs {
			result[addr] = true
		}
	} else {
		for addr := range GetMaccPerms() {
			result[addr] = true
		}
	}
	return result
}
