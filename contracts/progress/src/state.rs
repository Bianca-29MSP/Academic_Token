// src/state.rs
use cosmwasm_schema::cw_serde;
use cosmwasm_std::Addr;
use cw_storage_plus::{Item, Map};

/// Main contract state
#[cw_serde]
pub struct State {
    pub owner: Addr,
    pub total_students: u64,
    pub total_institutions: u64,
    pub total_progress_records: u64,
    pub analytics_enabled: bool,
}

/// Configuration for progress tracking and analytics
#[cw_serde]
pub struct ProgressConfig {
    pub update_frequency: UpdateFrequency,
    pub analytics_depth: AnalyticsDepth,
    pub retention_period_days: u32,
    pub dashboard_refresh_hours: u32,
    pub benchmark_computation_enabled: bool,
}

/// How frequently progress is updated
#[cw_serde]
pub enum UpdateFrequency {
    RealTime,    // Update immediately on changes
    Daily,       // Update once per day
    Weekly,      // Update weekly
    Monthly,     // Update monthly
    OnDemand,    // Update only when requested
}

/// Depth of analytics computation
#[cw_serde]
pub enum AnalyticsDepth {
    Basic,       // Essential metrics only
    Standard,    // Standard academic metrics
    Advanced,    // Advanced analytics with predictions
    Comprehensive, // Full analytics suite
}

/// Complete student progress tracking
#[cw_serde]
pub struct StudentProgress {
    pub student_id: String,
    pub institution_id: String,
    pub course_id: String,
    pub curriculum_id: String,
    
    // Academic standing
    pub current_semester: u32,
    pub total_semesters_enrolled: u32,
    pub academic_status: AcademicStatus,
    
    // Credit tracking
    pub total_credits_completed: u32,
    pub total_credits_required: u32,
    pub credits_in_progress: u32,
    
    // Grade performance
    pub gpa: String,
    pub cgpa: String, // Cumulative GPA
    pub grade_distribution: GradeDistribution,
    
    // Subject tracking
    pub completed_subjects: Vec<CompletedSubjectProgress>,
    pub current_subjects: Vec<CurrentSubjectProgress>,
    pub failed_subjects: Vec<FailedSubjectProgress>,
    
    // Milestone tracking
    pub milestones_achieved: Vec<AcademicMilestone>,
    pub upcoming_milestones: Vec<UpcomingMilestone>,
    
    // Performance metrics
    pub performance_metrics: PerformanceMetrics,
    
    // Predictions and recommendations
    pub graduation_forecast: GraduationForecast,
    pub risk_factors: Vec<RiskFactor>,
    pub improvement_recommendations: Vec<String>,
    
    // Timestamps
    pub enrollment_date: String,
    pub last_updated: u64,
    pub next_milestone_date: Option<String>,
}

/// Student's academic status
#[cw_serde]
pub enum AcademicStatus {
    Active,           // Currently enrolled and in good standing
    Probation,        // Academic probation
    Suspended,        // Temporarily suspended
    Graduated,        // Successfully graduated
    Withdrawn,        // Voluntarily withdrawn
    Expelled,         // Expelled from institution
    Transfer,         // Transferred to another institution
}

/// Distribution of grades across subjects
#[cw_serde]
pub struct GradeDistribution {
    pub excellent_count: u32,    // A grades (90-100)
    pub good_count: u32,         // B grades (80-89)
    pub satisfactory_count: u32, // C grades (70-79)
    pub poor_count: u32,         // D grades (60-69)
    pub fail_count: u32,         // F grades (<60)
}

/// Progress for a completed subject
#[cw_serde]
pub struct CompletedSubjectProgress {
    pub subject_id: String,
    pub title: String,
    pub credits: u32,
    pub final_grade: u32,
    pub letter_grade: String,
    pub completion_date: String,
    pub semester_taken: u32,
    pub attempt_number: u32,
    pub study_hours_logged: Option<u32>,
    pub difficulty_rating: Option<u32>,
    pub satisfaction_rating: Option<u32>,
    pub nft_token_id: String,
    pub competencies_gained: Vec<String>,
}

