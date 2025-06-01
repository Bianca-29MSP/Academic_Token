package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgRequestEquivalence{}

func NewMsgRequestEquivalence(creator string, sourceSubjectId string, targetInstitution string, targetSubjectId string) *MsgRequestEquivalence {
	return &MsgRequestEquivalence{
		Creator:           creator,
		SourceSubjectId:   sourceSubjectId,
		TargetInstitution: targetInstitution,
		TargetSubjectId:   targetSubjectId,
	}
}

func (msg *MsgRequestEquivalence) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Additional validation
	if msg.SourceSubjectId == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "source_subject_id cannot be empty")
	}

	if msg.TargetInstitution == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "target_institution cannot be empty")
	}

	if msg.TargetSubjectId == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "target_subject_id cannot be empty")
	}

	return nil
}

// GetSigners implements the sdk.Msg interface
func (msg *MsgRequestEquivalence) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// Route implements the sdk.Msg interface
func (msg *MsgRequestEquivalence) Route() string {
	return ModuleName
}

// Type implements the sdk.Msg interface
func (msg *MsgRequestEquivalence) Type() string {
	return "request_equivalence"
}

// GetSignBytes implements the sdk.Msg interface
func (msg *MsgRequestEquivalence) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
