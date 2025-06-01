package types

// ============================================================================
// CORE TYPES (Non-contract related)
// ============================================================================

// Note: AcademicProgress struct is defined in academic_progress.pb.go (generated from protobuf)

// SubjectProgress represents progress for a specific subject
type SubjectProgress struct {
	SubjectId        string `json:"subject_id"`
	Status           string `json:"status"` // "available", "in_progress", "completed", "locked"
	Credits          uint64 `json:"credits"`
	PrerequisitesMet bool   `json:"prerequisites_met"`
	Grade            uint64 `json:"grade,omitempty"`
	CompletionDate   string `json:"completion_date,omitempty"`
	NftTokenId       string `json:"nft_token_id,omitempty"`
}

// StudentStats represents statistical data about a student's academic progress
type StudentStats struct {
	TotalSubjects        uint64  `json:"total_subjects"`
	CompletedSubjects    uint64  `json:"completed_subjects"`
	InProgressSubjects   uint64  `json:"in_progress_subjects"`
	AvailableSubjects    uint64  `json:"available_subjects"`
	CompletionPercentage float64 `json:"completion_percentage"`
	CurrentGPA           float64 `json:"current_gpa"`
	TotalCredits         uint64  `json:"total_credits"`
	RemainingCredits     uint64  `json:"remaining_credits"`
}

// GraduationRequirement represents a requirement for graduation
type GraduationRequirement struct {
	Id            string   `json:"id"`
	Type          string   `json:"type"` // "credits", "subjects", "gpa", "special"
	Description   string   `json:"description"`
	RequiredValue uint64   `json:"required_value"`
	CurrentValue  uint64   `json:"current_value"`
	IsSatisfied   bool     `json:"is_satisfied"`
	SubjectIds    []string `json:"subject_ids,omitempty"`
	MinimumGrade  uint64   `json:"minimum_grade,omitempty"`
}

// TransferCredits represents credits transferred from another institution
type TransferCredits struct {
	Id                string `json:"id"`
	StudentId         string `json:"student_id"`
	SourceInstitution string `json:"source_institution"`
	SourceSubjectId   string `json:"source_subject_id"`
	SourceSubjectName string `json:"source_subject_name"`
	TargetSubjectId   string `json:"target_subject_id"`
	Credits           uint64 `json:"credits"`
	Grade             uint64 `json:"grade"`
	EquivalenceId     string `json:"equivalence_id"`
	ApprovalDate      string `json:"approval_date"`
	ApprovedBy        string `json:"approved_by"`
	Status            string `json:"status"` // "pending", "approved", "rejected"
}

// StudentWarning represents academic warnings or alerts for a student
type StudentWarning struct {
	Id          string `json:"id"`
	StudentId   string `json:"student_id"`
	Type        string `json:"type"`     // "low_gpa", "missing_prerequisites", "overdue", "graduation_risk"
	Severity    string `json:"severity"` // "low", "medium", "high", "critical"
	Message     string `json:"message"`
	CreatedDate string `json:"created_date"`
	IsRead      bool   `json:"is_read"`
	IsResolved  bool   `json:"is_resolved"`
}

// ContractAnalysisResult represents the result from equivalence contract analysis
type ContractAnalysisResult struct {
	EquivalenceId      string `json:"equivalence_id"`
	SourceSubjectId    string `json:"source_subject_id"`
	TargetSubjectId    string `json:"target_subject_id"`
	EquivalencePercent string `json:"equivalence_percent"`
	AnalysisMetadata   string `json:"analysis_metadata"`
	ContractAddress    string `json:"contract_address"`
	ContractVersion    string `json:"contract_version"`
	AnalysisHash       string `json:"analysis_hash"`
	ProcessingTime     string `json:"processing_time"`
	Success            bool   `json:"success"`
	ErrorMessage       string `json:"error_message"`
}
