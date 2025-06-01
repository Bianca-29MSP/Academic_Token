package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateInstitution{}

// NewMsgUpdateInstitution creates a new message for updating an institution
func NewMsgUpdateInstitution(creator string, index string, name string, address string, isAuthorized string) *MsgUpdateInstitution {
	return &MsgUpdateInstitution{
		Creator:      creator,
		Index:        index,
		Name:         name,
		Address:      address,
		IsAuthorized: isAuthorized,
	}
}

// ValidateBasic performs basic validation of the message
func (msg *MsgUpdateInstitution) ValidateBasic() error {
	// Validate the creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate institution index
	if msg.Index == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "institution index cannot be empty")
	}

	// Validate institution name
	if msg.Name == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "institution name cannot be empty")
	}

	// Validate institution address
	if msg.Address == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "institution address cannot be empty")
	}

	// Validate authorization status
	if msg.IsAuthorized != "true" && msg.IsAuthorized != "false" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "isAuthorized must be either 'true' or 'false'")
	}

	return nil
}