/// Progress for a currently enrolled subject
#[cw_serde]
pub struct CurrentSubjectProgress {
    pub subject_id: String,
    pub title: String,
    pub credits: u32,
    pub enrollment_date: String,
    pub current_grade: Option<u32>,
    pub attendance_rate: Option<u32>,
    pub assignments_completed: u32,
    pub assignments_total: u32,
    pub study_hours_logged: u32,
    pub predicted_final_grade: Option<u32>,
    pub risk_level: RiskLevel,
    pub last_activity_date: Option<String>,
}

/// Progress for a failed subject
#[cw_serde]
pub struct FailedSubjectProgress {
    pub subject_id: String,
    pub title: String,
    pub credits: u32,
    pub final_grade: u32,
    pub completion_date: String,
    pub semester_taken: u32,
    pub attempt_number: u32,
    pub failure_reasons: Vec<String>,
    pub remediation_required: bool,
    pub retry_available: bool,
    pub retry_recommendations: Vec<String>,
}

/// Academic milestones achieved
#[cw_serde]
pub struct AcademicMilestone {
    pub milestone_id: String,
    pub milestone_type: MilestoneType,
    pub title: String,
    pub description: String,
    pub achieved_date: String,
    pub semester_achieved: u32,
    pub credits_at_achievement: u32,
    pub gpa_at_achievement: String,
    pub significance: MilestoneSignificance,
}

/// Upcoming milestones to track
#[cw_serde]
pub struct UpcomingMilestone {
    pub milestone_id: String,
    pub milestone_type: MilestoneType,
    pub title: String,
    pub description: String,
    pub estimated_date: String,
    pub requirements: Vec<String>,
    pub progress_percentage: u32,
    pub probability: u32, // Probability of achieving on time
}

/// Types of academic milestones
#[cw_serde]
pub enum MilestoneType {
    FirstSemester,        // Completion of first semester
    QuarterProgress,      // 25% of program complete
    HalfwayPoint,         // 50% of program complete
    ThreeQuarterProgress, // 75% of program complete
    FinalSemester,        // Entering final semester
    Graduation,           // Program completion
    HighGPA,              // Achieving high GPA threshold
    PerfectSemester,      // Perfect GPA for a semester
    ResearchCompletion,   // Completion of research project
    InternshipCompletion, // Completion of internship
    CertificationEarned,  // Professional certification
    AcademicAward,        // Academic achievement award
    Custom(String),       // Institution-specific milestone
}

/// Significance levels for milestones
#[cw_serde]
pub enum MilestoneSignificance {
    Minor,     // Small achievement
    Standard,  // Regular milestone
    Major,     // Important achievement
    Critical,  // Graduation-critical milestone
}

/// Risk levels for current subjects
#[cw_serde]
pub enum RiskLevel {
    Low,       // Student performing well
    Medium,    // Some concerns
    High,      // At risk of failing
    Critical,  // Likely to fail without intervention
}

/// Comprehensive performance metrics
#[cw_serde]
pub struct PerformanceMetrics {
    // Academic performance
    pub success_rate: u32,               // Percentage of subjects passed
    pub retry_rate: u32,                 // Percentage of subjects retaken
    pub grade_trend: GradeTrend,         // Trending up, down, or stable
    pub consistency_score: u32,          // How consistent grades are
    
    // Engagement metrics
    pub attendance_average: u32,         // Average attendance rate
    pub assignment_completion_rate: u32, // Percentage of assignments completed
    pub study_hours_per_credit: u32,     // Study efficiency
    
    // Progression metrics
    pub credits_per_semester_avg: u32,   // Average credits per semester
    pub time_to_completion_ratio: u32,   // Actual vs expected time ratio
    pub prerequisite_success_rate: u32,  // Success in prerequisite chains
    
    // Comparative metrics
    pub percentile_rank_gpa: u32,        // GPA percentile in cohort
    pub percentile_rank_progress: u32,   // Progress percentile in cohort
    
