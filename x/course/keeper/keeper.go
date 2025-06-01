package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"academictoken/x/course/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string

		// Module dependencies
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

	institutionKeeper types.InstitutionKeeper,
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

// SetCourse stores a course in the KV store
func (k Keeper) SetCourse(ctx sdk.Context, course types.Course) error {
	store := k.storeService.OpenKVStore(ctx)

	// Validate course data before storing
	if course.Index == "" {
		return fmt.Errorf("course index cannot be empty")
	}

	if course.Name == "" {
		return fmt.Errorf("course name cannot be empty")
	}

	if course.Institution == "" {
		return fmt.Errorf("course institution cannot be empty")
	}

	// Verify that the institution exists first
	_, found := k.institutionKeeper.GetInstitution(ctx, course.Institution)
	if !found {
		return fmt.Errorf("institution %s does not exist", course.Institution)
	}

	// Verify that the institution exists and is authorized
	if !k.institutionKeeper.IsInstitutionAuthorized(ctx, course.Institution) {
		return fmt.Errorf("institution %s is not authorized", course.Institution)
	}

	// Serialize course data
	bz := k.cdc.MustMarshal(&course)

	// Store using course key
	key := types.CourseKey(course.Index)

	// Debug: verificar a chave sendo usada
	k.logger.Info("SetCourse debug",
		"courseIndex", course.Index,
		"key", string(key),
		"keyHex", fmt.Sprintf("%x", key),
		"dataLength", len(bz),
	)

	return store.Set(key, bz)
}

// GetCourse retrieves a course by its index
func (k Keeper) GetCourse(ctx sdk.Context, index string) (val types.Course, found bool) {
	store := k.storeService.OpenKVStore(ctx)

	key := types.CourseKey(index)
	bz, err := store.Get(key)
	if err != nil || bz == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(bz, &val)
	return val, true
}

// RemoveCourse removes a course from the store
func (k Keeper) RemoveCourse(ctx sdk.Context, index string) error {
	store := k.storeService.OpenKVStore(ctx)

	key := types.CourseKey(index)
	return store.Delete(key)
}

// GetAllCourse returns all courses
func (k Keeper) GetAllCourse(ctx sdk.Context) (list []types.Course) {
	store := k.storeService.OpenKVStore(ctx)

	// Debug: verificar o prefixo
	prefix := types.KeyPrefix(types.CourseKeyPrefix)
	prefixEnd := storetypes.PrefixEndBytes(prefix)

	k.logger.Info("GetAllCourse debug",
		"prefix", string(prefix),
		"prefixHex", fmt.Sprintf("%x", prefix),
		"prefixEnd", fmt.Sprintf("%x", prefixEnd),
	)

	iterator, err := store.Iterator(prefix, prefixEnd)
	if err != nil {
		k.logger.Error("failed to create iterator for courses", "error", err)
		return list
	}
	defer iterator.Close()

	// Debug: verificar se o iterator tem dados
	count := 0
	for ; iterator.Valid(); iterator.Next() {
		count++
		key := iterator.Key()
		value := iterator.Value()

		k.logger.Info("GetAllCourse found item",
			"count", count,
			"key", string(key),
			"keyHex", fmt.Sprintf("%x", key),
			"valueLength", len(value),
		)

		var val types.Course
		k.cdc.MustUnmarshal(value, &val)
		list = append(list, val)
	}

	k.logger.Info("GetAllCourse result", "totalFound", count, "listLength", len(list))
	return
}

// GetCoursesByInstitution returns all courses for a specific institution
func (k Keeper) GetCoursesByInstitution(ctx sdk.Context, institutionIndex string) (list []types.Course) {
	allCourses := k.GetAllCourse(ctx)

	for _, course := range allCourses {
		if course.Institution == institutionIndex {
			list = append(list, course)
		}
	}

	return
}

