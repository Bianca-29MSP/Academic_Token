package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateCourse{}

func NewMsgUpdateCourse(
	creator string,
	index string,
	name string,
	description string,
	totalCredits string,
) *MsgUpdateCourse {
	return &MsgUpdateCourse{
		Creator:      creator,
		Index:        index,
		Name:         name,
		Description:  description,
		TotalCredits: totalCredits,
	}
}

func (msg *MsgUpdateCourse) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Index == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index cannot be empty")
	}

	if msg.Name == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "name cannot be empty")
	}

	if msg.TotalCredits == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "total credits cannot be empty")
	}

	return nil
}
