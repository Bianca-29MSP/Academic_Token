package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCompleteSubject{}

const TypeMsgCompleteSubject = "complete_subject"

func NewMsgCompleteSubject(
	creator string,
	studentId string,
	subjectId string,
	grade uint32,
	completionDate string,
	semester string,
	professorSignature string,
) *MsgCompleteSubject {
	return &MsgCompleteSubject{
		Creator:            creator,
		StudentId:          studentId,
		SubjectId:          subjectId,
		Grade:              grade,
		CompletionDate:     completionDate,
		Semester:           semester,
		ProfessorSignature: professorSignature,
	}
}

func (msg *MsgCompleteSubject) Type() string {
	return TypeMsgCompleteSubject
}

func (msg *MsgCompleteSubject) Route() string {
	return RouterKey
}

func (msg *MsgCompleteSubject) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCompleteSubject) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCompleteSubject) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.StudentId == "" {
		return errorsmod.Wrap(ErrInvalidBasicValidation, "student ID cannot be empty")
	}

	if msg.SubjectId == "" {
		return errorsmod.Wrap(ErrInvalidBasicValidation, "subject ID cannot be empty")
	}

	if msg.Grade > 10000 { // Assuming grade is stored as integer (e.g., 85.5 -> 8550)
		return errorsmod.Wrap(ErrInvalidBasicValidation, "grade cannot exceed 100.00")
	}

	if msg.CompletionDate == "" {
		return errorsmod.Wrap(ErrInvalidBasicValidation, "completion date cannot be empty")
	}

	if msg.Semester == "" {
		return errorsmod.Wrap(ErrInvalidBasicValidation, "semester cannot be empty")
	}

	return nil
}
