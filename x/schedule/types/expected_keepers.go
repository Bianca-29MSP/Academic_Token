package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SubjectContent defines the locally defined version of Subject content
// to avoid direct dependency on subject module
type SubjectContent struct {
	Index              string
	SubjectId          string
	Institution        string
	Title              string
	Code               string
	WorkloadHours      uint64
	Credits            uint64
	SubjectType        string // "Required" or "Elective"
	KnowledgeArea      string
	PrerequisiteGroups []PrerequisiteGroup
	DifficultyLevel    string // "Easy", "Medium", "Hard"
}

// PrerequisiteGroup defines a group of prerequisites with logical operation
type PrerequisiteGroup struct {
	GroupType                string // "ALL" for AND, "ANY" for OR
	SubjectIds               []string
	MinimumCredits           uint64
	MinimumCompletedSubjects uint64
}

// SubjectKeeper defines the expected interface for the Subject module
type SubjectKeeper interface {
	// GetSubject returns a subject by its ID
	GetSubject(ctx sdk.Context, subjectID string) (SubjectContent, bool)

	// CheckPrerequisites checks if prerequisites are met for a subject
	CheckPrerequisites(ctx sdk.Context, studentID string, subjectID string) (bool, []string, error)

	// GetSubjectsByArea returns subjects from a specific knowledge area
	GetSubjectsByArea(ctx sdk.Context, area string) []SubjectContent
}

// StudentAcademicTree defines the locally defined version of student academic progress
// to avoid direct dependency on student module
type StudentAcademicTree struct {
	Index               string
	Student             string
	Institution         string
	CourseId            string
	CurriculumVersion   string
	CompletedTokens     []string
	InProgressTokens    []string
	AvailableTokens     []string
	TotalCredits        uint64
	TotalCompletedHours uint64
	CoefficientGPA      float64
}

// StudentKeeper defines the expected interface for the Student module
type StudentKeeper interface {
	// GetAcademicTree returns the academic progress tree for a student
	GetAcademicTree(ctx sdk.Context, studentID string, courseID string) (StudentAcademicTree, bool)

	// GetCompletedSubjects returns a list of completed subject IDs
	GetCompletedSubjects(ctx sdk.Context, studentID string) []string

	// GetInProgressSubjects returns a list of in-progress subject IDs
	GetInProgressSubjects(ctx sdk.Context, studentID string) []string
}

// CurriculumTree defines the locally defined version of curriculum tree
// to avoid direct dependency on curriculum module
type CurriculumTree struct {
	Index             string
	CourseId          string
	Version           string
	RequiredSubjects  []string
	ElectiveSubjects  []string
	SemesterStructure []CurriculumSemester
	ElectiveGroups    []ElectiveGroup
}

// CurriculumSemester defines the structure of a semester in the curriculum
type CurriculumSemester struct {
	SemesterNumber uint64
	SubjectIds     []string
}

// ElectiveGroup defines a group of elective subjects
type ElectiveGroup struct {
	GroupId             string
	Name                string
	SubjectIds          []string
	MinSubjectsRequired uint64
	CreditsRequired     uint64
	KnowledgeArea       string
}

// CurriculumKeeper defines the expected interface for the Curriculum module
type CurriculumKeeper interface {
	// GetCurriculumTree returns a curriculum tree by course ID and version
	GetCurriculumTree(ctx sdk.Context, courseID string, version string) (CurriculumTree, bool)

	// GetCurrentCurriculumVersion returns the current version of a curriculum
	GetCurrentCurriculumVersion(ctx sdk.Context, courseID string) string
}

// ParamSubspace defines the expected Subspace interface for parameters
type ParamSubspace interface {
	Get(ctx context.Context, key []byte, ptr interface{})
	Set(ctx context.Context, key []byte, param interface{})
}

// AccountKeeper defines the expected interface for the Account module
// Kept for compatibility with Cosmos SDK
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
}

// BankKeeper defines the expected interface for the Bank module
// Kept for compatibility with Cosmos SDK
type BankKeeper interface {
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
}
