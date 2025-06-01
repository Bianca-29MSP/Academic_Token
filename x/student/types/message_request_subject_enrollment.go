package types

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgRequestSubjectEnrollment{}

func NewMsgRequestSubjectEnrollment(creator string, student string, subjectId string) *MsgRequestSubjectEnrollment {
	return &MsgRequestSubjectEnrollment{
		Creator:   creator,
		Student:   student,
		SubjectId: subjectId,
	}
}

func (msg *MsgRequestSubjectEnrollment) ValidateBasic() error {
	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate student ID
	if strings.TrimSpace(msg.Student) == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "student ID cannot be empty")
	}

	// Validate subject ID
	if strings.TrimSpace(msg.SubjectId) == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "subject ID cannot be empty")
	}

	return nil
}

// GetSigners returns the expected signers for the message
func (msg *MsgRequestSubjectEnrollment) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}
