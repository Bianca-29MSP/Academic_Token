package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"academictoken/x/institution/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string

		accountKeeper types.AccountKeeper
		bankKeeper    types.BankKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,

	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		authority:    authority,
		logger:       logger,

		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
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

// SetInstitution stores an institution in the KV store
func (k Keeper) SetInstitution(ctx sdk.Context, institution types.Institution) error {
	store := k.storeService.OpenKVStore(ctx)

	// Validate institution data before storing
	if institution.Index == "" {
		return fmt.Errorf("institution index cannot be empty")
	}

	if institution.Name == "" {
		return fmt.Errorf("institution name cannot be empty")
	}

	if institution.Creator == "" {
		return fmt.Errorf("institution creator cannot be empty")
	}

	// Serialize institution data
	bz := k.cdc.MustMarshal(&institution)

	// Store using institution key
	key := types.InstitutionKey(institution.Index)
	return store.Set(key, bz)
}

// GetInstitution retrieves an institution by its index
func (k Keeper) GetInstitution(ctx sdk.Context, index string) (val types.Institution, found bool) {
	store := k.storeService.OpenKVStore(ctx)

	key := types.InstitutionKey(index)
	bz, err := store.Get(key)
	if err != nil || bz == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(bz, &val)
	return val, true
}

// RemoveInstitution removes an institution from the store
func (k Keeper) RemoveInstitution(ctx sdk.Context, index string) error {
	store := k.storeService.OpenKVStore(ctx)

	key := types.InstitutionKey(index)
	return store.Delete(key)
}