    // Predictive metrics
    pub graduation_probability: u32,     // Probability of graduation
    pub time_extension_probability: u32, // Probability of needing extra time
}

/// Grade trend analysis
#[cw_serde]
pub enum GradeTrend {
    StronglyIncreasing,   // Grades improving significantly
    ModeratelyIncreasing, // Grades slowly improving
    Stable,               // Grades consistent
    ModeratelyDecreasing, // Grades slowly declining
    StronglyDecreasing,   // Grades declining significantly
    Volatile,             // Grades highly variable
}

/// Graduation forecast
#[cw_serde]
pub struct GraduationForecast {
    pub estimated_graduation_date: String,
    pub confidence_level: u32,
    pub remaining_semesters: u32,
    pub remaining_credits: u32,
    pub on_track: bool,
    pub potential_delays: Vec<DelayFactor>,
    pub acceleration_opportunities: Vec<String>,
}

/// Factors that could delay graduation
#[cw_serde]
pub struct DelayFactor {
    pub factor_type: DelayFactorType,
    pub description: String,
    pub impact_semesters: u32,
    pub mitigation_strategies: Vec<String>,
}

/// Types of delay factors
#[cw_serde]
pub enum DelayFactorType {
    AcademicPerformance,  // Poor grades requiring retakes
    PrerequisiteChain,    // Prerequisite sequence delays
    SubjectAvailability,  // Subjects not offered when needed
    PersonalCircumstances, // Student personal issues
    InstitutionalChanges, // Curriculum or policy changes
    TransferCredits,      // Transfer credit issues
}

/// Risk factors that could impact success
#[cw_serde]
pub struct RiskFactor {
    pub risk_type: RiskFactorType,
    pub severity: RiskSeverity,
    pub description: String,
    pub early_warning_indicators: Vec<String>,
    pub intervention_recommendations: Vec<String>,
    pub monitoring_frequency: MonitoringFrequency,
}

/// Types of risk factors
#[cw_serde]
pub enum RiskFactorType {
    AcademicPerformance,  // Declining grades
    Attendance,           // Poor attendance
    Engagement,           // Low engagement
    Prerequisites,        // Struggling with prerequisites
    Workload,            // Overwhelming course load
    TimeManagement,      // Poor time management
    FinancialStress,     // Financial difficulties
    HealthIssues,        // Health concerns
    SocialFactors,       // Social integration issues
}

/// Severity levels for risk factors
#[cw_serde]
#[derive(Eq, Hash)]
pub enum RiskSeverity {
    Low,       // Minor concern
    Moderate,  // Needs attention
    High,      // Urgent intervention needed
    Critical,  // Immediate action required
}

/// How frequently to monitor risk factors
#[cw_serde]
pub enum MonitoringFrequency {
    Weekly,
    Biweekly,
    Monthly,
    Quarterly,
}

/// Institution-wide analytics
#[cw_serde]
pub struct InstitutionAnalytics {
    pub institution_id: String,
    pub analytics_period: AnalyticsPeriod,
    
    // Overall metrics
    pub total_students: u32,
    pub active_students: u32,
    pub graduated_students: u32,
    pub dropout_rate: u32,
    
    // Performance metrics
    pub average_gpa: String,
    pub grade_distribution: InstitutionGradeDistribution,
    pub completion_rates: CompletionRates,
    
    // Progression metrics
    pub average_time_to_graduation: u32,
    pub on_track_percentage: u32,
    pub risk_distribution: RiskDistribution,
    
    // Comparative benchmarks
    pub benchmark_comparisons: Vec<BenchmarkComparison>,
    
    // Trends
    pub enrollment_trends: Vec<EnrollmentTrend>,
    pub performance_trends: Vec<PerformanceTrend>,
    
    // Subject-specific analytics
    pub subject_analytics: Vec<SubjectAnalytics>,
    
    // Generated timestamp
    pub generated_timestamp: u64,
    pub next_update_due: u64,
}

