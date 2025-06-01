use cosmwasm_std::StdError;
use thiserror::Error;

#[derive(Error, Debug, PartialEq)]
pub enum ContractError {
    #[error("{0}")]
    Std(#[from] StdError),

    #[error("Unauthorized: {reason}")]
    Unauthorized { reason: String },

    #[error("NFT not found: {token_id}")]
    NFTNotFound { token_id: String },

    #[error("Student not found: {student_id}")]
    StudentNotFound { student_id: String },

    #[error("Invalid NFT type: {nft_type}")]
    InvalidNFTType { nft_type: String },

    #[error("Subject already completed: {subject_id}")]
    SubjectAlreadyCompleted { subject_id: String },

    #[error("Degree already issued: {degree_id}")]
    DegreeAlreadyIssued { degree_id: String },

    #[error("Invalid grade: {grade}")]
    InvalidGrade { grade: String },

    #[error("IPFS content not found: {ipfs_link}")]
    IPFSContentNotFound { ipfs_link: String },

    #[error("Invalid signature: {reason}")]
    InvalidSignature { reason: String },

    #[error("Token already exists: {token_id}")]
    TokenAlreadyExists { token_id: String },

    #[error("Invalid token URI: {uri}")]
    InvalidTokenURI { uri: String },

    #[error("Validation failed: {reason}")]
    ValidationFailed { reason: String },

    #[error("Contract integration error: {reason}")]
    ContractIntegrationError { reason: String },

    #[error("Metadata parsing error: {reason}")]
    MetadataParsingError { reason: String },
}