// GetAllInstitution returns all institutions
func (k Keeper) GetAllInstitution(ctx sdk.Context) (list []types.Institution) {
	store := k.storeService.OpenKVStore(ctx)
	iterator, err := store.Iterator(types.KeyPrefix(types.InstitutionKeyPrefix), storetypes.PrefixEndBytes(types.KeyPrefix(types.InstitutionKeyPrefix)))
	if err != nil {
		k.logger.Error("failed to create iterator for institutions", "error", err)
		return list
	}
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Institution
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetNextInstitutionIndex generates a new institution index
func (k Keeper) GetNextInstitutionIndex(ctx sdk.Context) string {
	store := k.storeService.OpenKVStore(ctx)

	// Get the current count
	countKey := []byte(types.InstitutionCountKey)
	countBytes, err := store.Get(countKey)
	var count uint64 = 1

	if err == nil && countBytes != nil {
		count = sdk.BigEndianToUint64(countBytes) + 1
	}

	// Store the new count
	countBytes = sdk.Uint64ToBigEndian(count)
	if err := store.Set(countKey, countBytes); err != nil {
		k.logger.Error("failed to update institution counter", "error", err)
	}

	return fmt.Sprintf("institution-%d", count)
}

// IsInstitutionAuthorized checks if an institution is authorized
func (k Keeper) IsInstitutionAuthorized(ctx sdk.Context, index string) bool {
	institution, found := k.GetInstitution(ctx, index)
	if !found {
		return false
	}

	// Check if institution is authorized (assuming "true" string means authorized)
	return institution.IsAuthorized == "true"
}

// IsInstitutionCreator checks if the given address is the original creator of the institution
func (k Keeper) IsInstitutionCreator(ctx sdk.Context, institutionIndex string, address string) bool {
	institution, found := k.GetInstitution(ctx, institutionIndex)
	if !found {
		return false
	}

	return institution.Creator == address
}

// CanUpdateInstitution checks if an address can update an institution
func (k Keeper) CanUpdateInstitution(ctx sdk.Context, institutionIndex string, updaterAddress string) bool {
	// Authority (governance) can always update
	if updaterAddress == k.GetAuthority() {
		return true
	}

	// Original creator can update their own institution
	return k.IsInstitutionCreator(ctx, institutionIndex, updaterAddress)
}

// CanAuthorizeInstitution checks if an address can change authorization status
func (k Keeper) CanAuthorizeInstitution(ctx sdk.Context, institutionIndex string, updaterAddress string) bool {
	// Authority (governance) can always authorize
	if updaterAddress == k.GetAuthority() {
		return true
	}

	// TEMPORARY: Allow the creator to authorize their own institution
	// In production, only governance should be able to do this
	return k.IsInstitutionCreator(ctx, institutionIndex, updaterAddress)
}

// GetInstitutionCount returns the total count of institutions
func (k Keeper) GetInstitutionCount(ctx sdk.Context) uint64 {
	store := k.storeService.OpenKVStore(ctx)

	countKey := []byte(types.InstitutionCountKey)
	countBytes, err := store.Get(countKey)
	if err != nil || countBytes == nil {
		return 0
	}

	return sdk.BigEndianToUint64(countBytes)
}

// Query methods for gRPC

// Institution returns institution by index
func (k Keeper) Institution(c context.Context, req *types.QueryGetInstitutionRequest) (*types.QueryGetInstitutionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetInstitution(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetInstitutionResponse{Institution: val}, nil
}

// InstitutionAll returns all institutions with pagination
func (k Keeper) InstitutionAll(c context.Context, req *types.QueryAllInstitutionRequest) (*types.QueryAllInstitutionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var institutions []types.Institution
	ctx := sdk.UnwrapSDKContext(c)

	store := k.storeService.OpenKVStore(ctx)
	institutionStore := prefix.NewStore(runtime.KVStoreAdapter(store), types.KeyPrefix(types.InstitutionKeyPrefix))

	pageRes, err := query.Paginate(institutionStore, req.Pagination, func(key []byte, value []byte) error {
		var institution types.Institution
		if err := k.cdc.Unmarshal(value, &institution); err != nil {
			return err
		}

		institutions = append(institutions, institution)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllInstitutionResponse{Institution: institutions, Pagination: pageRes}, nil
}

// AuthorizedInstitutions returns only authorized institutions
func (k Keeper) AuthorizedInstitutions(c context.Context, req *types.QueryAuthorizedInstitutionsRequest) (*types.QueryAuthorizedInstitutionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	allInstitutions := k.GetAllInstitution(ctx)

	var authorizedInstitutions []types.Institution
	for _, institution := range allInstitutions {
		if institution.IsAuthorized == "true" {
			authorizedInstitutions = append(authorizedInstitutions, institution)
		}
	}

	// Simple pagination for filtered results
	pageRes := &query.PageResponse{
		Total: uint64(len(authorizedInstitutions)),
	}

	return &types.QueryAuthorizedInstitutionsResponse{
		Institution: authorizedInstitutions,
		Pagination:  pageRes,
	}, nil
}

// InstitutionCount returns total count of institutions
func (k Keeper) InstitutionCount(c context.Context, req *types.QueryInstitutionCountRequest) (*types.QueryInstitutionCountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	count := k.GetInstitutionCount(ctx)

	return &types.QueryInstitutionCountResponse{Count: count}, nil
}

// GetAuthorizedInstitutions returns all authorized institutions
func (k Keeper) GetAuthorizedInstitutions(ctx sdk.Context) []types.Institution {
	allInstitutions := k.GetAllInstitution(ctx)
	var authorized []types.Institution

	for _, inst := range allInstitutions {
		if inst.IsAuthorized == "true" {
			authorized = append(authorized, inst)
		}
	}

	return authorized
}

// InstitutionExistsByName checks if an institution with the given name already exists
func (k Keeper) InstitutionExistsByName(ctx sdk.Context, name string) bool {
	allInstitutions := k.GetAllInstitution(ctx)
	for _, institution := range allInstitutions {
		if institution.Name == name {
			return true
		}
	}
	return false
}

// InstitutionExistsByAddress checks if an institution with the given address already exists
func (k Keeper) InstitutionExistsByAddress(ctx sdk.Context, address string) bool {
	allInstitutions := k.GetAllInstitution(ctx)
	for _, institution := range allInstitutions {
		if institution.Address == address {
			return true
		}
	}
	return false
}

// InstitutionExists checks if an institution exists by its index
func (k Keeper) InstitutionExists(ctx sdk.Context, index string) bool {
	_, found := k.GetInstitution(ctx, index)
	return found
}

var _ types.QueryServer = Keeper{}
