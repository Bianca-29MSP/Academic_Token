package types

import (
	"fmt"
	"time"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// Parameter store keys
var (
	KeyMaxCreditsPerSemester        = []byte("MaxCreditsPerSemester")
	KeyMaxPlannedSemesters          = []byte("MaxPlannedSemesters")
	KeyRecommendationWeight         = []byte("RecommendationWeight")
	KeyIPFSTimeout                  = []byte("IPFSTimeout")
	KeyMinimumGradeForProgress      = []byte("MinimumGradeForProgress")
	KeyAllowedRecommendationTypes   = []byte("AllowedRecommendationTypes")
	KeyDefaultDifficultyLevels      = []byte("DefaultDifficultyLevels")
	KeyMaxStudyPlansPerStudent      = []byte("MaxStudyPlansPerStudent")
	KeyRecommendationScoreThreshold = []byte("RecommendationScoreThreshold")
	KeyDefaultSemesterDuration      = []byte("DefaultSemesterDuration")
)

// Default parameter values
const (
	DefaultMaxCreditsPerSemester        uint64        = 24
	DefaultMaxPlannedSemesters          uint64        = 16
	DefaultRecommendationWeight         float32       = 0.75
	DefaultIPFSTimeout                  time.Duration = 30 * time.Second
	DefaultMinimumGradeForProgress      float32       = 6.0
	DefaultMaxStudyPlansPerStudent      uint64        = 5
	DefaultRecommendationScoreThreshold float32       = 0.6
	DefaultDefaultSemesterDuration      uint64        = 6 // months
)

// Default parameter values for arrays
var (
	DefaultAllowedRecommendationTypes = []string{
		RecommendationTypeNext,
		RecommendationTypeElective,
		RecommendationTypePrereq,
		RecommendationTypeOptimal,
		RecommendationTypeIntensive,
	}
	DefaultDefaultDifficultyLevels = []string{
		DifficultyEasy,
		DifficultyMedium,
		DifficultyHard,
	}
)

// ParamKeyTable returns the parameter key table for use with the keeper
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance with default values
func NewParams() Params {
	return Params{
		MaxCreditsPerSemester:        DefaultMaxCreditsPerSemester,
		MaxPlannedSemesters:          DefaultMaxPlannedSemesters,
		RecommendationWeight:         DefaultRecommendationWeight,
		IpfsTimeout:                  DefaultIPFSTimeout.String(),
		MinimumGradeForProgress:      DefaultMinimumGradeForProgress,
		AllowedRecommendationTypes:   DefaultAllowedRecommendationTypes,
		DefaultDifficultyLevels:      DefaultDefaultDifficultyLevels,
		MaxStudyPlansPerStudent:      DefaultMaxStudyPlansPerStudent,
		RecommendationScoreThreshold: DefaultRecommendationScoreThreshold,
		DefaultSemesterDuration:      DefaultDefaultSemesterDuration,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs returns the parameter set pairs
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMaxCreditsPerSemester, &p.MaxCreditsPerSemester, validateMaxCreditsPerSemester),
		paramtypes.NewParamSetPair(KeyMaxPlannedSemesters, &p.MaxPlannedSemesters, validateMaxPlannedSemesters),
		paramtypes.NewParamSetPair(KeyRecommendationWeight, &p.RecommendationWeight, validateRecommendationWeight),
		paramtypes.NewParamSetPair(KeyIPFSTimeout, &p.IpfsTimeout, validateIPFSTimeout),
		paramtypes.NewParamSetPair(KeyMinimumGradeForProgress, &p.MinimumGradeForProgress, validateMinimumGradeForProgress),
		paramtypes.NewParamSetPair(KeyAllowedRecommendationTypes, &p.AllowedRecommendationTypes, validateAllowedRecommendationTypes),
		paramtypes.NewParamSetPair(KeyDefaultDifficultyLevels, &p.DefaultDifficultyLevels, validateDefaultDifficultyLevels),
		paramtypes.NewParamSetPair(KeyMaxStudyPlansPerStudent, &p.MaxStudyPlansPerStudent, validateMaxStudyPlansPerStudent),
		paramtypes.NewParamSetPair(KeyRecommendationScoreThreshold, &p.RecommendationScoreThreshold, validateRecommendationScoreThreshold),
		paramtypes.NewParamSetPair(KeyDefaultSemesterDuration, &p.DefaultSemesterDuration, validateDefaultSemesterDuration),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateMaxCreditsPerSemester(p.MaxCreditsPerSemester); err != nil {
		return err
	}
	if err := validateMaxPlannedSemesters(p.MaxPlannedSemesters); err != nil {
		return err
	}
	if err := validateRecommendationWeight(p.RecommendationWeight); err != nil {
		return err
	}
	if err := validateIPFSTimeout(p.IpfsTimeout); err != nil {
		return err
	}
	if err := validateMinimumGradeForProgress(p.MinimumGradeForProgress); err != nil {
		return err
	}
	if err := validateAllowedRecommendationTypes(p.AllowedRecommendationTypes); err != nil {
		return err
	}
	if err := validateDefaultDifficultyLevels(p.DefaultDifficultyLevels); err != nil {
		return err
	}
	if err := validateMaxStudyPlansPerStudent(p.MaxStudyPlansPerStudent); err != nil {
		return err
	}
	if err := validateRecommendationScoreThreshold(p.RecommendationScoreThreshold); err != nil {
		return err
	}
	if err := validateDefaultSemesterDuration(p.DefaultSemesterDuration); err != nil {
		return err
	}
	return nil
}

// Validation functions

func validateMaxCreditsPerSemester(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("max credits per semester must be positive: %d", v)
	}

	if v > 50 {
		return fmt.Errorf("max credits per semester too high: %d", v)
	}

	return nil
}

