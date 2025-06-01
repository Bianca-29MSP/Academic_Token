package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"academictoken/x/student/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string

		// External module keepers
		accountKeeper     types.AccountKeeper
		bankKeeper        types.BankKeeper
		institutionKeeper types.InstitutionKeeper
		courseKeeper      types.CourseKeeper
		curriculumKeeper  types.CurriculumKeeper
		subjectKeeper     types.SubjectKeeper
		tokenDefKeeper    types.TokenDefKeeper
		academicNFTKeeper types.AcademicNFTKeeper

		// Contract integration components
		wasmMsgServer       types.WasmMsgServer
		wasmQuerier         types.WasmQuerier
		contractIntegration *ContractIntegration
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	institutionKeeper types.InstitutionKeeper,
	courseKeeper types.CourseKeeper,
	curriculumKeeper types.CurriculumKeeper,
	subjectKeeper types.SubjectKeeper,
	tokenDefKeeper types.TokenDefKeeper,
	academicNFTKeeper types.AcademicNFTKeeper,
	wasmMsgServer types.WasmMsgServer,
	wasmQuerier types.WasmQuerier,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	keeper := Keeper{
		cdc:               cdc,
		storeService:      storeService,
		authority:         authority,
		logger:            logger,
		accountKeeper:     accountKeeper,
		bankKeeper:        bankKeeper,
		institutionKeeper: institutionKeeper,
		courseKeeper:      courseKeeper,
		curriculumKeeper:  curriculumKeeper,
		subjectKeeper:     subjectKeeper,
		tokenDefKeeper:    tokenDefKeeper,
		academicNFTKeeper: academicNFTKeeper,
		wasmMsgServer:     wasmMsgServer,
		wasmQuerier:       wasmQuerier,
	}

	// Initialize contract integration
	keeper.contractIntegration = NewContractIntegration(&keeper, wasmMsgServer, wasmQuerier)

	return keeper
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// ============================================================================
// GETTER METHODS FOR EXTERNAL KEEPERS
// ============================================================================

// GetAccountKeeper returns the account keeper
func (k Keeper) GetAccountKeeper() types.AccountKeeper {
	return k.accountKeeper
}

// GetBankKeeper returns the bank keeper
func (k Keeper) GetBankKeeper() types.BankKeeper {
	return k.bankKeeper
}

// GetInstitutionKeeper returns the institution keeper
func (k Keeper) GetInstitutionKeeper() types.InstitutionKeeper {
	return k.institutionKeeper
}

// GetCourseKeeper returns the course keeper
func (k Keeper) GetCourseKeeper() types.CourseKeeper {
	return k.courseKeeper
}

// GetCurriculumKeeper returns the curriculum keeper
func (k Keeper) GetCurriculumKeeper() types.CurriculumKeeper {
	return k.curriculumKeeper
}

// GetSubjectKeeper returns the subject keeper
func (k Keeper) GetSubjectKeeper() types.SubjectKeeper {
	return k.subjectKeeper
}

// GetTokenDefKeeper returns the tokendef keeper
func (k Keeper) GetTokenDefKeeper() types.TokenDefKeeper {
	return k.tokenDefKeeper
}

// GetAcademicNFTKeeper returns the academicnft keeper
func (k Keeper) GetAcademicNFTKeeper() types.AcademicNFTKeeper {
	return k.academicNFTKeeper
}

// GetContractIntegration returns the contract integration instance
func (k Keeper) GetContractIntegration() interface{} {
	return k.contractIntegration
}

// GetWasmMsgServer returns the wasm message server
func (k Keeper) GetWasmMsgServer() types.WasmMsgServer {
	return k.wasmMsgServer
}

// GetWasmQuerier returns the wasm querier
func (k Keeper) GetWasmQuerier() types.WasmQuerier {
	return k.wasmQuerier
}

// ============================================================================
// INTERNAL STUDENT CRUD OPERATIONS
// ============================================================================

// setStudent set a specific student in the store from its index
func (k Keeper) setStudent(ctx sdk.Context, student types.Student) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StudentKeyPrefix))
	b := k.cdc.MustMarshal(&student)
	store.Set(types.KeyPrefix(student.Index), b)
}

