package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateDegreeContract = "update_degree_contract"

var _ sdk.Msg = &MsgUpdateDegreeContract{}

func NewMsgUpdateDegreeContract(
	authority string,
	newContractAddress string,
	contractVersion string,
	migrationReason string,
) *MsgUpdateDegreeContract {
	return &MsgUpdateDegreeContract{
		Authority:          authority,
		NewContractAddress: newContractAddress,
		ContractVersion:    contractVersion,
		MigrationReason:    migrationReason,
	}
}

func (msg *MsgUpdateDegreeContract) Route() string {
	return RouterKey
}

func (msg *MsgUpdateDegreeContract) Type() string {
	return TypeMsgUpdateDegreeContract
}

func (msg *MsgUpdateDegreeContract) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgUpdateDegreeContract) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateDegreeContract) ValidateBasic() error {
	// Validate authority address
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	// Validate new contract address
	if msg.NewContractAddress == "" {
		return errorsmod.Wrap(ErrInvalidContractAddress, "new contract address cannot be empty")
	}

	// Validate contract version
	if msg.ContractVersion == "" {
		return errorsmod.Wrap(ErrInvalidDegreeRequest, "contract version cannot be empty")
	}

	// Migration reason is optional, so no validation needed

	return nil
}
