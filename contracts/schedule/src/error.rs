// src/error.rs
use cosmwasm_std::StdError;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum ContractError {
    #[error("{0}")]
    Std(#[from] StdError),

    #[error("Unauthorized: only {owner} can execute this action")]
    Unauthorized { owner: String },

    #[error("Student record not found for ID: {student_id}")]
    StudentNotFound { student_id: String },

    #[error("Subject not found for ID: {subject_id}")]
    SubjectNotFound { subject_id: String },

    #[error("Curriculum not found for ID: {curriculum_id}")]
    CurriculumNotFound { curriculum_id: String },

    #[error("Academic path not found for ID: {path_id}")]
    AcademicPathNotFound { path_id: String },

    #[error("Invalid schedule configuration: {reason}")]
    InvalidScheduleConfig { reason: String },

    #[error("Prerequisites not met for subject: {subject_id}")]
    PrerequisitesNotMet { subject_id: String },

    #[error("Maximum subjects per semester exceeded: {max} allowed")]
    MaxSubjectsExceeded { max: u32 },

    #[error("IPFS content not found for link: {ipfs_link}")]
    IpfsContentNotFound { ipfs_link: String },

    #[error("Schedule generation failed: {reason}")]
    ScheduleGenerationFailed { reason: String },

    #[error("Academic path optimization failed: {reason}")]
    OptimizationFailed { reason: String },
}