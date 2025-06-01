package keeper

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types" // Add store types import
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"academictoken/x/schedule/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger
		authority    string

		// Expected keepers for interactions with other modules
		subjectKeeper    types.SubjectKeeper
		studentKeeper    types.StudentKeeper
		curriculumKeeper types.CurriculumKeeper
		accountKeeper    types.AccountKeeper
		bankKeeper       types.BankKeeper

		// Add wasmKeeper for contract integration (for future use)
		// wasmKeeper       wasmkeeper.Keeper  // Uncomment when ready for contracts
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	subjectKeeper types.SubjectKeeper,
	studentKeeper types.StudentKeeper,
	curriculumKeeper types.CurriculumKeeper,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	// wasmKeeper wasmkeeper.Keeper, // Uncomment when ready for contracts
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:              cdc,
		storeService:     storeService,
		authority:        authority,
		logger:           logger,
		subjectKeeper:    subjectKeeper,
		studentKeeper:    studentKeeper,
		curriculumKeeper: curriculumKeeper,
		accountKeeper:    accountKeeper,
		bankKeeper:       bankKeeper,
	}
}

func (k Keeper) GetAuthority() string {
	return k.authority
}

func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// CreateStudyPlan creates a new study plan for a student
func (k Keeper) CreateStudyPlan(ctx sdk.Context, msg *types.MsgCreateStudyPlan) (string, error) {
	// Generate unique study plan ID
	studyPlanId := fmt.Sprintf("sp_%s_%d", msg.Student, time.Now().Unix())

	// Create the study plan with only fields that exist in proto
	studyPlan := types.StudyPlan{
		Index:            studyPlanId,
		Student:          msg.Student,
		CreationDate:     time.Now().Format(time.RFC3339),
		CompletionTarget: msg.CompletionTarget,
		AdditionalNotes:  msg.AdditionalNotes,
		Status:           "draft",
		PlannedSemesters: []*types.PlannedSemester{}, // Note: é []*PlannedSemester nos protos
	}

	// Store the study plan
	k.SetStudyPlan(ctx, studyPlan)

	return studyPlanId, nil
}

// UpdateStudyPlanStatus updates the status of a study plan
func (k Keeper) UpdateStudyPlanStatus(ctx sdk.Context, msg *types.MsgUpdateStudyPlanStatus) error {
	studyPlan, found := k.GetStudyPlan(ctx, msg.StudyPlanId)
	if !found {
		return types.ErrStudyPlanNotFound
	}

	// Update only the status field
	studyPlan.Status = msg.Status

	// Store updated study plan
	k.SetStudyPlan(ctx, studyPlan)

	return nil
}

// AddPlannedSemester adds a planned semester to a study plan
func (k Keeper) AddPlannedSemester(ctx sdk.Context, msg *types.MsgAddPlannedSemester) error {
	studyPlan, found := k.GetStudyPlan(ctx, msg.StudyPlanId)
	if !found {
		return types.ErrStudyPlanNotFound
	}

	// Create planned semester with only proto fields
	plannedSemester := &types.PlannedSemester{
		SemesterCode:    msg.SemesterCode,
		PlannedSubjects: msg.PlannedSubjects,
		TotalCredits:    strconv.FormatUint(uint64(len(msg.PlannedSubjects)*3), 10),  // Simple calculation
		TotalHours:      strconv.FormatUint(uint64(len(msg.PlannedSubjects)*60), 10), // Simple calculation
		Status:          "planned",
	}

	// Add to study plan
	studyPlan.PlannedSemesters = append(studyPlan.PlannedSemesters, plannedSemester)

	// Store updated study plan
	k.SetStudyPlan(ctx, studyPlan)

	return nil
}

// CreateSubjectRecommendation creates subject recommendations for a student
func (k Keeper) CreateSubjectRecommendation(ctx sdk.Context, msg *types.MsgCreateSubjectRecommendation) (string, error) {
	// Generate unique recommendation ID
	recommendationId := fmt.Sprintf("sr_%s_%d", msg.Student, time.Now().Unix())

	// Create recommendation subjects with only proto fields
	var recommendedSubjects []*types.RecommendedSubject
	// Split recommendation metadata if it contains subject IDs (comma-separated)
	if msg.RecommendationMetadata != "" {
		// For now, create a single recommendation based on metadata
		recommendedSubject := &types.RecommendedSubject{
			SubjectId:          msg.RecommendationMetadata, // Use metadata as single subject ID
			RecommendationRank: "1",
			Reason:             "Generated recommendation",
			IsRequired:         "false",
			SemesterAlignment:  msg.RecommendationSemester,
			DifficultyLevel:    "medium",
		}
		recommendedSubjects = append(recommendedSubjects, recommendedSubject)
	}

	// Create the subject recommendation with only proto fields
	subjectRecommendation := types.SubjectRecommendation{
		Index:                  recommendationId,
		Student:                msg.Student,
		RecommendationSemester: msg.RecommendationSemester,
		RecommendationMetadata: msg.RecommendationMetadata,
		GeneratedDate:          time.Now().Format(time.RFC3339),
		RecommendedSubjects:    recommendedSubjects,
	}

	// Store the recommendation
	k.SetSubjectRecommendation(ctx, subjectRecommendation)

	return recommendationId, nil
}

