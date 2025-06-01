pub mod contract;
pub mod error;
pub mod msg;
pub mod state;

pub use crate::error::ContractError;

// Re-export important types for easier access
pub use crate::msg::{ExecuteMsg, InstantiateMsg, QueryMsg, ValidateDegreeRequirementsResponse};
pub use crate::state::{Config, CurriculumRequirements, DegreeValidationResult};