/// Time period for analytics
#[cw_serde]
pub struct AnalyticsPeriod {
    pub period_type: PeriodType,  // ← Agora está definido
    pub start_date: String,
    pub end_date: String,
}

/// Types of analytics periods
#[cw_serde]
pub enum PeriodType {
    AcademicYear,  
    CalendarYear, 
    Semester,
    Quarter,
}

// Adicionar o tipo Priority que está faltando
#[cw_serde]
pub enum Priority {
    Critical,
    High,
    Medium,
    Low,
}

/// Institution-wide grade distribution
#[cw_serde]
pub struct InstitutionGradeDistribution {
    pub excellent_percentage: u32,
    pub good_percentage: u32,
    pub satisfactory_percentage: u32,
    pub poor_percentage: u32,
    pub fail_percentage: u32,
}

/// Completion rates across different metrics
#[cw_serde]
pub struct CompletionRates {
    pub overall_completion_rate: u32,
    pub completion_rate_by_course: Vec<CourseCompletionRate>,
    pub completion_rate_by_semester: Vec<SemesterCompletionRate>,
}

/// Completion rate for a specific course
#[cw_serde]
pub struct CourseCompletionRate {
    pub course_id: String,
    pub course_title: String,
    pub completion_rate: u32,
    pub average_time_semesters: u32,
}

/// Completion rate for a specific semester
#[cw_serde]
pub struct SemesterCompletionRate {
    pub semester_id: String,
    pub completion_rate: u32,
    pub dropout_rate: u32,
}

/// Distribution of students by risk level
#[cw_serde]
pub struct RiskDistribution {
    pub low_risk_count: u32,
    pub medium_risk_count: u32,
    pub high_risk_count: u32,
    pub critical_risk_count: u32,
}

/// Benchmark comparison data
#[cw_serde]
pub struct BenchmarkComparison {
    pub metric_name: String,
    pub institution_value: u32,
    pub benchmark_value: u32,
    pub percentile_rank: u32,
    pub comparison_group: String,
}

/// Enrollment trend data
#[cw_serde]
pub struct EnrollmentTrend {
    pub period: String,
    pub new_enrollments: u32,
    pub total_enrolled: u32,
    pub dropout_count: u32,
    pub transfer_in_count: u32,
    pub transfer_out_count: u32,
}

/// Trend direction for analytics
#[cw_serde]
pub enum TrendDirection {
    Improving,
    Stable,
    Declining,
}

/// Performance trend data
#[cw_serde]
pub struct PerformanceTrend {
    pub period: String,
    pub average_gpa: String,
    pub completion_rate: u32,
    pub success_rate: u32,
    pub trend_direction: TrendDirection,
}

/// Analytics for individual subjects
#[cw_serde]
pub struct SubjectAnalytics {
    pub subject_id: String,
    pub subject_title: String,
    pub enrollment_count: u32,
    pub completion_rate: u32,
    pub average_grade: u32,
    pub difficulty_rating: u32,
    pub student_satisfaction: u32,
    pub prerequisite_success_correlation: u32,
    pub retry_rate: u32,
    pub common_failure_reasons: Vec<String>,
}

/// Student dashboard data
#[cw_serde]
pub struct StudentDashboard {
    pub student_id: String,
    pub last_updated: u64,
    
    // Quick stats
    pub quick_stats: QuickStats,
    
    // Current status
    pub current_status: CurrentStatus,
    
    // Performance visualizations
    pub grade_history: Vec<GradeHistoryPoint>,
   pub progress_timeline: Vec<ProgressTimelineEvent>,
   
   // Alerts and notifications
   pub urgent_alerts: Vec<Alert>,
   pub upcoming_deadlines: Vec<Deadline>,
   pub recommendations: Vec<Recommendation>,
   
   // Achievements and milestones
   pub recent_achievements: Vec<Achievement>,
   pub progress_towards_goals: Vec<GoalProgress>,
   
   // Analytics insights
   pub performance_insights: Vec<PerformanceInsight>,
   pub risk_warnings: Vec<RiskWarning>,
   