// BuildRecommendations generates automatic subject recommendations for a student (internal logic)
func (k Keeper) BuildRecommendations(ctx sdk.Context, studentId string, semesterCode string) ([]types.SubjectRecommendation, error) {
	var recommendations []types.SubjectRecommendation

	// Simple recommendation generation
	recommendationId := fmt.Sprintf("sr_%s_%d", studentId, time.Now().Unix())

	recommendation := types.SubjectRecommendation{
		Index:                  recommendationId,
		Student:                studentId,
		RecommendationSemester: semesterCode,
		RecommendationMetadata: "auto-generated",
		GeneratedDate:          time.Now().Format(time.RFC3339),
		RecommendedSubjects:    []*types.RecommendedSubject{},
	}

	recommendations = append(recommendations, recommendation)
	return recommendations, nil
}

// ============================================================================
// STORAGE OPERATIONS
// ============================================================================

func (k Keeper) SetStudyPlan(ctx sdk.Context, studyPlan types.StudyPlan) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.KeyPrefix(types.StudyPlanPrefix + studyPlan.Index)
	value := k.cdc.MustMarshal(&studyPlan)
	store.Set(key, value)
}

func (k Keeper) GetStudyPlan(ctx sdk.Context, studyPlanId string) (types.StudyPlan, bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.KeyPrefix(types.StudyPlanPrefix + studyPlanId)
	value := store.Get(key)
	if value == nil {
		return types.StudyPlan{}, false
	}

	var studyPlan types.StudyPlan
	k.cdc.MustUnmarshal(value, &studyPlan)
	return studyPlan, true
}

func (k Keeper) SetPlannedSemester(ctx sdk.Context, plannedSemester types.PlannedSemester) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.KeyPrefix(types.PlannedSemesterPrefix + plannedSemester.SemesterCode)
	value := k.cdc.MustMarshal(&plannedSemester)
	store.Set(key, value)
}

func (k Keeper) GetPlannedSemester(ctx sdk.Context, plannedSemesterId string) (types.PlannedSemester, bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.KeyPrefix(types.PlannedSemesterPrefix + plannedSemesterId)
	value := store.Get(key)
	if value == nil {
		return types.PlannedSemester{}, false
	}

	var plannedSemester types.PlannedSemester
	k.cdc.MustUnmarshal(value, &plannedSemester)
	return plannedSemester, true
}

func (k Keeper) SetSubjectRecommendation(ctx sdk.Context, recommendation types.SubjectRecommendation) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.KeyPrefix(types.SubjectRecommendationPrefix + recommendation.Index)
	value := k.cdc.MustMarshal(&recommendation)
	store.Set(key, value)
}

func (k Keeper) GetSubjectRecommendation(ctx sdk.Context, recommendationId string) (types.SubjectRecommendation, bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.KeyPrefix(types.SubjectRecommendationPrefix + recommendationId)
	value := store.Get(key)
	if value == nil {
		return types.SubjectRecommendation{}, false
	}

	var recommendation types.SubjectRecommendation
	k.cdc.MustUnmarshal(value, &recommendation)
	return recommendation, true
}

// ============================================================================
// QUERY METHODS
// ============================================================================

func (k Keeper) GetAllStudyPlans(ctx sdk.Context, pagination *query.PageRequest) ([]types.StudyPlan, *query.PageResponse, error) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	var studyPlans []types.StudyPlan
	pageRes, err := query.Paginate(store, pagination, func(key []byte, value []byte) error {
		if !bytes.HasPrefix(key, []byte(types.StudyPlanPrefix)) {
			return nil
		}
		var studyPlan types.StudyPlan
		if err := k.cdc.Unmarshal(value, &studyPlan); err != nil {
			return err
		}
		studyPlans = append(studyPlans, studyPlan)
		return nil
	})

	return studyPlans, pageRes, err
}

func (k Keeper) GetStudyPlansByStudent(ctx sdk.Context, studentId string) []types.StudyPlan {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	var studyPlans []types.StudyPlan

	iterator := storetypes.KVStorePrefixIterator(store, []byte(types.StudyPlanPrefix))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var studyPlan types.StudyPlan
		if err := k.cdc.Unmarshal(iterator.Value(), &studyPlan); err != nil {
			continue
		}
		if studyPlan.Student == studentId {
			studyPlans = append(studyPlans, studyPlan)
		}
	}

	return studyPlans
}

func (k Keeper) GetSubjectRecommendationsByStudent(ctx sdk.Context, studentId string) []types.SubjectRecommendation {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	var recommendations []types.SubjectRecommendation

	iterator := storetypes.KVStorePrefixIterator(store, []byte(types.SubjectRecommendationPrefix))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var recommendation types.SubjectRecommendation
		if err := k.cdc.Unmarshal(iterator.Value(), &recommendation); err != nil {
			continue
		}
		if recommendation.Student == studentId {
			recommendations = append(recommendations, recommendation)
		}
	}

	return recommendations
}
