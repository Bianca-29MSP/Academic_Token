package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgValidateDegreeRequirements = "validate_degree_requirements"

var _ sdk.Msg = &MsgValidateDegreeRequirements{}

func NewMsgValidateDegreeRequirements(
	creator string,
	degreeRequestId string,
	contractAddress string,
	validationParameters string,
) *MsgValidateDegreeRequirements {
	return &MsgValidateDegreeRequirements{
		Creator:              creator,
		DegreeRequestId:      degreeRequestId,
		ContractAddress:      contractAddress,
		ValidationParameters: validationParameters,
	}
}

func (msg *MsgValidateDegreeRequirements) Route() string {
	return RouterKey
}

func (msg *MsgValidateDegreeRequirements) Type() string {
	return TypeMsgValidateDegreeRequirements
}

func (msg *MsgValidateDegreeRequirements) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgValidateDegreeRequirements) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgValidateDegreeRequirements) ValidateBasic() error {
	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate degree request ID
	if msg.DegreeRequestId == "" {
		return errorsmod.Wrap(ErrDegreeRequestNotFound, "degree request ID cannot be empty")
	}

	// Validate contract address
	if msg.ContractAddress == "" {
		return errorsmod.Wrap(ErrInvalidContractAddress, "contract address cannot be empty")
	}

	return nil
}
