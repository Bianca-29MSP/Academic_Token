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

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		"http://localhost:5001", // default IPFS gateway
		true,                    // IPFS enabled by default
		"",                      // empty admin
	)
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
