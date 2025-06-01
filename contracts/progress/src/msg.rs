// src/msg.rs
use cosmwasm_schema::{cw_serde, QueryResponses};
use crate::state::{
    State, ProgressConfig, StudentProgress, InstitutionAnalytics, StudentDashboard,
    UpdateFrequency, AnalyticsDepth, AcademicStatus, RiskFactorType, RiskSeverity,
    AnalyticsPeriod,
};

/// Instantiate message
#[cw_serde]
pub struct InstantiateMsg {
    pub owner: Option<String>,
    pub analytics_enabled: bool,
    pub update_frequency: UpdateFrequency,
    pub analytics_depth: AnalyticsDepth,
}

/// Execute messages
#[cw_serde]
pub enum ExecuteMsg {
    /// Update contract configuration
    UpdateConfig {
        new_owner: Option<String>,
        analytics_enabled: Option<bool>,
        update_frequency: Option<UpdateFrequency>,
        analytics_depth: Option<AnalyticsDepth>,
        retention_period_days: Option<u32>,
    },
    
    /// Register or update student progress
    UpdateStudentProgress {
        student_progress: StudentProgress,
        force_analytics_refresh: Option<bool>,
    },
    
    /// Batch update multiple student progress records
    BatchUpdateStudentProgress {
        student_progress_list: Vec<StudentProgress>,
        analytics_refresh_mode: AnalyticsRefreshMode,
    },
    
    /// Record subject completion
    RecordSubjectCompletion {
        student_id: String,
        subject_id: String,
        final_grade: u32,
        completion_date: String,
        study_hours: Option<u32>,
        difficulty_rating: Option<u32>,
        satisfaction_rating: Option<u32>,
        nft_token_id: String,
    },
    
    /// Record subject enrollment
    RecordSubjectEnrollment {
        student_id: String,
        subject_id: String,
        enrollment_date: String,
        expected_completion: String,
    },
    
    /// Update current subject progress
    UpdateCurrentSubjectProgress {
        student_id: String,
        subject_id: String,
        current_grade: Option<u32>,
        attendance_rate: Option<u32>,
        assignments_completed: u32,
        study_hours_logged: u32,
    },
    
    /// Record milestone achievement
    RecordMilestoneAchievement {
        student_id: String,
        milestone_id: String,
        achievement_date: String,
        notes: Option<String>,
    },
    
    /// Add risk factor for student
    AddRiskFactor {
        student_id: String,
        risk_type: RiskFactorType,
        severity: RiskSeverity,
        description: String,
        intervention_recommendations: Vec<String>,
    },
    
    /// Remove risk factor
    RemoveRiskFactor {
        student_id: String,
        risk_type: RiskFactorType,
    },
    
    /// Generate institution analytics
    GenerateInstitutionAnalytics {
        institution_id: String,
        period: AnalyticsPeriod,
        force_refresh: Option<bool>,
    },
    
    /// Generate student dashboard
    GenerateStudentDashboard {
        student_id: String,
        include_peer_comparison: Option<bool>,
        force_refresh: Option<bool>,
    },
    
    /// Batch generate dashboards for institution
    BatchGenerateDashboards {
        institution_id: String,
        student_ids: Option<Vec<String>>, // If None, generate for all students
        include_analytics: bool,
    },
    
    /// Run analytics computation for all institutions
    RunGlobalAnalytics {
        period: AnalyticsPeriod,
        institutions: Option<Vec<String>>, // If None, process all
    },
    
    /// Update academic status
    UpdateAcademicStatus {
        student_id: String,
        new_status: AcademicStatus,
        effective_date: String,
        reason: Option<String>,
    },
    
    /// Archive old data (for storage management)
    ArchiveOldData {
        cutoff_date: String,
        dry_run: Option<bool>,
    },
    
    /// Reset student progress (admin function)
    ResetStudentProgress {
        student_id: String,
        reason: String,
        backup: bool,
    },
}

/// Mode for refreshing analytics during batch operations
#[cw_serde]
pub enum AnalyticsRefreshMode {
    None,           // Don't refresh analytics
    Individual,     // Refresh for each student individually
    Batch,          // Refresh once at the end for all affected institutions
    Background,     // Schedule for background processing
}

