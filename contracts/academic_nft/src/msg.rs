use cosmwasm_schema::{cw_serde, QueryResponses};
use crate::state::{NFTType, SubjectNFT, DegreeNFT, NFTMetadata, StudentNFTCollection, NFTStatistics, Config};

// Instantiate message
#[cw_serde]
pub struct InstantiateMsg {
    pub admin: Option<String>,
    pub minter: String,                    // AcademicNFT module address
    pub ipfs_gateway: String,
    pub collection_name: String,
    pub collection_symbol: String,
}

// Execute messages
#[cw_serde]
pub enum ExecuteMsg {
    /// Mint subject completion NFT
    MintSubjectNFT {
        student_id: String,
        subject_data: SubjectCompletionData,
        metadata: NFTMetadata,
        validation_hash: String,
    },
    
    /// Mint degree NFT
    MintDegreeNFT {
        student_id: String,
        degree_data: DegreeCompletionData,
        metadata: NFTMetadata,
        validation_hash: String,
        signatures: Vec<String>,
    },
    
    /// Transfer NFT ownership
    TransferNft {
        recipient: String,
        token_id: String,
    },
    
    /// Approve NFT for transfer
    Approve {
        spender: String,
        token_id: String,
        expires: Option<cw721::Expiration>,
    },
    
    /// Revoke approval
    Revoke {
        spender: String,
        token_id: String,
    },
    
    /// Approve all tokens for operator
    ApproveAll {
        operator: String,
        expires: Option<cw721::Expiration>,
    },
    
    /// Revoke all approvals for operator
    RevokeAll {
        operator: String,
    },
    
    /// Update NFT metadata (admin only)
    UpdateMetadata {
        token_id: String,
        metadata: NFTMetadata,
    },
    
    /// Burn NFT (for corrections - admin only)
    Burn {
        token_id: String,
        reason: String,
    },
    
    /// Update contract configuration
    UpdateConfig {
        admin: Option<String>,
        minter: Option<String>,
        ipfs_gateway: Option<String>,
    },
    
    /// Cache IPFS content for better performance
    CacheIpfsContent {
        ipfs_link: String,
        content: String,
    },
    
    /// Validate NFT data before minting (internal use)
    ValidateNFTData {
        nft_type: NFTType,
        data: String, // JSON serialized data
    },
}

// Query messages
#[cw_serde]
#[derive(QueryResponses)]
pub enum QueryMsg {
    /// Get contract configuration
    #[returns(ConfigResponse)]
    GetConfig {},
    
    /// Get NFT owner
    #[returns(cw721::OwnerOfResponse)]
    OwnerOf {
        token_id: String,
        include_expired: Option<bool>,
    },
    
    /// Get NFT approval info
    #[returns(cw721::ApprovalResponse)]
    Approval {
        token_id: String,
        spender: String,
        include_expired: Option<bool>,
    },
    
    /// Get all approvals for token
    #[returns(cw721::ApprovalsResponse)]
    Approvals {
        token_id: String,
        include_expired: Option<bool>,
    },
    
    /// Get all tokens owned by address
    #[returns(cw721::TokensResponse)]
    Tokens {
        owner: String,
        start_after: Option<String>,
        limit: Option<u32>,
    },
    
    /// Get all tokens in collection
    #[returns(cw721::TokensResponse)]
    AllTokens {
        start_after: Option<String>,
        limit: Option<u32>,
    },
    
    /// Get NFT info
    #[returns(cw721::NftInfoResponse<NFTMetadata>)]
    NftInfo {
        token_id: String,
    },
    
    /// Get all NFT info (owner + metadata)
    #[returns(cw721::AllNftInfoResponse<NFTMetadata>)]
    AllNftInfo {
        token_id: String,
        include_expired: Option<bool>,
    },
    
    /// Get collection info
    #[returns(cw721::ContractInfoResponse)]
    ContractInfo {},
    
    /// Get subject NFT details
    #[returns(SubjectNFTResponse)]
    GetSubjectNFT {
        token_id: String,
    },
    
    /// Get degree NFT details
    #[returns(DegreeNFTResponse)]
    GetDegreeNFT {
        token_id: String,
    },
    
    /// Get student's NFT collection
    #[returns(StudentCollectionResponse)]
    GetStudentCollection {
        student_id: String,
    },
    
    /// Get NFTs by student and type
    #[returns(NFTsByTypeResponse)]
    GetNFTsByType {
        student_id: String,
        nft_type: NFTType,
        limit: Option<u32>,
    },
    
    /// Get NFT statistics
    #[returns(StatisticsResponse)]
    GetStatistics {},
    
    /// Get NFTs by institution
    #[returns(InstitutionNFTsResponse)]
    GetNFTsByInstitution {
        institution_id: String,
        nft_type: Option<NFTType>,
        limit: Option<u32>,
        start_after: Option<String>,
    },
    
    /// Check if NFT exists
    #[returns(bool)]
    NFTExists {
        token_id: String,
    },
    
    /// Get cached IPFS content
    #[returns(String)]
    GetCachedContent {
        ipfs_link: String,
    },
    
    /// Verify NFT authenticity
    #[returns(VerificationResponse)]
    VerifyNFT {
        token_id: String,
    },
}

// Helper structs for messages
#[cw_serde]
pub struct SubjectCompletionData {
    pub subject_id: String,
    pub institution_id: String,
    pub course_id: String,
    pub subject_name: String,
    pub credits: u32,
    pub final_grade: u32,
    pub completion_date: String,
    pub semester: String,
    pub academic_year: String,
    pub instructor: Option<String>,
}

#[cw_serde]
pub struct DegreeCompletionData {
    pub institution_id: String,
    pub course_id: String,
    pub degree_name: String,
    pub degree_type: String,
    pub major: String,
    pub minor: Option<String>,
    pub graduation_date: String,
    pub final_gpa: String,
    pub total_credits: u32,
    pub honors: Option<String>,
}

// Response structs
#[cw_serde]
pub struct ConfigResponse {
    pub config: Config,
}

#[cw_serde]
pub struct SubjectNFTResponse {
    pub nft: Option<SubjectNFT>,
}

#[cw_serde]
pub struct DegreeNFTResponse {
    pub nft: Option<DegreeNFT>,
}

#[cw_serde]
pub struct StudentCollectionResponse {
    pub collection: Option<StudentNFTCollection>,
    pub total_value: String, // Could represent academic value/score
}

#[cw_serde]
pub struct NFTsByTypeResponse {
    pub nfts: Vec<String>, // Token IDs
    pub total_count: u32,
}

#[cw_serde]
pub struct StatisticsResponse {
    pub stats: NFTStatistics,
}

#[cw_serde]
pub struct InstitutionNFTsResponse {
    pub nfts: Vec<InstitutionNFTSummary>,
    pub total_count: u32,
}

#[cw_serde]
pub struct InstitutionNFTSummary {
    pub token_id: String,
    pub student_id: String,
    pub nft_type: NFTType,
    pub issued_date: String,
    pub name: String,
}

#[cw_serde]
pub struct VerificationResponse {
    pub is_valid: bool,
    pub issued_by: String,
    pub validation_hash: String,
    pub issue_date: u64,
    pub verification_details: Vec<String>,
}
