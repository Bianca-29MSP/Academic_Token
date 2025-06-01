use cosmwasm_std::StdError;
use thiserror::Error;

#[derive(Error, Debug, PartialEq)]
pub enum ContractError {
    #[error("{0}")]
    Std(#[from] StdError),

    #[error("Unauthorized")]
    Unauthorized {},

    #[error("Invalid GPA format")]
    InvalidGPA {},

    #[error("Curriculum not found")]
    CurriculumNotFound {},

    #[error("Validation not found")]
    ValidationNotFound {},
}