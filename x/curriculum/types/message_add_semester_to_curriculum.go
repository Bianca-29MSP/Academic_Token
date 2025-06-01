package types

import (
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgAddSemesterToCurriculum{}

// NewMsgAddSemesterToCurriculum creates a new message for adding a semester to a curriculum
func NewMsgAddSemesterToCurriculum(creator string, curriculumIndex string, semesterNumber string) *MsgAddSemesterToCurriculum {
	// Convert semester number from string to uint64
	semesterNum, _ := strconv.ParseUint(semesterNumber, 10, 64)

	return &MsgAddSemesterToCurriculum{
		Creator:         creator,
		CurriculumIndex: curriculumIndex,
		SemesterNumber:  semesterNum,
	}
}

// ValidateBasic performs basic validation of the message
func (msg *MsgAddSemesterToCurriculum) ValidateBasic() error {
	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate curriculum index
	if msg.CurriculumIndex == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "curriculum index cannot be empty")
	}

	// Validate semester number
	if msg.SemesterNumber == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "semester number must be greater than 0")
	}

	return nil
}
