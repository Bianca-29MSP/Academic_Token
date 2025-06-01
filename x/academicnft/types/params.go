package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// Parameter store keys
var (
	KeyIPFSGateway = []byte("IpfsGateway")
	KeyIPFSEnabled = []byte("IpfsEnabled")
	KeyAdmin       = []byte("Admin")
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	ipfsGateway string,
	ipfsEnabled bool,
	admin string,
) Params {
	return Params{
		IpfsGateway: ipfsGateway,
		IpfsEnabled: ipfsEnabled,
		Admin:       admin,
	}
}

// CreateDefaultPassiveModeConfig creates default passive mode configuration
func CreateDefaultPassiveModeConfig() PassiveModeConfig {
	return PassiveModeConfig{
		Enabled:             true,
		RequireContractAuth: true,
		AuthorizedContracts: []string{},
		AllowDirectMinting:  false,
		VerificationLevel:   "strict",
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		"http://localhost:5001", // default IPFS gateway
		true,                    // IPFS enabled by default
		"",                      // empty admin
	)
}

// DefaultParamsWithPassiveMode returns default params
func DefaultParamsWithPassiveMode() Params {
	return DefaultParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyIPFSGateway, &p.IpfsGateway, validateString),
		paramtypes.NewParamSetPair(KeyIPFSEnabled, &p.IpfsEnabled, validateBool),
		paramtypes.NewParamSetPair(KeyAdmin, &p.Admin, validateString),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateString(p.IpfsGateway); err != nil {
		return err
	}
	if err := validateString(p.Admin); err != nil {
		return err
	}
	return nil
}

// GetPassiveModeConfig returns the default passive mode configuration
// Note: In a full implementation, this configuration might be stored separately
// or retrieved from a different source than the basic Params
func (p Params) GetPassiveModeConfig() PassiveModeConfig {
	return CreateDefaultPassiveModeConfig()
}

// ============================================================================
// VALIDATION FUNCTIONS
// ============================================================================

func validateString(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

// validateStringSlice validates a slice of strings
func validateStringSlice(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// Validate each contract address in the slice
	for _, addr := range v {
		if addr == "" {
			return fmt.Errorf("empty contract address not allowed")
		}
		// Add more specific address validation if needed
	}

	return nil
}

// validateVerificationLevel validates verification level values
func validateVerificationLevel(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	validLevels := []string{"strict", "moderate", "permissive"}
	for _, level := range validLevels {
		if v == level {
			return nil
		}
	}

	return fmt.Errorf("invalid verification level: %s. Valid options: %v", v, validLevels)
}

// ValidatePassiveModeConfig validates a PassiveModeConfig
func ValidatePassiveModeConfig(config PassiveModeConfig) error {
	if err := validateStringSlice(config.AuthorizedContracts); err != nil {
		return err
	}
	if err := validateVerificationLevel(config.VerificationLevel); err != nil {
		return err
	}
	return nil
}
