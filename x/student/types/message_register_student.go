package types

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgRegisterStudent{}

func NewMsgRegisterStudent(creator string, name string, address string) *MsgRegisterStudent {
	return &MsgRegisterStudent{
		Creator: creator,
		Name:    name,
		Address: address,
	}
}

func (msg *MsgRegisterStudent) ValidateBasic() error {
	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate student name
	if strings.TrimSpace(msg.Name) == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "student name cannot be empty")
	}

	if len(msg.Name) < 2 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "student name must be at least 2 characters long")
	}

	if len(msg.Name) > 100 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "student name cannot exceed 100 characters")
	}

	// Validate student address
	_, err = sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid student address (%s)", err)
	}

	// Ensure creator and student address are the same (self-registration)
	if msg.Creator != msg.Address {
		return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only the student can register themselves")
	}

	return nil
}

// GetSigners returns the expected signers for the message
func (msg *MsgRegisterStudent) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}
