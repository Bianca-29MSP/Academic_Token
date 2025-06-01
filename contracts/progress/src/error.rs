// src/error.rs
use cosmwasm_std::StdError;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum ContractError {
    #[error("{0}")]
    Std(#[from] StdError),

    #[error("Unauthorized: only {owner} can execute this action")]
    Unauthorized { owner: String },

    #[error("Student not found for ID: {student_id}")]
    StudentNotFound { student_id: String },

    #[error("Institution not found for ID: {institution_id}")]
    InstitutionNotFound { institution_id: String },

    #[error("Course not found for ID: {course_id}")]
    CourseNotFound { course_id: String },

    #[error("Progress record not found")]
    ProgressRecordNotFound,

    #[error("Invalid analytics configuration: {reason}")]
    InvalidAnalyticsConfig { reason: String },

    #[error("Analytics computation failed: {reason}")]
    AnalyticsComputationFailed { reason: String },

    #[error("Insufficient data for analysis: {reason}")]
    InsufficientData { reason: String },

    #[error("Invalid time range: start {start} must be before end {end}")]
    InvalidTimeRange { start: String, end: String },

    #[error("Dashboard generation failed: {reason}")]
    DashboardGenerationFailed { reason: String },

    #[error("Metric calculation failed for {metric}: {reason}")]
    MetricCalculationFailed { metric: String, reason: String },
}