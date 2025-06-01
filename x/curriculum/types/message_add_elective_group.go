package types

import (
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgAddElectiveGroup{}

// NewMsgAddElectiveGroup creates a new message for adding an elective group
func NewMsgAddElectiveGroup(
	creator string,
	curriculumIndex string,
	name string,
	description string,
	minSubjectsRequired string,
	creditsRequired string,
	knowledgeArea string,
) *MsgAddElectiveGroup {
	// Convert string numeric values to uint64
	minSubjects, _ := strconv.ParseUint(minSubjectsRequired, 10, 64)
	credits, _ := strconv.ParseUint(creditsRequired, 10, 64)

	return &MsgAddElectiveGroup{
		Creator:             creator,
		CurriculumIndex:     curriculumIndex,
		Name:                name,
		Description:         description,
		MinSubjectsRequired: minSubjects,
		CreditsRequired:     credits,
		KnowledgeArea:       knowledgeArea,
	}
}

// ValidateBasic performs basic validation of the message
func (msg *MsgAddElectiveGroup) ValidateBasic() error {
	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate required fields
	if msg.CurriculumIndex == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "curriculum index cannot be empty")
	}

	if msg.Name == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "name cannot be empty")
	}

	// Validate numeric fields
	if msg.MinSubjectsRequired == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "minimum subjects required must be greater than 0")
	}

	if msg.CreditsRequired == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "credits required must be greater than 0")
	}

	return nil
}
