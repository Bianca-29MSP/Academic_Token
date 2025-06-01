// src/msg.rs
use cosmwasm_schema::{cw_serde, QueryResponses};
use crate::state::{
    ScheduleRecommendation, AcademicPath, StudentProgress, SubjectScheduleInfo,
    AvailableSubjects, RecommendationAlgorithm, OptimizationCriteria, SchedulePreferences,
    State, ScheduleConfig,
};
use crate::ipfs::SubjectContent;

/// Instantiate message
#[cw_serde]
pub struct InstantiateMsg {
    pub owner: Option<String>,
    pub ipfs_gateway: String,
    pub max_subjects_per_semester: u32,
    pub recommendation_algorithm: RecommendationAlgorithm,
}

/// Execute messages
#[cw_serde]
pub enum ExecuteMsg {
    /// Update contract configuration
    UpdateConfig {
        ipfs_gateway: Option<String>,
        max_subjects_per_semester: Option<u32>,
        recommendation_algorithm: Option<RecommendationAlgorithm>,
        new_owner: Option<String>,
    },
    
    /// Register or update student progress
    RegisterStudentProgress {
        student_progress: StudentProgress,
    },
    
    /// Update student's schedule preferences
    UpdateStudentPreferences {
        student_id: String,
        preferences: SchedulePreferences,
    },
    
    /// Register subject schedule information
    RegisterSubjectScheduleInfo {
        subject_info: SubjectScheduleInfo,
    },
    
    /// Batch register multiple subjects
    BatchRegisterSubjects {
        subjects: Vec<SubjectScheduleInfo>,
    },
    
    /// Generate schedule recommendation for next semester
    GenerateScheduleRecommendation {
        student_id: String,
        target_semester: u32,
        force_refresh: Option<bool>,
        custom_preferences: Option<SchedulePreferences>,
    },
    
    /// Create complete academic path to graduation
    CreateAcademicPath {
        student_id: String,
        path_name: String,
        optimization_criteria: OptimizationCriteria,
        target_graduation_semester: Option<u32>,
    },
    
    /// Optimize existing academic path
    OptimizeAcademicPath {
        path_id: String,
        optimization_criteria: OptimizationCriteria,
        preserve_current_semester: Option<bool>,
    },
    
    /// Update academic path (manual modifications)
    UpdateAcademicPath {
        path_id: String,
        semester_number: u32,
        new_subjects: Vec<String>,
        notes: Option<String>,
    },
    
    /// Activate academic path (set as current)
    ActivateAcademicPath {
        path_id: String,
    },
    
    /// Mark subject as completed in student progress
    CompleteSubject {
        student_id: String,
        subject_id: String,
        grade: u32,
        completion_date: String,
        difficulty_rating: Option<u32>,
        workload_rating: Option<u32>,
        nft_token_id: String,
    },
    
    /// Enroll student in subject for current semester
    EnrollInSubject {
        student_id: String,
        subject_id: String,
        enrollment_date: String,
        expected_completion: String,
    },
    
    /// Cache IPFS content for subject analysis
    CacheIpfsContent {
        ipfs_link: String,
        content: SubjectContent,
    },
    
    /// Generate alternative recommendations
    GenerateAlternativeRecommendations {
        student_id: String,
        target_semester: u32,
        excluded_subjects: Vec<String>,
    },
    
    /// Simulate schedule with what-if scenarios
    SimulateSchedule {
        student_id: String,
        hypothetical_completions: Vec<String>,
        target_semester: u32,
    },
}

/// Query messages
#[cw_serde]
#[derive(QueryResponses)]
pub enum QueryMsg {
    /// Get contract state
    #[returns(StateResponse)]
    GetState {},
    
    /// Get schedule configuration
    #[returns(ScheduleConfigResponse)]
    GetConfig {},
    
    /// Get student progress
    #[returns(StudentProgressResponse)]
    GetStudentProgress { student_id: String },
    
    /// Get subject schedule information
    #[returns(SubjectScheduleInfoResponse)]
    GetSubjectScheduleInfo { subject_id: String },
    
    /// Get schedule recommendation for semester
    #[returns(ScheduleRecommendationResponse)]
    GetScheduleRecommendation {
        student_id: String,
        semester: u32,
    },
    
    /// Get academic path
    #[returns(AcademicPathResponse)]
    GetAcademicPath { path_id: String },
    
    /// List student's academic paths
    #[returns(StudentPathsResponse)]
    GetStudentPaths {
        student_id: String,
        include_inactive: Option<bool>,
    },
    
