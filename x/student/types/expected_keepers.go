package types

import (
	"context"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
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

// WasmMsgServer defines the expected interface for WASM message server.
type WasmMsgServer interface {
	ExecuteContract(context.Context, *wasmtypes.MsgExecuteContract) (*wasmtypes.MsgExecuteContractResponse, error)
}

// WasmQuerier defines the expected interface for WASM querier.
type WasmQuerier interface {
	SmartContractState(context.Context, *wasmtypes.QuerySmartContractStateRequest) (*wasmtypes.QuerySmartContractStateResponse, error)
}

// InstitutionKeeper defines the methods from the Institution module that Student needs
type InstitutionKeeper interface {
	GetInstitution(ctx sdk.Context, index string) (Institution, bool)
	IsInstitutionAuthorized(ctx sdk.Context, institutionIndex string) bool
}

// CourseKeeper defines the methods from the Course module that Student needs
type CourseKeeper interface {
	GetCourse(ctx sdk.Context, index string) (Course, bool)
	GetCoursesByInstitution(ctx sdk.Context, institutionIndex string) []Course
}

// CurriculumKeeper defines the methods from the Curriculum module that Student needs
type CurriculumKeeper interface {
	GetCurriculumTree(ctx context.Context, index string) (CurriculumTree, bool)
	GetCurriculumTreesByCourse(ctx context.Context, courseId string) []CurriculumTree
}

// SubjectKeeper defines the methods from the Subject module that Student needs
type SubjectKeeper interface {
	GetSubject(ctx sdk.Context, subjectId string) (SubjectContent, bool)
	SubjectExists(ctx sdk.Context, index string) bool                                                                                           // ✅ Existe no Subject keeper
	CheckPrerequisitesViaContract(ctx sdk.Context, studentID string, subjectID string) (bool, []string, error)                                  // ✅ Existe no Subject keeper
	CheckEquivalenceViaContract(ctx sdk.Context, sourceSubjectID string, targetSubjectID string, forceRecalculate bool) (uint64, string, error) // ✅ Existe no Subject keeper
}

// TokenDefKeeper defines the methods from the TokenDef module that Student needs
type TokenDefKeeper interface {
	GetTokenDefinitionByIndex(ctx sdk.Context, index string) (TokenDefinition, bool)
	GetTokenDefinitionsBySubject(ctx sdk.Context, subjectId string) []TokenDefinition
}

// AcademicNFTKeeper defines the methods from the AcademicNFT module that Student needs
type AcademicNFTKeeper interface {
	GetSubjectTokenInstance(ctx sdk.Context, tokenInstanceId string) (SubjectTokenInstance, bool)
	GetStudentTokenInstances(ctx sdk.Context, studentAddress string) ([]SubjectTokenInstance, error)
}

// AcademicNFTMsgServer defines message server methods that Student might need to call
type AcademicNFTMsgServer interface {
	MintSubjectToken(goCtx context.Context, msg *MsgMintSubjectToken) (*MsgMintSubjectTokenResponse, error)          // ✅ Existe como MsgServer
	VerifyTokenInstance(goCtx context.Context, msg *MsgVerifyTokenInstance) (*MsgVerifyTokenInstanceResponse, error) // ✅ Existe como MsgServer
}

// Types needed from the Institution module
type Institution struct {
	Index        string
	Name         string
	Address      string
	IsAuthorized bool
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

// Types needed from the Curriculum module
type CurriculumTree struct {
	Index                  string
	CourseId               string
	Version                string
	RequiredSubjects       []string
	ElectiveMin            uint64
	ElectiveSubjects       []string
	TotalWorkloadHours     uint64
	GraduationRequirements GraduationRequirements
}

type GraduationRequirements struct {
	TotalCreditsRequired    uint64
	MinGPA                  float64
	RequiredElectiveCredits uint64
	RequiredActivities      []string
	MinimumTimeYears        float64
	MaximumTimeYears        float64
}

// Types needed from the Subject module
type SubjectContent struct {
	Index         string
	SubjectId     string
	Institution   string
	Title         string
	Code          string
	WorkloadHours uint64
	Credits       uint64
	Description   string
	ContentHash   string // Hash do conteúdo completo no IPFS
	SubjectType   string
	KnowledgeArea string
	IPFSLink      string // Link para o conteúdo completo no IPFS
}

// Types needed from the TokenDef module
type TokenDefinition struct {
	Index          string
	SubjectId      string
	TokenName      string
	TokenSymbol    string
	Description    string
	TokenType      string
	IsTransferable bool
	IsBurnable     bool
	MaxSupply      uint64
	ImageUri       string
	ContentHash    string
	IPFSLink       string
}

// Types needed from the AcademicNFT module
type SubjectTokenInstance struct {
	TokenInstanceId    string
	TokenDefId         string
	Student            string
	CompletionDate     string
	Grade              string
	IssuerInstitution  string
	Semester           string
	ProfessorSignature string
	MintedAt           string
	IsValid            bool
}

// Message types needed for AcademicNFT operations
type MsgMintSubjectToken struct {
	Creator            string
	TokenDefId         string
	Student            string
	CompletionDate     string
	Grade              string
	IssuerInstitution  string
	Semester           string
	ProfessorSignature string
}

type MsgMintSubjectTokenResponse struct {
	TokenInstanceId string
}

type MsgVerifyTokenInstance struct {
	Creator         string
	TokenInstanceId string
}

type MsgVerifyTokenInstanceResponse struct {
	IsValid bool
	Message string
}