// getStudentByIndex returns a student from its index
func (k Keeper) getStudentByIndex(ctx sdk.Context, index string) (val types.Student, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StudentKeyPrefix))

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// getAllStudents returns all student
func (k Keeper) getAllStudents(ctx sdk.Context) (list []types.Student) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StudentKeyPrefix))
	iterator := store.Iterator(nil, nil)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Student
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// getStudentByAddress returns a student by their address
func (k Keeper) getStudentByAddress(ctx sdk.Context, address string) (val types.Student, found bool) {
	students := k.getAllStudents(ctx)
	for _, student := range students {
		if student.Address == address {
			return student, true
		}
	}
	return val, false
}

// ============================================================================
// INTERNAL STUDENT ENROLLMENT CRUD OPERATIONS
// ============================================================================

// setStudentEnrollment set a specific studentEnrollment in the store from its index
func (k Keeper) setStudentEnrollment(ctx sdk.Context, studentEnrollment types.StudentEnrollment) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StudentEnrollmentKeyPrefix))
	b := k.cdc.MustMarshal(&studentEnrollment)
	store.Set(types.KeyPrefix(studentEnrollment.Index), b)
}

// getStudentEnrollment returns a studentEnrollment from its index
func (k Keeper) getStudentEnrollment(ctx sdk.Context, index string) (val types.StudentEnrollment, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StudentEnrollmentKeyPrefix))

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// getAllStudentEnrollments returns all studentEnrollment
func (k Keeper) getAllStudentEnrollments(ctx sdk.Context) (list []types.StudentEnrollment) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StudentEnrollmentKeyPrefix))
	iterator := store.Iterator(nil, nil)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.StudentEnrollment
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// getEnrollmentsByStudentId returns all enrollments for a specific student
func (k Keeper) getEnrollmentsByStudentId(ctx sdk.Context, studentId string) (list []types.StudentEnrollment) {
	enrollments := k.getAllStudentEnrollments(ctx)
	for _, enrollment := range enrollments {
		if enrollment.Student == studentId {
			list = append(list, enrollment)
		}
	}
	return
}

// getEnrollmentsByInstitution returns all enrollments for a specific institution
func (k Keeper) getEnrollmentsByInstitution(ctx sdk.Context, institutionId string) (list []types.StudentEnrollment) {
	enrollments := k.getAllStudentEnrollments(ctx)
	for _, enrollment := range enrollments {
		if enrollment.Institution == institutionId {
			list = append(list, enrollment)
		}
	}
	return
}

// getEnrollmentsByCourse returns all enrollments for a specific course
func (k Keeper) getEnrollmentsByCourse(ctx sdk.Context, courseId string) (list []types.StudentEnrollment) {
	enrollments := k.getAllStudentEnrollments(ctx)
	for _, enrollment := range enrollments {
		if enrollment.CourseId == courseId {
			list = append(list, enrollment)
		}
	}
	return
}

// ============================================================================
// INTERNAL STUDENT ACADEMIC TREE CRUD OPERATIONS
// ============================================================================

// setStudentAcademicTree set a specific studentAcademicTree in the store from its index
func (k Keeper) setStudentAcademicTree(ctx sdk.Context, studentAcademicTree types.StudentAcademicTree) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StudentAcademicTreeKeyPrefix))
	b := k.cdc.MustMarshal(&studentAcademicTree)
	store.Set(types.KeyPrefix(studentAcademicTree.Index), b)
}

/*
// removeStudentAcademicTree removes a studentAcademicTree from the store
func (k Keeper) removeStudentAcademicTree(ctx sdk.Context, index string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StudentAcademicTreeKeyPrefix))
	store.Delete(types.KeyPrefix(index))
}*/

// getAllStudentAcademicTrees returns all studentAcademicTree
func (k Keeper) getAllStudentAcademicTrees(ctx sdk.Context) (list []types.StudentAcademicTree) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StudentAcademicTreeKeyPrefix))
	iterator := store.Iterator(nil, nil)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.StudentAcademicTree
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// getAcademicTreeByStudent returns the academic tree for a specific student
func (k Keeper) getAcademicTreeByStudent(ctx sdk.Context, studentId string) (val types.StudentAcademicTree, found bool) {
	trees := k.getAllStudentAcademicTrees(ctx)
	for _, tree := range trees {
		if tree.Student == studentId {
			return tree, true
		}
	}
	return val, false
}

// ============================================================================
// COUNTER MANAGEMENT FUNCTIONS
// ============================================================================

