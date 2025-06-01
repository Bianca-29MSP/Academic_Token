package schedule

import (
	"context"
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	// this line is used by starport scaffolding # 1

	modulev1 "academictoken/api/academictoken/schedule/module"
	"academictoken/x/schedule/keeper"

	//"academictoken/x/schedule/simulation"
	"academictoken/x/schedule/types"
)

var (
	_ module.AppModuleBasic      = (*AppModule)(nil)
	_ module.AppModuleSimulation = (*AppModule)(nil)
	_ module.HasGenesis          = (*AppModule)(nil)
	_ module.HasInvariants       = (*AppModule)(nil)
	_ module.HasConsensusVersion = (*AppModule)(nil)

	_ appmodule.AppModule       = (*AppModule)(nil)
	_ appmodule.HasBeginBlocker = (*AppModule)(nil)
	_ appmodule.HasEndBlocker   = (*AppModule)(nil)
)

// ----------------------------------------------------------------------------
// AppModuleBasic
// ----------------------------------------------------------------------------

// AppModuleBasic implements the AppModuleBasic interface that defines the
// independent methods a Cosmos SDK module needs to implement.
type AppModuleBasic struct {
	cdc codec.BinaryCodec
}

func NewAppModuleBasic(cdc codec.BinaryCodec) AppModuleBasic {
	return AppModuleBasic{cdc: cdc}
}

// Name returns the name of the module as a string.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers the amino codec for the module, which is used
// to marshal and unmarshal structs to/from []byte in order to persist them in the module's KVStore.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {}

// RegisterInterfaces registers a module's interface types and their concrete implementations as proto.Message.
func (a AppModuleBasic) RegisterInterfaces(reg cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(reg)
}

// DefaultGenesis returns a default GenesisState for the module, marshalled to json.RawMessage.
// The default GenesisState need to be defined by the module developer and is primarily used for testing.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesis())
}

// ValidateGenesis used to validate the GenesisState, given in its json.RawMessage form.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return genState.Validate()
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule implements the AppModule interface that defines the inter-dependent methods that modules need to implement
type AppModule struct {
	AppModuleBasic

	keeper           keeper.Keeper
	accountKeeper    types.AccountKeeper
	bankKeeper       types.BankKeeper
	subjectKeeper    types.SubjectKeeper
	studentKeeper    types.StudentKeeper
	curriculumKeeper types.CurriculumKeeper
}

func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	subjectKeeper types.SubjectKeeper,
	studentKeeper types.StudentKeeper,
	curriculumKeeper types.CurriculumKeeper,
) AppModule {
	return AppModule{
		AppModuleBasic:   NewAppModuleBasic(cdc),
		keeper:           keeper,
		accountKeeper:    accountKeeper,
		bankKeeper:       bankKeeper,
		subjectKeeper:    subjectKeeper,
		studentKeeper:    studentKeeper,
		curriculumKeeper: curriculumKeeper,
	}
}

// RegisterServices registers a gRPC query service to respond to the module-specific gRPC queries
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)
}

// RegisterInvariants registers the invariants of the module. If an invariant deviates from its predicted value, the InvariantRegistry triggers appropriate logic (most often the chain will be halted)
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// InitGenesis performs the module's genesis initialization. It returns no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) {
	var genState types.GenesisState
	// Initialize global index to index in genesis state
	cdc.MustUnmarshalJSON(gs, &genState)

	InitGenesis(ctx, am.keeper, genState)
}

// ExportGenesis returns the module's exported genesis state as raw JSON bytes.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genState := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(genState)
}

// ConsensusVersion is a sequence number for state-breaking change of the module.
// It should be incremented on each consensus-breaking change introduced by the module.
// To avoid wrong/empty versions, the initial version should be set to 1.
func (AppModule) ConsensusVersion() uint64 { return 1 }

// BeginBlock contains the logic that is automatically triggered at the beginning of each block.
// The begin block implementation is optional.
func (am AppModule) BeginBlock(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Check for study plans that need automatic recommendations
	am.processAutomaticRecommendations(sdkCtx)

	// Check for study plans approaching deadlines
	am.checkStudyPlanDeadlines(sdkCtx)

	return nil
}

// EndBlock contains the logic that is automatically triggered at the end of each block.
// The end block implementation is optional.
func (am AppModule) EndBlock(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Process pending schedule optimizations
	am.processScheduleOptimizations(sdkCtx)

	// Update recommendation scores based on student progress
	am.updateRecommendationScores(sdkCtx)

	return nil
}

// processAutomaticRecommendations generates automatic recommendations for active study plans
func (am AppModule) processAutomaticRecommendations(ctx sdk.Context) {
	// Get all active study plans
	studyPlans, _, err := am.keeper.GetAllStudyPlans(ctx, nil)
	if err != nil {
		am.keeper.Logger().Error("Failed to get study plans for automatic recommendations", "error", err)
		return
	}

	for _, studyPlan := range studyPlans {
		if studyPlan.Status == types.StudyPlanStatusActive {
			// Generate recommendations for this study plan
			_, err := am.keeper.BuildRecommendations(ctx, studyPlan.Student, "")
			if err != nil {
				am.keeper.Logger().Error("Failed to generate automatic recommendations",
					"study_plan_id", studyPlan.Index,
					"student", studyPlan.Student,
					"error", err)
			}
		}
	}
}

// checkStudyPlanDeadlines checks for study plans approaching completion targets
func (am AppModule) checkStudyPlanDeadlines(ctx sdk.Context) {
	// This would implement logic to check for study plans approaching their completion targets
	// and emit events or notifications accordingly
	_ = ctx // Suppress unused parameter warning
	am.keeper.Logger().Debug("Checking study plan deadlines")
}

// processScheduleOptimizations processes any pending schedule optimizations
func (am AppModule) processScheduleOptimizations(ctx sdk.Context) {
	// This would implement batch processing of schedule optimizations
	_ = ctx // Suppress unused parameter warning
	am.keeper.Logger().Debug("Processing schedule optimizations")
}

// updateRecommendationScores updates recommendation scores based on student progress
func (am AppModule) updateRecommendationScores(ctx sdk.Context) {
	// This would implement logic to update recommendation scores based on
	// recent student academic performance and progress
	_ = ctx // Suppress unused parameter warning
	am.keeper.Logger().Debug("Updating recommendation scores")
}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

// ----------------------------------------------------------------------------
// App Wiring Setup
// ----------------------------------------------------------------------------

func init() {
	appmodule.Register(
		&modulev1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	StoreService store.KVStoreService
	Cdc          codec.Codec
	Config       *modulev1.Module
	Logger       log.Logger

	AccountKeeper    types.AccountKeeper
	BankKeeper       types.BankKeeper
	SubjectKeeper    types.SubjectKeeper
	StudentKeeper    types.StudentKeeper
	CurriculumKeeper types.CurriculumKeeper
}

type ModuleOutputs struct {
	depinject.Out

	ScheduleKeeper keeper.Keeper
	Module         appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	// default to governance authority if not provided
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}

	k := keeper.NewKeeper(
		in.Cdc,
		in.StoreService,
		in.Logger,
		authority.String(),
		in.SubjectKeeper,
		in.StudentKeeper,
		in.CurriculumKeeper,
		in.AccountKeeper,
		in.BankKeeper,
	)

	m := NewAppModule(
		in.Cdc,
		k,
		in.AccountKeeper,
		in.BankKeeper,
		in.SubjectKeeper,
		in.StudentKeeper,
		in.CurriculumKeeper,
	)

	return ModuleOutputs{ScheduleKeeper: k, Module: m}
}