/// Query messages
#[cw_serde]
#[derive(QueryResponses)]
pub enum QueryMsg {
    /// Get contract state
    #[returns(StateResponse)]
    GetState {},
    
    /// Get progress configuration
    #[returns(ProgressConfigResponse)]
    GetConfig {},
    
    /// Get student progress
    #[returns(StudentProgressResponse)]
    GetStudentProgress { student_id: String },
    
    /// Get student dashboard
    #[returns(StudentDashboardResponse)]
    GetStudentDashboard { 
        student_id: String,
        include_peer_comparison: Option<bool>,
    },
    
    /// Get institution analytics
    #[returns(InstitutionAnalyticsResponse)]
    GetInstitutionAnalytics { 
        institution_id: String,
        period: Option<AnalyticsPeriod>,
    },
    
    /// Get students by institution
    #[returns(StudentsListResponse)]
    GetStudentsByInstitution {
        institution_id: String,
        status_filter: Option<AcademicStatus>,
        limit: Option<u32>,
        start_after: Option<String>,
    },
    
    /// Get students by course
    #[returns(StudentsListResponse)]
    GetStudentsByCourse {
        course_id: String,
        status_filter: Option<AcademicStatus>,
        limit: Option<u32>,
        start_after: Option<String>,
    },
    
    /// Get students at risk
    #[returns(AtRiskStudentsResponse)]
    GetAtRiskStudents {
        institution_id: Option<String>,
        risk_level_filter: Option<RiskSeverity>,
        limit: Option<u32>,
    },
    
    /// Get top performers
    #[returns(TopPerformersResponse)]
    GetTopPerformers {
        institution_id: Option<String>,
        course_id: Option<String>,
        metric: PerformanceMetric,
        limit: Option<u32>,
    },
    
    /// Get graduation forecast
    #[returns(GraduationForecastResponse)]
    GetGraduationForecast {
        student_id: String,
        scenario: ForecastScenario,
    },
    
    /// Get comparative analytics
    #[returns(ComparativeAnalyticsResponse)]
    GetComparativeAnalytics {
        institution_id: String,
        comparison_group: ComparisonGroup,
        metrics: Vec<String>,
    },
    
    /// Get progress trends
    #[returns(ProgressTrendsResponse)]
    GetProgressTrends {
        entity_id: String, // student_id or institution_id
        entity_type: EntityType,
        period: AnalyticsPeriod,
        metrics: Vec<String>,
    },
    
    /// Get subject performance analytics
    #[returns(SubjectPerformanceResponse)]
    GetSubjectPerformance {
        subject_id: String,
        institution_id: Option<String>,
        period: Option<AnalyticsPeriod>,
    },
    
    /// Get cohort analysis
    #[returns(CohortAnalysisResponse)]
    GetCohortAnalysis {
        institution_id: String,
        course_id: String,
        enrollment_year: u32,
        analysis_type: CohortAnalysisType,
    },
    
    /// Get predictive insights
    #[returns(PredictiveInsightsResponse)]
    GetPredictiveInsights {
        student_id: String,
        prediction_horizon_semesters: u32,
    },
    
    /// Search students
    #[returns(StudentSearchResponse)]
    SearchStudents {
        query: StudentSearchQuery,
        limit: Option<u32>,
        start_after: Option<String>,
    },
    
    /// Get analytics summary
    #[returns(AnalyticsSummaryResponse)]
    GetAnalyticsSummary {
        scope: AnalyticsScope,
        quick_stats_only: Option<bool>,
    },
}

/// Performance metrics for ranking
#[cw_serde]
pub enum PerformanceMetric {
    GPA,
    ProgressSpeed,
    ConsistencyScore,
    ImprovementRate,
    CompletionRate,
}

/// Forecast scenarios
#[cw_serde]
pub enum ForecastScenario {
    Current,      // Based on current trajectory
    Optimistic,   // Best case scenario
    Pessimistic,  // Worst case scenario
    Intervention, // With recommended interventions
}

/// Comparison groups for analytics
#[cw_serde]
pub enum ComparisonGroup {
    NationalAverage,
    InstitutionType,
    SimilarPrograms,
    HistoricalData,
    CustomGroup(Vec<String>),
}

