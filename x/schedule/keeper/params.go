package keeper

import (
	"context"
	"time"

	"academictoken/x/schedule/types"
)

// HARDCODED SCHEDULE CONFIGURATION
const (
	// Schedule recommendation contract address
	SCHEDULE_CONTRACT_ADDRESS = "cosmos1schedulecontract123456789abcdef"

	// Contract version for compatibility checks
	SCHEDULE_CONTRACT_VERSION = "1.0.0"

	// Contract execution gas limit
	SCHEDULE_CONTRACT_GAS_LIMIT = uint64(800000) // Higher limit for complex algorithms

	// IPFS configuration
	IPFS_GATEWAY_URL = "https://ipfs.io/ipfs/"
	IPFS_TIMEOUT     = 30 * time.Second

	// Schedule limits and defaults
	MAX_CREDITS_PER_SEMESTER       = uint64(24)
	MAX_PLANNED_SEMESTERS          = uint64(16)
	RECOMMENDATION_WEIGHT          = float32(0.75)
	MINIMUM_GRADE_FOR_PROGRESS     = float32(6.0)
	MAX_STUDY_PLANS_PER_STUDENT    = uint64(5)
	RECOMMENDATION_SCORE_THRESHOLD = float32(0.6)
	DEFAULT_SEMESTER_DURATION      = uint64(6) // months
)

// Default arrays
var (
	ALLOWED_RECOMMENDATION_TYPES = []string{
		types.RecommendationTypeNext,
		types.RecommendationTypeElective,
		types.RecommendationTypePrereq,
		types.RecommendationTypeOptimal,
		types.RecommendationTypeIntensive,
	}

	DEFAULT_DIFFICULTY_LEVELS = []string{
		types.DifficultyEasy,
		types.DifficultyMedium,
		types.DifficultyHard,
	}
)

// GetScheduleContractAddress returns the hardcoded contract address
func (k Keeper) GetScheduleContractAddress(ctx context.Context) string {
	return SCHEDULE_CONTRACT_ADDRESS
}

// GetScheduleContractVersion returns the hardcoded contract version
func (k Keeper) GetScheduleContractVersion(ctx context.Context) string {
	return SCHEDULE_CONTRACT_VERSION
}

// GetScheduleContractGasLimit returns the hardcoded gas limit
func (k Keeper) GetScheduleContractGasLimit(ctx context.Context) uint64 {
	return SCHEDULE_CONTRACT_GAS_LIMIT
}

// GetIPFSGatewayURL returns the hardcoded IPFS gateway URL
func (k Keeper) GetIPFSGatewayURL(ctx context.Context) string {
	return IPFS_GATEWAY_URL
}

// GetIPFSTimeout returns the hardcoded IPFS timeout
func (k Keeper) GetIPFSTimeout(ctx context.Context) time.Duration {
	return IPFS_TIMEOUT
}

// GetMaxCreditsPerSemester returns the hardcoded max credits per semester
func (k Keeper) GetMaxCreditsPerSemester(ctx context.Context) uint64 {
	return MAX_CREDITS_PER_SEMESTER
}

// GetMaxPlannedSemesters returns the hardcoded max planned semesters
func (k Keeper) GetMaxPlannedSemesters(ctx context.Context) uint64 {
	return MAX_PLANNED_SEMESTERS
}

// GetRecommendationWeight returns the hardcoded recommendation weight
func (k Keeper) GetRecommendationWeight(ctx context.Context) float32 {
	return RECOMMENDATION_WEIGHT
}

// GetMinimumGradeForProgress returns the hardcoded minimum grade for progress
func (k Keeper) GetMinimumGradeForProgress(ctx context.Context) float32 {
	return MINIMUM_GRADE_FOR_PROGRESS
}

// GetMaxStudyPlansPerStudent returns the hardcoded max study plans per student
func (k Keeper) GetMaxStudyPlansPerStudent(ctx context.Context) uint64 {
	return MAX_STUDY_PLANS_PER_STUDENT
}

// GetRecommendationScoreThreshold returns the hardcoded recommendation score threshold
func (k Keeper) GetRecommendationScoreThreshold(ctx context.Context) float32 {
	return RECOMMENDATION_SCORE_THRESHOLD
}

// GetDefaultSemesterDuration returns the hardcoded default semester duration
func (k Keeper) GetDefaultSemesterDuration(ctx context.Context) uint64 {
	return DEFAULT_SEMESTER_DURATION
}

// GetAllowedRecommendationTypes returns the hardcoded allowed recommendation types
func (k Keeper) GetAllowedRecommendationTypes(ctx context.Context) []string {
	return ALLOWED_RECOMMENDATION_TYPES
}

// GetDefaultDifficultyLevels returns the hardcoded default difficulty levels
func (k Keeper) GetDefaultDifficultyLevels(ctx context.Context) []string {
	return DEFAULT_DIFFICULTY_LEVELS
}

// DEPRECATED: Remove after migration
func (k Keeper) GetParams(ctx context.Context) types.Params {
	return types.Params{
		MaxCreditsPerSemester:        MAX_CREDITS_PER_SEMESTER,
		MaxPlannedSemesters:          MAX_PLANNED_SEMESTERS,
		RecommendationWeight:         RECOMMENDATION_WEIGHT,
		IpfsTimeout:                  IPFS_TIMEOUT.String(),
		MinimumGradeForProgress:      MINIMUM_GRADE_FOR_PROGRESS,
		AllowedRecommendationTypes:   ALLOWED_RECOMMENDATION_TYPES,
		DefaultDifficultyLevels:      DEFAULT_DIFFICULTY_LEVELS,
		MaxStudyPlansPerStudent:      MAX_STUDY_PLANS_PER_STUDENT,
		RecommendationScoreThreshold: RECOMMENDATION_SCORE_THRESHOLD,
		DefaultSemesterDuration:      DEFAULT_SEMESTER_DURATION,
	}
}

// DEPRECATED: Remove after migration
func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	// No-op for hardcoded params
	return nil
}
