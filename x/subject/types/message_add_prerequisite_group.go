package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Message type constants
const (
	TypeMsgAddPrerequisiteGroup = "add_prerequisite_group"
)

var _ sdk.Msg = &MsgAddPrerequisiteGroup{}

// NewMsgAddPrerequisiteGroup creates a new message for adding a prerequisite group
func NewMsgAddPrerequisiteGroup(
	creator string,
	subjectId string,
	groupType string,
	minimumCredits uint64,
	minimumCompletedSubjects uint64,
	subjectIds []string,
) *MsgAddPrerequisiteGroup {
	return &MsgAddPrerequisiteGroup{
		Creator:                  creator,
		SubjectId:                subjectId,
		GroupType:                groupType,
		MinimumCredits:           minimumCredits,
		MinimumCompletedSubjects: minimumCompletedSubjects,
		SubjectIds:               subjectIds,
	}
}

func (msg *MsgAddPrerequisiteGroup) Route() string {
	return ModuleName
}

func (msg *MsgAddPrerequisiteGroup) Type() string {
	return TypeMsgAddPrerequisiteGroup
}

func (msg *MsgAddPrerequisiteGroup) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddPrerequisiteGroup) GetSignBytes() []byte {
	bz := sdk.MustSortJSON([]byte(msg.String()))
	return bz
}

func (msg *MsgAddPrerequisiteGroup) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}

	if msg.SubjectId == "" {
		return fmt.Errorf("subject_id cannot be empty")
	}

	if msg.GroupType == "" {
		return fmt.Errorf("group_type cannot be empty")
	}

	return nil
}
