package types

const (
	// ModuleName defines the module name
	ModuleName = "schedule"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_schedule"
)

var (
	ParamsKey = []byte("p_schedule")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// Storage prefixes
const (
	StudyPlanPrefix              = "StudyPlan/value/"
	StudyPlanCountPrefix         = "StudyPlan/count/"
	PlannedSemesterPrefix        = "PlannedSemester/value/"
	PlannedSemesterCountPrefix   = "PlannedSemester/count/"
	SubjectRecommendationPrefix  = "SubjectRecommendation/value/"
	SubjectRecommendationCountPrefix = "SubjectRecommendation/count/"
	RecommendedSubjectPrefix     = "RecommendedSubject/value/"
)

// Events
const (
	EventTypeStudyPlanCreated           = "study_plan_created"
	EventTypeStudyPlanUpdated           = "study_plan_updated"
	EventTypeStudyPlanStatusChanged     = "study_plan_status_changed"
	EventTypePlannedSemesterAdded       = "planned_semester_added"
	EventTypePlannedSemesterUpdated     = "planned_semester_updated"
	EventTypeSubjectRecommendationCreated = "subject_recommendation_created"
	EventTypeSubjectRecommendationUpdated = "subject_recommendation_updated"
	EventTypeRecommendationGenerated    = "recommendation_generated"
	EventTypeScheduleOptimized          = "schedule_optimized"
)

// Event attributes
const (
	AttributeKeyStudyPlanID        = "study_plan_id"
	AttributeKeyStudentID          = "student_id"
	AttributeKeyCourseID           = "course_id"
	AttributeKeyInstitutionID      = "institution_id"
	AttributeKeySemesterNumber     = "semester_number"
	AttributeKeySubjectID          = "subject_id"
	AttributeKeyRecommendationType = "recommendation_type"
	AttributeKeyPriority           = "priority"
	AttributeKeyCredits            = "credits"
	AttributeKeyWorkloadHours      = "workload_hours"
	AttributeKeyDifficultyLevel    = "difficulty_level"
	AttributeKeyRecommendationScore = "recommendation_score"
	AttributeKeyStatus             = "status"
	AttributeKeyReason             = "reason"
	AttributeValueCategory         = ModuleName
)

// Study Plan Status
const (
	StudyPlanStatusDraft     = "draft"
	StudyPlanStatusActive    = "active"
	StudyPlanStatusCompleted = "completed"
	StudyPlanStatusArchived  = "archived"
)

// Recommendation Types
const (
	RecommendationTypeNext       = "next_semester"
	RecommendationTypeElective   = "elective"
	RecommendationTypePrereq     = "prerequisite"
	RecommendationTypeOptimal    = "optimal_path"
	RecommendationTypeIntensive  = "intensive"
)

// Priority Levels
const (
	PriorityHigh   = "high"
	PriorityMedium = "medium"
	PriorityLow    = "low"
)

// Difficulty Levels
const (
	DifficultyEasy   = "easy"
	DifficultyMedium = "medium"
	DifficultyHard   = "hard"
)
