package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// Parameter store keys
var (
	KeyContractAddress = []byte("ContractAddress")
	KeyContractVersion = []byte("ContractVersion")
)

// NewParams creates a new Params instance
func NewParams(contractAddress string, contractVersion string) Params {
	return Params{
		ContractAddress: contractAddress,
		ContractVersion: contractVersion,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams("", "1.0.0")
}

// ParamSetPairs implements the ParamSet interface
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyContractAddress, &p.ContractAddress, validateContractAddress),
		paramtypes.NewParamSetPair(KeyContractVersion, &p.ContractVersion, validateContractVersion),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateContractAddress(p.ContractAddress); err != nil {
		return err
	}
	if err := validateContractVersion(p.ContractVersion); err != nil {
		return err
	}
	return nil
}

// Validation functions
func validateContractAddress(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	// Allow empty string for initial setup
	return nil
}

func validateContractVersion(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == "" {
		return fmt.Errorf("contract version cannot be empty")
	}
	return nil
}
