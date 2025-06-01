use cosmwasm_schema::cw_serde;
use cosmwasm_std::Addr;
use cw_storage_plus::{Item, Map};

// ===================== CORE STATE =====================

#[cw_serde]
pub struct State {
    pub owner: Addr,
    pub total_equivalences: u64,
    pub total_analyses: u64,
    pub auto_approval_threshold: u32,
}

// ===================== SUBJECT TYPES =====================

#[cw_serde]
pub struct SubjectInfo {
    pub subject_id: String,
    pub title: String,
    pub institution_id: String,
    pub credits: u32,
    pub ipfs_link: String,
    pub content_hash: String,
    pub metadata: SubjectMetadata,
}

#[cw_serde]
pub struct SubjectMetadata {
    pub level: AcademicLevel,
    pub department: String,
    pub workload_hours: u32,
    pub semester: u32,
    pub language: String,
}

#[cw_serde]
pub enum AcademicLevel {
    Undergraduate,
    Graduate,
    Postgraduate,
}

// ===================== EQUIVALENCE TYPES =====================

#[cw_serde]
pub struct Equivalence {
    pub id: String,
    pub source_subject: SubjectInfo,
    pub target_subject: SubjectInfo,
    pub equivalence_type: EquivalenceType,
    pub similarity_percentage: u32,
    pub status: EquivalenceStatus,
    pub analysis_method: AnalysisMethod,
    pub created_timestamp: u64,
    pub approved_timestamp: Option<u64>,
    pub approved_by: Option<Addr>,
    pub notes: String,
    pub confidence_score: u32,
}

#[cw_serde]
pub enum EquivalenceType {
    Full,
    Partial,
    Conditional,
    None,
}

#[cw_serde]
pub enum EquivalenceStatus {
    Pending,
    Analyzing,
    UnderReview,
    Approved,
    Rejected,
}

#[cw_serde]
pub enum AnalysisMethod {
    Automatic,
    Manual,
    Hybrid,
    Institutional,
}

// ===================== ANALYSIS TYPES =====================

#[cw_serde]
pub struct AnalysisResult {
    pub equivalence_id: String,
    pub content_similarity: u32,
    pub structure_similarity: u32,
    pub credit_compatibility: u32,
    pub level_compatibility: u32,
    pub overall_score: u32,
    pub recommendation: EquivalenceType,
    pub analysis_details: String,
    pub analyzed_timestamp: u64,
}

#[cw_serde]
pub struct DetailedAnalysisResult {
    pub base_result: AnalysisResult,
    pub quality_metrics: QualityMetrics,
    pub language_compatibility: u32,
    pub content_depth_score: u32,
    pub prerequisite_alignment: u32,
    pub learning_outcome_alignment: u32,
    pub workload_compatibility: u32,
    pub bibliography_overlap: u32,
    pub recommendations: Vec<String>,
}

#[cw_serde]
pub struct QualityMetrics {
    pub overall_quality: u32,
    pub content_completeness: u32,
    pub analysis_confidence: u32,
    pub data_availability: u32,
    pub language_quality: u32,
}

// ===================== TRANSFER TYPES =====================

#[cw_serde]
pub struct TransferRequest {
    pub id: String,
    pub student_id: String,
    pub source_institution: String,
    pub target_institution: String,
    pub completed_subjects: Vec<String>,
    pub requested_equivalences: Vec<String>,
    pub approved_equivalences: Vec<String>,
    pub status: TransferStatus,
    pub submitted_timestamp: u64,
    pub processed_timestamp: Option<u64>,
    pub processed_by: Option<Addr>,
    pub notes: String,
}

#[cw_serde]
pub enum TransferStatus {
    Pending,
    Processing,
    PartiallyApproved,
    FullyApproved,
    Rejected,
}

// ===================== STORAGE =====================

pub const STATE: Item<State> = Item::new("state");
pub const EQUIVALENCES: Map<&str, Equivalence> = Map::new("equivalences");
pub const EQUIVALENCE_INDEX: Map<(&str, &str), String> = Map::new("equivalence_index");
pub const ANALYSIS_RESULTS: Map<&str, AnalysisResult> = Map::new("analysis_results");
pub const DETAILED_ANALYSIS_RESULTS: Map<&str, DetailedAnalysisResult> = Map::new("detailed_analysis_results");
pub const TRANSFER_REQUESTS: Map<&str, TransferRequest> = Map::new("transfer_requests");
pub const STUDENT_TRANSFERS: Map<&str, Vec<String>> = Map::new("student_transfers");

// IPFS Cache - now stores multilingual content
use crate::ipfs::MultilingualSyllabusContent;
pub const IPFS_CACHE: Map<&str, MultilingualSyllabusContent> = Map::new("ipfs_cache");