/// Entity types for trend analysis
#[cw_serde]
pub enum EntityType {
    Student,
    Institution,
    Course,
    Subject,
}

/// Types of cohort analysis
#[cw_serde]
pub enum CohortAnalysisType {
    RetentionRates,
    ProgressRates,
    GraduationRates,
    GradeDistribution,
    TimeToCompletion,
}

/// Student search query parameters
#[cw_serde]
pub struct StudentSearchQuery {
    pub name_contains: Option<String>,
    pub institution_id: Option<String>,
    pub course_id: Option<String>,
    pub status: Option<AcademicStatus>,
    pub gpa_min: Option<String>,
    pub gpa_max: Option<String>,
    pub credits_min: Option<u32>,
    pub credits_max: Option<u32>,
    pub risk_level: Option<RiskSeverity>,
    pub enrollment_year: Option<u32>,
}

/// Scope for analytics summary
#[cw_serde]
pub enum AnalyticsScope {
    Global,
    Institution(String),
    Course(String),
    Student(String),
}

// Response structs
#[cw_serde]
pub struct StateResponse {
    pub state: State,
}

#[cw_serde]
pub struct ProgressConfigResponse {
    pub config: ProgressConfig,
}

#[cw_serde]
pub struct StudentProgressResponse {
    pub progress: StudentProgress,
}

#[cw_serde]
pub struct StudentDashboardResponse {
    pub dashboard: StudentDashboard,
}

#[cw_serde]
pub struct InstitutionAnalyticsResponse {
    pub analytics: InstitutionAnalytics,
}

#[cw_serde]
pub struct StudentsListResponse {
    pub students: Vec<StudentSummary>,
    pub total_count: u32,
    pub has_more: bool,
}

#[cw_serde]
pub struct StudentSummary {
    pub student_id: String,
    pub institution_id: String,
    pub course_id: String,
    pub current_semester: u32,
    pub gpa: String,
    pub academic_status: AcademicStatus,
    pub completion_percentage: u32,
}

#[cw_serde]
pub struct AtRiskStudentsResponse {
    pub at_risk_students: Vec<AtRiskStudent>,
    pub summary: RiskSummary,
}

#[cw_serde]
pub struct AtRiskStudent {
    pub student_id: String,
    pub risk_level: RiskSeverity,
    pub risk_factors: Vec<RiskFactorType>,
    pub intervention_priority: u32,
    pub last_updated: u64,
}

#[cw_serde]
pub struct RiskSummary {
    pub total_at_risk: u32,
    pub critical_risk_count: u32,
    pub high_risk_count: u32,
    pub medium_risk_count: u32,
    pub intervention_success_rate: u32,
}

#[cw_serde]
pub struct TopPerformersResponse {
    pub top_performers: Vec<TopPerformer>,
    pub metric_used: PerformanceMetric,
    pub benchmark_value: String,
}

#[cw_serde]
pub struct TopPerformer {
    pub student_id: String,
    pub metric_value: String,
    pub percentile_rank: u32,
    pub achievement_highlights: Vec<String>,
}

#[cw_serde]
pub struct GraduationForecastResponse {
    pub forecast: crate::state::GraduationForecast,
    pub scenario_comparisons: Vec<ScenarioComparison>,
    pub confidence_factors: Vec<String>,
}

#[cw_serde]
pub struct ScenarioComparison {
    pub scenario: ForecastScenario,
    pub estimated_date: String,
    pub probability: u32,
    pub key_assumptions: Vec<String>,
}

#[cw_serde]
pub struct ComparativeAnalyticsResponse {
    pub institution_metrics: Vec<MetricComparison>,
    pub percentile_rankings: Vec<PercentileRanking>,
    pub improvement_opportunities: Vec<String>,
}

#[cw_serde]
pub struct MetricComparison {
    pub metric_name: String,
    pub institution_value: String,
    pub comparison_value: String,
    pub performance_ratio: String,
    pub trend: TrendDirection,
}

pub use crate::state::TrendDirection;

#[cw_serde]
pub struct PercentileRanking {
    pub metric_name: String,
    pub percentile: u32,
    pub rank_category: RankCategory,
}