// GetStudentCount get the total number of students
func (k Keeper) GetStudentCount(ctx sdk.Context) uint64 {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	byteKey := types.KeyPrefix(types.StudentCounterKey)
	bz := store.Get(byteKey)

	if bz == nil {
		return 0
	}

	return sdk.BigEndianToUint64(bz)
}

// SetStudentCount set the total number of students
func (k Keeper) SetStudentCount(ctx sdk.Context, count uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	byteKey := types.KeyPrefix(types.StudentCounterKey)
	bz := sdk.Uint64ToBigEndian(count)
	store.Set(byteKey, bz)
}

// GetStudentEnrollmentCount get the total number of student enrollments
func (k Keeper) GetStudentEnrollmentCount(ctx sdk.Context) uint64 {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	byteKey := types.KeyPrefix(types.StudentEnrollmentCounterKey)
	bz := store.Get(byteKey)

	if bz == nil {
		return 0
	}

	return sdk.BigEndianToUint64(bz)
}

// SetStudentEnrollmentCount set the total number of student enrollments
func (k Keeper) SetStudentEnrollmentCount(ctx sdk.Context, count uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	byteKey := types.KeyPrefix(types.StudentEnrollmentCounterKey)
	bz := sdk.Uint64ToBigEndian(count)
	store.Set(byteKey, bz)
}

// GetStudentAcademicTreeCount get the total number of student academic trees
func (k Keeper) GetStudentAcademicTreeCount(ctx sdk.Context) uint64 {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	byteKey := types.KeyPrefix(types.StudentAcademicTreeCounterKey)
	bz := store.Get(byteKey)

	if bz == nil {
		return 0
	}

	return sdk.BigEndianToUint64(bz)
}

// SetStudentAcademicTreeCount set the total number of student academic trees
func (k Keeper) SetStudentAcademicTreeCount(ctx sdk.Context, count uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	byteKey := types.KeyPrefix(types.StudentAcademicTreeCounterKey)
	bz := sdk.Uint64ToBigEndian(count)
	store.Set(byteKey, bz)
}

// ============================================================================
// UTILITY FUNCTIONS
// ============================================================================

// AppendStudent appends a student in the store with a new id and update the count
func (k Keeper) AppendStudent(ctx sdk.Context, student types.Student) (uint64, error) {
	count := k.GetStudentCount(ctx)
	student.Index = fmt.Sprintf("%d", count)

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StudentKeyPrefix))
	appendedValue := k.cdc.MustMarshal(&student)
	store.Set(types.KeyPrefix(student.Index), appendedValue)

	k.SetStudentCount(ctx, count+1)
	return count, nil
}

// AppendStudentEnrollment appends a student enrollment in the store with a new id and update the count
func (k Keeper) AppendStudentEnrollment(ctx sdk.Context, studentEnrollment types.StudentEnrollment) (uint64, error) {
	count := k.GetStudentEnrollmentCount(ctx)
	studentEnrollment.Index = fmt.Sprintf("%d", count)

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StudentEnrollmentKeyPrefix))
	appendedValue := k.cdc.MustMarshal(&studentEnrollment)
	store.Set(types.KeyPrefix(studentEnrollment.Index), appendedValue)

	k.SetStudentEnrollmentCount(ctx, count+1)
	return count, nil
}

// AppendStudentAcademicTree appends a student academic tree in the store with a new id and update the count
func (k Keeper) AppendStudentAcademicTree(ctx sdk.Context, studentAcademicTree types.StudentAcademicTree) (uint64, error) {
	count := k.GetStudentAcademicTreeCount(ctx)
	studentAcademicTree.Index = fmt.Sprintf("%d", count)

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StudentAcademicTreeKeyPrefix))
	appendedValue := k.cdc.MustMarshal(&studentAcademicTree)
	store.Set(types.KeyPrefix(studentAcademicTree.Index), appendedValue)

	k.SetStudentAcademicTreeCount(ctx, count+1)
	return count, nil
}

// ============================================================================
// VALIDATION FUNCTIONS
// ============================================================================

// ValidateStudent validates student data
func (k Keeper) ValidateStudent(ctx sdk.Context, student types.Student) error {
	if _, found := k.getStudentByAddress(ctx, student.Address); found {
		return types.ErrStudentAlreadyExists
	}
	return nil
}

