package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ============================================================================
// MESSAGE EXTENSIONS FOR PASSIVE MODE
// ============================================================================

// ExtendedMsgMintSubjectToken provides additional fields for passive mode
// This extends the protobuf message with fields needed for contract authorization
type ExtendedMsgMintSubjectToken struct {
	*MsgMintSubjectToken
	SubjectId                string `json:"subject_id,omitempty"`
	ContractAuthorizationHash string `json:"contract_authorization_hash,omitempty"`
}

// NewExtendedMsgMintSubjectToken creates a new extended message for passive mode
func NewExtendedMsgMintSubjectToken(
	creator string,
	tokenDefId string,
	student string,
	completionDate string,
	grade string,
	issuerInstitution string,
	semester string,
	professorSignature string,
	subjectId string,
	contractAuthorizationHash string,
) *ExtendedMsgMintSubjectToken {
	base := &MsgMintSubjectToken{
		Creator:            creator,
		TokenDefId:         tokenDefId,
		Student:            student,
		CompletionDate:     completionDate,
		Grade:              grade,
		IssuerInstitution:  issuerInstitution,
		Semester:           semester,
		ProfessorSignature: professorSignature,
	}
	
	return &ExtendedMsgMintSubjectToken{
		MsgMintSubjectToken:       base,
		SubjectId:                subjectId,
		ContractAuthorizationHash: contractAuthorizationHash,
	}
}

// IsContractAuthorized returns true if this message includes contract authorization
func (msg *ExtendedMsgMintSubjectToken) IsContractAuthorized() bool {
	return msg.SubjectId != "" && msg.ContractAuthorizationHash != ""
}

// ValidateBasic validates the extended message
func (msg *ExtendedMsgMintSubjectToken) ValidateBasic() error {
	// First validate the base message
	if err := msg.MsgMintSubjectToken.ValidateBasic(); err != nil {
		return err
	}
	
	// If passive mode fields are provided, both must be present
	if msg.SubjectId != "" || msg.ContractAuthorizationHash != "" {
		if msg.SubjectId == "" {
			return ErrMissingContractAuthorization.Wrap("subjectId required for contract-authorized minting")
		}
		if msg.ContractAuthorizationHash == "" {
			return ErrMissingContractAuthorization.Wrap("contractAuthorizationHash required for contract-authorized minting")
		}
	}
	
	return nil
}

// ToBaseMessage returns the base protobuf message
func (msg *ExtendedMsgMintSubjectToken) ToBaseMessage() *MsgMintSubjectToken {
	return msg.MsgMintSubjectToken
}

// ============================================================================
// MESSAGE HELPERS FOR PASSIVE MODE COMPATIBILITY
// ============================================================================

// WithPassiveModeFields adds passive mode fields to a standard message
func (msg *MsgMintSubjectToken) WithPassiveModeFields(subjectId, authHash string) *ExtendedMsgMintSubjectToken {
	return &ExtendedMsgMintSubjectToken{
		MsgMintSubjectToken:       msg,
		SubjectId:                subjectId,
		ContractAuthorizationHash: authHash,
	}
}

// GetSubjectId returns the subject ID if available
func (msg *MsgMintSubjectToken) GetSubjectId() string {
	// For compatibility, check if this message has been extended
	// In practice, this would be handled by the caller
	return ""
}

// GetContractAuthorizationHash returns the authorization hash if available
func (msg *MsgMintSubjectToken) GetContractAuthorizationHash() string {
	// For compatibility, check if this message has been extended
	// In practice, this would be handled by the caller
	return ""
}

// ============================================================================
// PASSIVE MODE MESSAGE PROCESSING HELPERS
// ============================================================================

// ProcessPassiveModeMessage processes a message for passive mode
func ProcessPassiveModeMessage(msg *MsgMintSubjectToken, subjectId, authHash string) *ExtendedMsgMintSubjectToken {
	return &ExtendedMsgMintSubjectToken{
		MsgMintSubjectToken:       msg,
		SubjectId:                subjectId,
		ContractAuthorizationHash: authHash,
	}
}

// ExtractPassiveModeData extracts passive mode data from context or message
func ExtractPassiveModeData(ctx sdk.Context, msg *MsgMintSubjectToken) (subjectId, authHash string, isPassive bool) {
	// Check for passive mode data in transaction context
	// This is a placeholder for actual implementation
	
	// In practice, this would check:
	// 1. Transaction metadata
	// 2. Message annotations
	// 3. Context values set by contract calls
	
	return "", "", false
}

// ============================================================================
// CONTRACT AUTHORIZATION VALIDATION
// ============================================================================

// ValidateContractAuthorization validates contract authorization for a message
func ValidateContractAuthorization(extMsg *ExtendedMsgMintSubjectToken) error {
	if !extMsg.IsContractAuthorized() {
		return ErrMissingContractAuthorization
	}
	
	// Validate subject ID format
	if extMsg.SubjectId == "" {
		return ErrMissingContractAuthorization.Wrap("subject ID cannot be empty")
	}
	
	// Validate authorization hash format
	if len(extMsg.ContractAuthorizationHash) != 64 { // SHA-256 hex string
		return ErrInvalidContractAuthorization.Wrap("invalid authorization hash format")
	}
	
	return nil
}

// ============================================================================
// MESSAGE CONVERSION UTILITIES
// ============================================================================

// ConvertToExtended converts a standard message to extended if passive mode data exists
func ConvertToExtended(msg *MsgMintSubjectToken, ctx sdk.Context) (*ExtendedMsgMintSubjectToken, bool) {
	// Check if passive mode data is available in context
	subjectId, authHash, isPassive := ExtractPassiveModeData(ctx, msg)
	
	if isPassive {
		return ProcessPassiveModeMessage(msg, subjectId, authHash), true
	}
	
	return nil, false
}

// ConvertFromExtended converts an extended message back to standard format
func ConvertFromExtended(extMsg *ExtendedMsgMintSubjectToken) *MsgMintSubjectToken {
	return extMsg.MsgMintSubjectToken
}
