package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// Parameter store keys - HARDCODED KEYS
var (
	KeyEquivalenceContractAddress = []byte("EquivalenceContractAddress")
	KeyIPFSGateway               = []byte("IPFSGateway") 
	KeyIPFSEnabled               = []byte("IPFSEnabled")
	KeyMinApprovalThreshold      = []byte("MinApprovalThreshold")
	KeyMaxAnalysisRetries        = []byte("MaxAnalysisRetries")
	KeyAnalysisTimeoutSeconds    = []byte("AnalysisTimeoutSeconds")
	KeyRequireContractAuth       = []byte("RequireContractAuth")
	KeyAdmin                     = []byte("Admin")
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance with HARDCODED values
func NewParams() Params {
	return Params{
		// Note: Since protobuf Params is empty, we'll store hardcoded values in functions
		// and return empty Params struct for compatibility
	}
}

// DefaultParams returns HARDCODED default parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet with HARDCODED keys
// Note: Since protobuf Params is empty, we'll use empty ParamSetPairs for now
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		// When protobuf params are added, uncomment these:
		// paramtypes.NewParamSetPair(KeyEquivalenceContractAddress, &p.EquivalenceContractAddress, validateString),
		// paramtypes.NewParamSetPair(KeyIPFSGateway, &p.IpfsGateway, validateString),
		// paramtypes.NewParamSetPair(KeyIPFSEnabled, &p.IpfsEnabled, validateBool),
		// paramtypes.NewParamSetPair(KeyMinApprovalThreshold, &p.MinApprovalThreshold, validateString),
		// paramtypes.NewParamSetPair(KeyMaxAnalysisRetries, &p.MaxAnalysisRetries, validateUint64),
		// paramtypes.NewParamSetPair(KeyAnalysisTimeoutSeconds, &p.AnalysisTimeoutSeconds, validateUint64),
		// paramtypes.NewParamSetPair(KeyRequireContractAuth, &p.RequireContractAuth, validateBool),
		// paramtypes.NewParamSetPair(KeyAdmin, &p.Admin, validateString),
	}
}

// Validate validates the HARDCODED params
func (p Params) Validate() error {
	// Since protobuf Params is empty, validation always passes
	return nil
}

// ============================================================================
// HARDCODED PARAMETER GETTERS (Since protobuf Params is empty)
// ============================================================================

// GetHardcodedEquivalenceContractAddress returns the hardcoded contract address
func GetHardcodedEquivalenceContractAddress() string {
	return "cosmos1mockequivalencecontract"
}

// GetHardcodedIPFSGateway returns the hardcoded IPFS gateway
func GetHardcodedIPFSGateway() string {
	return "http://localhost:5001"
}

// IsHardcodedIPFSEnabled returns if IPFS is enabled (hardcoded true)
func IsHardcodedIPFSEnabled() bool {
	return true
}

// GetHardcodedMinApprovalThreshold returns the hardcoded minimum approval threshold
func GetHardcodedMinApprovalThreshold() string {
	return "70.0"
}

// GetHardcodedMaxAnalysisRetries returns the hardcoded max analysis retries
func GetHardcodedMaxAnalysisRetries() uint64 {
	return 3
}

// GetHardcodedAnalysisTimeoutSeconds returns the hardcoded analysis timeout
func GetHardcodedAnalysisTimeoutSeconds() uint64 {
	return 300
}

// IsHardcodedContractAuthRequired returns if contract auth is required (hardcoded true)
func IsHardcodedContractAuthRequired() bool {
	return true
}

// GetHardcodedAdmin returns the hardcoded admin address (empty)
func GetHardcodedAdmin() string {
	return ""
}

// ============================================================================
// HELPER FUNCTIONS FOR EQUIVALENCE MODULE (FIXED)
// ============================================================================

// GetDefaultContractConfig returns default contract configuration
func GetDefaultContractConfig() map[string]interface{} {
	return map[string]interface{}{
		"contract_address":        GetHardcodedEquivalenceContractAddress(),
		"ipfs_gateway":           GetHardcodedIPFSGateway(),
		"ipfs_enabled":           IsHardcodedIPFSEnabled(),
		"min_approval_threshold": GetHardcodedMinApprovalThreshold(),
		"max_analysis_retries":   GetHardcodedMaxAnalysisRetries(),
		"analysis_timeout":       GetHardcodedAnalysisTimeoutSeconds(),
		"require_contract_auth":  IsHardcodedContractAuthRequired(),
		"admin":                  GetHardcodedAdmin(),
	}
}

// GetDefaultAnalysisConfig returns default analysis configuration - FIXED SIGNATURE
func GetDefaultAnalysisConfig() AnalysisConfig {
	return AnalysisConfig{
		SimilarityAlgorithm:    DefaultSimilarityAlgorithm,
		ConfidenceThreshold:    DefaultConfidenceThreshold,
		TextSimilarityWeight:   DefaultTextSimilarityWeight,
		TopicSimilarityWeight:  DefaultTopicSimilarityWeight,
		DefaultLanguage:        DefaultLanguage,
		AnalysisMode:           DefaultAnalysisMode,
		MaxRetries:             GetHardcodedMaxAnalysisRetries(),
		TimeoutSeconds:         GetHardcodedAnalysisTimeoutSeconds(),
	}
}

// IsValidEquivalencePercent validates equivalence percentage
func IsValidEquivalencePercent(percent string) bool {
	if percent == "" {
		return false
	}
	// Additional validation can be added here
	return true
}

// IsValidContractAddress validates contract address format
func IsValidContractAddress(address string) bool {
	if address == "" {
		return false
	}
	// Additional validation can be added here
	return true
}