func validateMaxPlannedSemesters(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("max planned semesters must be positive: %d", v)
	}

	if v > 50 {
		return fmt.Errorf("max planned semesters too high: %d", v)
	}

	return nil
}

func validateRecommendationWeight(i interface{}) error {
	v, ok := i.(float32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0.0 || v > 1.0 {
		return fmt.Errorf("recommendation weight must be between 0.0 and 1.0: %f", v)
	}

	return nil
}

func validateIPFSTimeout(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	_, err := time.ParseDuration(v)
	if err != nil {
		return fmt.Errorf("invalid IPFS timeout duration: %s", v)
	}

	return nil
}

func validateMinimumGradeForProgress(i interface{}) error {
	v, ok := i.(float32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0.0 || v > 10.0 {
		return fmt.Errorf("minimum grade for progress must be between 0.0 and 10.0: %f", v)
	}

	return nil
}

func validateAllowedRecommendationTypes(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if len(v) == 0 {
		return fmt.Errorf("allowed recommendation types cannot be empty")
	}

	validTypes := map[string]bool{
		RecommendationTypeNext:      true,
		RecommendationTypeElective:  true,
		RecommendationTypePrereq:    true,
		RecommendationTypeOptimal:   true,
		RecommendationTypeIntensive: true,
	}

	for _, recType := range v {
		if !validTypes[recType] {
			return fmt.Errorf("invalid recommendation type: %s", recType)
		}
	}

	return nil
}

func validateDefaultDifficultyLevels(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if len(v) == 0 {
		return fmt.Errorf("default difficulty levels cannot be empty")
	}

	validLevels := map[string]bool{
		DifficultyEasy:   true,
		DifficultyMedium: true,
		DifficultyHard:   true,
	}

	for _, level := range v {
		if !validLevels[level] {
			return fmt.Errorf("invalid difficulty level: %s", level)
		}
	}

	return nil
}

func validateMaxStudyPlansPerStudent(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("max study plans per student must be positive: %d", v)
	}

	if v > 20 {
		return fmt.Errorf("max study plans per student too high: %d", v)
	}

	return nil
}

func validateRecommendationScoreThreshold(i interface{}) error {
	v, ok := i.(float32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0.0 || v > 1.0 {
		return fmt.Errorf("recommendation score threshold must be between 0.0 and 1.0: %f", v)
	}

	return nil
}

func validateDefaultSemesterDuration(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("default semester duration must be positive: %d", v)
	}

	if v > 12 {
		return fmt.Errorf("default semester duration too high (max 12 months): %d", v)
	}

	return nil
}
