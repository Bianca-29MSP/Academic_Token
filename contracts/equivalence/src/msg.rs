use cosmwasm_schema::{cw_serde, QueryResponses};
use crate::state::{
    Equivalence, AnalysisResult, TransferRequest, SubjectInfo, 
    EquivalenceType, AnalysisMethod, State, QualityMetrics,
};
use crate::ipfs::{MultilingualSyllabusContent, Language, ContentQualityAssessment}; 

/// Instantiate message
#[cw_serde]
pub struct InstantiateMsg {
    pub owner: Option<String>,
    pub auto_approval_threshold: Option<u32>, // Default: 85%
}

/// Execute messages
#[cw_serde]
pub enum ExecuteMsg {
    /// Register a new equivalence between subjects
    RegisterEquivalence {
        source_subject: SubjectInfo,
        target_subject: SubjectInfo,
        analysis_method: AnalysisMethod,
        notes: Option<String>,
    },
    
    /// Analyze equivalence using AI/algorithm
    AnalyzeEquivalence {
        equivalence_id: String,
        force_reanalysis: Option<bool>,
    },
    
    /// Manually approve/reject equivalence
    ApproveEquivalence {
        equivalence_id: String,
        equivalence_type: EquivalenceType,
        similarity_percentage: u32,
        notes: Option<String>,
    },

    CacheIpfsContent {
        ipfs_link: String,
        content: MultilingualSyllabusContent,
    },
    
    /// Enhanced analysis with options
    AnalyzeEquivalenceEnhanced {
        equivalence_id: String,
        analysis_options: AnalysisOptions,
    }, 
    /// Submit student transfer request
    SubmitTransferRequest {
        student_id: String,
        source_institution: String,
        target_institution: String,
        completed_subjects: Vec<String>,
        requested_equivalences: Vec<String>,
    },
    
    /// Process transfer request
    ProcessTransferRequest {
        transfer_id: String,
        approved_equivalences: Vec<String>,
        notes: Option<String>,
    },
    
    /// Batch register multiple equivalences
    BatchRegisterEquivalences {
        equivalences: Vec<EquivalenceRegistration>,
    },
    
    /// Update contract configuration
    UpdateConfig {
        auto_approval_threshold: Option<u32>,
        new_owner: Option<String>,
    },
}

/// Query messages
#[cw_serde]
#[derive(QueryResponses)]
pub enum QueryMsg {
    /// Get contract state
    #[returns(StateResponse)]
    GetState {},
    
    /// Get equivalence by ID
    #[returns(EquivalenceResponse)]
    GetEquivalence { equivalence_id: String },
    
    /// Find equivalence between two subjects
    #[returns(EquivalenceResponse)]
    FindEquivalence {
        source_subject_id: String,
        target_subject_id: String,
    },
    
    /// List equivalences for an institution
    #[returns(EquivalencesResponse)]
    ListEquivalencesByInstitution {
        institution_id: String,
        limit: Option<u32>,
        start_after: Option<String>,
    },
    
    /// Get analysis result
    #[returns(AnalysisResponse)]
    GetAnalysisResult { analysis_id: String },
    
    /// Get detailed analysis result
    #[returns(crate::state::DetailedAnalysisResult)]
    GetDetailedAnalysisResult { analysis_id: String },
    
    /// Get transfer request
    #[returns(TransferResponse)]
    GetTransferRequest { transfer_id: String },
    
    /// List student transfers
    #[returns(TransfersResponse)]
    ListStudentTransfers {
        student_id: String,
        limit: Option<u32>,
    },
    
    /// Check if subjects are equivalent
    #[returns(EquivalenceCheckResponse)]
    CheckEquivalence {
        source_subject_id: String,
        target_subject_id: String,
        minimum_similarity: Option<u32>,
    },
    
    /// Get equivalence statistics
    #[returns(StatisticsResponse)]
    GetStatistics {
        institution_id: Option<String>,
    },
    
    /// Debug: List all equivalence IDs
    #[returns(Vec<String>)]
    DebugEquivalences {},
}

// Helper structs for messages
#[cw_serde]
pub struct EquivalenceRegistration {
    pub source_subject: SubjectInfo,
    pub target_subject: SubjectInfo,
    pub analysis_method: AnalysisMethod,
    pub notes: Option<String>,
}


#[cw_serde]
pub struct EnhancedAnalysisResponse {
    pub equivalence_id: String,
    pub analysis_result: AnalysisResult,
    pub quality_metrics: QualityMetrics,
    pub language_compatibility: u32,
    pub recommendations: Vec<String>,
    pub source_quality_assessment: Option<ContentQualityAssessment>,
    pub target_quality_assessment: Option<ContentQualityAssessment>,
    pub confidence_factors: ConfidenceFactors,
}

#[cw_serde]
pub struct ConfidenceFactors {
    pub content_available: bool,
    pub languages_match: bool,
    pub quality_threshold_met: bool,
    pub auto_translated: bool,
    pub comprehensive_analysis: bool,
}
#[cw_serde]
pub struct AnalysisOptions {
    pub use_semantic_analysis: bool,
    pub check_prerequisites: bool,
    pub analyze_learning_outcomes: bool,
    pub minimum_content_depth: u32,
    pub language_preference: Option<Language>,
}

// Response structs
#[cw_serde]
pub struct StateResponse {
    pub state: State,
}

#[cw_serde]
pub struct EquivalenceResponse {
    pub equivalence: Option<Equivalence>,
}

#[cw_serde]
pub struct EquivalencesResponse {
    pub equivalences: Vec<Equivalence>,
    pub total: u64,
}

#[cw_serde]
pub struct AnalysisResponse {
    pub analysis: Option<AnalysisResult>,
}

#[cw_serde]
pub struct TransferResponse {
    pub transfer: Option<TransferRequest>,
}

#[cw_serde]
pub struct TransfersResponse {
    pub transfers: Vec<TransferRequest>,
    pub total: u64,
}

#[cw_serde]
pub struct EquivalenceCheckResponse {
    pub is_equivalent: bool,
    pub equivalence_type: Option<EquivalenceType>,
    pub similarity_percentage: u32,
    pub equivalence_id: Option<String>,
}

#[cw_serde]
pub struct StatisticsResponse {
    pub total_equivalences: u64,
    pub approved_equivalences: u64,
    pub pending_equivalences: u64,
    pub average_similarity: u32,
    pub total_transfers: u64,
    pub successful_transfers: u64,
}
