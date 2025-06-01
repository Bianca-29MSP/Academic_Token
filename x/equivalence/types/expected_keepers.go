package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SubjectContent defines the locally defined version of Subject content
// to avoid direct dependency on subject module
type SubjectContent struct {
	Index         string
	SubjectId     string
	Institution   string
	Title         string
	Code          string
	WorkloadHours uint64
	Credits       uint64
	ContentHash   string
	TopicUnits    []string
	Keywords      []string
	KnowledgeArea string
	IPFSLink      string
}

// SubjectKeeper defines the expected interface for the Subject module
type SubjectKeeper interface {
	// GetSubject returns a subject by its ID
	GetSubject(ctx sdk.Context, subjectID string) (SubjectContent, bool)

	// GetSubjectsByInstitution returns subjects offered by an institution
	GetSubjectsByInstitution(ctx sdk.Context, institutionID string) []SubjectContent

	// CompareSubjectContent compares two subjects for similarity
	CompareSubjectContent(ctx sdk.Context, sourceSubjectID string, targetSubjectID string) (float64, error)
}

// Institution defines the locally defined version of Institution
// to avoid direct dependency on institution module
type Institution struct {
	Index        string
	Name         string
	IsAuthorized bool
}

// InstitutionKeeper defines the expected interface for the Institution module
type InstitutionKeeper interface {
	// GetInstitution returns an institution by its ID
	GetInstitution(ctx sdk.Context, institutionID string) (Institution, bool)

	// IsInstitutionAuthorized checks if an institution is authorized
	IsInstitutionAuthorized(ctx sdk.Context, institutionID string) bool
}

// WasmKeeper defines the expected interface for CosmWasm integration
type WasmKeeper interface {
	// QuerySmart executes a smart contract query
	QuerySmart(ctx sdk.Context, contractAddr sdk.AccAddress, req []byte) ([]byte, error)
	
	// Execute executes a smart contract
	Execute(ctx sdk.Context, contractAddr sdk.AccAddress, caller sdk.AccAddress, msg []byte, coins sdk.Coins) ([]byte, error)
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
