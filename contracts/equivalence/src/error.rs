use cosmwasm_std::StdError;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum ContractError {
    #[error("{0}")]
    Std(#[from] StdError),

    #[error("Unauthorized")]
    Unauthorized {},

    #[error("Self equivalence not allowed for subject {subject_id}")]
    SelfEquivalenceNotAllowed { subject_id: String },

    #[error("Equivalence already exists between {source_id} and {target_id}")]
    EquivalenceAlreadyExists { source_id: String, target_id: String },

    #[error("Equivalence {equivalence_id} already analyzed")]
    EquivalenceAlreadyAnalyzed { equivalence_id: String },

    #[error("Invalid similarity percentage: {percentage}%")]
    InvalidSimilarityPercentage { percentage: u32 },

    #[error("Content integrity error for subject {subject_id}")]
    ContentIntegrityError { subject_id: String },

    #[error("Content processing error: {reason}")]
    ContentProcessingError { reason: String },

    #[error("Invalid content: {reason}")]
    InvalidContent { reason: String },

    #[error("Storage error: {reason}")]
    StorageError { reason: String },

    #[error("Language not supported: {language}")]
    LanguageNotSupported { language: String },

    #[error("Translation quality too low: {confidence}%")]
    TranslationQualityTooLow { confidence: u32 },

    #[error("Insufficient data for analysis")]
    InsufficientDataForAnalysis {},

    #[error("Analysis timeout")]
    AnalysisTimeout {},
}