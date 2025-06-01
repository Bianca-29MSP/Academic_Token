package types

const (
	// ModuleName defines the module name
	ModuleName = "degree"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_degree"

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

// Storage prefixes
var (
	// DegreePrefix is the prefix to retrieve all Degree
	DegreePrefix = []byte("degree/value/")

	// DegreeRequestPrefix is the prefix to retrieve all DegreeRequest
	DegreeRequestPrefix = []byte("degree_request/value/")

	// DegreeCountPrefix tracks the total count of degrees
	DegreeCountPrefix = []byte("degree/count/")

	// DegreeRequestCountPrefix tracks the total count of degree requests
	DegreeRequestCountPrefix = []byte("degree_request/count/")

	// DegreeByStudentPrefix indexes degrees by student ID
	DegreeByStudentPrefix = []byte("degree/by_student/")

	// DegreeRequestByStudentPrefix indexes degree requests by student ID
	DegreeRequestByStudentPrefix = []byte("degree_request/by_student/")

	// DegreeByCurriculumPrefix indexes degrees by curriculum ID
	DegreeByCurriculumPrefix = []byte("degree/by_curriculum/")

	// DegreeByStatusPrefix indexes degrees by status
	DegreeByStatusPrefix = []byte("degree/by_status/")

	// DegreeRequestByStatusPrefix indexes degree requests by status
	DegreeRequestByStatusPrefix = []byte("degree_request/by_status/")
)

// Event types
const (
	EventTypeDegreeRequested    = "degree_requested"
	EventTypeDegreeValidated    = "degree_validated"
	EventTypeDegreeIssued       = "degree_issued"
	EventTypeDegreeRejected     = "degree_rejected"
	EventTypeContractUpdated    = "degree_contract_updated"
	EventTypeValidationStarted  = "degree_validation_started"
	EventTypeValidationComplete = "degree_validation_complete"
)

// Event attributes
const (
	AttributeKeyDegreeID               = "degree_id"
	AttributeKeyDegreeRequestID        = "degree_request_id"
	AttributeKeyStudentID              = "student_id"
	AttributeKeyCurriculumID           = "curriculum_id"
	AttributeKeyInstitutionID          = "institution_id"
	AttributeKeyContractAddress        = "contract_address"
	AttributeKeyDegreeStatus           = "degree_status"
	AttributeKeyValidationScore        = "validation_score"
	AttributeKeyIssueDate              = "issue_date"
	AttributeKeyGPA                    = "gpa"
	AttributeKeyTotalCredits           = "total_credits"
	AttributeKeyDegreeType             = "degree_type"
	AttributeKeyNFTTokenID             = "nft_token_id"
	AttributeKeyIPFSHash               = "ipfs_hash"
	AttributeKeyRejectionReason        = "rejection_reason"
	AttributeKeyExpectedGraduationDate = "expected_graduation_date"
	AttributeKeyValidationPassed       = "validation_passed"
	AttributeKeyFinalGPA               = "final_gpa"
)

// Degree statuses
const (
	DegreeStatusRequested  = "requested"
	DegreeStatusValidating = "validating"
	DegreeStatusValidated  = "validated"
	DegreeStatusIssued     = "issued"
	DegreeStatusRejected   = "rejected"
	DegreeStatusCancelled  = "cancelled"
)

// Degree request statuses
const (
	DegreeRequestStatusPending          = "pending"
	DegreeRequestStatusProcessing       = "processing"
	DegreeRequestStatusApproved         = "approved"
	DegreeRequestStatusRejected         = "rejected"
	DegreeRequestStatusCancelled        = "cancelled"
	DegreeRequestStatusValidated        = "validated"
	DegreeRequestStatusValidationFailed = "validation_failed"
)

// Default values
const (
	DefaultMinimumGPA       = "2.0"
	DefaultValidationPeriod = 7 * 24 * 3600 // 7 days in seconds
	DefaultContractVersion  = "1.0.0"
)
