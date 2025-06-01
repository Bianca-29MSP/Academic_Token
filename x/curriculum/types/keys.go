package types

const (
	// ModuleName defines the module name
	ModuleName = "curriculum"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_curriculum"
)

var (
	ParamsKey = []byte("p_curriculum")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// CurriculumTreeKey returns the store key to retrieve a CurriculumTree from the index fields
func CurriculumTreeKey(index string) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}

// Event types for curriculum module
const (
	EventTypeCurriculumTreeCreated     = "curriculum_tree_created"
	EventTypeSemesterAdded             = "semester_added"
	EventTypeElectiveGroupAdded        = "elective_group_added"
	EventTypeGraduationRequirementsSet = "graduation_requirements_set"

	AttributeKeyIndex           = "index"
	AttributeKeyCourseId        = "course_id"
	AttributeKeyCreator         = "creator"
	AttributeKeyCurriculumIndex = "curriculum_index"
	AttributeKeySemesterNumber  = "semester_number"
	AttributeKeyGroupName       = "group_name"
)

// Event definitions
type EventCurriculumTreeCreated struct {
	Index    string `json:"index"`
	CourseId string `json:"course_id"`
	Creator  string `json:"creator"`
}

type EventSemesterAdded struct {
	CurriculumIndex string `json:"curriculum_index"`
	SemesterNumber  uint64 `json:"semester_number"`
	Creator         string `json:"creator"`
}

type EventElectiveGroupAdded struct {
	CurriculumIndex string `json:"curriculum_index"`
	GroupName       string `json:"group_name"`
	Creator         string `json:"creator"`
}

type EventGraduationRequirementsSet struct {
	CurriculumIndex string `json:"curriculum_index"`
	Creator         string `json:"creator"`
}
