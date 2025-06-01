use cosmwasm_schema::cw_serde;
use cosmwasm_std::Addr;
use cw_storage_plus::{Item, Map};

// Contract configuration
#[cw_serde]
pub struct Config {
    pub admin: Addr,
    pub minter: Addr,           // Usually the AcademicNFT module address
    pub ipfs_gateway: String,
    pub collection_name: String,
    pub collection_symbol: String,
}

// NFT Types
#[cw_serde]
pub enum NFTType {
    SubjectCompletion,
    Degree,
    Certificate,
    Achievement,
}

// Subject NFT data
#[cw_serde]
pub struct SubjectNFT {
    pub token_id: String,
    pub student_id: String,
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
    pub nft_metadata: NFTMetadata,
}

// Degree NFT data
#[cw_serde]
pub struct DegreeNFT {
    pub token_id: String,
    pub student_id: String,
    pub institution_id: String,
    pub course_id: String,
    pub degree_name: String,
    pub degree_type: String, // Bachelor, Master, PhD, etc.
    pub major: String,
    pub minor: Option<String>,
    pub graduation_date: String,
    pub final_gpa: String,
    pub total_credits: u32,
    pub honors: Option<String>,
    pub validation_hash: String,
    pub signatures: Vec<String>,
    pub nft_metadata: NFTMetadata,
}

// Common NFT metadata structure
#[cw_serde]
pub struct NFTMetadata {
    pub name: String,
    pub description: String,
    pub image: String,           // IPFS link to NFT image
    pub external_url: Option<String>,
    pub animation_url: Option<String>,
    pub attributes: Vec<NFTAttribute>,
    pub ipfs_metadata_link: String, // Complete metadata in IPFS
    pub created_at: u64,
    pub updated_at: u64,
}

// NFT attributes for metadata
#[cw_serde]
pub struct NFTAttribute {
    pub trait_type: String,
    pub value: String,
    pub display_type: Option<String>,
}

// Student NFT collection summary
#[cw_serde]
pub struct StudentNFTCollection {
    pub student_id: String,
    pub subject_nfts: Vec<String>,    // Token IDs
    pub degree_nfts: Vec<String>,     // Token IDs
    pub total_credits: u32,
    pub current_gpa: String,
    pub completion_percentage: u32,
}

// NFT statistics
#[cw_serde]
pub struct NFTStatistics {
    pub total_subject_nfts: u64,
    pub total_degree_nfts: u64,
    pub total_students: u64,
    pub total_institutions: u64,
    pub nfts_by_type: Vec<(NFTType, u64)>,
}

// Validation data for NFT minting
#[cw_serde]
pub struct MintValidation {
    pub validator_address: String,
    pub validation_hash: String,
    pub timestamp: u64,
    pub signature: String,
}

// Storage keys
pub const CONFIG: Item<Config> = Item::new("config");
pub const NFT_COUNTER: Item<u64> = Item::new("nft_counter");

// Maps for NFT data
pub const SUBJECT_NFTS: Map<&str, SubjectNFT> = Map::new("subject_nfts");
pub const DEGREE_NFTS: Map<&str, DegreeNFT> = Map::new("degree_nfts");

// Student collections
pub const STUDENT_COLLECTIONS: Map<&str, StudentNFTCollection> = Map::new("student_collections");

// Token to owner mapping (for CW721 compatibility)
pub const TOKEN_OWNERS: Map<&str, Addr> = Map::new("token_owners");

// Token to approval mapping
pub const TOKEN_APPROVALS: Map<&str, Addr> = Map::new("token_approvals");

// Operator approvals (owner -> operator -> approved)
pub const OPERATOR_APPROVALS: Map<(&str, &str), bool> = Map::new("operator_approvals");

// Statistics
pub const NFT_STATS: Item<NFTStatistics> = Item::new("nft_stats");

// IPFS content cache for metadata
pub const IPFS_CACHE: Map<&str, String> = Map::new("ipfs_cache");

// Validation records
pub const VALIDATIONS: Map<&str, MintValidation> = Map::new("validations");
