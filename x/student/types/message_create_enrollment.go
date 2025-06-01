package types

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateEnrollment{}

func NewMsgCreateEnrollment(creator string, student string, institution string, courseId string) *MsgCreateEnrollment {
	return &MsgCreateEnrollment{
		Creator:     creator,
		Student:     student,
		Institution: institution,
		CourseId:    courseId,
	}
}

func (msg *MsgCreateEnrollment) ValidateBasic() error {
	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate student ID
	if strings.TrimSpace(msg.Student) == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "student ID cannot be empty")
	}

	// Validate institution ID
	if strings.TrimSpace(msg.Institution) == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "institution ID cannot be empty")
	}

	// Validate course ID
	if strings.TrimSpace(msg.CourseId) == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "course ID cannot be empty")
	}

	return nil
}

// GetSigners returns the expected signers for the message
func (msg *MsgCreateEnrollment) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}
