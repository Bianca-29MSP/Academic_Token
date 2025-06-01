use cosmwasm_schema::cw_serde;
use cosmwasm_std::Addr;
use cw_storage_plus::{Item, Map};
use crate::ipfs::SubjectContent;

/// Main contract state
#[cw_serde]
pub struct State {
    pub owner: Addr,
    pub total_subjects: u64,
    pub total_verifications: u64,
}

/// Prerequisite group for a subject
#[cw_serde]
pub struct PrerequisiteGroup {
    pub id: String,
    pub subject_id: String,
    pub group_type: GroupType,
    pub minimum_credits: u64,
    pub minimum_completed_subjects: u64,
    pub subject_ids: Vec<String>,
    pub logic: LogicType,
    pub priority: u32,
    pub confidence: u32,
    pub ipfs_link: Option<String>, // Added IPFS link
}

/// Type of prerequisite group
#[cw_serde]
pub enum GroupType {
    #[serde(rename = "all")]
    All,      // All subjects in the group must be completed
    #[serde(rename = "any")]
    Any,      // At least one subject must be completed
    #[serde(rename = "minimum")]
    Minimum,  // Minimum number of subjects must be completed
    #[serde(rename = "none")]
    None,     // No prerequisites
}

/// Logic type for combining groups
#[cw_serde]
pub enum LogicType {
    #[serde(rename = "and")]
    And,       // All groups must be satisfied
    #[serde(rename = "or")]
    Or,        // At least one group must be satisfied
    #[serde(rename = "xor")]
    Xor,       // Exactly one group must be satisfied
    #[serde(rename = "threshold")]
    Threshold, // Custom threshold logic
    #[serde(rename = "none")]
    None,      // No logic needed
}

/// Student completion record
#[cw_serde]
pub struct StudentRecord {
    pub student_id: String,
    pub completed_subjects: Vec<CompletedSubject>,
    pub total_credits: u64,
}

/// Completed subject information
#[cw_serde]
pub struct CompletedSubject {
    pub subject_id: String,
    pub credits: u64,
    pub completion_date: String,
    pub grade: u32, // Changed back to u32 (grade * 100, e.g., 85.5 -> 8550)
    pub nft_token_id: String,
    pub ipfs_link: Option<String>, // Added IPFS link
}

/// Result of prerequisite verification
#[cw_serde]
pub struct VerificationResult {
    pub can_enroll: bool,
    pub missing_prerequisites: Vec<String>,
    pub satisfied_groups: Vec<String>,
    pub unsatisfied_groups: Vec<String>,
    pub verification_timestamp: u64,
    pub details: String,
    pub used_ipfs_content: bool, // Added IPFS usage indicator
}

// Storage items
pub const STATE: Item<State> = Item::new("state");

// Maps
// subject_id -> Vec<PrerequisiteGroup>
pub const PREREQUISITES: Map<&str, Vec<PrerequisiteGroup>> = Map::new("prerequisites");

// student_address -> StudentRecord
pub const STUDENT_RECORDS: Map<&str, StudentRecord> = Map::new("student_records");

// verification_id -> VerificationResult
pub const VERIFICATIONS: Map<&str, VerificationResult> = Map::new("verifications");

// IPFS content cache: ipfs_link -> SubjectContent
pub const IPFS_CACHE: Map<&str, SubjectContent> = Map::new("ipfs_cache");
