package types

const (
	// ModuleName defines the module name
	ModuleName = "student"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_student"
)

const (
	ParamsKey = "params"
)

// Student store keys
const (
	StudentKeyPrefix  = "Student/value/"
	StudentCounterKey = "Student/count/"
)

// StudentEnrollment store keys
const (
	StudentEnrollmentKeyPrefix  = "StudentEnrollment/value/"
	StudentEnrollmentCounterKey = "StudentEnrollment/count/"
)

// StudentAcademicTree store keys
const (
	StudentAcademicTreeKeyPrefix  = "StudentAcademicTree/value/"
	StudentAcademicTreeCounterKey = "StudentAcademicTree/count/"
)

// SubjectEnrollment store keys
const (
	SubjectEnrollmentKeyPrefix  = "SubjectEnrollment/value/"
	SubjectEnrollmentCounterKey = "SubjectEnrollment/count/"
)

// KeyPrefix returns the key prefix for the given string
func KeyPrefix(p string) []byte {
	return []byte(p)
}

// Event types
const (
	EventTypeRegisterStudent          = "register_student"
	EventTypeCreateEnrollment         = "create_enrollment"
	EventTypeUpdateEnrollmentStatus   = "update_enrollment_status"
	EventTypeRequestSubjectEnrollment = "request_subject_enrollment"
	EventTypeUpdateAcademicTree       = "update_academic_tree"
)

// Event attribute keys
const (
	AttributeKeyStudentIndex   = "student_index"
	AttributeKeyStudentAddress = "student_address"
	AttributeKeyStudentName    = "student_name"

	AttributeKeyEnrollmentId   = "enrollment_id"
	AttributeKeyStudent        = "student"
	AttributeKeyInstitution    = "institution"
	AttributeKeyCourseId       = "course_id"
	AttributeKeyEnrollmentDate = "enrollment_date"
	AttributeKeyStatus         = "status"
	AttributeKeyOldStatus      = "old_status"
	AttributeKeyNewStatus      = "new_status"
	AttributeKeyUpdatedBy      = "updated_by"

	AttributeKeySubjectId      = "subject_id"
	AttributeKeyAcademicTreeId = "academic_tree_id"

	AttributeKeyCompletedTokens  = "completed_tokens"
	AttributeKeyInProgressTokens = "in_progress_tokens"
	AttributeKeyAvailableTokens  = "available_tokens"
)
