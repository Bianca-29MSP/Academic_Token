package keeper

import (
	"context"
	"fmt"
	"strconv"

	"cosmossdk.io/core/store"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	coursekeeper "academictoken/x/course/keeper"
	"academictoken/x/curriculum/types"
	subjectkeeper "academictoken/x/subject/keeper"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		storeKey     storetypes.StoreKey
		logger       log.Logger
		paramstore   paramtypes.Subspace

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string

		// Dependencies on other modules
		courseKeeper  *coursekeeper.Keeper
		subjectKeeper *subjectkeeper.Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	storeKey storetypes.StoreKey,
	logger log.Logger,
	ps paramtypes.Subspace,
	courseKeeper *coursekeeper.Keeper,
	subjectKeeper *subjectkeeper.Keeper,
	authority string,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:           cdc,
		storeService:  storeService,
		storeKey:      storeKey,
		authority:     authority,
		logger:        logger,
		paramstore:    ps,
		courseKeeper:  courseKeeper,
		subjectKeeper: subjectKeeper,
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

// GetCourseKeeper returns the course keeper
func (k Keeper) GetCourseKeeper() *coursekeeper.Keeper {
	return k.courseKeeper
}

// GetSubjectKeeper returns the subject keeper
func (k Keeper) GetSubjectKeeper() *subjectkeeper.Keeper {
	return k.subjectKeeper
}

// ============================================================================
// CurriculumTree CRUD Operations
// ============================================================================

// SetCurriculumTree stores a curriculum tree in the KV store
func (k Keeper) SetCurriculumTree(ctx context.Context, curriculumTree types.CurriculumTree) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix("CurriculumTree/"))
	b := k.cdc.MustMarshal(&curriculumTree)
	store.Set(types.CurriculumTreeKey(curriculumTree.Index), b)
}

// GetCurriculumTree returns a curriculum tree from its index (context.Context version)
func (k Keeper) GetCurriculumTree(ctx context.Context, index string) (val types.CurriculumTree, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix("CurriculumTree/"))

	b := store.Get(types.CurriculumTreeKey(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetCurriculumTree returns a curriculum tree from its index (sdk.Context version)
func (k Keeper) GetCurriculumTreeSDK(ctx sdk.Context, index string) (val types.CurriculumTree, found bool) {
	return k.GetCurriculumTree(sdk.WrapSDKContext(ctx), index)
}

// RemoveCurriculumTree removes a curriculum tree from the store
func (k Keeper) RemoveCurriculumTree(ctx context.Context, index string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix("CurriculumTree/"))
	store.Delete(types.CurriculumTreeKey(index))
}

// GetAllCurriculumTree returns all curriculum trees
func (k Keeper) GetAllCurriculumTree(ctx context.Context) (list []types.CurriculumTree) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix("CurriculumTree/"))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.CurriculumTree
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetCurriculumTreeCount gets the total number of curriculum trees
func (k Keeper) GetCurriculumTreeCount(ctx context.Context) uint64 {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix("CurriculumTreeCount/"))
	byteKey := types.KeyPrefix("CurriculumTreeCount/")
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	count, err := strconv.ParseUint(string(bz), 10, 64)
	if err != nil {
		// Panic because the count should be always formattable to uint64
		panic("cannot decode count")
	}

	return count
}

// SetCurriculumTreeCount sets the total number of curriculum trees
func (k Keeper) SetCurriculumTreeCount(ctx context.Context, count uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix("CurriculumTreeCount/"))
	byteKey := types.KeyPrefix("CurriculumTreeCount/")
	bz := []byte(strconv.FormatUint(count, 10))
	store.Set(byteKey, bz)
}

// AppendCurriculumTree appends a curriculum tree in the store with a new id and update the count
func (k Keeper) AppendCurriculumTree(ctx context.Context, curriculumTree types.CurriculumTree) (uint64, error) {
	// Validate that course exists
	if err := k.ValidateCourseExists(ctx, curriculumTree.CourseId); err != nil {
		return 0, err
	}

	// Validate that all subjects exist
	if err := k.ValidateSubjectsExist(ctx, curriculumTree.RequiredSubjects); err != nil {
		return 0, err
	}
	if err := k.ValidateSubjectsExist(ctx, curriculumTree.ElectiveSubjects); err != nil {
		return 0, err
	}

	// Validate semester structure subjects
	for _, semester := range curriculumTree.SemesterStructure {
		if err := k.ValidateSubjectsExist(ctx, semester.SubjectIds); err != nil {
			return 0, err
		}
	}

	// Validate elective groups subjects
	for _, group := range curriculumTree.ElectiveGroups {
		if err := k.ValidateSubjectsExist(ctx, group.SubjectIds); err != nil {
			return 0, err
		}
	}

	// Create the curriculum tree
	count := k.GetCurriculumTreeCount(ctx)

	// Set the ID of the appended value
	curriculumTree.Index = strconv.FormatUint(count, 10)

	k.SetCurriculumTree(ctx, curriculumTree)
	k.SetCurriculumTreeCount(ctx, count+1)

	return count, nil
}

