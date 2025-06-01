package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgVerifyTokenInstance{}

func NewMsgVerifyTokenInstance(creator string, tokenInstanceId string) *MsgVerifyTokenInstance {
	return &MsgVerifyTokenInstance{
		Creator:         creator,
		TokenInstanceId: tokenInstanceId,
	}
}

func (msg *MsgVerifyTokenInstance) ValidateBasic() error {
	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate required fields
	if msg.TokenInstanceId == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "tokenInstanceId cannot be empty")
	}

	return nil
}
