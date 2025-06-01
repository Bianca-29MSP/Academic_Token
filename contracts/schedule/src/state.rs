// src/state.rs
use cosmwasm_schema::cw_serde;
use cosmwasm_std::Addr;
use cw_storage_plus::{Item, Map};

/// Main contract state
#[cw_serde]
pub struct State {
    pub owner: Addr,
    pub ipfs_gateway: String,
    pub total_students: u64,
    pub total_recommendations: u64,
    pub total_paths: u64,
}

/// Configuration parameters for schedule generation
#[cw_serde]
pub struct ScheduleConfig {
    pub max_subjects_per_semester: u32,
    pub recommendation_algorithm: RecommendationAlgorithm,
    pub optimization_weights: OptimizationWeights,
    pub default_preferences: DefaultSchedulePreferences,
}

/// Available recommendation algorithms
#[cw_serde]
pub enum RecommendationAlgorithm {
    Basic,              // Simple prerequisite-based
    Balanced,           // Considers workload and difficulty
    OptimalPath,        // Shortest path to graduation
    MachineLearning,    // ML-based recommendations
}

/// Weights for optimization criteria
#[cw_serde]
pub struct OptimizationWeights {
    pub graduation_speed: u32,     // Weight for faster graduation
    pub workload_balance: u32,     // Weight for balanced semesters
    pub difficulty_distribution: u32, // Weight for difficulty spread
    pub subject_availability: u32, // Weight for subject scheduling
    pub student_preferences: u32,  // Weight for student preferences
}

/// Student's academic progress and history
#[cw_serde]
pub struct StudentProgress {
    pub student_id: String,
    pub institution_id: String,
    pub course_id: String,
    pub curriculum_id: String,
    pub current_semester: u32,
    pub completed_subjects: Vec<CompletedSubject>,
    pub current_subjects: Vec<EnrolledSubject>,
    pub total_credits: u32,
    pub total_credits_required: u32,
    pub gpa: String,
    pub preferences: SchedulePreferences,
    pub academic_standing: AcademicStanding,
    pub expected_graduation: Option<String>,
}

/// Completed subject with performance data
#[cw_serde]
pub struct CompletedSubject {
    pub subject_id: String,
    pub credits: u32,
    pub grade: u32,
    pub completion_date: String,
    pub semester_taken: u32,
    pub difficulty_rating: Option<u32>, // Student's subjective difficulty rating
    pub workload_rating: Option<u32>,   // Student's workload assessment
    pub nft_token_id: String,
}

/// Currently enrolled subject
#[cw_serde]
pub struct EnrolledSubject {
    pub subject_id: String,
    pub credits: u32,
    pub enrollment_date: String,
    pub expected_completion: String,
    pub current_grade: Option<u32>,
}

/// Student's schedule preferences
#[cw_serde]
pub struct SchedulePreferences {
    pub max_subjects_per_semester: u32,
    pub preferred_study_intensity: StudyIntensity,
    pub priority_subjects: Vec<String>,
    pub avoid_subjects: Vec<String>,
    pub preferred_days: Vec<DayOfWeek>,
    pub preferred_times: Vec<TimeSlot>,
    pub balance_theory_practice: bool,
    pub prefer_prerequisites_early: bool,
    pub graduation_target: Option<String>,
}

/// Default preferences for new students
#[cw_serde]
pub struct DefaultSchedulePreferences {
    pub max_subjects_per_semester: u32,
    pub study_intensity: StudyIntensity,
    pub balance_theory_practice: bool,
    pub prefer_prerequisites_early: bool,
}

/// Study intensity levels
#[cw_serde]
pub enum StudyIntensity {
    Light,      // 3-4 subjects per semester
    Moderate,   // 5-6 subjects per semester
    Intensive,  // 7-8 subjects per semester
    Maximum,    // 8+ subjects per semester
}

/// Academic standing
#[cw_serde]
pub enum AcademicStanding {
    Excellent,     // GPA >= 9.0
    Good,          // GPA >= 8.0
    Satisfactory,  // GPA >= 7.0
    Probation,     // GPA < 7.0
}

/// Days of the week
#[cw_serde]
pub enum DayOfWeek {
    Monday,
    Tuesday,
    Wednesday,
    Thursday,
    Friday,
    Saturday,
    Sunday,
}

/// Time slots for classes
#[cw_serde]
pub enum TimeSlot {
    EarlyMorning,  // 7:00-9:00
    Morning,       // 9:00-12:00
    Afternoon,     // 13:00-17:00
    Evening,       // 18:00-22:00
    Night,         // 22:00+
}

/// Subject information for scheduling
#[cw_serde]
pub struct SubjectScheduleInfo {
    pub subject_id: String,
    pub title: String,
    pub credits: u32,
    pub prerequisites: Vec<String>,
    pub corequisites: Vec<String>,
    pub semester_offered: Vec<u32>,
    pub max_students: Option<u32>,
    pub difficulty_level: u32,
    pub workload_hours: u32,
    pub is_elective: bool,
    pub department: String,
    pub professor: Option<String>,
    pub schedule_info: ClassSchedule,
    pub ipfs_link: Option<String>,
}

/// Class schedule information
#[cw_serde]
pub struct ClassSchedule {
    pub days: Vec<DayOfWeek>,
    pub time_slots: Vec<TimeSlot>,
    pub location: Option<String>,
    pub online_option: bool,
}

