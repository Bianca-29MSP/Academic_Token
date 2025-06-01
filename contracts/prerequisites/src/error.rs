use cosmwasm_std::StdError;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum ContractError {
    #[error("{0}")]
    Std(#[from] StdError),

    #[error("Unauthorized")]
    Unauthorized {},

    #[error("Invalid prerequisite: {reason}")]
    InvalidPrerequisite { reason: String },

    #[error("Subject already completed: {subject_id}")]
    SubjectAlreadyCompleted { subject_id: String },

    #[error("Student not found: {student_id}")]
    StudentNotFound { student_id: String },

    #[error("Subject not found: {subject_id}")]
    SubjectNotFound { subject_id: String },

    #[error("Invalid confidence value: must be between 0 and 1")]
    InvalidConfidence {},

    #[error("Custom Error val: {val:?}")]
    CustomError { val: String },
    
    #[error("IPFS content not found: {ipfs_link}")]
    IpfsContentNotFound { ipfs_link: String },
    
    #[error("Invalid data: {reason}")]
    InvalidData { reason: String },
    
    #[error("Storage error: {reason}")]
    StorageError { reason: String },
}