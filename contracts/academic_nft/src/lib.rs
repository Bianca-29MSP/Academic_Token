pub mod contract;
pub mod error;
pub mod execute;
pub mod ipfs;
pub mod msg;
pub mod query;
pub mod state;

#[cfg(test)]
mod tests;

pub use crate::error::ContractError;

// Re-export important types for easier access
pub use crate::msg::{ExecuteMsg, InstantiateMsg, QueryMsg};
pub use crate::state::{Config, NFTMetadata, SubjectNFT, DegreeNFT};
