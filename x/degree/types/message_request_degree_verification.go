package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRequestDegree = "request_degree"

var _ sdk.Msg = &MsgRequestDegree{}

func NewMsgRequestDegree(
	creator string,
	studentId string,
	institutionId string,
	curriculumId string,
	expectedGraduationDate string,
) *MsgRequestDegree {
	return &MsgRequestDegree{
		Creator:                creator,
		StudentId:              studentId,
		InstitutionId:          institutionId,
		CurriculumId:           curriculumId,
		ExpectedGraduationDate: expectedGraduationDate,
	}
}

func (msg *MsgRequestDegree) Route() string {
	return RouterKey
}

func (msg *MsgRequestDegree) Type() string {
	return TypeMsgRequestDegree
}

func (msg *MsgRequestDegree) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRequestDegree) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRequestDegree) ValidateBasic() error {
	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate student ID
	if msg.StudentId == "" {
		return errorsmod.Wrap(ErrInvalidDegreeRequest, "student ID cannot be empty")
	}

	// Validate institution ID
	if msg.InstitutionId == "" {
		return errorsmod.Wrap(ErrInvalidDegreeRequest, "institution ID cannot be empty")
	}

	// Validate curriculum ID
	if msg.CurriculumId == "" {
		return errorsmod.Wrap(ErrInvalidDegreeRequest, "curriculum ID cannot be empty")
	}

	// Validate expected graduation date
	if msg.ExpectedGraduationDate == "" {
		return errorsmod.Wrap(ErrInvalidDegreeRequest, "expected graduation date cannot be empty")
	}

	return nil
}
