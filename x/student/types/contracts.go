package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ============================================================================
// CONTRACT INTEGRATION TYPES
// ============================================================================

// CompletedSubject represents a subject completed by a student
type CompletedSubject struct {
	SubjectId      string `json:"subject_id"`
	Credits        uint64 `json:"credits"`
	CompletionDate string `json:"completion_date"`
	Grade          uint64 `json:"grade"`          // Grade as integer (e.g., 85.5 -> 8550 for precision)
	NftTokenId     string `json:"nft_token_id"`
}

// PrerequisiteGroup represents a group of prerequisites for a subject
type PrerequisiteGroup struct {
	Id                        string   `json:"id"`
	SubjectId                 string   `json:"subject_id"`
	GroupType                 string   `json:"group_type"`                   // "all", "any", "minimum", "none"
	MinimumCredits            uint64   `json:"minimum_credits"`
	MinimumCompletedSubjects  uint64   `json:"minimum_completed_subjects"`
	SubjectIds                []string `json:"subject_ids"`
	Logic                     string   `json:"logic"`                        // "and", "or", "xor", "threshold", "none"
	Priority                  uint64   `json:"priority"`
	Confidence                float64  `json:"confidence"`
}

// EquivalenceResult represents the result of an equivalence analysis
type EquivalenceResult struct {
	EquivalenceId        string `json:"equivalence_id"`
	SimilarityPercentage uint32 `json:"similarity_percentage"`
	Status               string `json:"status"`                 // "pending", "analyzing", "approved", "rejected", "under_review"
	EquivalenceType      string `json:"equivalence_type"`       // "full", "partial", "conditional", "none"
	Notes                string `json:"notes"`
}

// SubjectEnrollmentRequest represents a request to enroll in a subject
type SubjectEnrollmentRequest struct {
	StudentId              string   `json:"student_id"`
	SubjectId              string   `json:"subject_id"`
	CanEnroll              bool     `json:"can_enroll"`
	MissingPrerequisites   []string `json:"missing_prerequisites"`
	VerificationTimestamp  uint64   `json:"verification_timestamp"`
}

// StudentAcademicProgress represents the academic progress tracking
type StudentAcademicProgress struct {
	StudentId                    string            `json:"student_id"`
	CourseId                     string            `json:"course_id"`
	InstitutionId                string            `json:"institution_id"`
	CompletedSubjects            []CompletedSubject `json:"completed_subjects"`
	InProgressSubjects           []string          `json:"in_progress_subjects"`
	AvailableSubjects            []string          `json:"available_subjects"`
	TotalCreditsCompleted        uint64            `json:"total_credits_completed"`
	RequiredCreditsRemaining     uint64            `json:"required_credits_remaining"`
	ElectiveCreditsCompleted     uint64            `json:"elective_credits_completed"`
	CurrentSemester              uint32            `json:"current_semester"`
	EstimatedGraduationSemester  string            `json:"estimated_graduation_semester"`
}

// ContractCallResult represents the result of a contract call
type ContractCallResult struct {
	Success       bool   `json:"success"`
	ErrorMessage  string `json:"error_message,omitempty"`
	TransactionId string `json:"transaction_id,omitempty"`
	BlockHeight   int64  `json:"block_height"`
	GasUsed       uint64 `json:"gas_used"`
}

// ============================================================================
// CONTRACT MESSAGE TYPES
// ============================================================================

// PrerequisiteCheckRequest represents a request to check prerequisites
type PrerequisiteCheckRequest struct {
	StudentId string `json:"student_id"`
	SubjectId string `json:"subject_id"`
}

// PrerequisiteCheckResponse represents the response from prerequisite check
type PrerequisiteCheckResponse struct {
	CanEnroll            bool     `json:"can_enroll"`
	MissingPrerequisites []string `json:"missing_prerequisites"`
	SatisfiedGroups      []string `json:"satisfied_groups"`
	UnsatisfiedGroups    []string `json:"unsatisfied_groups"`
	VerificationTimestamp uint64  `json:"verification_timestamp"`
	Details              string   `json:"details"`
}

// EquivalenceCheckRequest represents a request to check equivalence
type EquivalenceCheckRequest struct {
	SourceSubjectId  string `json:"source_subject_id"`
	TargetSubjectId  string `json:"target_subject_id"`
	ForceRecalculate bool   `json:"force_recalculate"`
}

// EquivalenceCheckResponse represents the response from equivalence check
type EquivalenceCheckResponse struct {
	IsEquivalent         bool   `json:"is_equivalent"`
	EquivalenceType      string `json:"equivalence_type"`
	SimilarityPercentage uint32 `json:"similarity_percentage"`
	EquivalenceId        string `json:"equivalence_id"`
	Status               string `json:"status"`
	Notes                string `json:"notes"`
}

// ============================================================================
// VALIDATION HELPER TYPES
// ============================================================================

