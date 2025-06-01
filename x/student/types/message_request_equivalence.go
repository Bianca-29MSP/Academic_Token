package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgRequestEquivalence{}

const TypeMsgRequestEquivalence = "request_equivalence"

func NewMsgRequestEquivalence(
	creator string,
	studentId string,
	sourceSubjectId string,
	targetSubjectId string,
	reason string,
) *MsgRequestEquivalence {
	return &MsgRequestEquivalence{
		Creator:         creator,
		StudentId:       studentId,
		SourceSubjectId: sourceSubjectId,
		TargetSubjectId: targetSubjectId,
		Reason:          reason,
	}
}

func (msg *MsgRequestEquivalence) Type() string {
	return TypeMsgRequestEquivalence
}

func (msg *MsgRequestEquivalence) Route() string {
	return RouterKey
}

func (msg *MsgRequestEquivalence) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRequestEquivalence) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRequestEquivalence) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.StudentId == "" {
		return errorsmod.Wrap(ErrInvalidBasicValidation, "student ID cannot be empty")
	}

	if msg.SourceSubjectId == "" {
		return errorsmod.Wrap(ErrInvalidBasicValidation, "source subject ID cannot be empty")
	}

	if msg.TargetSubjectId == "" {
		return errorsmod.Wrap(ErrInvalidBasicValidation, "target subject ID cannot be empty")
	}

	if msg.SourceSubjectId == msg.TargetSubjectId {
		return errorsmod.Wrap(ErrInvalidBasicValidation, "source and target subjects cannot be the same")
	}

	return nil
}