#[cw_serde]
pub enum RankCategory {
    TopTier,     // 90-100th percentile
    HighPerforming, // 75-89th percentile
    AboveAverage,   // 60-74th percentile
    Average,        // 40-59th percentile
    BelowAverage,   // 25-39th percentile
    LowPerforming,  // 10-24th percentile
    Bottom,         // 0-9th percentile
}

#[cw_serde]
pub struct ProgressTrendsResponse {
    pub trend_data: Vec<TrendDataPoint>,
    pub trend_analysis: TrendAnalysis,
    pub forecasted_values: Vec<ForecastPoint>,
}

#[cw_serde]
pub struct TrendDataPoint {
    pub period: String,
    pub value: String,
    pub period_change: String,
    pub cumulative_change: String,
}

#[cw_serde]
pub struct TrendAnalysis {
    pub overall_direction: TrendDirection,
    pub volatility: VolatilityLevel,
    pub seasonal_patterns: Vec<String>,
    pub significant_events: Vec<SignificantEvent>,
}

#[cw_serde]
pub enum VolatilityLevel {
    Low,
    Moderate,
    High,
    Extreme,
}

#[cw_serde]
pub struct SignificantEvent {
    pub date: String,
    pub event_type: String,
    pub impact_description: String,
    pub impact_magnitude: u32,
}

#[cw_serde]
pub struct ForecastPoint {
    pub period: String,
    pub predicted_value: String,
    pub confidence_interval_low: String,
    pub confidence_interval_high: String,
    pub confidence_level: u32,
}

#[cw_serde]
pub struct SubjectPerformanceResponse {
    pub subject_analytics: crate::state::SubjectAnalytics,
    pub grade_distribution: Vec<GradeDistributionPoint>,
    pub success_predictors: Vec<SuccessPredictor>,
    pub improvement_recommendations: Vec<String>,
}

#[cw_serde]
pub struct GradeDistributionPoint {
    pub grade_range: String,
    pub student_count: u32,
    pub percentage: u32,
}

#[cw_serde]
pub struct SuccessPredictor {
    pub predictor_name: String,
    pub correlation_strength: u32,
    pub description: String,
}

#[cw_serde]
pub struct CohortAnalysisResponse {
    pub cohort_metrics: CohortMetrics,
    pub retention_analysis: RetentionAnalysis,
    pub performance_distribution: PerformanceDistribution,
    pub comparative_insights: Vec<String>,
}

#[cw_serde]
pub struct CohortMetrics {
    pub cohort_size: u32,
    pub current_active: u32,
    pub graduated: u32,
    pub withdrawn: u32,
    pub average_time_to_completion: u32,
}

#[cw_serde]
pub struct RetentionAnalysis {
    pub semester_retention_rates: Vec<SemesterRetention>,
    pub dropout_patterns: Vec<DropoutPattern>,
    pub retention_predictors: Vec<String>,
}

#[cw_serde]
pub struct SemesterRetention {
    pub semester: u32,
    pub retention_rate: u32,
    pub dropout_count: u32,
    pub key_factors: Vec<String>,
}

#[cw_serde]
pub struct DropoutPattern {
    pub semester_range: String,
    pub dropout_percentage: u32,
    pub common_reasons: Vec<String>,
}

#[cw_serde]
pub struct PerformanceDistribution {
    pub gpa_quartiles: GpaQuartiles,
    pub credit_completion_rates: Vec<CompletionRateDistribution>,
    pub time_to_graduation_distribution: Vec<TimeDistribution>,
}

#[cw_serde]
pub struct GpaQuartiles {
    pub q1: String,
    pub q2_median: String,
    pub q3: String,
    pub mean: String,
    pub standard_deviation: String,
}

#[cw_serde]
pub struct CompletionRateDistribution {
    pub rate_range: String,
    pub student_count: u32,
    pub percentage: u32,
}

#[cw_serde]
pub struct TimeDistribution {
    pub time_range_semesters: String,
    pub student_count: u32,
    pub percentage: u32,
}

#[cw_serde]
pub struct PredictiveInsightsResponse {
    pub predictions: Vec<Prediction>,
    pub risk_assessments: Vec<RiskAssessment>,
    pub opportunity_identifications: Vec<Opportunity>,
    pub recommended_interventions: Vec<Intervention>,
}