/// Schedule recommendation for a specific semester
#[cw_serde]
pub struct ScheduleRecommendation {
    pub student_id: String,
    pub target_semester: u32,
    pub recommended_subjects: Vec<RecommendedSubject>,
    pub alternative_subjects: Vec<RecommendedSubject>,
    pub total_credits: u32,
    pub estimated_workload: u32,
    pub difficulty_score: u32,
    pub completion_percentage: String,
    pub notes: Vec<String>,
    pub confidence_score: u32,
    pub generated_timestamp: u64,
    pub algorithm_used: RecommendationAlgorithm,
}

/// Individual subject recommendation
#[cw_serde]
pub struct RecommendedSubject {
    pub subject_id: String,
    pub title: String,
    pub credits: u32,
    pub priority: Priority,
    pub recommendation_reason: RecommendationReason,
    pub prerequisites_met: bool,
    pub estimated_difficulty: u32,
    pub estimated_workload: u32,
    pub schedule_conflicts: Vec<String>,
    pub alternative_semesters: Vec<u32>,
}

/// Priority levels for subject recommendations
#[cw_serde]
pub enum Priority {
    Critical,    // Must take this semester
    High,        // Should take this semester
    Medium,      // Good to take this semester
    Low,         // Can defer to later
    Optional,    // Elective suggestion
}

/// Reasons for recommending a subject
#[cw_serde]
pub enum RecommendationReason {
    MandatoryForGraduation,
    PrerequisiteForFutureSubjects,
    OptimalSequencing,
    WorkloadBalancing,
    StudentPreference,
    ProfessorAvailability,
    ScheduleOptimization,
    GraduationSpeedOptimization,
}

/// Complete academic path from current state to graduation
#[cw_serde]
pub struct AcademicPath {
    pub path_id: String,
    pub student_id: String,
    pub path_name: String,
    pub semesters: Vec<SemesterPlan>,
    pub total_duration_semesters: u32,
    pub total_credits: u32,
    pub expected_graduation_date: String,
    pub optimization_criteria: OptimizationCriteria,
    pub path_metrics: PathMetrics,
    pub created_timestamp: u64,
    pub last_updated: u64,
    pub is_optimized: bool,
    pub status: PathStatus,
}

/// Plan for a specific semester
#[cw_serde]
pub struct SemesterPlan {
    pub semester_number: u32,
    pub subjects: Vec<PlannedSubject>,
    pub total_credits: u32,
    pub estimated_difficulty: u32,
    pub estimated_workload: u32,
    pub notes: Vec<String>,
    pub flexibility_score: u32, // How easy it is to modify this semester
}

/// Subject planned for a specific semester
#[cw_serde]
pub struct PlannedSubject {
    pub subject_id: String,
    pub title: String,
    pub credits: u32,
    pub is_mandatory: bool,
    pub backup_options: Vec<String>, // Alternative subjects if this one isn't available
    pub placement_reason: String,
}

/// Optimization criteria for academic path
#[cw_serde]
pub enum OptimizationCriteria {
    Fastest,          // Minimize total semesters
    Balanced,         // Balance workload across semesters
    EasiestFirst,     // Take easier subjects first
    PrerequisiteOptimal, // Optimize prerequisite chains
    Custom(OptimizationWeights), // Custom weight configuration
}

/// Metrics for evaluating academic paths
#[cw_serde]
pub struct PathMetrics {
    pub total_duration_semesters: u32,
    pub average_credits_per_semester: u32,
    pub workload_variance: u32,
    pub difficulty_progression_score: u32,
    pub prerequisite_efficiency: u32,
    pub schedule_flexibility: u32,
    pub risk_factors: Vec<String>,
}

/// Status of academic path
#[cw_serde]
pub enum PathStatus {
    Draft,       // Being created/modified
    Active,      // Currently being followed
    Completed,   // Student graduated
    Abandoned,   // Student changed path
    Outdated,    // Curriculum changed
}

/// Available subjects for a semester (query result)
#[cw_serde]
pub struct AvailableSubjects {
    pub semester: u32,
    pub subjects: Vec<AvailableSubject>,
    pub total_count: u32,
}

/// Individual available subject
#[cw_serde]
pub struct AvailableSubject {
    pub subject_id: String,
    pub title: String,
    pub credits: u32,
    pub can_enroll: bool,
    pub missing_prerequisites: Vec<String>,
    pub conflicts: Vec<String>,
    pub recommendation_priority: Option<Priority>,
}

// Storage items
pub const STATE: Item<State> = Item::new("state");
pub const SCHEDULE_CONFIG: Item<ScheduleConfig> = Item::new("schedule_config");

// Maps
pub const STUDENT_PROGRESS: Map<&str, StudentProgress> = Map::new("student_progress");
pub const SUBJECT_SCHEDULE_INFO: Map<&str, SubjectScheduleInfo> = Map::new("subject_schedule_info");
pub const SCHEDULE_RECOMMENDATIONS: Map<(&str, u32), ScheduleRecommendation> = Map::new("schedule_recommendations");
pub const ACADEMIC_PATHS: Map<&str, AcademicPath> = Map::new("academic_paths");
pub const STUDENT_PATHS: Map<&str, Vec<String>> = Map::new("student_paths"); // student_id -> path_ids

// IPFS cache for subject content
use crate::ipfs::SubjectContent;
pub const IPFS_CACHE: Map<&str, SubjectContent> = Map::new("ipfs_cache");