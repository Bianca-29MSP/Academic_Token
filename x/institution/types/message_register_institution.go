package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgRegisterInstitution{}

// NewMsgRegisterInstitution creates a new message for registering an institution
func NewMsgRegisterInstitution(creator string, name string, address string) *MsgRegisterInstitution {
	return &MsgRegisterInstitution{
		Creator: creator,
		Name:    name,
		Address: address,
	}
}

// ValidateBasic performs basic validation of the message
func (msg *MsgRegisterInstitution) ValidateBasic() error {
	// Validate the creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate institution name
	if msg.Name == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "institution name cannot be empty")
	}

	// Validate institution address
	if msg.Address == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "institution address cannot be empty")
	}

	return nil
}