// GetNextCourseIndex generates a new course index
func (k Keeper) GetNextCourseIndex(ctx sdk.Context) string {
	store := k.storeService.OpenKVStore(ctx)

	// Get the current count
	countKey := []byte(types.CourseCountKey)
	countBytes, err := store.Get(countKey)
	var count uint64 = 1

	if err == nil && countBytes != nil {
		count = sdk.BigEndianToUint64(countBytes) + 1
	}

	// Store the new count
	countBytes = sdk.Uint64ToBigEndian(count)
	if err := store.Set(countKey, countBytes); err != nil {
		k.logger.Error("failed to update course counter", "error", err)
	}

	return fmt.Sprintf("course-%d", count)
}

// GetCourseCount returns the total count of courses
func (k Keeper) GetCourseCount(ctx sdk.Context) uint64 {
	store := k.storeService.OpenKVStore(ctx)

	countKey := []byte(types.CourseCountKey)
	countBytes, err := store.Get(countKey)
	if err != nil || countBytes == nil {
		return 0
	}

	return sdk.BigEndianToUint64(countBytes)
}

// CanCreateCourse checks if an address can create a course for an institution
func (k Keeper) CanCreateCourse(ctx sdk.Context, creatorAddress string, institutionIndex string) bool {
	// Check if institution is authorized
	if !k.institutionKeeper.IsInstitutionAuthorized(ctx, institutionIndex) {
		return false
	}

	// Get the institution
	institution, found := k.institutionKeeper.GetInstitution(ctx, institutionIndex)
	if !found {
		return false
	}

	// For now, only the institution's creator can create courses
	// This can be extended later to support multiple authorized users per institution
	return institution.Creator == creatorAddress
}

// CanUpdateCourse checks if an address can update a course
func (k Keeper) CanUpdateCourse(ctx sdk.Context, courseIndex string, updaterAddress string) bool {
	// Authority (governance) can always update
	if updaterAddress == k.GetAuthority() {
		return true
	}

	// Get the course
	course, found := k.GetCourse(ctx, courseIndex)
	if !found {
		return false
	}

	// Check if the updater is authorized for the institution
	return k.CanCreateCourse(ctx, updaterAddress, course.Institution)
}

// Query methods for gRPC

// Course returns course by index
func (k Keeper) Course(c context.Context, req *types.QueryGetCourseRequest) (*types.QueryGetCourseResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetCourse(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetCourseResponse{Course: val}, nil
}

// CourseAll returns all courses with pagination - FIXED
func (k Keeper) CourseAll(c context.Context, req *types.QueryAllCourseRequest) (*types.QueryAllCourseResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	// Use the working GetAllCourse method
	allCourses := k.GetAllCourse(ctx)

	// Simple pagination response
	pageRes := &query.PageResponse{
		Total: uint64(len(allCourses)),
	}

	return &types.QueryAllCourseResponse{Course: allCourses, Pagination: pageRes}, nil
}

// CoursesByInstitution returns courses for a specific institution - FIXED
func (k Keeper) CoursesByInstitution(c context.Context, req *types.QueryCoursesByInstitutionRequest) (*types.QueryCoursesByInstitutionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	// Verify institution exists first
	_, found := k.institutionKeeper.GetInstitution(ctx, req.InstitutionIndex)
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("institution %s not found", req.InstitutionIndex))
	}

	// Get courses for the institution
	courses := k.GetCoursesByInstitution(ctx, req.InstitutionIndex)

	// Simple pagination for filtered results
	pageRes := &query.PageResponse{
		Total: uint64(len(courses)),
	}

	return &types.QueryCoursesByInstitutionResponse{
		Course:     courses,
		Pagination: pageRes,
	}, nil
}

// CourseCount returns total count of courses
func (k Keeper) CourseCount(c context.Context, req *types.QueryCourseCountRequest) (*types.QueryCourseCountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	count := k.GetCourseCount(ctx)

	return &types.QueryCourseCountResponse{Count: count}, nil
}

// CourseExists checks if a course exists by its index
func (k Keeper) CourseExists(ctx sdk.Context, index string) bool {
	_, found := k.GetCourse(ctx, index)
	return found
}

// HasCourse is an alias for CourseExists (required by Subject module interface)
func (k Keeper) HasCourse(ctx sdk.Context, index string) bool {
	return k.CourseExists(ctx, index)
}

var _ types.QueryServer = Keeper{}
