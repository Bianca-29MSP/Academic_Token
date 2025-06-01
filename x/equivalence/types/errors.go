package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/equivalence module sentinel errors
var (
	ErrInvalidSigner              = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrEquivalenceNotFound        = sdkerrors.Register(ModuleName, 1101, "equivalence not found")
	ErrEquivalenceAlreadyExists   = sdkerrors.Register(ModuleName, 1102, "equivalence already exists")
	ErrInvalidEquivalenceStatus   = sdkerrors.Register(ModuleName, 1103, "invalid equivalence status")
	ErrInvalidEquivalencePercent  = sdkerrors.Register(ModuleName, 1104, "invalid equivalence percentage")
	ErrSubjectNotFound           = sdkerrors.Register(ModuleName, 1105, "subject not found")
	ErrInstitutionNotFound       = sdkerrors.Register(ModuleName, 1106, "institution not found")
	ErrInstitutionNotAuthorized  = sdkerrors.Register(ModuleName, 1107, "institution not authorized")
	ErrContractAnalysisFailed    = sdkerrors.Register(ModuleName, 1108, "contract analysis failed")
	ErrInvalidContractAddress    = sdkerrors.Register(ModuleName, 1109, "invalid contract address")
	ErrAnalysisIntegrityFailed   = sdkerrors.Register(ModuleName, 1110, "analysis integrity verification failed")
	ErrInvalidAnalysisMetadata   = sdkerrors.Register(ModuleName, 1111, "invalid analysis metadata")
	ErrEquivalenceUpdateFailed   = sdkerrors.Register(ModuleName, 1112, "equivalence update failed")
	ErrBatchRequestFailed        = sdkerrors.Register(ModuleName, 1113, "batch request failed")
	ErrInvalidRequest            = sdkerrors.Register(ModuleName, 1114, "invalid request")
	
	// Additional errors for contract integration
	ErrContractNotAvailable      = sdkerrors.Register(ModuleName, 1115, "equivalence contract not available")
	ErrContractCallFailed        = sdkerrors.Register(ModuleName, 1116, "contract call failed")
	ErrInvalidContractResponse   = sdkerrors.Register(ModuleName, 1117, "invalid contract response")
	ErrContractAuthRequired      = sdkerrors.Register(ModuleName, 1118, "contract authorization required")
	ErrInvalidAnalysisHash       = sdkerrors.Register(ModuleName, 1119, "invalid analysis hash")
	ErrAnalysisTimeout           = sdkerrors.Register(ModuleName, 1120, "analysis timeout")
	ErrMaxRetriesExceeded        = sdkerrors.Register(ModuleName, 1121, "maximum analysis retries exceeded")
	ErrInvalidAnalysisParams     = sdkerrors.Register(ModuleName, 1122, "invalid analysis parameters")
	ErrContractVersionMismatch   = sdkerrors.Register(ModuleName, 1123, "contract version mismatch")
	ErrIPFSAccessFailed          = sdkerrors.Register(ModuleName, 1124, "IPFS access failed")
	ErrInvalidIPFSHash           = sdkerrors.Register(ModuleName, 1125, "invalid IPFS hash")
	ErrSubjectContentNotFound    = sdkerrors.Register(ModuleName, 1126, "subject content not found")
	ErrInvalidSubjectData        = sdkerrors.Register(ModuleName, 1127, "invalid subject data")
	ErrAnalysisConfigInvalid     = sdkerrors.Register(ModuleName, 1128, "analysis configuration invalid")
	ErrPermissionDenied          = sdkerrors.Register(ModuleName, 1129, "permission denied")
)

// ============================================================================
// ERROR HELPER FUNCTIONS
// ============================================================================

// IsContractError checks if an error is related to contract operations
func IsContractError(err error) bool {
	return sdkerrors.IsOf(err, ErrContractNotAvailable) ||
		sdkerrors.IsOf(err, ErrContractCallFailed) ||
		sdkerrors.IsOf(err, ErrInvalidContractResponse) ||
		sdkerrors.IsOf(err, ErrContractAuthRequired) ||
		sdkerrors.IsOf(err, ErrContractVersionMismatch)
}

// IsAnalysisError checks if an error is related to analysis operations
func IsAnalysisError(err error) bool {
	return sdkerrors.IsOf(err, ErrContractAnalysisFailed) ||
		sdkerrors.IsOf(err, ErrInvalidAnalysisMetadata) ||
		sdkerrors.IsOf(err, ErrInvalidAnalysisHash) ||
		sdkerrors.IsOf(err, ErrAnalysisTimeout) ||
		sdkerrors.IsOf(err, ErrMaxRetriesExceeded) ||
		sdkerrors.IsOf(err, ErrInvalidAnalysisParams) ||
		sdkerrors.IsOf(err, ErrAnalysisConfigInvalid)
}

// IsIPFSError checks if an error is related to IPFS operations
func IsIPFSError(err error) bool {
	return sdkerrors.IsOf(err, ErrIPFSAccessFailed) ||
		sdkerrors.IsOf(err, ErrInvalidIPFSHash)
}

// IsValidationError checks if an error is related to validation
func IsValidationError(err error) bool {
	return sdkerrors.IsOf(err, ErrInvalidEquivalencePercent) ||
		sdkerrors.IsOf(err, ErrInvalidEquivalenceStatus) ||
		sdkerrors.IsOf(err, ErrInvalidContractAddress) ||
		sdkerrors.IsOf(err, ErrInvalidSubjectData) ||
		sdkerrors.IsOf(err, ErrInvalidRequest)
}

// GetErrorCode returns the error code for an error
func GetErrorCode(err error) uint32 {
	if sdkErr, ok := err.(*sdkerrors.Error); ok {
		return sdkErr.ABCICode()
	}
	return 0
}

// WrapContractError wraps a contract-related error with additional context
func WrapContractError(err error, context string) error {
	return sdkerrors.Wrapf(ErrContractCallFailed, "%s: %v", context, err)
}

// WrapAnalysisError wraps an analysis-related error with additional context
func WrapAnalysisError(err error, context string) error {
	return sdkerrors.Wrapf(ErrContractAnalysisFailed, "%s: %v", context, err)
}

// WrapIPFSError wraps an IPFS-related error with additional context
func WrapIPFSError(err error, context string) error {
	return sdkerrors.Wrapf(ErrIPFSAccessFailed, "%s: %v", context, err)
}
