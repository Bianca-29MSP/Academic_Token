package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateStudyPlanStatus{}

func NewMsgUpdateStudyPlanStatus(creator string, studyPlanId string, status string) *MsgUpdateStudyPlanStatus {
	return &MsgUpdateStudyPlanStatus{
		Creator:     creator,
		StudyPlanId: studyPlanId,
		Status:      status, // Use Status instead of NewStatus
	}
}

func (msg *MsgUpdateStudyPlanStatus) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.StudyPlanId == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "study plan ID cannot be empty")
	}

	if msg.Status == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "status cannot be empty")
	}

	return nil
}
