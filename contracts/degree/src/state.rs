use cosmwasm_schema::cw_serde;
use cosmwasm_std::{Addr, Timestamp, Uint128};
use cw_storage_plus::{Item, Map};

#[cw_serde]
pub struct Config {
    pub admin: Addr,
    pub module_address: Addr,
}

#[cw_serde]
pub struct DegreeValidationRequest {
    pub student_id: String,
    pub curriculum_id: String,
    pub institution_id: String,
    pub request_date: Timestamp,
    pub expected_graduation_date: String,
}

#[cw_serde]
pub struct DegreeValidationResult {
    pub is_valid: bool,
    pub validation_score: String,
    pub message: String,
    pub requirements_met: Vec<String>,
    pub missing_requirements: Vec<String>,
    pub total_credits_completed: Uint128,
    pub gpa: String,
}

#[cw_serde]
pub struct CurriculumRequirements {
    pub curriculum_id: String,
    pub min_credits: Uint128,
    pub required_subjects: Vec<String>,
    pub min_gpa: Option<String>,
    pub additional_requirements: Vec<String>,
}

pub const CONFIG: Item<Config> = Item::new("config");
pub const CURRICULUM_REQUIREMENTS: Map<&str, CurriculumRequirements> = Map::new("curriculum_requirements");
pub const VALIDATION_CACHE: Map<(&str, &str), DegreeValidationResult> = Map::new("validation_cache");