   // Comparative data
   pub peer_comparisons: Option<PeerComparison>,
}

/// Quick statistics for dashboard
#[cw_serde]
pub struct QuickStats {
   pub current_gpa: String,
   pub credits_completed: u32,
   pub credits_remaining: u32,
   pub completion_percentage: u32,
   pub semesters_remaining: u32,
   pub current_semester_subjects: u32,
}

/// Current academic status
#[cw_serde]
pub struct CurrentStatus {
   pub academic_standing: AcademicStatus,
   pub current_semester: u32,
   pub enrollment_status: EnrollmentStatus,
   pub next_milestone: Option<String>,
   pub days_to_next_milestone: Option<u32>,
}

/// Enrollment status details
#[cw_serde]
pub enum EnrollmentStatus {
   FullTime,
   PartTime,
   OnLeave,
   Graduated,
   Withdrawn,
}

/// Point in grade history for visualization
#[cw_serde]
pub struct GradeHistoryPoint {
   pub semester: u32,
   pub gpa: String,
   pub credits_earned: u32,
   pub subjects_completed: u32,
   pub trend_indicator: TrendDirection,
}

/// Event in progress timeline
#[cw_serde]
pub struct ProgressTimelineEvent {
   pub date: String,
   pub event_type: TimelineEventType,
   pub title: String,
   pub description: String,
   pub impact_level: ImpactLevel,
}

/// Types of timeline events
#[cw_serde]
pub enum TimelineEventType {
   Enrollment,
   SubjectCompletion,
   MilestoneAchievement,
   GradeImprovement,
   Warning,
   Achievement,
   StatusChange,
}

/// Impact level of events
#[cw_serde]
pub enum ImpactLevel {
   Low,
   Medium,
   High,
   Critical,
}

/// Alert for immediate attention
#[cw_serde]
pub struct Alert {
   pub alert_id: String,
   pub alert_type: AlertType,
   pub severity: AlertSeverity,
   pub title: String,
   pub message: String,
   pub action_required: bool,
   pub deadline: Option<String>,
   pub related_subject: Option<String>,
}

/// Types of alerts
#[cw_serde]
pub enum AlertType {
   GradeWarning,
   AttendanceWarning,
   DeadlineMissing,
   PrerequisiteIssue,
   RegistrationRequired,
   PaymentDue,
   DocumentRequired,
   AcademicProbation,
}

/// Alert severity levels
#[cw_serde]
pub enum AlertSeverity {
   Info,
   Warning,
   Critical,
   Emergency,
}

/// Upcoming deadline
#[cw_serde]
pub struct Deadline {
   pub deadline_id: String,
   pub title: String,
   pub description: String,
   pub due_date: String,
   pub days_remaining: i32,
   pub subject_id: Option<String>,
   pub deadline_type: DeadlineType,
   pub completion_status: CompletionStatus,
}

/// Types of deadlines
#[cw_serde]
pub enum DeadlineType {
   Assignment,
   Exam,
   Registration,
   Payment,
   Documentation,
   Application,
   Thesis,
}

/// Completion status of deadlines
#[cw_serde]
pub enum CompletionStatus {
   NotStarted,
   InProgress,
   Completed,
   Overdue,
}

/// Recommendation for student
#[cw_serde]
pub struct Recommendation {
   pub recommendation_id: String,
   pub category: RecommendationCategory,
   pub priority: Priority,
   pub title: String,
   pub description: String,
   pub action_steps: Vec<String>,
   pub expected_benefit: String,
   pub effort_level: EffortLevel,
}

/// Categories of recommendations
#[cw_serde]
pub enum RecommendationCategory {
   AcademicImprovement,
   StudyStrategy,
   TimeManagement,
   SubjectSelection,
   CareerDevelopment,
   HealthWellness,
   ResourceUtilization,
}

/// Effort level required for recommendations
#[cw_serde]
pub enum EffortLevel {
   Low,     // Easy to implement
   Medium,  // Moderate effort required
   High,    // Significant effort required
}

