package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
}

// InstitutionKeeper defines the expected interface for the Institution module.
type InstitutionKeeper interface {
	GetInstitution(ctx sdk.Context, index string) (Institution, bool)
	IsInstitutionAuthorized(ctx sdk.Context, index string) bool
}

// CourseKeeper defines the expected interface for the Course module.
type CourseKeeper interface {
	GetCourse(ctx sdk.Context, index string) (Course, bool)
	HasCourse(ctx sdk.Context, index string) bool
}

// SubjectKeeper defines the expected interface for the Subject module.
type SubjectKeeper interface {
	GetSubject(ctx sdk.Context, index string) (SubjectContent, bool)
	HasSubject(ctx sdk.Context, index string) bool
}

// Types needed from the Institution module
type Institution struct {
	Index        string
	Address      string
	Name         string
	IsAuthorized string
	Creator      string
}

// Types needed from the Course module
type Course struct {
	Index        string
	Institution  string
	Name         string
	Code         string
	Description  string
	TotalCredits string
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
	IpfsLink      string
	Creator       string
}