// ValidationError represents validation errors from contracts
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// ContractState represents the current state of contracts
type ContractState struct {
	PrerequisitesContractAddr    string `json:"prerequisites_contract_addr"`
	EquivalenceContractAddr      string `json:"equivalence_contract_addr"`
	AcademicProgressContractAddr string `json:"academic_progress_contract_addr"`
	DegreeContractAddr           string `json:"degree_contract_addr"`
	NftMintingContractAddr       string `json:"nft_minting_contract_addr"`
	IsConfigured                 bool   `json:"is_configured"`
	LastUpdated                  string `json:"last_updated"`
}

// ============================================================================
// CONTRACT MESSAGE TYPES FOR ACADEMIC OPERATIONS
// ============================================================================

// SubjectCompletionRequest represents a request to validate subject completion
type SubjectCompletionRequest struct {
	StudentId      string `json:"student_id"`
	SubjectId      string `json:"subject_id"`
	Grade          uint64 `json:"grade"`
	CompletionDate string `json:"completion_date"`
	Semester       string `json:"semester"`
	Credits        uint64 `json:"credits"`
	Institution    string `json:"institution"`
}

// SubjectCompletionResult represents the result of subject completion validation
type SubjectCompletionResult struct {
	Success                   bool             `json:"success"`
	Message                   string           `json:"message"`
	UpdatedCompletedSubjects  []string         `json:"updated_completed_subjects"`
	UpdatedInProgressSubjects []string         `json:"updated_in_progress_subjects"`
	UpdatedProgress           AcademicProgress `json:"updated_progress"`
	ShouldCheckGraduation     bool             `json:"should_check_graduation"`
}

// GraduationEligibilityResult represents graduation eligibility check result
type GraduationEligibilityResult struct {
	IsEligible                  bool     `json:"is_eligible"`
	Message                     string   `json:"message"`
	EstimatedGraduationDate     string   `json:"estimated_graduation_date"`
	RequiredCreditsRemaining    uint64   `json:"required_credits_remaining"`
	RequiredSubjectsRemaining   []string `json:"required_subjects_remaining"`
	MissingElectiveCredits      uint64   `json:"missing_elective_credits"`
	GPARequirementMet           bool     `json:"gpa_requirement_met"`
	TimeframeRequirementMet     bool     `json:"timeframe_requirement_met"`
}

// NFTMintingRequest represents a request to authorize NFT minting
type NFTMintingRequest struct {
	StudentAddress    string           `json:"student_address"`
	SubjectId         string           `json:"subject_id"`
	Grade             uint64           `json:"grade"`
	CompletionDate    string           `json:"completion_date"`
	Semester          string           `json:"semester"`
	IssuerInstitution string           `json:"issuer_institution"`
	ProgressData      AcademicProgress `json:"progress_data"`
}

// NFTMintingResult represents the result of NFT minting authorization
type NFTMintingResult struct {
	Success                  bool   `json:"success"`
	TokenInstanceId          string `json:"token_instance_id"`
	Message                  string `json:"message"`
	MetadataHash            string `json:"metadata_hash"`
	ContractAuthorizationHash string `json:"contract_authorization_hash"`
}

// DegreeValidationRequest represents a request to validate degree requirements
type DegreeValidationRequest struct {
	StudentId     string   `json:"student_id"`
	CurriculumId  string   `json:"curriculum_id"`
	InstitutionId string   `json:"institution_id"`
	FinalGPA      string   `json:"final_gpa"`
	TotalCredits  uint64   `json:"total_credits"`
	Signatures    []string `json:"signatures"`
	RequestedDate string   `json:"requested_date"`
}

// DegreeValidationResult represents the result of degree validation
type DegreeValidationResult struct {
	IsValid           bool     `json:"is_valid"`
	Message           string   `json:"message"`
	DegreeType        string   `json:"degree_type"`
	CurriculumVersion string   `json:"curriculum_version"`
	ValidationHash    string   `json:"validation_hash"`
	RequirementsMet   []string `json:"requirements_met"`
	MissingRequirements []string `json:"missing_requirements"`
}

// DegreeNFTMintingRequest represents a request to mint degree NFT
type DegreeNFTMintingRequest struct {
	StudentId      string                  `json:"student_id"`
	CurriculumId   string                  `json:"curriculum_id"`
	InstitutionId  string                  `json:"institution_id"`
	DegreeType     string                  `json:"degree_type"`
	FinalGPA       string                  `json:"final_gpa"`
	TotalCredits   uint64                  `json:"total_credits"`
	ValidationData DegreeValidationResult  `json:"validation_data"`
	IssueDate      string                  `json:"issue_date"`
}

// DegreeNFTMintingResult represents the result of degree NFT minting
type DegreeNFTMintingResult struct {
	Success          bool   `json:"success"`
	TokenId          string `json:"token_id"`
	MetadataIPFSHash string `json:"metadata_ipfs_hash"`
	Message          string `json:"message"`
}

