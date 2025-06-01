package types

import "fmt"

var (
	ErrInvalidIndex       = fmt.Errorf("invalid index")
	ErrInvalidSubjectId   = fmt.Errorf("invalid subject ID")
	ErrInvalidInstitution = fmt.Errorf("invalid institution")
	ErrInvalidCourseId    = fmt.Errorf("invalid course ID")
)

// NewSubjectContent creates a new instance of SubjectContent
func NewSubjectContent(
	index string,
	subjectId string,
	institution string,
	courseId string,
	title string,
	code string,
	workloadHours uint64,
	credits uint64,
	description string,
	contentHash string,
	subjectType string,
	knowledgeArea string,
	ipfsLink string,
	creator string,
) SubjectContent {
	return SubjectContent{
		Index:         index,
		SubjectId:     subjectId,
		Institution:   institution,
		CourseId:      courseId,
		Title:         title,
		Code:          code,
		WorkloadHours: workloadHours,
		Credits:       credits,
		Description:   description,
		ContentHash:   contentHash,
		SubjectType:   subjectType,
		KnowledgeArea: knowledgeArea,
		IpfsLink:      ipfsLink,
		Creator:       creator,
	}
}

// ValidateBasic performs a basic validation of SubjectContent fields
func (s SubjectContent) ValidateBasic() error {
	if s.Index == "" {
		return ErrInvalidIndex
	}
	if s.SubjectId == "" {
		return ErrInvalidSubjectId
	}
	if s.Institution == "" {
		return ErrInvalidInstitution
	}
	if s.CourseId == "" {
		return ErrInvalidCourseId
	}

	return nil
}

// IsComplete checks if all required fields are filled
func (s SubjectContent) IsComplete() bool {
	return s.Index != "" &&
		s.SubjectId != "" &&
		s.Institution != "" &&
		s.CourseId != "" &&
		s.Title != "" &&
		s.Code != "" &&
		s.WorkloadHours > 0 &&
		s.Credits > 0
}

// HasIPFSContent checks if the subject has valid IPFS references
func (s SubjectContent) HasIPFSContent() bool {
	return s.ContentHash != "" && s.IpfsLink != ""
}