    /// Get available subjects for enrollment
    #[returns(AvailableSubjectsResponse)]
    GetAvailableSubjects {
        student_id: String,
        semester: u32,
        include_electives: Option<bool>,
    },
    
    /// Get optimal academic path recommendations
    #[returns(OptimalPathResponse)]
    GetOptimalPath {
        student_id: String,
        criteria: OptimizationCriteria,
        max_paths: Option<u32>,
    },
    
    /// Get graduation timeline analysis
    #[returns(GraduationTimelineResponse)]
    GetGraduationTimeline {
        student_id: String,
        path_id: Option<String>,
    },
    
    /// Get subject sequence recommendations
    #[returns(SubjectSequenceResponse)]
    GetSubjectSequence {
        student_id: String,
        target_subjects: Vec<String>,
    },
    
    /// Get workload analysis for semester
    #[returns(WorkloadAnalysisResponse)]
    GetWorkloadAnalysis {
        student_id: String,
        semester: u32,
        proposed_subjects: Vec<String>,
    },
    
    /// Check IPFS cache status
    #[returns(IpfsCacheStatusResponse)]
    GetIpfsCacheStatus { ipfs_link: String },
    
    /// Get cached IPFS content
    #[returns(SubjectContent)]
    GetCachedContent { ipfs_link: String },
    
    /// Get schedule statistics
    #[returns(ScheduleStatisticsResponse)]
    GetScheduleStatistics {
        student_id: Option<String>,
        institution_id: Option<String>,
    },
}

// Response structs
#[cw_serde]
pub struct StateResponse {
    pub state: State,
}

#[cw_serde]
pub struct ScheduleConfigResponse {
    pub config: ScheduleConfig,
}

#[cw_serde]
pub struct StudentProgressResponse {
    pub progress: StudentProgress,
}

#[cw_serde]
pub struct SubjectScheduleInfoResponse {
    pub subject_info: SubjectScheduleInfo,
}

#[cw_serde]
pub struct ScheduleRecommendationResponse {
    pub recommendation: Option<ScheduleRecommendation>,
}

#[cw_serde]
pub struct AcademicPathResponse {
    pub path: Option<AcademicPath>,
}

#[cw_serde]
pub struct StudentPathsResponse {
    pub paths: Vec<AcademicPath>,
    pub active_path_id: Option<String>,
}

#[cw_serde]
pub struct AvailableSubjectsResponse {
    pub available_subjects: AvailableSubjects,
}

#[cw_serde]
pub struct OptimalPathResponse {
    pub paths: Vec<AcademicPath>,
    pub recommendation_notes: Vec<String>,
}

#[cw_serde]
pub struct GraduationTimelineResponse {
    pub current_progress_percentage: String,
    pub estimated_graduation_semester: u32,
    pub estimated_graduation_date: String,
    pub remaining_credits: u32,
    pub remaining_subjects: Vec<String>,
    pub critical_path_subjects: Vec<String>,
    pub risk_factors: Vec<String>,
}

#[cw_serde]
pub struct SubjectSequenceResponse {
    pub recommended_sequence: Vec<SemesterSubjects>,
    pub total_semesters: u32,
    pub notes: Vec<String>,
}

#[cw_serde]
pub struct SemesterSubjects {
    pub semester: u32,
    pub subjects: Vec<String>,
}

#[cw_serde]
pub struct WorkloadAnalysisResponse {
    pub total_credits: u32,
    pub estimated_hours_per_week: u32,
    pub difficulty_score: u32,
    pub workload_rating: WorkloadRating,
    pub recommendations: Vec<String>,
    pub subject_breakdown: Vec<SubjectWorkload>,
}

#[cw_serde]
pub enum WorkloadRating {
    Light,
    Moderate,
    Heavy,
    Excessive,
}

#[cw_serde]
pub struct SubjectWorkload {
    pub subject_id: String,
    pub credits: u32,
    pub estimated_hours: u32,
    pub difficulty_contribution: u32,
}

#[cw_serde]
pub struct IpfsCacheStatusResponse {
    pub is_cached: bool,
    pub ipfs_link: String,
    pub cache_timestamp: Option<u64>,
}

#[cw_serde]
pub struct ScheduleStatisticsResponse {
    pub total_students: u64,
    pub total_recommendations_generated: u64,
    pub total_paths_created: u64,
    pub average_time_to_graduation: Option<u32>,
    pub most_recommended_subjects: Vec<PopularSubject>,
    pub optimization_success_rate: u32,
}

#[cw_serde]
pub struct PopularSubject {
    pub subject_id: String,
    pub recommendation_count: u32,
    pub completion_rate: u32,
}