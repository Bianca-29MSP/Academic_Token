package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// Methods for accounts
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	// Methods for bank operations
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(ctx sdk.Context, key []byte, ptr interface{})
	Set(ctx sdk.Context, key []byte, param interface{})
	GetParamSet(ctx sdk.Context, ps paramtypes.ParamSet)
	SetParamSet(ctx sdk.Context, ps paramtypes.ParamSet)
}

// Student represents a simplified student for avoiding circular imports
type Student struct {
	Id            string
	InstitutionId string
	Status        string
}

// AcademicRecord represents a simplified academic record
type AcademicRecord struct {
	StudentId         string
	CompletedCredits  uint64
	GPA               string
	CompletedSubjects []string
}

// Curriculum represents a simplified curriculum
type Curriculum struct {
	Id               string
	InstitutionId    string
	RequiredCredits  uint64
	RequiredSubjects []string
}

// DegreeNFTData represents data for creating degree NFT
type DegreeNFTData struct {
	StudentId     string
	CurriculumId  string
	InstitutionId string
	DegreeType    string
	IssueDate     string
	GPA           string
	TotalCredits  uint64
	IPFSHash      string
}

// AcademicNFT represents a simplified academic NFT
type AcademicNFT struct {
	TokenId   string
	Owner     string
	Metadata  string
	TokenType string
}

// StudentKeeper defines the expected student keeper interface
type StudentKeeper interface {
	GetStudent(ctx sdk.Context, id string) (Student, bool)
	GetStudentByIndex(ctx sdk.Context, index string) (Student, bool)
	GetStudentAcademicRecord(ctx sdk.Context, studentId string) (AcademicRecord, bool)
	GetAcademicTreeByStudent(ctx sdk.Context, studentId string) (interface{}, bool)
	GetStudentsByInstitution(ctx sdk.Context, institutionId string) []Student
	ValidateStudentExists(ctx sdk.Context, studentId string) error
	GetStudentGPA(ctx sdk.Context, studentId string) (string, error)
	GetStudentTotalCredits(ctx sdk.Context, studentId string) (uint64, error)
	GetCompletedSubjects(ctx sdk.Context, studentId string) ([]string, error)
	GetContractIntegration() interface{} // Returns ContractIntegration from student module
}

// CurriculumKeeper defines the expected curriculum keeper interface
type CurriculumKeeper interface {
	GetCurriculum(ctx sdk.Context, id string) (Curriculum, bool)
	GetCurriculumTree(ctx sdk.Context, id string) (interface{}, bool)
	ValidateCurriculumRequirements(ctx sdk.Context, curriculumId string, completedSubjects []string) error
	GetCurriculumRequiredCredits(ctx sdk.Context, curriculumId string) (uint64, error)
	GetCurriculumRequiredSubjects(ctx sdk.Context, curriculumId string) ([]string, error)
	GetCurriculumsByInstitution(ctx sdk.Context, institutionId string) []Curriculum
}

// AcademicNFTKeeper defines the expected academic NFT keeper interface
type AcademicNFTKeeper interface {
	MintDegreeNFT(ctx sdk.Context, recipient string, degreeData DegreeNFTData) (string, error)
	GetNFTByTokenID(ctx sdk.Context, tokenId string) (AcademicNFT, bool)
	TransferNFT(ctx sdk.Context, from string, to string, tokenId string) error
	BurnNFT(ctx sdk.Context, tokenId string) error
	ValidateNFTOwnership(ctx sdk.Context, owner string, tokenId string) error
}

// WasmKeeper defines the expected wasm keeper interface for CosmWasm integration
type WasmKeeper interface {
	Sudo(ctx sdk.Context, contractAddress sdk.AccAddress, msg []byte) ([]byte, error)
	QuerySmart(ctx context.Context, contractAddr sdk.AccAddress, req []byte) ([]byte, error)
	Execute(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress, msg []byte, coins sdk.Coins) ([]byte, error)
}

// InstitutionKeeper defines the expected institution keeper interface
type InstitutionKeeper interface {
	GetInstitution(ctx sdk.Context, id string) (interface{}, bool) // Using interface{} to avoid circular import
	ValidateInstitutionAuthorization(ctx sdk.Context, institutionId string) error
	IsAuthorizedToIssueDegrees(ctx sdk.Context, institutionId string) bool
}
