package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Message type constants
const (
	TypeMsgUpdateSubjectContent = "update_subject_content"
)

var _ sdk.Msg = &MsgUpdateSubjectContent{}

func NewMsgUpdateSubjectContent(
	creator string,
	subjectId string,
	objectives []string,
	topicUnits []string,
	methodologies []string,
	evaluationMethods []string,
	bibliographyBasic []string,
	bibliographyComplementary []string,
	keywords []string,
) *MsgUpdateSubjectContent {
	return &MsgUpdateSubjectContent{
		Creator:                   creator,
		SubjectId:                 subjectId,
		Objectives:                objectives,
		TopicUnits:                topicUnits,
		Methodologies:             methodologies,
		EvaluationMethods:         evaluationMethods,
		BibliographyBasic:         bibliographyBasic,
		BibliographyComplementary: bibliographyComplementary,
		Keywords:                  keywords,
	}
}

func (msg *MsgUpdateSubjectContent) Route() string {
	return ModuleName
}

func (msg *MsgUpdateSubjectContent) Type() string {
	return TypeMsgUpdateSubjectContent
}

func (msg *MsgUpdateSubjectContent) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateSubjectContent) GetSignBytes() []byte {
	bz := sdk.MustSortJSON([]byte(msg.String()))
	return bz
}

func (msg *MsgUpdateSubjectContent) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}

	if msg.SubjectId == "" {
		return fmt.Errorf("subject_id cannot be empty")
	}

	return nil
}
