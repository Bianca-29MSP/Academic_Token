use cosmwasm_schema::cw_serde;

/// IPFS metadata structure for academic NFTs
#[cw_serde]
pub struct IPFSMetadata {
    pub name: String,
    pub description: String,
    pub image: String,
    pub external_url: Option<String>,
    pub animation_url: Option<String>,
    pub attributes: Vec<IPFSAttribute>,
    pub background_color: Option<String>,
    pub youtube_url: Option<String>,
}

/// NFT attributes in IPFS metadata
#[cw_serde]
pub struct IPFSAttribute {
    pub trait_type: String,
    pub value: String,
    pub display_type: Option<String>,
    pub max_value: Option<u32>,
}

/// Subject completion metadata for IPFS
#[cw_serde]
pub struct SubjectIPFSMetadata {
    pub base_metadata: IPFSMetadata,
    pub academic_data: SubjectAcademicData,
    pub verification: VerificationData,
}

/// Academic data for subject NFT
#[cw_serde]
pub struct SubjectAcademicData {
    pub subject_id: String,
    pub subject_name: String,
    pub subject_code: String,
    pub institution_name: String,
    pub course_name: String,
    pub credits: u32,
    pub final_grade: u32,
    pub grade_scale: String, // "0-100", "A-F", etc.
    pub completion_date: String,
    pub semester: String,
    pub academic_year: String,
    pub instructor: Option<String>,
    pub syllabus_ipfs: Option<String>, // Link to full syllabus
    pub prerequisites: Vec<String>,
    pub learning_outcomes: Vec<String>,
}

/// Degree completion metadata for IPFS
#[cw_serde]
pub struct DegreeIPFSMetadata {
    pub base_metadata: IPFSMetadata,
    pub academic_data: DegreeAcademicData,
    pub verification: VerificationData,
}

/// Academic data for degree NFT
#[cw_serde]
pub struct DegreeAcademicData {
    pub degree_name: String,
    pub degree_type: String, // Bachelor, Master, PhD, etc.
    pub major: String,
    pub minor: Option<String>,
    pub institution_name: String,
    pub graduation_date: String,
    pub final_gpa: String,
    pub gpa_scale: String,
    pub total_credits: u32,
    pub honors: Option<String>,
    pub thesis_title: Option<String>,
    pub advisor: Option<String>,
    pub completed_subjects: Vec<CompletedSubjectSummary>,
    pub certifications: Vec<String>,
}

/// Summary of completed subject for degree metadata
#[cw_serde]
pub struct CompletedSubjectSummary {
    pub subject_id: String,
    pub subject_name: String,
    pub credits: u32,
    pub grade: u32,
    pub semester: String,
    pub nft_token_id: Option<String>,
}

/// Verification data for both subject and degree NFTs
#[cw_serde]
pub struct VerificationData {
    pub issued_by: String,
    pub issuer_authority: String,
    pub validation_hash: String,
    pub blockchain_network: String,
    pub contract_address: String,
    pub token_id: String,
    pub issue_timestamp: u64,
    pub signatures: Vec<AuthoritySignature>,
    pub verification_url: String,
}

/// Authority signature for verification
#[cw_serde]
pub struct AuthoritySignature {
    pub authority_name: String,
    pub authority_role: String,
    pub signature: String,
    pub public_key: String,
    pub timestamp: u64,
}

/// IPFS content for subject syllabus
#[cw_serde]
pub struct SubjectSyllabusContent {
    pub subject_info: SubjectInfo,
    pub content: SyllabusContent,
    pub multilingual_support: Option<MultilingualContent>,
}

/// Basic subject information
#[cw_serde]
pub struct SubjectInfo {
    pub subject_id: String,
    pub subject_code: String,
    pub subject_name: String,
    pub institution_id: String,
    pub course_id: String,
    pub credits: u32,
    pub level: String, // Undergraduate, Graduate, etc.
    pub department: String,
}

/// Detailed syllabus content
#[cw_serde]
pub struct SyllabusContent {
    pub description: String,
    pub objectives: Vec<String>,
    pub learning_outcomes: Vec<String>,
    pub topics: Vec<TopicDetail>,
    pub assessment_methods: Vec<AssessmentMethod>,
    pub bibliography: Vec<Reference>,
    pub prerequisites: Vec<String>,
    pub corequisites: Vec<String>,
    pub workload: WorkloadInfo,
}

/// Topic detail in syllabus
#[cw_serde]
pub struct TopicDetail {
    pub topic_name: String,
    pub description: String,
    pub week: Option<u32>,
    pub hours: Option<u32>,
    pub subtopics: Vec<String>,
}

/// Assessment method
#[cw_serde]
pub struct AssessmentMethod {
    pub method_type: String, // Exam, Project, Assignment, etc.
    pub weight_percentage: u32,
    pub description: String,
    pub due_date: Option<String>,
}

/// Reference/Bibliography entry
#[cw_serde]
pub struct Reference {
    pub reference_type: String, // Book, Article, Website, etc.
    pub title: String,
    pub authors: Vec<String>,
    pub publication_year: Option<u32>,
    pub publisher: Option<String>,
    pub isbn: Option<String>,
    pub url: Option<String>,
    pub is_required: bool,
}

