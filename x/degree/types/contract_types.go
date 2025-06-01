package types

// ============================================================================
// CONTRACT INTEGRATION TYPES FOR DEGREE MODULE
// ============================================================================

// DegreeValidationRequest represents a request to validate degree requirements via contract
type DegreeValidationRequest struct {
	StudentId     string   `json:"student_id"`
	CurriculumId  string   `json:"curriculum_id"`
	InstitutionId string   `json:"institution_id"`
	FinalGPA      string   `json:"final_gpa"`
	TotalCredits  uint64   `json:"total_credits"`
	Signatures    []string `json:"signatures"`
	RequestedDate string   `json:"requested_date"`
}

// DegreeValidationResult represents the result from Degree Validation Contract
type DegreeValidationResult struct {
	IsValid             bool     `json:"is_valid"`
	Message             string   `json:"message"`
	DegreeType          string   `json:"degree_type"`
	CurriculumVersion   string   `json:"curriculum_version"`
	ValidationHash      string   `json:"validation_hash"`
	RequirementsMet     []string `json:"requirements_met"`
	MissingRequirements []string `json:"missing_requirements"`
}

// DegreeNFTMintingRequest represents a request to mint degree NFT via contract
type DegreeNFTMintingRequest struct {
	StudentId      string                 `json:"student_id"`
	CurriculumId   string                 `json:"curriculum_id"`
	InstitutionId  string                 `json:"institution_id"`
	DegreeType     string                 `json:"degree_type"`
	FinalGPA       string                 `json:"final_gpa"`
	TotalCredits   uint64                 `json:"total_credits"`
	ValidationData DegreeValidationResult `json:"validation_data"`
	IssueDate      string                 `json:"issue_date"`
}

// DegreeNFTMintingResult represents the result from Degree NFT Minting Contract
type DegreeNFTMintingResult struct {
	Success          bool   `json:"success"`
	TokenId          string `json:"token_id"`
	MetadataIPFSHash string `json:"metadata_ipfs_hash"`
	Message          string `json:"message"`
}

// DegreeVerificationRequest represents a request to verify a degree via contract
type DegreeVerificationRequest struct {
	DegreeId        string `json:"degree_id"`
	TokenId         string `json:"token_id"`
	StudentId       string `json:"student_id"`
	InstitutionId   string `json:"institution_id"`
	VerifierAddress string `json:"verifier_address"`
}

// DegreeVerificationResult represents the result from Degree Verification Contract
type DegreeVerificationResult struct {
	IsValid          bool   `json:"is_valid"`
	IsAuthentic      bool   `json:"is_authentic"`
	Message          string `json:"message"`
	VerificationHash string `json:"verification_hash"`
	VerifiedAt       string `json:"verified_at"`
}

// ContractDegreeNFTData represents the NFT data for a degree (used by contracts)
// This is separate from the DegreeNFTData in expected_keepers.go to avoid conflicts
type ContractDegreeNFTData struct {
	StudentId     string `json:"student_id"`
	CurriculumId  string `json:"curriculum_id"`
	InstitutionId string `json:"institution_id"`
	DegreeType    string `json:"degree_type"`
	IssueDate     string `json:"issue_date"`
	GPA           string `json:"gpa"`
	TotalCredits  uint64 `json:"total_credits"`
}

// ============================================================================
// CONTRACT-SPECIFIC CONSTANTS
// ============================================================================

// Additional degree statuses specific to contract operations
const (
	DegreeStatusRevoked   = "revoked"
	DegreeStatusSuspended = "suspended"
)

// Additional event types specific to contract operations
const (
	EventTypeDegreeVerified = "degree_verified"
)
