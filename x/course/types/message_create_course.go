package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateCourse{}

func NewMsgCreateCourse(
	creator string,
	institution string,
	name string,
	code string,
	description string,
	totalCredits string,
	degreeLevel string,
) *MsgCreateCourse {
	return &MsgCreateCourse{
		Creator:      creator,
		Institution:  institution,
		Name:         name,
		Code:         code,
		Description:  description,
		TotalCredits: totalCredits,
		DegreeLevel:  degreeLevel,
	}
}

func (msg *MsgCreateCourse) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Institution == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "institution cannot be empty")
	}

	if msg.Name == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "name cannot be empty")
	}

	if msg.Code == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "code cannot be empty")
	}

	if msg.TotalCredits == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "total credits cannot be empty")
	}

	if msg.DegreeLevel == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "degree level cannot be empty")
	}

	return nil
}
