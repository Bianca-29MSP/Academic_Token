package types

import "fmt"

var (
	ErrInvalidGroupType = fmt.Errorf("invalid group type")
)

// NewPrerequisiteGroup creates a new instance of PrerequisiteGroup
func NewPrerequisiteGroup(
	id string,
	subjectId string,
	groupType string,
	minimumCredits uint64,
	minimumCompletedSubjects uint64,
	subjectIds []string,
) PrerequisiteGroup {
	return PrerequisiteGroup{
		Id:                       id,
		SubjectId:                subjectId,
		GroupType:                groupType,
		MinimumCredits:           minimumCredits,
		MinimumCompletedSubjects: minimumCompletedSubjects,
		SubjectIds:               subjectIds,
	}
}

// ValidateBasic performs a basic validation of PrerequisiteGroup fields
func (pg PrerequisiteGroup) ValidateBasic() error {
	if pg.Id == "" {
		return fmt.Errorf("prerequisite group id cannot be empty")
	}

	if pg.SubjectId == "" {
		return fmt.Errorf("subject id cannot be empty")
	}

	if pg.GroupType == "" {
		return ErrInvalidGroupType
	}

	if pg.GroupType != "ALL" && pg.GroupType != "ANY" && pg.GroupType != "CREDITS" {
		return ErrInvalidGroupType
	}

	// Add more validations as needed
	return nil
}