#[cw_serde]
pub struct Prediction {
    pub prediction_type: PredictionType,
    pub predicted_outcome: String,
    pub confidence_level: u32,
    pub time_horizon: String,
    pub key_factors: Vec<String>,
}

#[cw_serde]
pub enum PredictionType {
    GraduationProbability,
    NextSemesterGPA,
    SubjectSuccess,
    TimeToGraduation,
    RiskOfDropout,
}

#[cw_serde]
pub struct RiskAssessment {
    pub risk_factor: RiskFactorType,
    pub probability: u32,
    pub impact_severity: RiskSeverity,
    pub early_indicators: Vec<String>,
    pub mitigation_strategies: Vec<String>,
}

#[cw_serde]
pub struct Opportunity {
    pub opportunity_type: OpportunityType,
    pub description: String,
    pub potential_benefit: String,
    pub action_required: Vec<String>,
    pub success_probability: u32,
}

#[cw_serde]
pub enum OpportunityType {
    AcademicImprovement,
    AcceleratedCompletion,
    SkillDevelopment,
    CareerAdvancement,
    LeadershipRole,
    ResearchOpportunity,
}

#[cw_serde]
pub struct Intervention {
    pub intervention_id: String,
    pub intervention_type: InterventionType,
    pub target_issue: String,
    pub recommended_actions: Vec<String>,
    pub expected_outcome: String,
    pub success_probability: u32,
    pub effort_required: EffortLevel,
    pub timeline: String,
}

#[cw_serde]
pub enum InterventionType {
    AcademicSupport,
    Counseling,
    TutoringProgram,
    StudySkillsTraining,
    TimeManagementCoaching,
    PeerMentoring,
    ProfessionalDevelopment,
    HealthWellnessSupport,
}

#[cw_serde]
pub enum EffortLevel {
    Low,
    Medium,
    High,
    Intensive,
}

#[cw_serde]
pub struct StudentSearchResponse {
    pub students: Vec<StudentSearchResult>,
    pub total_matches: u32,
    pub search_metadata: SearchMetadata,
}

#[cw_serde]
pub struct StudentSearchResult {
    pub student_id: String,
    pub basic_info: StudentSummary,
    pub match_relevance: u32,
    pub highlighted_attributes: Vec<String>,
}

#[cw_serde]
pub struct SearchMetadata {
    pub query_time_ms: u32,
    pub filters_applied: Vec<String>,
    pub result_quality: SearchQuality,
}

#[cw_serde]
pub enum SearchQuality {
    Excellent,
    Good,
    Fair,
    Limited,
}

#[cw_serde]
pub struct AnalyticsSummaryResponse {
    pub summary: AnalyticsSummary,
    pub key_insights: Vec<KeyInsight>,
    pub actionable_recommendations: Vec<String>,
}

#[cw_serde]
pub struct AnalyticsSummary {
    pub scope: AnalyticsScope,
    pub period: AnalyticsPeriod,
    pub quick_stats: GlobalQuickStats,
    pub performance_overview: PerformanceOverview,
    pub trend_summary: TrendSummary,
}

#[cw_serde]
pub struct GlobalQuickStats {
    pub total_students: u32,
    pub active_students: u32,
    pub graduation_rate: u32,
    pub average_gpa: String,
    pub at_risk_percentage: u32,
}

#[cw_serde]
pub struct PerformanceOverview {
    pub top_performing_institutions: Vec<String>,
    pub improvement_leaders: Vec<String>,
    pub areas_needing_attention: Vec<String>,
}

#[cw_serde]
pub struct TrendSummary {
    pub enrollment_trend: TrendDirection,
    pub performance_trend: TrendDirection,
    pub completion_trend: TrendDirection,
    pub key_trend_drivers: Vec<String>,
}

#[cw_serde]
pub struct KeyInsight {
    pub insight_category: InsightCategory,
    pub title: String,
    pub description: String,
    pub supporting_metrics: Vec<String>,
    pub confidence_level: u32,
}

#[cw_serde]
pub enum InsightCategory {
    PerformancePattern,
    RiskIdentification,
    OpportunitySpotting,
    TrendAlert,
    BenchmarkComparison,
    PredictiveWarning,
}