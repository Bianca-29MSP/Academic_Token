package types

import (
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateCurriculumTree{}

// NewMsgCreateCurriculumTree creates a new message for creating a curriculum tree
func NewMsgCreateCurriculumTree(
	creator string,
	courseId string,
	version string,
	electiveMin string,
	totalWorkloadHours string,
) *MsgCreateCurriculumTree {
	// Convert string numeric values to uint64
	electiveMinUint, _ := strconv.ParseUint(electiveMin, 10, 64)
	totalWorkloadHoursUint, _ := strconv.ParseUint(totalWorkloadHours, 10, 64)

	return &MsgCreateCurriculumTree{
		Creator:            creator,
		CourseId:           courseId,
		Version:            version,
		ElectiveMin:        electiveMinUint,
		TotalWorkloadHours: totalWorkloadHoursUint,
	}
}

// ValidateBasic performs basic validation of the message
func (msg *MsgCreateCurriculumTree) ValidateBasic() error {
	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate required fields
	if msg.CourseId == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "course ID cannot be empty")
	}

	if msg.Version == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "version cannot be empty")
	}

	// Validate numeric fields
	if msg.TotalWorkloadHours == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "total workload hours must be greater than 0")
	}

	return nil
}