/// Achievement earned by student
#[cw_serde]
pub struct Achievement {
   pub achievement_id: String,
   pub title: String,
   pub description: String,
   pub earned_date: String,
   pub achievement_type: AchievementType,
   pub rarity: AchievementRarity,
   pub points_awarded: u32,
}

/// Types of achievements
#[cw_serde]
pub enum AchievementType {
   Academic,
   Participation,
   Improvement,
   Consistency,
   Leadership,
   Innovation,
   Community,
}

/// Rarity levels for achievements
#[cw_serde]
pub enum AchievementRarity {
   Common,
   Uncommon,
   Rare,
   Epic,
   Legendary,
}

/// Progress toward specific goals
#[cw_serde]
pub struct GoalProgress {
   pub goal_id: String,
   pub goal_title: String,
   pub goal_description: String,
   pub current_progress: u32,
   pub target_value: u32,
   pub progress_percentage: u32,
   pub estimated_completion_date: Option<String>,
   pub on_track: bool,
}

/// Performance insight for analytics
#[cw_serde]
pub struct PerformanceInsight {
   pub insight_id: String,
   pub insight_type: InsightType,
   pub title: String,
   pub description: String,
   pub supporting_data: Vec<String>,
   pub actionable_advice: Vec<String>,
   pub confidence_level: u32,
}

/// Types of insights
#[cw_serde]
pub enum InsightType {
   StrengthIdentification,
   WeaknessIdentification,
   TrendAnalysis,
   PredictiveWarning,
   OpportunityIdentification,
   BenchmarkComparison,
}

/// Risk warning for students
#[cw_serde]
pub struct RiskWarning {
   pub warning_id: String,
   pub risk_type: RiskFactorType,
   pub severity: RiskSeverity,
   pub probability: u32,
   pub description: String,
   pub early_indicators: Vec<String>,
   pub prevention_strategies: Vec<String>,
   pub support_resources: Vec<String>,
}

/// Peer comparison data
#[cw_serde]
pub struct PeerComparison {
   pub gpa_percentile: u32,
   pub progress_percentile: u32,
   pub credit_completion_rate_vs_peers: ComparisonResult,
   pub grade_trend_vs_peers: ComparisonResult,
   pub anonymous_peer_insights: Vec<String>,
}

/// Result of comparison with peers
#[cw_serde]
pub enum ComparisonResult {
   BelowAverage,
   Average,
   AboveAverage,
   TopTier,
}

// Storage items
pub const STATE: Item<State> = Item::new("state");
pub const PROGRESS_CONFIG: Item<ProgressConfig> = Item::new("progress_config");

// Maps
pub const STUDENT_PROGRESS: Map<&str, StudentProgress> = Map::new("student_progress");
pub const INSTITUTION_ANALYTICS: Map<&str, InstitutionAnalytics> = Map::new("institution_analytics");
pub const STUDENT_DASHBOARDS: Map<&str, StudentDashboard> = Map::new("student_dashboards");

// Index maps for efficient querying
pub const STUDENTS_BY_INSTITUTION: Map<&str, Vec<String>> = Map::new("students_by_institution");
pub const STUDENTS_BY_COURSE: Map<&str, Vec<String>> = Map::new("students_by_course");
pub const STUDENTS_BY_STATUS: Map<&str, Vec<String>> = Map::new("students_by_status");
pub const STUDENTS_BY_RISK_LEVEL: Map<&str, Vec<String>> = Map::new("students_by_risk");

// Historical data for trends
pub const PROGRESS_HISTORY: Map<(&str, u64), StudentProgress> = Map::new("progress_history");
pub const ANALYTICS_HISTORY: Map<(&str, u64), InstitutionAnalytics> = Map::new("analytics_history");

// Recommendations and paths
pub const RECOMMENDATIONS_BY_INSTITUTION: Map<&str, Vec<String>> = Map::new("rec_by_inst");
pub const PATHS_BY_INSTITUTION: Map<&str, Vec<String>> = Map::new("paths_by_inst");