// ValidateStudentEnrollment validates student enrollment data
func (k Keeper) ValidateStudentEnrollment(ctx sdk.Context, enrollment types.StudentEnrollment) error {
	if _, found := k.getStudentByIndex(ctx, enrollment.Student); !found {
		return types.ErrStudentNotFound
	}

	if _, found := k.institutionKeeper.GetInstitution(ctx, enrollment.Institution); !found {
		return types.ErrInvalidInstitution
	}

	if _, found := k.courseKeeper.GetCourse(ctx, enrollment.CourseId); !found {
		return types.ErrInvalidCourse
	}

	return nil
}

// ValidateStudentAcademicTree validates student academic tree data
func (k Keeper) ValidateStudentAcademicTree(ctx sdk.Context, tree types.StudentAcademicTree) error {
	if _, found := k.getStudentByIndex(ctx, tree.Student); !found {
		return types.ErrStudentNotFound
	}

	if _, found := k.institutionKeeper.GetInstitution(ctx, tree.Institution); !found {
		return types.ErrInvalidInstitution
	}

	if _, found := k.courseKeeper.GetCourse(ctx, tree.CourseId); !found {
		return types.ErrInvalidCourse
	}

	return nil
}

// ============================================================================
// PUBLIC METHODS FOR EXTERNAL USAGE
// ============================================================================

// GetStudentByAddress returns a student by address
func (k Keeper) GetStudentByAddress(ctx sdk.Context, address string) (types.Student, bool) {
	return k.getStudentByAddress(ctx, address)
}

// GetAcademicTreeByStudent returns academic tree by student - FIXED METHOD
func (k Keeper) GetAcademicTreeByStudent(ctx sdk.Context, studentIndex string) (interface{}, bool) {
	tree, found := k.getAcademicTreeByStudent(ctx, studentIndex)
	if !found {
		return nil, false
	}
	return tree, true
}

// GetAcademicTreeByStudentTyped returns academic tree by student with proper type
func (k Keeper) GetAcademicTreeByStudentTyped(ctx sdk.Context, studentIndex string) (types.StudentAcademicTree, bool) {
	return k.getAcademicTreeByStudent(ctx, studentIndex)
}

// SetStudentAcademicTree sets academic tree
func (k Keeper) SetStudentAcademicTree(ctx sdk.Context, academicTree types.StudentAcademicTree) {
	k.setStudentAcademicTree(ctx, academicTree)
}

// GetStudentEnrollments returns enrollments by student (for adapters)
func (k Keeper) GetStudentEnrollments(ctx sdk.Context, studentIndex string) ([]types.StudentEnrollment, error) {
	enrollments := k.getEnrollmentsByStudentId(ctx, studentIndex)
	return enrollments, nil
}

// GetStudentByIndex returns a student by index (for adapters)
func (k Keeper) GetStudentByIndex(ctx sdk.Context, index string) (types.Student, bool) {
	return k.getStudentByIndex(ctx, index)
}

// ============================================================================
// QUERYSERVER INTERFACE IMPLEMENTATION
// ============================================================================

// Params implements the QueryServer interface
func (k Keeper) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(sdkCtx)

	return &types.QueryParamsResponse{
		Params: params,
	}, nil
}

// ListStudents implements the QueryServer interface
func (k Keeper) ListStudents(ctx context.Context, req *types.QueryListStudentsRequest) (*types.QueryListStudentsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	var students []types.Student
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(sdkCtx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StudentKeyPrefix))

	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		var student types.Student
		if err := k.cdc.Unmarshal(value, &student); err != nil {
			return err
		}
		students = append(students, student)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to paginate: %w", err)
	}

	return &types.QueryListStudentsResponse{
		Students:   students,
		Pagination: pageRes,
	}, nil
}

// GetStudent implements the QueryServer interface
func (k Keeper) GetStudent(ctx context.Context, req *types.QueryGetStudentRequest) (*types.QueryGetStudentResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	student, found := k.getStudentByIndex(sdkCtx, req.StudentId)
	if !found {
		return nil, types.ErrStudentNotFound
	}

	return &types.QueryGetStudentResponse{
		Student: student,
	}, nil
}

