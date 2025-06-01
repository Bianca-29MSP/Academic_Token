package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}

// CourseKeeper defines the methods from the Course module that Curriculum needs
type CourseKeeper interface {
	GetCourse(ctx sdk.Context, index string) (Course, bool)
	CourseExists(ctx sdk.Context, index string) bool
}

// SubjectKeeper defines the methods from the Subject module that Curriculum needs
type SubjectKeeper interface {
	GetSubject(ctx sdk.Context, index string) (SubjectContent, bool)
	SubjectExists(ctx sdk.Context, index string) bool
}

// Types needed from the Course module
type Course struct {
	Index        string
	Institution  string
	Name         string
	Code         string
	Description  string
	TotalCredits string // Changed to string to match proto
	DegreeLevel  string
}

// Types needed from the Subject module
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
	ContentHash   string
	SubjectType   string
	KnowledgeArea string
	IPFSLink      string
}
