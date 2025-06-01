package types

const (
	// Message types
	TypeMsgMintSubjectToken    = "mint_subject_token"
	TypeMsgVerifyTokenInstance = "verify_token_instance"
	TypeMsgUpdateParams        = "update_params"
)

// ============================================================================
// PASSIVE MODULE EXTENDED TYPES
// ============================================================================

// ExtendedSubjectTokenInstance adds passive mode fields to SubjectTokenInstance
// This extends the protobuf-generated type with additional fields needed for passive mode
type ExtendedSubjectTokenInstance struct {
	SubjectTokenInstance      SubjectTokenInstance `json:"subject_token_instance"`
	SubjectId                 string               `json:"subject_id"`
	ContractAuthorizationHash string               `json:"contract_authorization_hash"`
	MintedByContract          bool                 `json:"minted_by_contract"`
	PassiveModeEnabled        bool                 `json:"passive_mode_enabled"`
}

// ExtendedMsgMintSubjectTokenResponse adds passive mode fields to the response
type ExtendedMsgMintSubjectTokenResponse struct {
	TokenInstanceId    string `json:"token_instance_id"`
	ContractAuthorized bool   `json:"contract_authorized"`
	AuthorizationHash  string `json:"authorization_hash"`
	PassiveExecution   bool   `json:"passive_execution"`
	SubjectId          string `json:"subject_id"`
}

// PassiveModeConfig contains configuration for passive mode operations
type PassiveModeConfig struct {
	Enabled             bool     `json:"enabled"`
	RequireContractAuth bool     `json:"require_contract_auth"`
	AuthorizedContracts []string `json:"authorized_contracts"`
	AllowDirectMinting  bool     `json:"allow_direct_minting"`
	VerificationLevel   string   `json:"verification_level"` // "strict", "moderate", "permissive"
}