// ListEnrollments implements the QueryServer interface
func (k Keeper) ListEnrollments(ctx context.Context, req *types.QueryListEnrollmentsRequest) (*types.QueryListEnrollmentsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	var enrollments []types.StudentEnrollment
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(sdkCtx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StudentEnrollmentKeyPrefix))

	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		var enrollment types.StudentEnrollment
		if err := k.cdc.Unmarshal(value, &enrollment); err != nil {
			return err
		}
		enrollments = append(enrollments, enrollment)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to paginate: %w", err)
	}

	return &types.QueryListEnrollmentsResponse{
		Enrollments: enrollments,
		Pagination:  pageRes,
	}, nil
}

// GetEnrollment implements the QueryServer interface
func (k Keeper) GetEnrollment(ctx context.Context, req *types.QueryGetEnrollmentRequest) (*types.QueryGetEnrollmentResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	enrollment, found := k.getStudentEnrollment(sdkCtx, req.EnrollmentId)
	if !found {
		return nil, types.ErrEnrollmentNotFound
	}

	return &types.QueryGetEnrollmentResponse{
		Enrollment: enrollment,
	}, nil
}

// GetEnrollmentsByStudent implements the QueryServer interface (gRPC query)
func (k Keeper) GetEnrollmentsByStudent(ctx context.Context, req *types.QueryGetEnrollmentsByStudentRequest) (*types.QueryGetEnrollmentsByStudentResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if _, found := k.getStudentByIndex(sdkCtx, req.StudentId); !found {
		return nil, types.ErrStudentNotFound
	}

	enrollments := k.getEnrollmentsByStudentId(sdkCtx, req.StudentId)

	startIndex := 0
	endIndex := len(enrollments)

	if req.Pagination != nil {
		if req.Pagination.Offset > 0 {
			startIndex = int(req.Pagination.Offset)
		}
		if req.Pagination.Limit > 0 && startIndex+int(req.Pagination.Limit) < endIndex {
			endIndex = startIndex + int(req.Pagination.Limit)
		}
	}

	if startIndex >= len(enrollments) {
		startIndex = len(enrollments)
		endIndex = len(enrollments)
	}
	if endIndex > len(enrollments) {
		endIndex = len(enrollments)
	}

	paginatedEnrollments := enrollments[startIndex:endIndex]

	return &types.QueryGetEnrollmentsByStudentResponse{
		Enrollments: paginatedEnrollments,
		Pagination: &query.PageResponse{
			NextKey: nil,
			Total:   uint64(len(enrollments)),
		},
	}, nil
}

// GetStudentProgress implements the QueryServer interface
func (k Keeper) GetStudentProgress(ctx context.Context, req *types.QueryGetStudentProgressRequest) (*types.QueryGetStudentProgressResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if _, found := k.getStudentByIndex(sdkCtx, req.StudentId); !found {
		return nil, types.ErrStudentNotFound
	}

	academicTree, found := k.getAcademicTreeByStudent(sdkCtx, req.StudentId)
	if !found {
		return &types.QueryGetStudentProgressResponse{
			Progress: types.AcademicProgress{
				RequiredCreditsCompleted:   0,
				ElectiveCreditsCompleted:   0,
				RequiredSubjectsPercentage: 0.0,
				CurrentSemester:            1,
				CurrentYear:                1,
				EnrollmentYears:            0.0,
				ElectivesByAreaCompleted:   make(map[string]uint64),
			},
		}, nil
	}

	var progress types.AcademicProgress
	if academicTree.AcademicProgress != nil {
		progress = *academicTree.AcademicProgress
	}

	return &types.QueryGetStudentProgressResponse{
		Progress: progress,
	}, nil
}

// GetStudentsByInstitution implements the QueryServer interface
func (k Keeper) GetStudentsByInstitution(ctx context.Context, req *types.QueryGetStudentsByInstitutionRequest) (*types.QueryGetStudentsByInstitutionResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if _, found := k.institutionKeeper.GetInstitution(sdkCtx, req.InstitutionId); !found {
		return nil, types.ErrInvalidInstitution
	}

	enrollments := k.getEnrollmentsByInstitution(sdkCtx, req.InstitutionId)

	studentMap := make(map[string]types.Student)
	for _, enrollment := range enrollments {
		if student, found := k.getStudentByIndex(sdkCtx, enrollment.Student); found {
			studentMap[student.Index] = student
		}
	}

	var students []types.Student
	for _, student := range studentMap {
		students = append(students, student)
	}

	startIndex := 0
	endIndex := len(students)

	if req.Pagination != nil {
		if req.Pagination.Offset > 0 {
			startIndex = int(req.Pagination.Offset)
		}
		if req.Pagination.Limit > 0 && startIndex+int(req.Pagination.Limit) < endIndex {
			endIndex = startIndex + int(req.Pagination.Limit)
		}
	}

	if startIndex >= len(students) {
		startIndex = len(students)
		endIndex = len(students)
	}
	if endIndex > len(students) {
		endIndex = len(students)
	}

	paginatedStudents := students[startIndex:endIndex]

	return &types.QueryGetStudentsByInstitutionResponse{
		Students: paginatedStudents,
		Pagination: &query.PageResponse{
			NextKey: nil,
			Total:   uint64(len(students)),
		},
	}, nil
}

