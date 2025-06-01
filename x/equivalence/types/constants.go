package types

// EquivalenceStatus constants for different states of equivalence analysis
const (
	EquivalenceStatusPending  = "pending"
	EquivalenceStatusApproved = "approved"
	EquivalenceStatusRejected = "rejected"
	EquivalenceStatusError    = "error"
)

// IsValidEquivalenceStatus checks if the status is valid
func IsValidEquivalenceStatus(status string) bool {
	return status == EquivalenceStatusPending ||
		status == EquivalenceStatusApproved ||
		status == EquivalenceStatusRejected ||
		status == EquivalenceStatusError
}

// ContractVersion constants for tracking different versions of analysis contracts
const (
	DefaultContractVersion = "v1.0.0"
)

// Validation constants
const (
	MaxEquivalencePercent = "100.00"
	MinEquivalencePercent = "0.00"
	
	// Minimum equivalence percentage to be considered approved
	MinApprovalThreshold = "75.00"
)

// Index generation helper
func GenerateEquivalenceIndex(sourceSubjectId, targetSubjectId string) string {
	return sourceSubjectId + "-" + targetSubjectId
}

// Event types for equivalence operations
const (
	EventTypeEquivalenceRequested = "equivalence_requested"
	EventTypeEquivalenceApproved  = "equivalence_approved"
	EventTypeEquivalenceRejected  = "equivalence_rejected"
	EventTypeEquivalenceUpdated   = "equivalence_updated"
)

// Event attribute keys
const (
	AttributeKeyEquivalenceId      = "equivalence_id"
	AttributeKeySourceSubject      = "source_subject"
	AttributeKeyTargetInstitution  = "target_institution"
	AttributeKeyTargetSubject      = "target_subject"
	AttributeKeyRequester          = "requester"
	AttributeKeyApprover           = "approver"
	AttributeKeyRejector           = "rejector"
	AttributeKeyStatus             = "status"
	AttributeKeyPreviousStatus     = "previous_status"
	AttributeKeyEquivalencePercent = "equivalence_percent"
	AttributeKeyAnalysisCount      = "analysis_count"
	AttributeKeyTimestamp          = "timestamp"
	AttributeKeyContractAddress    = "contract_address"
	AttributeKeyContractVersion    = "contract_version"
)

// ValidateEquivalencePercent validates that the equivalence percentage is within valid range
func ValidateEquivalencePercent(percent string) bool {
	if percent == "" {
		return true // Allow empty for pending analysis
	}
	// Add more sophisticated validation here if needed
	// For now, just check it's not obviously invalid
	return len(percent) > 0 && len(percent) <= 6 // e.g., "100.00"
}
