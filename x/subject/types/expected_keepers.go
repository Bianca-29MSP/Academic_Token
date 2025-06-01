package types

import (
	"context"

	coursetypes "academictoken/x/course/types"
	institutiontypes "academictoken/x/institution/types"

	"github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI
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

// StudentKeeper defines the methods from the Student module that Subject needs
type StudentKeeper interface {
	GetStudent(ctx sdk.Context, index string) (Student, bool)
	GetAcademicTree(ctx sdk.Context, index string) (StudentAcademicTree, bool)
	GetCompletedSubjects(ctx sdk.Context, studentId string) []string
}

// InstitutionKeeper defines the expected institution keeper
type InstitutionKeeper interface {
	GetInstitution(ctx sdk.Context, institutionID string) (institutiontypes.Institution, bool)
	IsInstitutionAuthorized(ctx sdk.Context, institutionID string) bool
}

// CourseKeeper defines the expected interface for the course keeper
type CourseKeeper interface {
	HasCourse(ctx sdk.Context, courseId string) bool
	GetCourse(ctx sdk.Context, courseId string) (coursetypes.Course, bool)
}

// WasmQuerier defines the expected interface for CosmWasm queries
type WasmQuerier interface {
	SmartContractState(ctx context.Context, req *types.QuerySmartContractStateRequest) (*types.QuerySmartContractStateResponse, error)
}

// Types needed from the Student module
type Student struct {
	Index         string
	Address       string
	Name          string
	EnrollmentIds []string
}

type StudentAcademicTree struct {
	Index            string
	Student          string
	CompletedTokens  []string
	InProgressTokens []string
	AvailableTokens  []string
}