// GetStudentsByCourse implements the QueryServer interface
func (k Keeper) GetStudentsByCourse(ctx context.Context, req *types.QueryGetStudentsByCourseRequest) (*types.QueryGetStudentsByCourseResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if _, found := k.courseKeeper.GetCourse(sdkCtx, req.CourseId); !found {
		return nil, types.ErrInvalidCourse
	}

	enrollments := k.getEnrollmentsByCourse(sdkCtx, req.CourseId)

	studentMap := make(map[string]types.Student)
	for _, enrollment := range enrollments {
		if student, found := k.getStudentByIndex(sdkCtx, enrollment.Student); found {
			studentMap[student.Index] = student
		}
	}

	var students []types.Student
	for _, student := range studentMap {
		students = append(students, student)
	}

	startIndex := 0
	endIndex := len(students)

	if req.Pagination != nil {
		if req.Pagination.Offset > 0 {
			startIndex = int(req.Pagination.Offset)
		}
		if req.Pagination.Limit > 0 && startIndex+int(req.Pagination.Limit) < endIndex {
			endIndex = startIndex + int(req.Pagination.Limit)
		}
	}

	if startIndex >= len(students) {
		startIndex = len(students)
		endIndex = len(students)
	}
	if endIndex > len(students) {
		endIndex = len(students)
	}

	paginatedStudents := students[startIndex:endIndex]

	return &types.QueryGetStudentsByCourseResponse{
		Students: paginatedStudents,
		Pagination: &query.PageResponse{
			NextKey: nil,
			Total:   uint64(len(students)),
		},
	}, nil
}

// GetStudentAcademicTree implements the QueryServer interface
func (k Keeper) GetStudentAcademicTree(ctx context.Context, req *types.QueryGetStudentAcademicTreeRequest) (*types.QueryGetStudentAcademicTreeResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if _, found := k.getStudentByIndex(sdkCtx, req.StudentId); !found {
		return nil, types.ErrStudentNotFound
	}

	academicTree, found := k.getAcademicTreeByStudent(sdkCtx, req.StudentId)
	if !found {
		return nil, types.ErrAcademicTreeNotFound
	}

	return &types.QueryGetStudentAcademicTreeResponse{
		AcademicTree: academicTree,
	}, nil
}

// CheckGraduationEligibility implements the QueryServer interface
func (k Keeper) CheckGraduationEligibility(ctx context.Context, req *types.QueryCheckGraduationEligibilityRequest) (*types.QueryCheckGraduationEligibilityResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if _, found := k.getStudentByIndex(sdkCtx, req.StudentId); !found {
		return nil, types.ErrStudentNotFound
	}

	academicTree, found := k.getAcademicTreeByStudent(sdkCtx, req.StudentId)
	if !found {
		return &types.QueryCheckGraduationEligibilityResponse{
			GraduationStatus: types.GraduationStatus{
				IsEligible:                  false,
				EstimatedGraduationSemester: "Unknown",
				GpaStatus:                   "No data available",
				TimeframeStatus:             "No data available",
				MissingRequirements:         []*types.MissingRequirement{},
			},
			IsEligible: false,
			Message:    "No academic tree found for student",
		}, nil
	}

	var graduationStatus types.GraduationStatus
	if academicTree.GraduationStatus != nil {
		graduationStatus = *academicTree.GraduationStatus
	}

	return &types.QueryCheckGraduationEligibilityResponse{
		GraduationStatus: graduationStatus,
		IsEligible:       graduationStatus.IsEligible,
		Message:          fmt.Sprintf("Graduation eligibility status for student %s", req.StudentId),
	}, nil
}
