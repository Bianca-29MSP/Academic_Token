package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgIssueDegree = "issue_degree"

var _ sdk.Msg = &MsgIssueDegree{}

func NewMsgIssueDegree(
	creator string,
	degreeRequestId string,
	finalGpa string,
	totalCredits uint64,
	signatures []string,
	additionalNotes string,
) *MsgIssueDegree {
	return &MsgIssueDegree{
		Creator:         creator,
		DegreeRequestId: degreeRequestId,
		FinalGpa:        finalGpa,
		TotalCredits:    totalCredits,
		Signatures:      signatures,
		AdditionalNotes: additionalNotes,
	}
}

func (msg *MsgIssueDegree) Route() string {
	return RouterKey
}

func (msg *MsgIssueDegree) Type() string {
	return TypeMsgIssueDegree
}

func (msg *MsgIssueDegree) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgIssueDegree) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgIssueDegree) ValidateBasic() error {
	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate degree request ID
	if msg.DegreeRequestId == "" {
		return errorsmod.Wrap(ErrDegreeRequestNotFound, "degree request ID cannot be empty")
	}

	// Validate final GPA
	if msg.FinalGpa == "" {
		return errorsmod.Wrap(ErrInvalidGrade, "final GPA cannot be empty")
	}

	// Validate total credits
	if msg.TotalCredits == 0 {
		return errorsmod.Wrap(ErrInsufficientCredits, "total credits must be greater than 0")
	}

	// Validate signatures
	if len(msg.Signatures) == 0 {
		return errorsmod.Wrap(ErrInvalidDegreeRequest, "at least one signature is required")
	}

	for i, signature := range msg.Signatures {
		if signature == "" {
			return errorsmod.Wrapf(ErrInvalidDegreeRequest, "signature at index %d cannot be empty", i)
		}
	}

	return nil
}
