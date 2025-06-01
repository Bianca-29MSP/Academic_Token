package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCancelDegreeRequest = "cancel_degree_request"

var _ sdk.Msg = &MsgCancelDegreeRequest{}

func NewMsgCancelDegreeRequest(
	creator string,
	degreeRequestId string,
	cancellationReason string,
) *MsgCancelDegreeRequest {
	return &MsgCancelDegreeRequest{
		Creator:            creator,
		DegreeRequestId:    degreeRequestId,
		CancellationReason: cancellationReason,
	}
}

func (msg *MsgCancelDegreeRequest) Route() string {
	return RouterKey
}

func (msg *MsgCancelDegreeRequest) Type() string {
	return TypeMsgCancelDegreeRequest
}

func (msg *MsgCancelDegreeRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCancelDegreeRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCancelDegreeRequest) ValidateBasic() error {
	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate degree request ID
	if msg.DegreeRequestId == "" {
		return errorsmod.Wrap(ErrDegreeRequestNotFound, "degree request ID cannot be empty")
	}

	// Validate cancellation reason
	if msg.CancellationReason == "" {
		return errorsmod.Wrap(ErrInvalidDegreeRequest, "cancellation reason cannot be empty")
	}

	return nil
}
