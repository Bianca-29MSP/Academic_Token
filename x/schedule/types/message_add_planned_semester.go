package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgAddPlannedSemester{}

func NewMsgAddPlannedSemester(creator string, studyPlanId string, semesterCode string, plannedSubjects []string) *MsgAddPlannedSemester {
	return &MsgAddPlannedSemester{
		Creator:         creator,
		StudyPlanId:     studyPlanId,
		SemesterCode:    semesterCode,
		PlannedSubjects: plannedSubjects,
	}
}

func (msg *MsgAddPlannedSemester) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.StudyPlanId == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "study plan ID cannot be empty")
	}

	if msg.SemesterCode == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "semester code cannot be empty")
	}

	if len(msg.PlannedSubjects) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "planned subjects cannot be empty")
	}

	return nil
}
