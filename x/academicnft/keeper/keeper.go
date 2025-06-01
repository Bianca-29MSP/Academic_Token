package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"academictoken/x/academicnft/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string

		// Expected keepers (can be nil during development)
		tokenDefKeeper    types.TokenDefKeeper
		studentKeeper     types.StudentKeeper
		institutionKeeper types.InstitutionKeeper
		accountKeeper     types.AccountKeeper
		bankKeeper        types.BankKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	tokenDefKeeper types.TokenDefKeeper,
	studentKeeper types.StudentKeeper,
	institutionKeeper types.InstitutionKeeper,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:               cdc,
		storeService:      storeService,
		authority:         authority,
		logger:            logger,
		tokenDefKeeper:    tokenDefKeeper,
		studentKeeper:     studentKeeper,
		institutionKeeper: institutionKeeper,
		accountKeeper:     accountKeeper,
		bankKeeper:        bankKeeper,
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

// ----------------------------------------------------------------------------
// Parameters Management
// ----------------------------------------------------------------------------

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.ParamsKey)
	if err != nil || bz == nil {
		return types.DefaultParams()
	}

	var params types.Params
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	store := k.storeService.OpenKVStore(ctx)
	bz := k.cdc.MustMarshal(&params)
	err := store.Set(types.ParamsKey, bz)
	return err
}

// ----------------------------------------------------------------------------
// CRUD Operations
// ----------------------------------------------------------------------------

// SetSubjectTokenInstance stores a token instance in the KV store
func (k Keeper) SetSubjectTokenInstance(ctx sdk.Context, tokenInstance types.SubjectTokenInstance) {
	store := k.storeService.OpenKVStore(ctx)
	b := k.cdc.MustMarshal(&tokenInstance)
	store.Set(types.SubjectTokenInstanceKey(tokenInstance.Index), b)

	// Also create indexes for efficient querying
	k.setStudentTokenIndex(ctx, tokenInstance.Student, tokenInstance.Index)
	k.setTokenDefInstanceIndex(ctx, tokenInstance.TokenDefId, tokenInstance.Index)
}

