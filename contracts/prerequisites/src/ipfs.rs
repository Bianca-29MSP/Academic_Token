use cosmwasm_schema::cw_serde;
use cosmwasm_std::{Deps, DepsMut, StdError, StdResult};
use serde::{Deserialize, Serialize};
use crate::state::IPFS_CACHE;

/// Complete subject content structure retrieved from IPFS for prerequisites analysis
#[cw_serde]
pub struct SubjectContent {
    pub title: String,
    pub code: String,
    pub description: String,
    pub objectives: Vec<String>,
    pub topics: Vec<TopicUnit>,
    pub competencies: Vec<String>,
    pub methodology: Vec<String>,
    pub assessment_methods: Vec<String>,
    pub workload_hours: u64,
    pub credits: u64,
    pub prerequisites_description: String,
    pub learning_outcomes: Vec<String>,
    pub knowledge_areas: Vec<String>,
    pub difficulty_level: String,
    pub language: String,
    pub institution: String,
    pub professor: String,
    pub semester: String,
    pub year: String,
    pub bibliography: Vec<String>,
    pub supplementary_materials: Vec<String>,
    pub practical_activities: Vec<String>,
    pub theoretical_hours: u64,
    pub practical_hours: u64,
    pub content_hash: String,
}

/// Topic or content unit within the subject
#[cw_serde]
pub struct TopicUnit {
    pub unit_number: u32,
    pub title: String,
    pub description: String,
    pub hours: u64,
    pub keywords: Vec<String>,
    pub subtopics: Vec<String>,
    pub learning_objectives: Vec<String>,
    pub required_knowledge: Vec<String>,
}

/// Fetch subject content from IPFS cache or return error
pub fn fetch_ipfs_content(deps: Deps, ipfs_link: &str) -> StdResult<SubjectContent> {
    // Try to load from cache first
    match IPFS_CACHE.may_load(deps.storage, ipfs_link)? {
        Some(content) => Ok(content),
        None => Err(StdError::generic_err(format!(
            "IPFS content not found in cache for link: {}. Please ensure content is cached first.", 
            ipfs_link
        ))),
    }
}

/// Cache IPFS content in the contract storage
pub fn cache_ipfs_content(
    deps: DepsMut,
    ipfs_link: &str,
    content: &SubjectContent,
) -> StdResult<()> {
    IPFS_CACHE.save(deps.storage, ipfs_link, content)
}

/// Check if IPFS content exists in cache
pub fn is_content_cached(deps: Deps, ipfs_link: &str) -> StdResult<bool> {
    Ok(IPFS_CACHE.may_load(deps.storage, ipfs_link)?.is_some())
}

/// Verify content integrity using hash
pub fn verify_content_integrity(content: &SubjectContent, expected_hash: &str) -> bool {
    calculate_content_hash(content) == expected_hash
}

/// Calculate hash of content for integrity verification
pub fn calculate_content_hash(content: &SubjectContent) -> String {
    use sha2::{Digest, Sha256};
    
    let content_str = format!(
        "{}|{}|{}|{}|{}|{}|{}",
        content.title,
        content.code,
        content.description,
        content.objectives.join(","),
        content.competencies.join(","),
        content.knowledge_areas.join(","),
        content.workload_hours
    );
    
    let mut hasher = Sha256::new();
    hasher.update(content_str.as_bytes());
    format!("{:x}", hasher.finalize())
}

/// Extract prerequisite information from IPFS content
pub fn extract_prerequisites_from_content(content: &SubjectContent) -> Vec<String> {
    let mut prerequisites = Vec::new();
    
    // Extract from prerequisites description
    if !content.prerequisites_description.is_empty() {
        // Simple parsing - could be more sophisticated
        let desc = content.prerequisites_description.to_lowercase();
        for competency in &content.competencies {
            if desc.contains(&competency.to_lowercase()) {
                prerequisites.push(competency.clone());
            }
        }
    }
    
    // Extract from required knowledge in topics
    for topic in &content.topics {
        prerequisites.extend(topic.required_knowledge.clone());
    }
    
    prerequisites.sort();
    prerequisites.dedup();
    prerequisites
}

/// Analyze prerequisite relationships between subjects using IPFS content
pub fn analyze_prerequisite_relationship(
    source_content: &SubjectContent,
    target_content: &SubjectContent,
) -> PrerequisiteAnalysis {
    let mut analysis = PrerequisiteAnalysis {
        is_prerequisite: false,
        confidence_score: 0,
        knowledge_overlap: 0.0,
        competency_overlap: 0.0,
        topic_progression: false,
        difficulty_progression: false,
        reasons: Vec::new(),
    };
    
    // Check knowledge areas overlap
    let source_areas: std::collections::HashSet<_> = source_content.knowledge_areas.iter().collect();
    let target_areas: std::collections::HashSet<_> = target_content.knowledge_areas.iter().collect();
    let overlap: Vec<_> = source_areas.intersection(&target_areas).collect();
    analysis.knowledge_overlap = overlap.len() as f64 / source_areas.len().max(1) as f64;
    
    // Check competencies overlap
    let source_comp: std::collections::HashSet<_> = source_content.competencies.iter().collect();
    let target_comp: std::collections::HashSet<_> = target_content.competencies.iter().collect();
    let comp_overlap: Vec<_> = source_comp.intersection(&target_comp).collect();
    analysis.competency_overlap = comp_overlap.len() as f64 / source_comp.len().max(1) as f64;
    
    // Check difficulty progression
    analysis.difficulty_progression = is_difficulty_progression(&source_content.difficulty_level, &target_content.difficulty_level);
    
    // Check if source competencies are mentioned in target prerequisites
    for competency in &source_content.competencies {
        if target_content.prerequisites_description.to_lowercase().contains(&competency.to_lowercase()) {
            analysis.reasons.push(format!("Competency '{}' is explicitly mentioned as prerequisite", competency));
        }
    }
    
    // Calculate confidence score
    let mut score = 0;
    if analysis.knowledge_overlap > 0.3 { score += 25; }
    if analysis.competency_overlap > 0.2 { score += 30; }
    if analysis.difficulty_progression { score += 20; }
    if !analysis.reasons.is_empty() { score += 25; }
    
    analysis.confidence_score = score;
    analysis.is_prerequisite = score >= 50;
    
    analysis
}

/// Result of prerequisite relationship analysis
#[cw_serde]
pub struct PrerequisiteAnalysis {
    pub is_prerequisite: bool,
    pub confidence_score: u32,
    pub knowledge_overlap: f64,
    pub competency_overlap: f64,
    pub topic_progression: bool,
    pub difficulty_progression: bool,
    pub reasons: Vec<String>,
}

/// Check if there's difficulty progression between subjects
fn is_difficulty_progression(source_level: &str, target_level: &str) -> bool {
    let difficulty_order = ["basic", "easy", "medium", "intermediate", "hard", "advanced"];
    
    let source_index = difficulty_order.iter().position(|&x| x == source_level.to_lowercase());
    let target_index = difficulty_order.iter().position(|&x| x == target_level.to_lowercase());
    
    match (source_index, target_index) {
        (Some(s), Some(t)) => s < t,
        _ => false,
    }
}
