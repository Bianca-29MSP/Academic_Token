package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Message type constants
const (
	TypeMsgCreateSubject = "create_subject"
)

var _ sdk.Msg = &MsgCreateSubject{}

// NewMsgCreateSubject creates a new message for creating a subject
func NewMsgCreateSubject(
	creator string,
	institution string,
	courseId string,
	title string,
	code string,
	workloadHours uint64,
	credits uint64,
	description string,
	subjectType string,
	knowledgeArea string,
	objectives []string,
	topicUnits []string,
) *MsgCreateSubject {
	return &MsgCreateSubject{
		Creator:       creator,
		Institution:   institution,
		CourseId:      courseId,
		Title:         title,
		Code:          code,
		WorkloadHours: workloadHours,
		Credits:       credits,
		Description:   description,
		SubjectType:   subjectType,
		KnowledgeArea: knowledgeArea,
		Objectives:    objectives,
		TopicUnits:    topicUnits,
	}
}

// Route implements the sdk.Msg interface
func (msg *MsgCreateSubject) Route() string {
	return ModuleName
}

// Type implements the sdk.Msg interface
func (msg *MsgCreateSubject) Type() string {
	return TypeMsgCreateSubject
}

// GetSigners implements the sdk.Msg interface
func (msg *MsgCreateSubject) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes implements the sdk.Msg interface
func (msg *MsgCreateSubject) GetSignBytes() []byte {
	bz := sdk.MustSortJSON([]byte(msg.String()))
	return bz
}

// ValidateBasic implements the sdk.Msg interface
func (msg *MsgCreateSubject) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}

	if msg.Institution == "" {
		return fmt.Errorf("institution cannot be empty")
	}

	if msg.CourseId == "" {
		return fmt.Errorf("course_id cannot be empty")
	}

	if msg.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}

	if msg.Code == "" {
		return fmt.Errorf("code cannot be empty")
	}

	if msg.Credits == 0 {
		return fmt.Errorf("credits must be greater than 0")
	}

	return nil
}