/// Workload information
#[cw_serde]
pub struct WorkloadInfo {
    pub total_hours: u32,
    pub lecture_hours: u32,
    pub lab_hours: Option<u32>,
    pub study_hours: u32,
    pub project_hours: Option<u32>,
    pub weekly_hours: u32,
}

/// Multilingual content support
#[cw_serde]
pub struct MultilingualContent {
    pub primary_language: String,
    pub translations: Vec<LanguageTranslation>,
}

/// Language translation
#[cw_serde]
pub struct LanguageTranslation {
    pub language_code: String, // en, pt, es, fr, etc.
    pub subject_name: String,
    pub description: String,
    pub objectives: Vec<String>,
    pub topics: Vec<String>,
}

/// IPFS utility functions
impl IPFSMetadata {
    /// Create metadata for subject completion NFT
    pub fn for_subject(
        subject_name: &str,
        institution_name: &str,
        final_grade: u32,
        completion_date: &str,
        image_ipfs: &str,
    ) -> Self {
        Self {
            name: format!("{} - Completion Certificate", subject_name),
            description: format!(
                "Certificate of completion for {} at {}. Final grade: {}. Completed on: {}",
                subject_name, institution_name, final_grade, completion_date
            ),
            image: image_ipfs.to_string(),
            external_url: None,
            animation_url: None,
            background_color: Some("#1e3a8a".to_string()), // Academic blue
            youtube_url: None,
            attributes: vec![
                IPFSAttribute {
                    trait_type: "Type".to_string(),
                    value: "Subject Completion".to_string(),
                    display_type: None,
                    max_value: None,
                },
                IPFSAttribute {
                    trait_type: "Grade".to_string(),
                    value: final_grade.to_string(),
                    display_type: Some("number".to_string()),
                    max_value: Some(100),
                },
                IPFSAttribute {
                    trait_type: "Institution".to_string(),
                    value: institution_name.to_string(),
                    display_type: None,
                    max_value: None,
                },
                IPFSAttribute {
                    trait_type: "Completion Date".to_string(),
                    value: completion_date.to_string(),
                    display_type: Some("date".to_string()),
                    max_value: None,
                },
            ],
        }
    }

    /// Create metadata for degree NFT
    pub fn for_degree(
        degree_name: &str,
        institution_name: &str,
        final_gpa: &str,
        graduation_date: &str,
        honors: Option<&str>,
        image_ipfs: &str,
    ) -> Self {
        let mut attributes = vec![
            IPFSAttribute {
                trait_type: "Type".to_string(),
                value: "Academic Degree".to_string(),
                display_type: None,
                max_value: None,
            },
            IPFSAttribute {
                trait_type: "Degree".to_string(),
                value: degree_name.to_string(),
                display_type: None,
                max_value: None,
            },
            IPFSAttribute {
                trait_type: "Institution".to_string(),
                value: institution_name.to_string(),
                display_type: None,
                max_value: None,
            },
            IPFSAttribute {
                trait_type: "GPA".to_string(),
                value: final_gpa.to_string(),
                display_type: Some("number".to_string()),
                max_value: Some(4), // Assuming 4.0 scale
            },
            IPFSAttribute {
                trait_type: "Graduation Date".to_string(),
                value: graduation_date.to_string(),
                display_type: Some("date".to_string()),
                max_value: None,
            },
        ];

        if let Some(honors) = honors {
            attributes.push(IPFSAttribute {
                trait_type: "Honors".to_string(),
                value: honors.to_string(),
                display_type: None,
                max_value: None,
            });
        }

        Self {
            name: format!("{} - {}", degree_name, institution_name),
            description: format!(
                "Academic degree certificate for {} from {}. GPA: {}. Graduated: {}{}",
                degree_name,
                institution_name,
                final_gpa,
                graduation_date,
                honors.map(|h| format!(" with {}", h)).unwrap_or_default()
            ),
            image: image_ipfs.to_string(),
            external_url: None,
            animation_url: None,
            background_color: Some("#059669".to_string()), // Academic green
            youtube_url: None,
            attributes,
        }
    }
}

/// Helper functions for IPFS integration
pub mod ipfs_utils {
    // use super::*;
    use sha2::{Sha256, Digest};

    /// Generate IPFS hash for content verification
    pub fn generate_content_hash(content: &str) -> String {
        let mut hasher = Sha256::new();
        hasher.update(content.as_bytes());
        format!("{:x}", hasher.finalize())
    }

    /// Create verification URL for NFT
    pub fn create_verification_url(
        contract_address: &str,
        token_id: &str,
        network: &str,
    ) -> String {
        format!(
            "https://academic-token.org/verify/{}/{}/{}",
            network, contract_address, token_id
        )
    }

    /// Validate IPFS link format
    pub fn validate_ipfs_link(ipfs_link: &str) -> bool {
        ipfs_link.starts_with("ipfs://") || ipfs_link.starts_with("https://ipfs.io/ipfs/")
    }

    /// Convert IPFS hash to gateway URL
    pub fn ipfs_to_gateway_url(ipfs_hash: &str, gateway: &str) -> String {
        let hash = ipfs_hash.strip_prefix("ipfs://").unwrap_or(ipfs_hash);
        format!("{}/ipfs/{}", gateway.trim_end_matches('/'), hash)
    }
}
