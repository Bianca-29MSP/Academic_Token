package types

import (
	"context"

	institutionmoduletypes "academictoken/x/institution/types"
	tokendefmoduletypes "academictoken/x/tokendef/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Institution defines the locally defined version of institution type
type Institution struct {
	Index        string
	Address      string
	Name         string
	IsAuthorized string // "true", "false", "pending", etc.
	Creator      string
}

// InstitutionKeeper defines the expected interface for the Institution module
type InstitutionKeeper interface {
	// GetInstitution returns an institution by ID - using real method signature
	GetInstitution(ctx sdk.Context, index string) (institutionmoduletypes.Institution, bool)

	// IsInstitutionAuthorized checks if institution is authorized - using real method
	IsInstitutionAuthorized(ctx sdk.Context, index string) bool
}

// TokenMetadata defines metadata for tokens
type TokenMetadata struct {
	Description string
	ImageUri    string
	// Simplified - we don't need full attributes here
}

// TokenDefinition defines the locally defined version of token definition
type TokenDefinition struct {
	Index          string
	TokenDefId     string
	SubjectId      string
	InstitutionId  string
	CourseId       string
	TokenName      string
	TokenSymbol    string
	TokenType      string
	IsTransferable bool
	IsBurnable     bool
	MaxSupply      uint64
	Metadata       TokenMetadata
	ContentHash    string
	IpfsLink       string
	Creator        string
	CreatedAt      string
}

// TokenDefKeeper defines the expected interface for the TokenDef module
type TokenDefKeeper interface {
	// GetTokenDefinitionByIndex returns a token definition by its index - using real method name
	GetTokenDefinitionByIndex(ctx sdk.Context, index string) (tokendefmoduletypes.TokenDefinition, bool)
}

// Student defines the locally defined version of student
type Student struct {
	Index       string
	Address     string
	Name        string
	Institution string
}

// StudentKeeper defines the expected interface for the Student module
type StudentKeeper interface {
	// GetStudentByAddress returns a student by blockchain address
	GetStudentByAddress(ctx sdk.Context, address sdk.AccAddress) (Student, bool)

	// AddCompletedSubject adds a completed subject token to student's record
	AddCompletedSubject(ctx sdk.Context, studentAddress string, tokenInstanceID string) error

	// ValidateStudentEligibility checks if student can receive tokens from institution
	ValidateStudentEligibility(ctx sdk.Context, studentAddress string, institutionID string) error
}

// ParamSubspace defines the expected Subspace interface for parameters
type ParamSubspace interface {
	Get(ctx context.Context, key []byte, ptr interface{})
	Set(ctx context.Context, key []byte, param interface{})
}

// AccountKeeper defines the expected interface for the Account module
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
}

// BankKeeper defines the expected interface for the Bank module
type BankKeeper interface {
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
}
