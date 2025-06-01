package types

import (
	"context"

	institutiontypes "academictoken/x/institution/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SubjectContent defines the locally defined version of Subject content
// Fields aligned with the Subject module's definition
type SubjectContent struct {
	Index         string
	SubjectId     string
	Institution   string
	CourseId      string
	Title         string
	Code          string
	WorkloadHours uint64
	Credits       uint64
	Description   string
	ContentHash   string // Added for hybrid storage support
	SubjectType   string
	KnowledgeArea string
	IPFSLink      string // Added for hybrid storage support
}

// InstitutionKeeper defines the expected interface for the Institution module
type InstitutionKeeper interface {
	// GetInstitution returns an institution by its ID
	GetInstitution(ctx sdk.Context, institutionID string) (institutiontypes.Institution, bool)

	// IsInstitutionAuthorized checks if an institution is authorized
	IsInstitutionAuthorized(ctx sdk.Context, institutionID string) bool

	// GetAuthorizedInstitutions returns all authorized institutions
	GetAuthorizedInstitutions(ctx sdk.Context) []institutiontypes.Institution

	// InstitutionExists checks if an institution exists
	InstitutionExists(ctx sdk.Context, institutionID string) bool
}

// SubjectKeeper defines the expected interface for the Subject module
type SubjectKeeper interface {
	// GetSubject returns a subject by its ID
	GetSubject(ctx sdk.Context, subjectID string) (SubjectContent, bool)

	// GetSubjectsByCourse returns all subjects for a specific course
	GetSubjectsByCourse(ctx sdk.Context, courseID string) []SubjectContent

	// SubjectExists checks if a subject exists
	SubjectExists(ctx sdk.Context, subjectID string) bool
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
