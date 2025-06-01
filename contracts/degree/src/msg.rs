// contracts/degree/src/msg.rs

use cosmwasm_schema::{cw_serde, QueryResponses};
use cosmwasm_std::Uint128;

// ============================================================================
// MESSAGES - Aligned with AcademicToken's contract_types.go
// ============================================================================

#[cw_serde]
pub struct InstantiateMsg {
    pub module_address: String,
}

#[cw_serde]
pub enum ExecuteMsg {
    // Admin functions
    SetCurriculumRequirements {
        curriculum_id: String,
        min_credits: Uint128,
        required_subjects: Vec<String>,
        min_gpa: Option<String>,
        additional_requirements: Vec<String>,
    },
    
    // Main validation function called by the degree module
    ValidateDegreeRequirements {
        student_id: String,
        curriculum_id: String,
        institution_id: String,
        final_gpa: String,
        total_credits: u64,
        completed_subjects: Vec<String>, // List of subject IDs completed
        signatures: Vec<String>,
        requested_date: String,
    },
    
    // Clear cached validation (admin only)
    ClearValidationCache {
        student_id: String,
        curriculum_id: String,
    },
}

#[cw_serde]
#[derive(QueryResponses)]
pub enum QueryMsg {
    #[returns(crate::state::Config)]
    GetConfig {},
    
    #[returns(crate::state::CurriculumRequirements)]
    GetCurriculumRequirements { curriculum_id: String },
    
    #[returns(crate::state::DegreeValidationResult)]
    GetValidationResult { 
        student_id: String,
        curriculum_id: String,
    },
    
    #[returns(Vec<String>)]
    GetMissingRequirements {
        student_id: String,
        curriculum_id: String,
        completed_subjects: Vec<String>,
    },
}

// Response types matching contract_types.go
#[cw_serde]
pub struct ValidateDegreeRequirementsResponse {
    pub is_valid: bool,
    pub message: String,
    pub degree_type: String,
    pub curriculum_version: String,
    pub validation_hash: String,
    pub requirements_met: Vec<String>,
    pub missing_requirements: Vec<String>,
}