package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateStudyPlan{}

func NewMsgCreateStudyPlan(creator string, student string, completionTarget string, additionalNotes string) *MsgCreateStudyPlan {
	return &MsgCreateStudyPlan{
		Creator:          creator,
		Student:          student,
		CompletionTarget: completionTarget,
		AdditionalNotes:  additionalNotes,
	}
}

func (msg *MsgCreateStudyPlan) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Student == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "student cannot be empty")
	}

	return nil
}