// SubjectRecommendationRequest represents a request for subject recommendations
type SubjectRecommendationRequest struct {
	StudentId           string   `json:"student_id"`
	CurrentSemester     string   `json:"current_semester"`
	TargetSemester      string   `json:"target_semester"`
	MaxSubjects         uint32   `json:"max_subjects"`
	MaxCredits          uint64   `json:"max_credits"`
	PreferredAreas      []string `json:"preferred_areas"`
	ExcludeSubjects     []string `json:"exclude_subjects"`
	PrioritizeElectives bool     `json:"prioritize_electives"`
}

// SubjectRecommendationResult represents the result of subject recommendations
type SubjectRecommendationResult struct {
	Success          bool                     `json:"success"`
	Recommendations  []RecommendedSubject     `json:"recommendations"`
	Message          string                   `json:"message"`
	AlgorithmUsed    string                   `json:"algorithm_used"`
	ConfidenceScore  float64                  `json:"confidence_score"`
}

// RecommendedSubject represents a recommended subject with context
type RecommendedSubject struct {
	SubjectId       string  `json:"subject_id"`
	SubjectName     string  `json:"subject_name"`
	Credits         uint64  `json:"credits"`
	RecommendationScore float64 `json:"recommendation_score"`
	Reason          string  `json:"reason"`
	PrerequisitesMet bool   `json:"prerequisites_met"`
	IsElective      bool    `json:"is_elective"`
	KnowledgeArea   string  `json:"knowledge_area"`
	Difficulty      string  `json:"difficulty"`
}

// ============================================================================
// EVENT ATTRIBUTE KEYS
// ============================================================================

const (
	EventTypeSubjectCompleted          = "subject_completed"
	EventTypePrerequisiteVerified      = "prerequisite_verified"
	EventTypeEquivalenceRequested      = "equivalence_requested"
	EventTypeContractCallExecuted      = "contract_call_executed"
	EventTypeAcademicProgressUpdated   = "academic_progress_updated"
	EventTypeDegreeValidated           = "degree_validated"
	EventTypeNFTMintingAuthorized      = "nft_minting_authorized"
	EventTypeSubjectRecommended        = "subject_recommended"
	
	AttributeKeyContractAddress        = "contract_address"
	AttributeKeyContractOperation      = "contract_operation"
	AttributeKeyContractResult         = "contract_result"
	AttributeKeyEquivalenceId          = "equivalence_id"
	AttributeKeyPrerequisiteResult     = "prerequisite_result"
	AttributeKeyCompletedSubjectId     = "completed_subject_id"
	AttributeKeyStudentProgress        = "student_progress"
	AttributeKeyNFTTokenId             = "nft_token_id"
	AttributeKeyRecommendationCount    = "recommendation_count"
	AttributeKeyDegreeValidationResult = "degree_validation_result"
)

// ============================================================================
// CONTRACT ERROR TYPES
// ============================================================================

// ContractError represents errors from contract calls
type ContractError struct {
	ContractAddress string `json:"contract_address"`
	Operation       string `json:"operation"`
	ErrorCode       string `json:"error_code"`
	ErrorMessage    string `json:"error_message"`
	BlockHeight     int64  `json:"block_height"`
}

// Error implements the error interface
func (e ContractError) Error() string {
	return fmt.Sprintf("contract error at %s: %s - %s", e.ContractAddress, e.ErrorCode, e.ErrorMessage)
}

// ============================================================================
// HELPER FUNCTIONS FOR CONTRACT RESPONSES
// ============================================================================

// NewContractError creates a new contract error
func NewContractError(contractAddr, operation, errorCode, errorMessage string, blockHeight int64) ContractError {
	return ContractError{
		ContractAddress: contractAddr,
		Operation:       operation,
		ErrorCode:       errorCode,
		ErrorMessage:    errorMessage,
		BlockHeight:     blockHeight,
	}
}

// ValidateContractResponse validates that a contract response is valid
func ValidateContractResponse(response interface{}) error {
	if response == nil {
		return fmt.Errorf("contract response is nil")
	}
	return nil
}

// ============================================================================
// CONTRACT INTEGRATION INTERFACES
// ============================================================================

// AcademicProgressContractInterface defines the interface for academic progress contract
type AcademicProgressContractInterface interface {
	ProcessSubjectCompletion(ctx sdk.Context, request SubjectCompletionRequest) (SubjectCompletionResult, error)
}

// DegreeContractInterface defines the interface for degree contract
type DegreeContractInterface interface {
	ValidateDegreeRequirements(ctx sdk.Context, request DegreeValidationRequest) (DegreeValidationResult, error)
	CheckGraduationEligibility(ctx sdk.Context, studentId string) (GraduationEligibilityResult, error)
}

// NFTMintingContractInterface defines the interface for NFT minting contract
type NFTMintingContractInterface interface {
	AuthorizeMinting(ctx sdk.Context, request NFTMintingRequest) (NFTMintingResult, error)
}
