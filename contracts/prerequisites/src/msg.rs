use cosmwasm_schema::{cw_serde, QueryResponses};
use crate::state::{PrerequisiteGroup, StudentRecord, VerificationResult};
use crate::ipfs::SubjectContent;

/// Instantiate message
#[cw_serde]
pub struct InstantiateMsg {
    pub owner: Option<String>,
}

/// Execute messages
#[cw_serde]
pub enum ExecuteMsg {
    /// Register prerequisites for a subject
    RegisterPrerequisites {
        subject_id: String,
        prerequisites: Vec<PrerequisiteGroup>,
    },
    
    /// Update student record with completed subjects
    UpdateStudentRecord {
        student_id: String,
        completed_subject: CompletedSubjectMsg,
    },
    
    /// Verify if a student can enroll in a subject
    VerifyEnrollment {
        student_id: String,
        subject_id: String,
    },
    
    /// Batch update prerequisites
    BatchRegisterPrerequisites {
        items: Vec<PrerequisiteRegistration>,
    },
    
    /// Update contract owner
    UpdateOwner {
        new_owner: String,
    },
    
    /// Cache IPFS content for better prerequisite analysis
    CacheIpfsContent {
        ipfs_link: String,
        content: SubjectContent,
    },
    
    /// Analyze prerequisite relationships using IPFS content
    AnalyzePrerequisiteRelationship {
        source_subject_id: String,
        target_subject_id: String,
        source_ipfs_link: String,
        target_ipfs_link: String,
    },
}

/// Query messages
#[cw_serde]
#[derive(QueryResponses)]
pub enum QueryMsg {
    /// Get prerequisites for a subject
    #[returns(PrerequisitesResponse)]
    GetPrerequisites { subject_id: String },
    
    /// Get student record
    #[returns(StudentRecordResponse)]
    GetStudentRecord { student_id: String },
    
    /// Check enrollment eligibility
    #[returns(VerificationResult)]
    CheckEligibility {
        student_id: String,
        subject_id: String,
    },
    
    /// Get verification history
    #[returns(VerificationHistoryResponse)]
    GetVerificationHistory {
        student_id: String,
        limit: Option<u32>,
    },
    
    /// Get contract state
    #[returns(StateResponse)]
    GetState {},
    
    /// Check if IPFS content is cached
    #[returns(IpfsCacheStatusResponse)]
    GetIpfsCacheStatus { ipfs_link: String },
    
    /// Get cached IPFS content
    #[returns(SubjectContent)]
    GetCachedContent { ipfs_link: String },
}

// Helper structs for messages
#[cw_serde]
pub struct CompletedSubjectMsg {
    pub subject_id: String,
    pub credits: u64,
    pub completion_date: String,
    pub grade: u32, // Grade as integer (e.g., 85.5 -> 8550 for precision)
    pub nft_token_id: String,
    pub ipfs_link: Option<String>, // Added IPFS link
}

#[cw_serde]
pub struct PrerequisiteRegistration {
    pub subject_id: String,
    pub prerequisites: Vec<PrerequisiteGroup>,
}

// Response structs
#[cw_serde]
pub struct PrerequisitesResponse {
    pub subject_id: String,
    pub prerequisites: Vec<PrerequisiteGroup>,
}

#[cw_serde]
pub struct StudentRecordResponse {
    pub record: StudentRecord,
}

#[cw_serde]
pub struct VerificationHistoryResponse {
    pub verifications: Vec<(String, VerificationResult)>,
}

#[cw_serde]
pub struct StateResponse {
    pub owner: String,
    pub total_subjects: u64,
    pub total_verifications: u64,
}

#[cw_serde]
pub struct IpfsCacheStatusResponse {
    pub is_cached: bool,
    pub ipfs_link: String,
}

#[cw_serde]
pub struct PrerequisiteAnalysisResponse {
    pub source_subject_id: String,
    pub target_subject_id: String,
    pub analysis: crate::ipfs::PrerequisiteAnalysis,
}