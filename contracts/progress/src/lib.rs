// src/lib.rs
pub mod contract;
pub mod error;
pub mod msg;
pub mod state;
pub mod execute;
pub mod query;
pub mod analytics;

#[cfg(test)]
mod tests;

pub use crate::error::ContractError;