// ValidateCourseExists checks if a course exists using the course keeper
func (k Keeper) ValidateCourseExists(ctx context.Context, courseId string) error {
	if k.courseKeeper == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "course keeper not available")
	}

	// Convert context.Context to sdk.Context
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	_, found := k.courseKeeper.GetCourse(sdkCtx, courseId)
	if !found {
		return errorsmod.Wrapf(sdkerrors.ErrKeyNotFound, "course with id %s not found", courseId)
	}

	return nil
}

// ValidateSubjectsExist checks if all subjects exist using the subject keeper
func (k Keeper) ValidateSubjectsExist(ctx context.Context, subjectIds []string) error {
	if k.subjectKeeper == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "subject keeper not available")
	}

	// Convert context.Context to sdk.Context
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	for _, subjectId := range subjectIds {
		_, found := k.subjectKeeper.GetSubject(sdkCtx, subjectId)
		if !found {
			return errorsmod.Wrapf(sdkerrors.ErrKeyNotFound, "subject with id %s not found", subjectId)
		}
	}

	return nil
}

// GetCurriculumTreesByCourse returns all curriculum trees for a specific course
func (k Keeper) GetCurriculumTreesByCourse(ctx context.Context, courseId string) (list []types.CurriculumTree) {
	allCurriculums := k.GetAllCurriculumTree(ctx)
	for _, curriculum := range allCurriculums {
		if curriculum.CourseId == courseId {
			list = append(list, curriculum)
		}
	}
	return
}

// AddSemesterToCurriculum adds a semester to an existing curriculum
func (k Keeper) AddSemesterToCurriculum(ctx context.Context, curriculumIndex string, semesterNumber uint64, subjectIds []string) error {
	curriculum, found := k.GetCurriculumTree(ctx, curriculumIndex)
	if !found {
		return errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "curriculum tree not found")
	}

	// Validate subjects exist
	if err := k.ValidateSubjectsExist(ctx, subjectIds); err != nil {
		return err
	}

	// Check if semester already exists
	for _, semester := range curriculum.SemesterStructure {
		if semester.SemesterNumber == fmt.Sprintf("%d", semesterNumber) {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "semester already exists")
		}
	}

	// Add new semester
	newSemester := &types.CurriculumSemester{
		SemesterNumber: fmt.Sprintf("%d", semesterNumber),
		SubjectIds:     subjectIds,
	}

	curriculum.SemesterStructure = append(curriculum.SemesterStructure, newSemester)
	k.SetCurriculumTree(ctx, curriculum)

	return nil
}

// AddElectiveGroupToCurriculum adds an elective group to an existing curriculum
func (k Keeper) AddElectiveGroupToCurriculum(ctx context.Context, curriculumIndex string, electiveGroup types.ElectiveGroup) error {
	curriculum, found := k.GetCurriculumTree(ctx, curriculumIndex)
	if !found {
		return errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "curriculum tree not found")
	}

	// Validate subjects exist
	if err := k.ValidateSubjectsExist(ctx, electiveGroup.SubjectIds); err != nil {
		return err
	}

	// Generate group ID
	electiveGroup.GroupId = fmt.Sprintf("%s-group-%d", curriculumIndex, len(curriculum.ElectiveGroups))

	curriculum.ElectiveGroups = append(curriculum.ElectiveGroups, &electiveGroup)
	k.SetCurriculumTree(ctx, curriculum)

	return nil
}

// SetGraduationRequirements sets graduation requirements for a curriculum
func (k Keeper) SetGraduationRequirements(ctx context.Context, curriculumIndex string, requirements types.GraduationRequirements) error {
	curriculum, found := k.GetCurriculumTree(ctx, curriculumIndex)
	if !found {
		return errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "curriculum tree not found")
	}

	curriculum.GraduationRequirements = &requirements
	k.SetCurriculumTree(ctx, curriculum)

	return nil
}

// ============================================================================
// Query Handlers
// ============================================================================

func (k Keeper) CurriculumTreeAll(goCtx context.Context, req *types.QueryAllCurriculumTreeRequest) (*types.QueryAllCurriculumTreeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var curriculumTrees []types.CurriculumTree

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(goCtx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix("CurriculumTree/"))

	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		var curriculumTree types.CurriculumTree
		if err := k.cdc.Unmarshal(value, &curriculumTree); err != nil {
			return err
		}

		curriculumTrees = append(curriculumTrees, curriculumTree)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllCurriculumTreeResponse{CurriculumTree: curriculumTrees, Pagination: pageRes}, nil
}

func (k Keeper) CurriculumTree(goCtx context.Context, req *types.QueryGetCurriculumTreeRequest) (*types.QueryGetCurriculumTreeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, found := k.GetCurriculumTree(goCtx, req.Index)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetCurriculumTreeResponse{CurriculumTree: val}, nil
}

func (k Keeper) CurriculumTreesByCourse(goCtx context.Context, req *types.QueryCurriculumTreesByCourseRequest) (*types.QueryCurriculumTreesByCourseResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	curriculums := k.GetCurriculumTreesByCourse(goCtx, req.CourseId)

	return &types.QueryCurriculumTreesByCourseResponse{CurriculumTrees: curriculums}, nil
}
