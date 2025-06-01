package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateSubjectRecommendation{}

func NewMsgCreateSubjectRecommendation(creator string, student string, recommendationSemester string, recommendationMetadata string) *MsgCreateSubjectRecommendation {
	return &MsgCreateSubjectRecommendation{
		Creator:                creator,
		Student:                student,
		RecommendationSemester: recommendationSemester,
		RecommendationMetadata: recommendationMetadata,
	}
}

func (msg *MsgCreateSubjectRecommendation) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Student == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "student cannot be empty")
	}

	if msg.RecommendationSemester == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "recommendation semester cannot be empty")
	}

	return nil
}