// GetSubjectTokenInstance retrieves a token instance by its index
func (k Keeper) GetSubjectTokenInstance(ctx sdk.Context, index string) (val types.SubjectTokenInstance, found bool) {
	store := k.storeService.OpenKVStore(ctx)
	b, err := store.Get(types.SubjectTokenInstanceKey(index))
	if err != nil || b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// HasSubjectTokenInstance checks if a token instance exists
func (k Keeper) HasSubjectTokenInstance(ctx sdk.Context, index string) bool {
	store := k.storeService.OpenKVStore(ctx)
	has, err := store.Has(types.SubjectTokenInstanceKey(index))
	if err != nil {
		return false
	}
	return has
}

// RemoveSubjectTokenInstance removes a token instance from the store
func (k Keeper) RemoveSubjectTokenInstance(ctx sdk.Context, index string) {
	store := k.storeService.OpenKVStore(ctx)

	// Get the token instance first to clean up indexes
	tokenInstance, found := k.GetSubjectTokenInstance(ctx, index)
	if !found {
		return
	}

	// Remove the main entry
	store.Delete(types.SubjectTokenInstanceKey(index))

	// Clean up indexes
	k.removeStudentTokenIndex(ctx, tokenInstance.Student, index)
	k.removeTokenDefInstanceIndex(ctx, tokenInstance.TokenDefId, index)
}

// GetAllSubjectTokenInstances returns all token instances
func (k Keeper) GetAllSubjectTokenInstances(ctx sdk.Context) (list []types.SubjectTokenInstance) {
	store := k.storeService.OpenKVStore(ctx)
	iterator, err := store.Iterator(types.SubjectTokenInstanceKeyPrefix, storetypes.PrefixEndBytes(types.SubjectTokenInstanceKeyPrefix))
	if err != nil {
		return list
	}

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.SubjectTokenInstance
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetStudentTokenInstances returns all token instances for a student
func (k Keeper) GetStudentTokenInstances(ctx sdk.Context, studentAddress string) ([]types.SubjectTokenInstance, error) {
	var tokens []types.SubjectTokenInstance

	// Get token instance IDs for this student from index
	tokenIndexes := k.getStudentTokenIndexes(ctx, studentAddress)

	for _, tokenIndex := range tokenIndexes {
		tokenInstance, found := k.GetSubjectTokenInstance(ctx, tokenIndex)
		if found {
			tokens = append(tokens, tokenInstance)
		}
	}

	return tokens, nil
}

// GetTokenDefInstancesInternal returns all instances of a specific TokenDef (internal method)
func (k Keeper) GetTokenDefInstancesInternal(ctx sdk.Context, tokenDefId string) ([]types.SubjectTokenInstance, error) {
	var tokens []types.SubjectTokenInstance

	// Get token instance IDs for this TokenDef from index
	tokenIndexes := k.getTokenDefInstanceIndexes(ctx, tokenDefId)

	for _, tokenIndex := range tokenIndexes {
		tokenInstance, found := k.GetSubjectTokenInstance(ctx, tokenIndex)
		if found {
			tokens = append(tokens, tokenInstance)
		}
	}

	return tokens, nil
}

// ----------------------------------------------------------------------------
// Business Logic
// ----------------------------------------------------------------------------

// ValidateTokenInstance performs comprehensive validation of a token instance
func (k Keeper) ValidateTokenInstance(ctx sdk.Context, tokenInstanceId string) (bool, error) {
	// Check if token instance exists
	tokenInstance, found := k.GetSubjectTokenInstance(ctx, tokenInstanceId)
	if !found {
		return false, fmt.Errorf("token instance %s not found", tokenInstanceId)
	}

	// Only validate with keepers that are available (not nil)
	if k.tokenDefKeeper != nil {
		// Use the real method name
		_, found := k.tokenDefKeeper.GetTokenDefinitionByIndex(ctx, tokenInstance.TokenDefId)
		if !found {
			return false, fmt.Errorf("token definition %s not found", tokenInstance.TokenDefId)
		}
	}

	if k.institutionKeeper != nil {
		// Use the real method
		if !k.institutionKeeper.IsInstitutionAuthorized(ctx, tokenInstance.IssuerInstitution) {
			return false, fmt.Errorf("institution %s is not authorized", tokenInstance.IssuerInstitution)
		}
	}

	if k.studentKeeper != nil {
		err := k.studentKeeper.ValidateStudentEligibility(ctx, tokenInstance.Student, tokenInstance.IssuerInstitution)
		if err != nil {
			return false, fmt.Errorf("invalid student eligibility: %v", err)
		}
	}

	return true, nil
}

// GenerateTokenInstanceID generates a unique ID for a new token instance
func (k Keeper) GenerateTokenInstanceID(ctx sdk.Context) string {
	// Simple counter-based ID generation
	// In production, you might want something more sophisticated
	count := k.GetTokenInstanceCount(ctx)
	k.SetTokenInstanceCount(ctx, count+1)
	return fmt.Sprintf("token-instance-%d", count+1)
}

// GetTokenInstanceCount gets the current count of token instances
func (k Keeper) GetTokenInstanceCount(ctx sdk.Context) uint64 {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.TokenInstanceCountKey)
	if err != nil || bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

// SetTokenInstanceCount sets the count of token instances
func (k Keeper) SetTokenInstanceCount(ctx sdk.Context, count uint64) {
	store := k.storeService.OpenKVStore(ctx)
	bz := sdk.Uint64ToBigEndian(count)
	store.Set(types.TokenInstanceCountKey, bz)
}

// ----------------------------------------------------------------------------
// Query Server Implementation
// ----------------------------------------------------------------------------

// Params returns the parameters of the module (Query Server method)
func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

// GetTokenInstance returns a specific token instance by ID
func (k Keeper) GetTokenInstance(goCtx context.Context, req *types.QueryGetTokenInstanceRequest) (*types.QueryGetTokenInstanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	tokenInstance, found := k.GetSubjectTokenInstance(ctx, req.TokenInstanceId)
	if !found {
		return nil, status.Error(codes.NotFound, "token instance not found")
	}

	return &types.QueryGetTokenInstanceResponse{TokenInstance: &tokenInstance}, nil
}

// GetStudentTokens returns all token instances for a student
func (k Keeper) GetStudentTokens(goCtx context.Context, req *types.QueryGetStudentTokensRequest) (*types.QueryGetStudentTokensResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get tokens for student
	tokens, err := k.GetStudentTokenInstances(ctx, req.StudentAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert to pointer slice for response
	var tokenInstances []*types.SubjectTokenInstance
	for i := range tokens {
		tokenInstances = append(tokenInstances, &tokens[i])
	}

	return &types.QueryGetStudentTokensResponse{
		TokenInstances: tokenInstances,
		Pagination:     nil, // Simple implementation without pagination
	}, nil
}

// GetTokenDefInstances returns all instances of a specific TokenDef (Query Server method)
func (k Keeper) GetTokenDefInstances(goCtx context.Context, req *types.QueryGetTokenDefInstancesRequest) (*types.QueryGetTokenDefInstancesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get instances for TokenDef
	tokens, err := k.GetTokenDefInstancesInternal(ctx, req.TokenDefId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert to pointer slice for response
	var tokenInstances []*types.SubjectTokenInstance
	for i := range tokens {
		tokenInstances = append(tokenInstances, &tokens[i])
	}

	return &types.QueryGetTokenDefInstancesResponse{
		TokenInstances: tokenInstances,
		Pagination:     nil, // Simple implementation without pagination
	}, nil
}

// VerifyTokenInstance verifies if a token instance exists and is valid
func (k Keeper) VerifyTokenInstance(goCtx context.Context, req *types.QueryVerifyTokenInstanceRequest) (*types.QueryVerifyTokenInstanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if token exists
	tokenInstance, exists := k.GetSubjectTokenInstance(ctx, req.TokenInstanceId)
	if !exists {
		return &types.QueryVerifyTokenInstanceResponse{
			Exists:        false,
			IsValid:       false,
			TokenInstance: nil,
		}, nil
	}

	// Validate token instance
	isValid, err := k.ValidateTokenInstance(ctx, req.TokenInstanceId)
	if err != nil {
		// Token exists but has validation errors
		return &types.QueryVerifyTokenInstanceResponse{
			Exists:        true,
			IsValid:       false,
			TokenInstance: &tokenInstance,
		}, nil
	}

	return &types.QueryVerifyTokenInstanceResponse{
		Exists:        true,
		IsValid:       isValid,
		TokenInstance: &tokenInstance,
	}, nil
}

// ----------------------------------------------------------------------------
// Helper methods for maintaining indexes
// ----------------------------------------------------------------------------

// setStudentTokenIndex adds a token instance to student's index
func (k Keeper) setStudentTokenIndex(ctx sdk.Context, studentAddress, tokenInstanceId string) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.StudentTokenIndexKey(studentAddress, tokenInstanceId)
	store.Set(key, []byte{1}) // Just a marker
}

// removeStudentTokenIndex removes a token instance from student's index
func (k Keeper) removeStudentTokenIndex(ctx sdk.Context, studentAddress, tokenInstanceId string) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.StudentTokenIndexKey(studentAddress, tokenInstanceId)
	store.Delete(key)
}

// getStudentTokenIndexes returns all token instance IDs for a student
func (k Keeper) getStudentTokenIndexes(ctx sdk.Context, studentAddress string) []string {
	store := k.storeService.OpenKVStore(ctx)
	prefix := types.StudentTokenIndexPrefix(studentAddress)
	iterator, err := store.Iterator(prefix, storetypes.PrefixEndBytes(prefix))
	if err != nil {
		return []string{}
	}
	defer iterator.Close()

	var tokenIndexes []string
	for ; iterator.Valid(); iterator.Next() {
		// Extract token instance ID from key
		key := iterator.Key()
		tokenInstanceId := string(key[len(prefix):])
		tokenIndexes = append(tokenIndexes, tokenInstanceId)
	}

	return tokenIndexes
}

// setTokenDefInstanceIndex adds a token instance to TokenDef's index
func (k Keeper) setTokenDefInstanceIndex(ctx sdk.Context, tokenDefId, tokenInstanceId string) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.TokenDefInstanceIndexKey(tokenDefId, tokenInstanceId)
	store.Set(key, []byte{1}) // Just a marker
}

// removeTokenDefInstanceIndex removes a token instance from TokenDef's index
func (k Keeper) removeTokenDefInstanceIndex(ctx sdk.Context, tokenDefId, tokenInstanceId string) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.TokenDefInstanceIndexKey(tokenDefId, tokenInstanceId)
	store.Delete(key)
}

// getTokenDefInstanceIndexes returns all token instance IDs for a TokenDef
func (k Keeper) getTokenDefInstanceIndexes(ctx sdk.Context, tokenDefId string) []string {
	store := k.storeService.OpenKVStore(ctx)
	prefix := types.TokenDefInstanceIndexPrefix(tokenDefId)
	iterator, err := store.Iterator(prefix, storetypes.PrefixEndBytes(prefix))
	if err != nil {
		return []string{}
	}
	defer iterator.Close()

	var tokenIndexes []string
	for ; iterator.Valid(); iterator.Next() {
		// Extract token instance ID from key
		key := iterator.Key()
		tokenInstanceId := string(key[len(prefix):])
		tokenIndexes = append(tokenIndexes, tokenInstanceId)
	}

	return tokenIndexes
}
