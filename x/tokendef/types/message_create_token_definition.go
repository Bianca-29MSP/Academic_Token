package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateTokenDefinition{}

func NewMsgCreateTokenDefinition(
	creator string,
	subjectId string,
	tokenName string,
	tokenSymbol string,
	description string,
	tokenType string,
	isTransferable bool,
	isBurnable bool,
	maxSupply uint64,
	imageUri string,
) *MsgCreateTokenDefinition {
	return &MsgCreateTokenDefinition{
		Creator:        creator,
		SubjectId:      subjectId,
		TokenName:      tokenName,
		TokenSymbol:    tokenSymbol,
		Description:    description,
		TokenType:      tokenType,
		IsTransferable: isTransferable,
		IsBurnable:     isBurnable,
		MaxSupply:      maxSupply,
		ImageUri:       imageUri,
		Attributes:     []*TokenAttributeInput{}, // Initialize empty slice
	}
}

func (msg *MsgCreateTokenDefinition) ValidateBasic() error {
	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate required fields
	if msg.SubjectId == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "subject_id cannot be empty")
	}

	if msg.TokenName == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "token_name cannot be empty")
	}

	if msg.TokenSymbol == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "token_symbol cannot be empty")
	}

	// Validate max supply
	if msg.MaxSupply == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "max_supply must be greater than 0")
	}

	// Validate attributes
	for i, attr := range msg.Attributes {
		if attr.TraitType == "" {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "attribute[%d] trait_type cannot be empty", i)
		}
		if attr.DisplayType == "" {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "attribute[%d] display_type cannot be empty", i)
		}
	}

	return nil
}
